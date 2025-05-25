from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    KEYCLOAK_URL: str = "http://keycloak:8080"
    KEYCLOAK_REALM: str = "testing"
    KEYCLOAK_CLIENT_ID: str = "user-service"
    KEYCLOAK_CLIENT_SECRET: str = "zfmm68hVZ45NSlumjejPw6Xhw56cToah"
    KEYCLOAK_ADMIN_USER: str = "admin"
    KEYCLOAK_ADMIN_PASSWORD: str = "admin"
    
    # Kafka settings
    KAFKA_BOOTSTRAP_SERVERS: str = "broker:29092"
    KAFKA_USER_TOPIC: str = "users"

    class Config:
        env_file = ".env"

settings = Settings()
