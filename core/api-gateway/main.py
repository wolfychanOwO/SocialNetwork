from fastapi import FastAPI, Request, HTTPException
import httpx
import jwt

app = FastAPI()
KEYCLOAK_URL = "http://keycloak:8080/realms/MyApp/protocol/openid-connect/certs"

SERVICES = {
    "user": "http://user-service:8001",
}

PUBLIC_ENDPOINTS = {
    "user": ["/api/v1/login", "/api/v1/register"]
}


def verify_token(token: str):
    try:
        headers = jwt.get_unverified_header(token)
        kid = headers["kid"]
        certs = httpx.get(KEYCLOAK_URL).json()["keys"]
        public_key = next(k["x5c"][0] for k in certs if k["kid"] == kid)
        
        payload = jwt.decode(token, public_key, algorithms=["RS256"])
        return payload
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")

@app.post("/{service}/{path:path}")
async def proxy(service: str, path: str, request: Request):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail="Service not found")

    if service in PUBLIC_ENDPOINTS and f"/{path}" in PUBLIC_ENDPOINTS[service]:
        token_data = None
    else:
        token = request.headers.get("Authorization")
        if not token or not token.startswith("Bearer "):
            raise HTTPException(status_code=401, detail="Token required")

        token_data = verify_token(token.split("Bearer ")[1])

    async with httpx.AsyncClient() as client:
        response = await client.post(f"{SERVICES[service]}/{path}", json=await request.json())
        return response.json()
    

@app.get("/{service}/{path:path}")
async def proxy(service: str, path: str, request: Request):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail="Service not found")

    token = request.headers.get("Authorization")
    if not token or not token.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Token required")

    token_data = verify_token(token.split("Bearer ")[1])

    async with httpx.AsyncClient() as client:
        headers = {"Authorization": request.headers.get("Authorization")}
        response = await client.get(f"{SERVICES[service]}/{path}", headers=headers)
        return response.json()

