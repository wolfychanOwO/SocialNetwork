from client.rpc.base import BaseGrpcClient
from generated.post_pb2_grpc import SocialServiceStub
from generated.post_pb2 import CreatePostRequest, UpdatePostRequest, DeleteRequest, GetRequest, ListRequest, LikePostRequest, CommentRequest, CommentsListRequest

class PostGrpcClient(BaseGrpcClient):
    def get_stub(self):
        return SocialServiceStub(self.channel)
    async def create_post(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(CreatePostRequest, request)
        response = await stub.CreatePost(request)
        return self.request_adapter.to_json(response)
    async def update_post(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(UpdatePostRequest, request)
        response = await stub.UpdatePost(request)
        return self.request_adapter.to_json(response)
    async def delete_post(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(DeleteRequest, request)
        response = await stub.DeletePost(request)
        return self.request_adapter.to_json(response)
    
    async def get_post(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetRequest, request)
        response = await stub.GetPost(request)
        return self.request_adapter.to_json(response)
    async def list_posts(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(ListRequest, request)
        response = await stub.ListPosts(request)
        return self.request_adapter.to_json(response)
    async def like_post(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(LikePostRequest, request)
        response = await stub.LikePost(request)
        return self.request_adapter.to_json(response)
    async def add_comment(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(CommentRequest, request)
        response = await stub.AddComment(request)
        return self.request_adapter.to_json(response)
    async def list_comments(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(CommentsListRequest, request)
        response = await stub.ListComments(request)
        return self.request_adapter.to_json(response)
