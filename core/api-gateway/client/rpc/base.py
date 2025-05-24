from abc import ABC, abstractmethod
import grpc


class BaseGrpcClient(ABC):
    def __init__(self, address: str, request_adapter=None):
        self.address = address
        self.channel = grpc.aio.insecure_channel(self.address)
        self.request_adapter = request_adapter
    
    @abstractmethod
    def get_stub(self):
        ...


