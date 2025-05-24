from fastapi import APIRouter, Depends, HTTPException, Header
from sqlalchemy.orm import Session

from app.service.auth import AuthService
from app.arch.methods import UserMethod
from app.models.schema import UserCreate, UserLogin, UserResponse

import jwt
import requests
from jose import jwt
from jose.exceptions import JOSEError

from app.config import settings
from app.service.kafka_producer import KafkaProducerService
from app.factories import get_db, get_kafka_producer

router = APIRouter()


def verify_token(authorization: str = Header(...)):
    try:
        jwks_url = f"{settings.KEYCLOAK_URL}/realms/{settings.KEYCLOAK_REALM}/protocol/openid-connect/certs"
        jwks = requests.get(jwks_url).json()
        return jwt.decode(
            authorization,
            jwks,
            algorithms=["RS256"],
            audience=settings.KEYCLOAK_CLIENT_ID,
            options={"verify_aud": False}
        )
    except JOSEError as e:
        raise HTTPException(status_code=401, detail=str(e))
    
    
@router.post("/register", response_model=UserResponse)
def register(user: UserCreate, db: Session = Depends(get_db), broker: KafkaProducerService = Depends(get_kafka_producer)):
    method = UserMethod(db)
    auth_service = AuthService(method, broker)
    try:
        return auth_service.register_user(user)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.post("/login")
def login(user: UserLogin, db: Session = Depends(get_db), broker: KafkaProducerService = Depends(get_kafka_producer)):
    method = UserMethod(db)
    auth_service = AuthService(method, broker)
    try:
        return auth_service.authenticate_user(user)
    except ValueError as e:
        raise HTTPException(status_code=401, detail=str(e))
    
@router.get("/validate", response_model=UserResponse)
def get_profile(token: dict = Depends(verify_token), db: Session = Depends(get_db)):
    method = UserMethod(db)
    print(token["sub"])
    user = method.get_user_by_keycloak_id(token["sub"])
    
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    return user
