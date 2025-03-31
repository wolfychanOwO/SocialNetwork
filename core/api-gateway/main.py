from fastapi import FastAPI, Request, HTTPException
import httpx
from fastapi.middleware.cors import CORSMiddleware
import grpc
from factory import GrpcFactory

app = FastAPI()


app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

SERVICES = {
    "user": "http://user-service:8001",
    "posts": "http://post-service:8002",
}


PUBLIC_ENDPOINTS = {
    "user": [
        "/api/v1/login",
        "/api/v1/register",
        "/api/v1/health"
    ],
    "posts": [
        "/api/v1/health",
        "list_comments",
    ]
}


async def validate_token(auth_header: str) -> dict:
    async with httpx.AsyncClient() as client:
        resp = await client.get("http://user-service:8001/api/v1/validate", headers={"Authorization": auth_header})

        if resp.status_code != 200:
            raise HTTPException(status_code=401, detail="Invalid token")
        return resp.json()

def is_public_endpoint(service: str, path: str) -> bool:
    if service not in PUBLIC_ENDPOINTS:
        return False
    
    full_path = f"/{path}" if not path.startswith("/") else path
    
    return any(
        full_path == public_path or 
        full_path.startswith(public_path + "/")
        for public_path in PUBLIC_ENDPOINTS[service]
    )

@app.api_route("/{service}/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
async def proxy(service: str, path: str, request: Request):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail="Service not found")

    is_public = is_public_endpoint(service, path)
    
    if not is_public:
        auth_header = request.headers.get("Authorization")
        if not auth_header:
            raise HTTPException(
                status_code=401,
                detail="Authorization header missing",
                headers={"WWW-Authenticate": "Bearer"}
            )

    target_url = f"{SERVICES[service]}/{path}"
    headers = {
        k: v for k, v in request.headers.items()
        if k.lower() not in ["host", "content-length"]
    }

    try:
        async with httpx.AsyncClient() as client:
            if request.method in ["POST", "PUT", "PATCH", "DELETE"]:
                content_type = request.headers.get("content-type", "")
                
                if "application/x-www-form-urlencoded" in content_type:
                    form_data = await request.form()
                    response = await client.request(
                        request.method,
                        target_url,
                        data=dict(form_data),
                        headers=headers,
                        timeout=30.0
                    )
                else:
                    try:
                        json_data = await request.json()
                    except:
                        json_data = None
                    
                    print('test')

                    if service == "posts":
                        if path not in PUBLIC_ENDPOINTS["posts"]:
                            print(path)
                            user = await validate_token(auth_header)
                            json_data['user_id'] = str(user['id'])
                            print('test')
                        grpc_client = GrpcFactory()
                        grpc_client = grpc_client.get_client("posts")
                        method = getattr(grpc_client, path, None)
                        print(method)
                        if method is None:
                            raise HTTPException(status_code=404, detail="method not found")
                        
                        response = await method(json_data)
                        print(response)
                        return response

                    response = await client.request(
                        request.method,
                        target_url,
                        json=json_data,
                        headers=headers,
                        timeout=30.0
                    )
            else:

                

                response = await client.request(
                    request.method,
                    target_url,
                    headers=headers,
                    timeout=30.0
                )

            return response.json()

    except httpx.ConnectError:
        raise HTTPException(
            status_code=503,
            detail=f"Service {service} unavailable"
        )
    except httpx.TimeoutException:
        raise HTTPException(
            status_code=504,
            detail=f"Service {service} timeout"
        )
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Internal server error: {str(e)}"
        )

@app.get("/health")
async def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
