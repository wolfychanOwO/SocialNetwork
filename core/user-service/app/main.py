from fastapi import FastAPI
from app.api.routes import router

app = FastAPI(
    title="User Service",
    description="Сервис управления пользователями с аутентификацией через Keycloak",
    version="1.0.0"
)

app.include_router(router, prefix="/api/v1")

from fastapi.middleware.cors import CORSMiddleware

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)

