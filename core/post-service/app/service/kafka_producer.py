from kafka import KafkaProducer
import json

class KafkaProducerService:
    _instance = None
    
    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            cls._instance = super(KafkaProducerService, cls).__new__(cls)
            cls._instance._initialized = False
        return cls._instance
    
    def __init__(self, bootstrap_servers='broker:29092'):
        if not self._initialized:
            self.producer = KafkaProducer(
                bootstrap_servers=bootstrap_servers,
                value_serializer=lambda v: json.dumps(v).encode('utf-8'),
                acks='all',
                retries=3
            )
            self._initialized = True
            print("Kafka producer initialized")
    
    def send(self, topic: str, message: dict):
        try:
            future = self.producer.send(topic, message)
            future.get(timeout=10)
            print(f"Message sent successfully to topic {topic}")
        except Exception as e:
            print(f"Failed to send message to topic {topic}: {str(e)}")
            raise
    
    def close(self):
        if self.producer:
            self.producer.close()
            print("Kafka producer closed") 
