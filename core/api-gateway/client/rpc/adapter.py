from google.protobuf.message import Message
from typing import Type
from google.protobuf.json_format import MessageToDict

class RequestAdapter:
    @staticmethod
    def from_json(proto_cls: Type[Message], json_data: dict) -> Message:
        return proto_cls(**json_data)
    
    @staticmethod
    def to_json(proto: Message) -> dict:
        return MessageToDict(proto)
