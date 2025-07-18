# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: post.proto
# Protobuf Python Version: 6.30.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    30,
    0,
    '',
    'post.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\npost.proto\x12\x06social\x1a\x1fgoogle/protobuf/timestamp.proto\"\x07\n\x05\x45mpty\",\n\rDeleteRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0f\n\x07user_id\x18\x02 \x01(\t\")\n\nGetRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0f\n\x07user_id\x18\x02 \x01(\t\"?\n\x0bListRequest\x12\x0c\n\x04page\x18\x01 \x01(\x05\x12\x11\n\tpage_size\x18\x02 \x01(\x05\x12\x0f\n\x07user_id\x18\x03 \x01(\t\"\xe3\x01\n\x04Post\x12\n\n\x02id\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x0f\n\x07user_id\x18\x04 \x01(\t\x12\x12\n\nis_private\x18\x05 \x01(\x08\x12\x0c\n\x04tags\x18\x06 \x03(\t\x12\x18\n\x10loyalty_platform\x18\x07 \x01(\t\x12.\n\ncreated_at\x18\x08 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12.\n\nupdated_at\x18\t \x01(\x0b\x32\x1a.google.protobuf.Timestamp\"\x84\x01\n\x11\x43reatePostRequest\x12\r\n\x05title\x18\x01 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x02 \x01(\t\x12\x0f\n\x07user_id\x18\x03 \x01(\t\x12\x12\n\nis_private\x18\x04 \x01(\x08\x12\x0c\n\x04tags\x18\x05 \x03(\t\x12\x18\n\x10loyalty_platform\x18\x06 \x01(\t\"\x90\x01\n\x11UpdatePostRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x12\n\nis_private\x18\x04 \x01(\x08\x12\x0c\n\x04tags\x18\x05 \x03(\t\x12\x18\n\x10loyalty_platform\x18\x06 \x01(\t\x12\x0f\n\x07user_id\x18\x07 \x01(\t\"*\n\x0cPostResponse\x12\x1a\n\x04post\x18\x01 \x01(\x0b\x32\x0c.social.Post\"?\n\x11ListPostsResponse\x12\x1b\n\x05posts\x18\x01 \x03(\x0b\x32\x0c.social.Post\x12\r\n\x05total\x18\x02 \x01(\x05\"3\n\x0fLikePostRequest\x12\x0f\n\x07post_id\x18\x01 \x01(\t\x12\x0f\n\x07user_id\x18\x02 \x01(\t\"x\n\x07\x43omment\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0f\n\x07post_id\x18\x02 \x01(\t\x12\x0f\n\x07user_id\x18\x03 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x04 \x01(\t\x12.\n\ncreated_at\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\"C\n\x0e\x43ommentRequest\x12\x0f\n\x07post_id\x18\x01 \x01(\t\x12\x0f\n\x07user_id\x18\x02 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x03 \x01(\t\"3\n\x0f\x43ommentResponse\x12 \n\x07\x63omment\x18\x01 \x01(\x0b\x32\x0f.social.Comment\"G\n\x13\x43ommentsListRequest\x12\x0f\n\x07post_id\x18\x01 \x01(\t\x12\x0c\n\x04page\x18\x02 \x01(\x05\x12\x11\n\tpage_size\x18\x03 \x01(\x05\"H\n\x14\x43ommentsListResponse\x12!\n\x08\x63omments\x18\x01 \x03(\x0b\x32\x0f.social.Comment\x12\r\n\x05total\x18\x02 \x01(\x05\x32\xf1\x03\n\rSocialService\x12=\n\nCreatePost\x12\x19.social.CreatePostRequest\x1a\x14.social.PostResponse\x12=\n\nUpdatePost\x12\x19.social.UpdatePostRequest\x1a\x14.social.PostResponse\x12\x32\n\nDeletePost\x12\x15.social.DeleteRequest\x1a\r.social.Empty\x12\x33\n\x07GetPost\x12\x12.social.GetRequest\x1a\x14.social.PostResponse\x12;\n\tListPosts\x12\x13.social.ListRequest\x1a\x19.social.ListPostsResponse\x12\x32\n\x08LikePost\x12\x17.social.LikePostRequest\x1a\r.social.Empty\x12=\n\nAddComment\x12\x16.social.CommentRequest\x1a\x17.social.CommentResponse\x12I\n\x0cListComments\x12\x1b.social.CommentsListRequest\x1a\x1c.social.CommentsListResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'post_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_EMPTY']._serialized_start=55
  _globals['_EMPTY']._serialized_end=62
  _globals['_DELETEREQUEST']._serialized_start=64
  _globals['_DELETEREQUEST']._serialized_end=108
  _globals['_GETREQUEST']._serialized_start=110
  _globals['_GETREQUEST']._serialized_end=151
  _globals['_LISTREQUEST']._serialized_start=153
  _globals['_LISTREQUEST']._serialized_end=216
  _globals['_POST']._serialized_start=219
  _globals['_POST']._serialized_end=446
  _globals['_CREATEPOSTREQUEST']._serialized_start=449
  _globals['_CREATEPOSTREQUEST']._serialized_end=581
  _globals['_UPDATEPOSTREQUEST']._serialized_start=584
  _globals['_UPDATEPOSTREQUEST']._serialized_end=728
  _globals['_POSTRESPONSE']._serialized_start=730
  _globals['_POSTRESPONSE']._serialized_end=772
  _globals['_LISTPOSTSRESPONSE']._serialized_start=774
  _globals['_LISTPOSTSRESPONSE']._serialized_end=837
  _globals['_LIKEPOSTREQUEST']._serialized_start=839
  _globals['_LIKEPOSTREQUEST']._serialized_end=890
  _globals['_COMMENT']._serialized_start=892
  _globals['_COMMENT']._serialized_end=1012
  _globals['_COMMENTREQUEST']._serialized_start=1014
  _globals['_COMMENTREQUEST']._serialized_end=1081
  _globals['_COMMENTRESPONSE']._serialized_start=1083
  _globals['_COMMENTRESPONSE']._serialized_end=1134
  _globals['_COMMENTSLISTREQUEST']._serialized_start=1136
  _globals['_COMMENTSLISTREQUEST']._serialized_end=1207
  _globals['_COMMENTSLISTRESPONSE']._serialized_start=1209
  _globals['_COMMENTSLISTRESPONSE']._serialized_end=1281
  _globals['_SOCIALSERVICE']._serialized_start=1284
  _globals['_SOCIALSERVICE']._serialized_end=1781
# @@protoc_insertion_point(module_scope)
