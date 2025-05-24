from client.rpc.post_client import PostGrpcClient
from client.rpc.adapter import RequestAdapter

class GrpcFactory:
    def __init__(self):
        self.services = {
            "posts": lambda: PostGrpcClient("post-service:50051", RequestAdapter),
        }

    def get_client(self, service_name: str):
        if service_name not in self.services:
            raise ValueError(f"Service {service_name} not found")
        return self.services[service_name]()
