from generated.stats_pb2_grpc import StatsServiceStub
from generated.stats_pb2 import GetPostStatsRequest, GetDynamicsRequest, GetTopRequest

from client.rpc.base import BaseGrpcClient


class StatisticGrpcClient(BaseGrpcClient):
    def get_stub(self):
        return StatsServiceStub(self.channel)
    async def get_post_stats(self, request):
        stub = self.get_stub()
        print('test1')
        request = self.request_adapter.from_json(GetPostStatsRequest, request)
        print('start test')
        responce = await stub.GetPostStats(request)
        print('end test')
        return self.request_adapter.to_json(responce)

    async def get_post_view_dynamics(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetDynamicsRequest, request)
        responce = await stub.GetPostViewsDynamics(request)
        return self.request_adapter.to_json(responce)

    async def get_post_likes_dynamics(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetDynamicsRequest, request)
        responce = await stub.GetPostLikesDynamics(request)
        return self.request_adapter.to_json(responce)

    async def get_post_comment_dynamics(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetDynamicsRequest, request)
        responce = await stub.GetPostCommentsDynamics(request)
        return self.request_adapter.to_json(responce)

    async def get_top_posts(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetTopRequest, request)
        responce = await stub.GetTopPosts(request)
        return self.request_adapter.to_json(responce)

    async def get_top_users(self, request):
        stub = self.get_stub()
        request = self.request_adapter.from_json(GetTopRequest, request)
        responce = await stub.GetTopUsers(request)
        return self.request_adapter.to_json(responce)
    
