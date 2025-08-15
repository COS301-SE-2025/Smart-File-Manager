# tests/conftest.py
import os, sys
import pytest
from concurrent import futures
import grpc

# ensure your `src/` is on the import path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from src import message_structure_pb2_grpc
from src.request_handler import RequestHandler


@pytest.fixture(scope="session")
def grpc_test_server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    message_structure_pb2_grpc.add_DirectoryServiceServicer_to_server(RequestHandler(), server)
    port = server.add_insecure_port("localhost:0")  # OS assigns free port
    server.start()

    channel = grpc.insecure_channel(f"localhost:{port}")
    stub = message_structure_pb2_grpc.DirectoryServiceStub(channel)

    yield stub  # This is used in the test function

    server.stop(None)
