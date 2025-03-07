from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
import jwt
import requests
from jose import jwt
from jose.exceptions import JOSEError

from app.arch.database import SessionLocal
from app.service.auth import AuthService
from app.arch.methods import UserMethods
from app.models.schema import UserCreate, UserLogin, UserResponse
from app.config import settings

router = APIRouter()

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

def verify_token(token: str):
    try:
        jwks_url = f"{settings.KEYCLOAK_URL}/realms/{settings.KEYCLOAK_REALM}/protocol/openid-connect/certs"
        jwks = requests.get(jwks_url).json()
        return jwt.decode(
            token,
            jwks,
            algorithms=["RS256"],
            audience=settings.KEYCLOAK_CLIENT_ID,
            options={"verify_aud": False}
        )
    except JOSEError as e:
        raise HTTPException(status_code=401, detail=str(e))
    
    
@router.post("/register", response_model=UserResponse)
def register(user: UserCreate, db: Session = Depends(get_db)):
    user_methods = UserMethods(db)
    auth_service = AuthService(user_methods)
    try:
        return auth_service.register_user(user)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.post("/login")
def login(user: UserLogin, db: Session = Depends(get_db)):
    user_methods = UserMethods(db)
    auth_service = AuthService(user_methods)
    try:
        return auth_service.authenticate_user(user)
    except ValueError as e:
        raise HTTPException(status_code=401, detail=str(e))
    
@router.get("/profile", response_model=UserResponse)
def get_profile(token: str = Depends(verify_token), db: Session = Depends(get_db)):
    user_methods = UserMethods(db)
    print(token["sub"])
    user = user_methods.get_user_by_keycloak_id(token["sub"])
    
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    return user

