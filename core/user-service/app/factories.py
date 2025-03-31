from app.service.kafka_producer import KafkaProducerService
from fastapi import Depends
from app.arch.database import SessionLocal
from app.config import settings

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

async def get_kafka_producer(broker_url: str = settings.KAFKA_BOOTSTRAP_SERVERS) -> KafkaProducerService:
    kafka_producer = KafkaProducerService(broker_url)
    return kafka_producer
