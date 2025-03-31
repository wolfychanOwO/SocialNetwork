from sqlalchemy.orm import Session

from app.models.models import User
from app.models.schema import UserCreate


class UserMethods:
    def __init__(self, db: Session):
        self.db = db

    def get_user_by_username(self, username: str):
        return self.db.query(User).filter(User.username == username).first()
    
    def get_user_by_keycloak_id(self, keycloak_id: str):
        return self.db.query(User).filter(User.keycloak_id == keycloak_id).first()

    def create_user(self, user: UserCreate, keycloak_id: str):
        user = User(**user.dict(), keycloak_id=keycloak_id)
        self.db.add(user)
        self.db.commit()
        self.db.refresh(user)
        return user

    def update_user(self, user: User):
        self.db.commit()
        self.db.refresh(user)
        return user
    
    def get_user_by_username(self, username: str):
        return self.db.query(User).filter(User.username == username).first()
