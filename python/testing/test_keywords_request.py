from typing import Dict, List
import grpc
import pytest
from concurrent import futures
import sys
import os
import time

# Add src to path temporarily so the generated grpc file can find message_structure_pb2
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from src import message_structure_pb2, message_structure_pb2_grpc
from src.message_structure_pb2 import Directory, File, DirectoryRequest, Keyword
from src.request_handler import RequestHandler


# <------ INTEGRATION TESTING ----->
# Fixture creates a gRPC DirectoryRequest from everything in python/testing/test_files_2 for repeated use
@pytest.fixture(scope="module")
def createDirectoryRequest():
    TEST_DIR = os.path.dirname(__file__)
    TEST_FILE_DIR = os.path.join(TEST_DIR, "test_files_3")

    def get_path(name):
        return os.path.join(TEST_FILE_DIR, name)

    files1 = [

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
        File(name="ENjoyment", original_path=get_path("ENjoyment.png")),
        File(name="Gantt chart", original_path=get_path("Gantt chart.png")),
        File(name="Gauteng", original_path=get_path("Gauteng.png")),
        File(name="most challanging", original_path=get_path("most challanging.png")),
        File(name="Most rewarding", original_path=get_path("Most rewarding.png")),
        File(name="Picture1", original_path=get_path("Picture1.png")),
        File(name="Picture2", original_path=get_path("Picture2.png")),
        File(name="Presentation speech", original_path=get_path("Presentation speech.docx")),
        File(name="Project Budget Form 2024", original_path=get_path("Project Budget Form 2024.pdf")),
        File(name="Taiichi ohno", original_path=get_path("Taiichi ohno.jpeg")),
        File(name="Week 3_Tutorial_2024_with Answers", original_path=get_path("Week 3_Tutorial_2024_with Answers.pdf")),
        File(name="Week 4_Tutorial_with answers", original_path=get_path("Week 4_Tutorial_with answers.pdf")),
        File(name="Week 5_Tutorial_2024_with answers", original_path=get_path("Week 5_Tutorial_2024_with answers.pdf"))

    ]

    root_dir = Directory(
        name = "test_files_3",
        path = get_path("test_files_3"),
        files = files1,
    )    

    req = DirectoryRequest(root=root_dir, requestType="KEYWORDS")
    yield req


def recHelper(curDir : Directory, kws : Dict[str, List[str]]):

    # Check response contains at least some keywords 
    for curFile in curDir.files:
        
        if not curFile.original_path.endswith(".png") and not curFile.original_path.endswith(".jpeg") and not curFile.original_path.endswith(".md"):
            assert len(curFile.keywords) > 0
            words = []
            for w in curFile.keywords:
                words.append(w)
            kws[curFile.name] = words
    
        # Recurisve call
        if len(curDir.directories) != 0:
            for dir in curDir.directories:
                recHelper(dir)

# Sends an actual directory and checks if metadata was correctly attached to filejjs
def test_send_real_dir(grpc_test_server, createDirectoryRequest):
    req = createDirectoryRequest  # Accessing req from the fixture
    start = time.time()
    response = grpc_test_server.SendDirectoryStructure(req)
    end = time.time()

    # Check if enough metadata was extracted
    kws = {}
    recHelper(response.root, kws)
    for file, words in kws.items():
        print(f"{file}: {words}")

    print(f"Method took: {end - start} seconds")
    # Check if response is well formed
    assert response.response_code == 200
    assert response.response_msg != "No file could be opened"


