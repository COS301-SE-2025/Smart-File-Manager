import grpc
import pytest
from concurrent import futures
import sys
import os

# Add src to path temporarily so the generated grpc file can find message_structure_pb2
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from src import message_structure_pb2, message_structure_pb2_grpc
from src.message_structure_pb2 import Directory, File, Tag, MetadataEntry, DirectoryRequest
from src.request_handler import RequestHandler

# Intgeration tests to see if gRPC server works
@pytest.fixture(scope="module")
def grpc_test_server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
    message_structure_pb2_grpc.add_DirectoryServiceServicer_to_server(RequestHandler(), server)
    port = server.add_insecure_port("localhost:0")  # OS assigns free port
    server.start()

    channel = grpc.insecure_channel(f"localhost:{port}")
    stub = message_structure_pb2_grpc.DirectoryServiceStub(channel)

    yield stub  # This is used in the test function

    server.stop(None)


def test_send_directory_structure(grpc_test_server):
    # Prepare sample request
    tag = Tag(name="testTag")
    metadata = MetadataEntry(key="author", value="testUser")

    test_file = File(
        name="test.pdf",
        original_path="/tmp/test.pdf",  # Simulate valid or invalid path
        new_path="/tmp/new_test.pdf",
        tags=[tag],
        metadata=[metadata],
    )

    test_dir = Directory(
        name="test_dir",
        path="/tmp/",
        files=[test_file],
        directories=[]
    )

    request = DirectoryRequest(root=test_dir)

    # Call gRPC
    response = grpc_test_server.SendDirectoryStructure(request)

    assert isinstance(response, message_structure_pb2.DirectoryResponse)
    assert response.root.name == "test_dir"
    assert len(response.root.files) == 1
