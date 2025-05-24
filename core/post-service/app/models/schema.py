from pydantic import BaseModel, EmailStr
from typing import Optional

class PostCreate(BaseModel):
    title: str
    content: str
    author_id: int

class PostResponse(PostCreate):
    id: int

