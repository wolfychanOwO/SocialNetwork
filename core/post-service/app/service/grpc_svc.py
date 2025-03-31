from app.generated import post_pb2, post_pb2_grpc
from app.service.post_svc import PostService
from google.protobuf.timestamp_pb2 import Timestamp


class PostRPCService(post_pb2_grpc.SocialServiceServicer):
    def __init__(self, post_service: PostService):
        self.post_service = post_service

    async def CreatePost(self, request, context):
        now = Timestamp()
        now.GetCurrentTime()
        
        post_data = {
            "title": request.title,
            "description": request.description,
            "user_id": request.user_id,
            "is_private": request.is_private,
            "tags": request.tags,
            "loyalty_platform": request.loyalty_platform,
            "created_at": now,
            "updated_at": now
        }
    
        created_post = await self.post_service.create_post(post_data)
    
        print('test2')
        return post_pb2.PostResponse(
            post=post_pb2.Post(
                id=str(created_post["id"]),
                title=created_post["title"],
                description=created_post["description"],
                user_id=str(created_post["user_id"]),
                is_private=created_post["is_private"],
                tags=created_post["tags"],
                loyalty_platform=created_post["loyalty_platform"],
                created_at=now,
                updated_at=now
            )
        )

    async def GetPost(self, request, context):
        post = await self.post_service.get_post_by_id(request.id)
        if not post:
            context.set_code(post_pb2.NOT_FOUND)
            context.set_details("Post not found")
            return post_pb2.PostResponse()
        
        return post_pb2.PostResponse(
            post=post_pb2.Post(
                id=str(post["id"]),
                title=post["title"],
                description=post["description"],
                user_id=str(post["user_id"]),
                is_private=post["is_private"],
                tags=post["tags"],
                loyalty_platform=post["loyalty_platform"],
                created_at=post["created_at"],
                updated_at=post["updated_at"]
            )
        )

    async def UpdatePost(self, request, context):
        post_data = {
            "id": request.id,
            "title": request.title,
            "description": request.description,
            "user_id": request.user_id,
            "is_private": request.is_private,
            "tags": request.tags,
            "loyalty_platform": request.loyalty_platform
        }

        updated_post = await self.post_service.update_post(post_data)

        if not updated_post:
            context.set_code(post_pb2.NOT_FOUND)
            context.set_details("Post not found for update")
            return post_pb2.PostResponse()

        return post_pb2.PostResponse(
            post=post_pb2.Post(
                id=str(updated_post["id"]),
                title=updated_post["title"],
                description=updated_post["description"],
                user_id=str(updated_post["user_id"]),
                is_private=updated_post["is_private"],
                tags=updated_post["tags"],
                loyalty_platform=updated_post["loyalty_platform"],
                created_at=updated_post["created_at"],
                updated_at=updated_post["updated_at"]
            )
        )

    def post_to_proto(self, post: dict) -> post_pb2.Post:
        created_at = Timestamp()
        created_at.FromDatetime(post["created_at"])

        updated_at = Timestamp()
        updated_at.FromDatetime(post["updated_at"])

        return post_pb2.Post(
            id=str(post["id"]),
            title=post["title"],
            description=post["description"],
            user_id=str(post["user_id"]),
            is_private=post["is_private"],
            tags=post["tags"],
            loyalty_platform=post["loyalty_platform"],
            created_at=created_at,
            updated_at=updated_at
        )

    async def DeletePost(self, request, context):
        success = await self.post_service.delete_post(request.id)

        if not success:
            context.set_code(post_pb2.NOT_FOUND)
            context.set_details("Post not found")
            return post_pb2.Empty()
        return post_pb2.Empty()
    
    async def ListPosts(self, request, context):
        result = await self.post_service.get_posts_paginated(
            page=request.page, page_size=request.page_size
        )

        posts_proto = [self.post_to_proto(post) for post in result["posts"]]

        return post_pb2.ListPostsResponse(
            posts=posts_proto,
            total=result["total"]
        )
    async def LikePost(self, request, context):
        post_id = request.post_id
        user_id = request.user_id

        success = await self.post_service.like_post(post_id, user_id)

        if not success:
            context.set_code(post_pb2.NOT_FOUND)
            context.set_details("Post not found")
            return post_pb2.Empty()

        return post_pb2.Empty()
    
    async def AddComment(self, request, context):
        
        comment_info = await self.post_service.add_comment(post_id=request.post_id, user_id=request.user_id, content=request.content)
        created_at = Timestamp()
        created_at.FromDatetime(comment_info["created_at"])
        return post_pb2.CommentResponse(
            comment = post_pb2.Comment(id=str(comment_info["id"]),
            post_id=str(comment_info["post_id"]),
            user_id=str(comment_info["user_id"]),
            content=comment_info["content"],
            created_at=created_at)
        )
    
    async def ListComments(self, request, context):

        comments_data = await self.post_service.get_comments_by_post(post_id=request.post_id, page=request.page, page_size=request.page_size)

        updated_at = Timestamp()
        return post_pb2.CommentsListResponse(
            comments=[
                post_pb2.Comment(
                    id=str(comment["id"]),
                    post_id=str(request.post_id),
                    user_id=str(comment["user_id"]),
                    content=comment["content"],
                    created_at=updated_at.FromDatetime(comment["created_at"])
                )
                for comment in comments_data["comments"]
            ],
            total=comments_data["total"]
        )

