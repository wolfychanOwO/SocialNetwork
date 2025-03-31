import grpc
from concurrent import futures
import asyncio
import grpc.aio
import time

from app.generated import post_pb2_grpc, post_pb2
from app.service.grpc_svc import PostRPCService
from app.service.post_svc import PostService
from app.arch.methods import PostMethod
from app.arch.database import SessionLocal
from grpc_reflection.v1alpha import reflection
from app.service.kafka_producer import KafkaProducerService

async def serve():
    server = grpc.aio.server()
    db = SessionLocal()
    kafka_producer = KafkaProducerService()
    post_service = PostService(PostMethod(db), kafka_producer)
    post_pb2_grpc.add_SocialServiceServicer_to_server(PostRPCService(post_service), server)

    SERVICE_NAMES = (
        post_pb2.DESCRIPTOR.services_by_name['SocialService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    
    server.add_insecure_port('[::]:50051')
    print("Starting server on port 50051...")
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    asyncio.run(serve())
