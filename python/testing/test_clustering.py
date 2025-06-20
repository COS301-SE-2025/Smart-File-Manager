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

# Create a fixture to automatically setup and tear down a grpc_test_server
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


# <------ INTEGRATION TESTING ----->
# Fixture creates a gRPC DirectoryRequest from everything in python/testing/test_files_2 for repeated use
@pytest.fixture(scope="module")
def createDirectoryRequest():
    TEST_DIR = os.path.dirname(__file__)
    TEST_FILE_DIR = os.path.join(TEST_DIR, "test_files_3")

    def get_path(name):
        return os.path.join(TEST_FILE_DIR, name)

    files1 =   [
        File(name="Apr8TODO.txt", original_path=get_path("Apr8TODO.txt")),
        File(name="Apr18 meeting.txt", original_path=get_path("Apr18 meeting.txt")),
        File(name="architecture_diagram.png", original_path=get_path("architecture_diagram.png")),
        File(name="Assignment2.pdf", original_path=get_path("Assignment2.pdf")),
        File(name="collection_page_wireframe.png", original_path=get_path("collection_page_wireframe.png")),
        File(name="COS 301 - Mini-Project - Demo 1 Instructions.pdf", original_path=get_path("COS 301 - Mini-Project - Demo 1 Instructions.pdf")),
        File(name="COS 301 - Mini-Project - Demo 2 Instructions.pdf", original_path=get_path("COS 301 - Mini-Project - Demo 2 Instructions.pdf")),
        File(name="COS122 Tutorial 4 Sept 7-8, 2023.pdf", original_path=get_path("COS122 Tutorial 4 Sept 7-8, 2023.pdf")),
        File(name="COS221 Assignment 1 2025.pdf", original_path=get_path("COS221 Assignment 1 2025.pdf")),
        File(name="cpp_api.md", original_path=get_path("cpp_api.md")),
        File(name="DeeBee.png", original_path=get_path("DeeBee.png")),
        File(name="Importing the Database.md", original_path=get_path("Importing the Database.md")),
        File(name="L01_Ch01a(1).pdf", original_path=get_path("L01_Ch01a(1).pdf")),
        File(name="L05_Ch02c.pdf", original_path=get_path("L05_Ch02c.pdf")),
        File(name="login_wireframe.png", original_path=get_path("login_wireframe.png")),
        File(name="MP Progress report.txt", original_path=get_path("MP Progress report.txt")),
        File(name="mp11_design_specification.md", original_path=get_path("mp11_design_specification.md")),
        File(name="mp11_requirement_spec.md", original_path=get_path("mp11_requirement_spec.md")),
        File(name="MPChecklist.txt", original_path=get_path("MPChecklist.txt")),
        File(name="Prac1Triggers.txt", original_path=get_path("Prac1Triggers.txt")),
        File(name="Screenshot_2025-02-26_at_15.36.48.png", original_path=get_path("Screenshot_2025-02-26_at_15.36.48.png")),
        File(name="statistics_page_wireframe.png", original_path=get_path("statistics_page_wireframe.png")),
        File(name="TODO mar30 Meeting.txt", original_path=get_path("TODO mar30 Meeting.txt")),
        File(name="Tututorial_2.pdf", original_path=get_path("Tututorial_2.pdf")),
        File(name="UseCase.png", original_path=get_path("UseCase.png")),
        File(name="init.py", original_path=get_path("init.py")),
        File(name="flamegraph.svg", original_path=get_path("flamegraph.svg"))
    ]

    root_dir = Directory(
        name = "test_files_3",
        path = get_path("test_files_3"),
        files = files1
    )    

    req = DirectoryRequest(root=root_dir)
    yield req


# Sends an actual directory and checks if metadata was correctly attached to files
def test_send_real_dir(grpc_test_server, createDirectoryRequest):
    req = createDirectoryRequest  # Accessing req from the fixture
    response = grpc_test_server.SendDirectoryStructure(req)