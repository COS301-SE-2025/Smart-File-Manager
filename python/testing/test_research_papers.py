import pytest
import sys
import os
import time

# Add src to path temporarily so the generated grpc file can find message_structure_pb2
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))


from src.message_structure_pb2 import Directory, File, Tag, DirectoryRequest


# <------ INTEGRATION TESTING ----->
@pytest.fixture(scope="module")
def createDirectoryRequest():
    TEST_DIR = os.path.dirname(__file__)
    TEST_FILE_DIR = os.path.join(TEST_DIR, "test_files_5")

    def get_path(name):
        return os.path.join(TEST_FILE_DIR, name)

    # gRPC objects for files we want to test
    files1 =   [
        File(name="astro_1", original_path=get_path("astro_1.pdf"), tags=[]),
        File(name="astro_2", original_path=get_path("astro_2.pdf"), tags=[]),
        File(name="astro_3", original_path=get_path("astro_3.pdf"), tags=[]),
        File(name="astro_4", original_path=get_path("astro_4.pdf"), tags=[]),
        File(name="astro_5", original_path=get_path("astro_5.pdf"), tags=[]),
        File(name="astro_6", original_path=get_path("astro_6.pdf"), tags=[]),
        File(name="econ_1", original_path=get_path("econ_1.pdf"), tags=[]),
        File(name="econ_2", original_path=get_path("econ_2.pdf"), tags=[]),
        File(name="econ_3", original_path=get_path("econ_3.pdf"), tags=[]),
        File(name="econ_4", original_path=get_path("econ_4.pdf"), tags=[]),
        File(name="econ_5", original_path=get_path("econ_5.pdf"), tags=[]),
        File(name="ee_1", original_path=get_path("ee_1.pdf"), tags=[]),
        File(name="ee_2", original_path=get_path("ee_2.pdf"), tags=[]),
        File(name="ee_3", original_path=get_path("ee_3.pdf"), tags=[]),
        File(name="ee_4", original_path=get_path("ee_4.pdf"), tags=[]),
        File(name="ee_5", original_path=get_path("ee_5.pdf"), tags=[]),
        File(name="math_1", original_path=get_path("math_1.pdf"), tags=[]),
        File(name="math_2", original_path=get_path("math_2.pdf"), tags=[]),
        File(name="math_3", original_path=get_path("math_3.pdf"), tags=[]),
        File(name="math_4", original_path=get_path("math_4.pdf"), tags=[]),
        File(name="math_5", original_path=get_path("math_5.pdf"), tags=[]),
    ]


    root_dir = Directory(
        name = "test_files_5",
        path = get_path("test_files_5"),
        files = files1
    )    

    req = DirectoryRequest(root=root_dir, requestType="CLUSTERING", serverSecret=os.environ["SFM_SERVER_SECRET"])
    yield req



# Sends an actual directory and checks if metadata was correctly attached to files
def test_send_real_dir(grpc_test_server, createDirectoryRequest):
    req = createDirectoryRequest  # Accessing req from the fixture

    start_first = time.time()
    response = grpc_test_server.SendDirectoryStructure(req)
    end_first = time.time()

    # start_second = time.time()
    # response2 = grpc_test_server.SendDirectoryStructure(req)
    # end_second = time.time()

    print("First took: " + str(end_first - start_first))
    # print("Second took: " + str(end_second - start_second))

    # Check if response contains all files

    # Check if response is well formed
    assert response.response_code == 200
    assert response.response_msg != "No file could be opened"