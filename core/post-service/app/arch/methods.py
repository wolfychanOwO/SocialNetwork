from sqlalchemy.orm import Session
from app.models.models import Post, PostLike, Comment
from sqlalchemy.exc import NoResultFound
from typing import Optional
from sqlalchemy import func
import datetime

class PostMethod:
    def __init__(self, db: Session):
        self.db = db

    async def create_post(self, post_data: dict) -> dict:
        
        db_post = Post(
            title=post_data["title"],
            description=post_data["description"],
            user_id=post_data["user_id"],
            is_private=post_data["is_private"],
            tags=post_data["tags"],
            loyalty_platform=post_data["loyalty_platform"]
        )
        
        self.db.add(db_post)
        self.db.commit()
        self.db.refresh(db_post)
        
        return {
            "id": db_post.id,
            "title": db_post.title,
            "description": db_post.description,
            "user_id": db_post.user_id,
            "is_private": db_post.is_private,
            "tags": db_post.tags,
            "loyalty_platform": db_post.loyalty_platform,
            "created_at": db_post.created_at,
            "updated_at": db_post.updated_at
        }

    async def get_post_by_id(self, post_id: str) -> Optional[dict]:
        try:
            post = self.db.query(Post).filter(Post.id == post_id).one()
            return {
                "id": post.id,
                "title": post.title,
                "description": post.description,
                "user_id": post.user_id,
                "is_private": post.is_private,
                "tags": post.tags,
                "loyalty_platform": post.loyalty_platform,
                "created_at": post.created_at,
                "updated_at": post.updated_at
            }
        except NoResultFound:
            return None

    async def update_post(self, post_data: dict) -> Optional[dict]:
        db_post = self.db.query(Post).filter(Post.id == post_data["id"]).first()
        if db_post:
            db_post.title = post_data["title"]
            db_post.description = post_data["description"]
            db_post.is_private = post_data["is_private"]
            db_post.tags = post_data["tags"]
            db_post.loyalty_platform = post_data["loyalty_platform"]
            self.db.commit()
            self.db.refresh(db_post)
            return {
                "id": db_post.id,
                "title": db_post.title,
                "description": db_post.description,
                "user_id": db_post.user_id,
                "is_private": db_post.is_private,
                "tags": db_post.tags,
                "loyalty_platform": db_post.loyalty_platform,
                "created_at": db_post.created_at,
                "updated_at": db_post.updated_at
            }
        return None

    async def delete_post(self, post_id: str) -> bool:
        db_post = self.db.query(Post).filter(Post.id == post_id).first()
        if db_post:
            self.db.delete(db_post)
            self.db.commit()
            return True
        return False
    

    async def get_posts_paginated(self, page: int = 1, page_size: int = 10) -> dict:
        query = self.db.query(Post)
        
        total_query = self.db.query(func.count(Post.id))
        total = total_query.scalar()

        offset = (page - 1) * page_size
        query = query.offset(offset).limit(page_size)
        
        posts = query.all()

        return {
            "posts": [
                {
                    "id": post.id,
                    "title": post.title,
                    "description": post.description,
                    "user_id": post.user_id,
                    "is_private": post.is_private,
                    "tags": post.tags,
                    "loyalty_platform": post.loyalty_platform,
                    "created_at": post.created_at,
                    "updated_at": post.updated_at
                }
                for post in posts
            ],
            "total": total
        }
    
    async def like_post(self, post_id: str, user_id: str) -> dict:
        
        post = self.db.query(Post).filter(Post.id == post_id).one()

        existing_like = self.db.query(PostLike).filter(PostLike.post_id == post_id, PostLike.user_id == user_id).first()
        if existing_like:
            raise ValueError("User has already liked this post.")
        
        new_like = PostLike(post_id=post_id, user_id=user_id)
        self.db.add(new_like)
        self.db.commit()
        
        return {
            "post_id": new_like.post_id,
            "user_id": new_like.user_id,
            "like_id": new_like.id
        }
    
    async def add_comment(self, post_id: str, user_id: str, content: str) -> dict:
        try:
            post = self.db.query(Post).filter(Post.id == post_id).one()

            comment = Comment(post_id=post_id, user_id=user_id, content=content)
            self.db.add(comment)
            self.db.commit()
            self.db.refresh(comment)

            return {
                "id": comment.id,
                "post_id": comment.post_id,
                "user_id": comment.user_id,
                "content": comment.content,
                "created_at": comment.created_at
            }
        except NoResultFound:
            raise ValueError("Post not found.")

    async def get_comments_by_post(self, post_id: str, page: int = 1, page_size: int = 10) -> dict:
        try:
            post = self.db.query(Post).filter(Post.id == post_id).one()

            query = self.db.query(Comment).filter(Comment.post_id == post_id)

            offset = (page - 1) * page_size
            comments = query.offset(offset).limit(page_size).all()

            total_comments = self.db.query(func.count(Comment.id)).filter(Comment.post_id == post_id).scalar()

            return {
                "comments": [
                    {
                        "id": comment.id,
                        "user_id": comment.user_id,
                        "content": comment.content,
                        "created_at": comment.created_at
                    }
                    for comment in comments
                ],
                "total": total_comments
            }
        except NoResultFound:
            raise ValueError("Post not found.")
