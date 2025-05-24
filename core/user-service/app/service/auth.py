# app/service/auth.py
import requests
from jose import jwt
from app.config import settings

from app.arch.methods import UserMethod
from app.models.schema import UserCreate, UserLogin

from passlib.context import CryptContext
from app.service.kafka_producer import KafkaProducerService
import asyncio

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

class AuthService:
    def __init__(self, method: UserMethod, broker: KafkaProducerService):
        self.method = method
        self.keycloak_openid_config = self._get_openid_config()
        self.broker = broker
        

    def _get_openid_config(self):
        url = f"{settings.KEYCLOAK_URL}/realms/{settings.KEYCLOAK_REALM}/.well-known/openid-configuration"
        return requests.get(url.replace('localhost', 'keycloak', 1)).json()

    def _get_admin_token(self):
        url = f"{settings.KEYCLOAK_URL}/realms/master/protocol/openid-connect/token"
        data = {
            "client_id": "admin-cli",
            "username": settings.KEYCLOAK_ADMIN_USER,
            "password": settings.KEYCLOAK_ADMIN_PASSWORD,
            "grant_type": "password"
        }
        response = requests.post(url.replace('localhost', 'keycloak', 1), data=data)
        return response.json()["access_token"]

    def register_user(self, user_data: UserCreate):
        admin_token = self._get_admin_token()
        url = f"{settings.KEYCLOAK_URL}/admin/realms/{settings.KEYCLOAK_REALM}/users"
        headers = {"Authorization": f"Bearer {admin_token}"}
        
        payload = {
            "username": user_data.username,
            "email": user_data.email,
            "firstName": user_data.username,
            "lastName": user_data.username,
            "enabled": True,
            "emailVerified": True,
            "requiredActions": [],
            "credentials": [{
                "type": "password",
                "value": user_data.hashed_password,
                "temporary": False
            }]
        }

        response = requests.post(url.replace('localhost', 'keycloak', 1), json=payload, headers=headers)

        if response.status_code != 201:
            raise ValueError(f"Failed to create user in Keycloak: {response.text}")

        location_header = response.headers.get("Location")
        if not location_header:
            raise ValueError("Keycloak did not return a Location header")

        keycloak_id = location_header.split("/")[-1]
        print(f"Created user in Keycloak with ID: {keycloak_id}")

        hashed_password = pwd_context.hash(user_data.hashed_password)
        user_data.hashed_password = hashed_password
        user = self.method.create_user(user_data, keycloak_id)
        self.broker.send("users", {"action": "create", "user": {"id": str(user.id), "username": user.username, "email": user.email, "created_at": str(user.created_at.isoformat())}})
        return user

    def authenticate_user(self, login_data: UserLogin):
        token_url = self.keycloak_openid_config["token_endpoint"]
        print(token_url)
        data = {
            "client_id": settings.KEYCLOAK_CLIENT_ID,
            "client_secret": settings.KEYCLOAK_CLIENT_SECRET,
            "username": login_data.username,
            "password": login_data.password,
            "grant_type": "password"
        }

        response = requests.post(token_url.replace('localhost', 'keycloak', 1), data=data)
        

        return response.json()
