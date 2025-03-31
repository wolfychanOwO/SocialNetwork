from fastapi import Depends
from app.arch.methods import UserMethod
from app.service.auth import AuthService
from app.service.kafka_producer import KafkaProducerService
from app.config import settings

async def get_kafka_producer():
    producer = KafkaProducerService(bootstrap_servers=settings.KAFKA_BOOTSTRAP_SERVERS)
    await producer.start()
    try:
        yield producer
    finally:
        await producer.stop()

def get_auth_service(
    user_method: UserMethod = Depends(),
    kafka_producer: KafkaProducerService = Depends(get_kafka_producer)
) -> AuthService:
    return AuthService(method=user_method, broker=kafka_producer) 
