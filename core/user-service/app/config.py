from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    KEYCLOAK_URL: str = "http://keycloak:8080"
    KEYCLOAK_REALM: str = "testing"
    KEYCLOAK_CLIENT_ID: str = "user-service"
    KEYCLOAK_CLIENT_SECRET: str = "Es0IqbYVisWsjKLQ3EqgnWY03ydu7fRR"
    KEYCLOAK_ADMIN_USER: str = "admin"
    KEYCLOAK_ADMIN_PASSWORD: str = "admin"

    class Config:
        env_file = ".env"

settings = Settings()

