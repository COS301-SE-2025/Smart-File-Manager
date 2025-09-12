import pytest
import sys
import os

# Add src to path temporarily so the generated grpc file can find message_structure_pb2
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from src.message_structure_pb2 import Directory, File, DirectoryRequest


""""
# <------ UNIT TESTING ------>
# Simply checks if a Directory response is returned
def test_send_directory_structure(grpc_test_server):
    # Prepare sample request
    tag = Tag(name="testTag")
    metadata = MetadataEntry(key="author", value="testUser")

    test_file = File(
        name="test.pdf",
        original_path="/tmp/test.pdf",  
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
"""

# <------ INTEGRATION TESTING ----->
# Fixture creates a gRPC DirectoryRequest from everything in python/testing/test_files_2 for repeated use
@pytest.fixture(scope="module")
def createDirectoryRequest():
    TEST_DIR = os.path.dirname(__file__)
    TEST_FILE_DIR = os.path.join(TEST_DIR, "test_files_2")

    def get_path(name):
        return os.path.join(TEST_FILE_DIR, name)

    holiday_file = File(
        name = "holiday.JPG",
        original_path = get_path("UserFiles/PersonalFiles/holiday.JPG")
    )

    vid_file = File(
        name = "myVideo.webm",
        original_path = get_path("UserFiles/PersonalFiles/myVideo.webm")
    )    

    img2_file = File(
        name = "thumbbig-708440.webp",
        original_path = get_path("UserFiles/PersonalFiles/thumbbig-708440.webp")
    )    
    pers_dir = Directory(
        name = "PersonalFiles",
        path = get_path("UserFiles/PersonalFiles"),
        files = [vid_file, img2_file, holiday_file]
    )

    pdf_file = File(
        name = "myPdf.pdf",
        original_path = get_path("UserFiles/UniFiles/myPdf.pdf")
    )    

    img_file = File(
        name = "myImg.jpg",
        original_path = get_path("UserFiles/UniFiles/myImg.jpg")
    )    

    uni_dir = Directory(
        name = "UniFiles",
        path =  get_path("UserFiles/UniFiles"),
        files = [img_file, pdf_file],
    )

    todo_file = File(
        name = "todo.docx",
        original_path=get_path("UserFiles/todo.docx")
    )    

    root_dir = Directory(
        name = "UserFiles",
        path = get_path("UserFiles"),
        files = [todo_file],
        directories= [uni_dir, pers_dir]
    )    

    req = DirectoryRequest(root=root_dir, requestType="METADATA", serverSecret=os.environ["SFM_SERVER_SECRET"])
    yield req


def recHelper(curDir : Directory):

    # Check response contains at least base metadata
    for curFile in curDir.files:

        assert len(curFile.metadata) >= 12 # Must extract at least base stats

        # Recurisve call
        if len(curDir.directories) != 0:
            for dir in curDir.directories:
                recHelper(dir)

# Sends an actual directory and checks if metadata was correctly attached to filejjs
def test_send_real_dir(grpc_test_server, createDirectoryRequest):
    req = createDirectoryRequest  # Accessing req from the fixture
    response = grpc_test_server.SendDirectoryStructure(req)

    # Check if enough metadata was extracted
    recHelper(response.root)

    # Check if response is well formed
    assert response.response_code == 200
    assert response.response_msg != "No file could be opened"


