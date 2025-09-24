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
    TEST_FILE_DIR = os.path.join(TEST_DIR, "test_files_3")

    def get_path(name):
        return os.path.join(TEST_FILE_DIR, name)

    # gRPC objects for files we want to test
    files1 =   [
        File(name="Apr8TODO.txt", original_path=get_path("Apr8TODO.txt"), tags=[]),
        File(name="Apr18 meeting.txt", original_path=get_path("Apr18 meeting.txt"), tags=[] ),
        File(name="architecture_diagram.png", original_path=get_path("architecture_diagram.png"), tags=[]),
        File(name="Assignment2.pdf", original_path=get_path("Assignment2.pdf"), tags=[]),
        File(name="collection_page_wireframe.png", original_path=get_path("collection_page_wireframe.png"), tags=[]),
        File(name="COS 301 - Mini-Project - Demo 1 Instructions.pdf", original_path=get_path("COS 301 - Mini-Project - Demo 1 Instructions.pdf"), tags=[]),
        File(name="COS 301 - Mini-Project - Demo 2 Instructions.pdf", original_path=get_path("COS 301 - Mini-Project - Demo 2 Instructions.pdf"), tags=[]),
        File(name="COS122 Tutorial 4 Sept 7-8, 2023.pdf", original_path=get_path("COS122 Tutorial 4 Sept 7-8, 2023.pdf"), tags=[]),
        File(name="COS221 Assignment 1 2025.pdf", original_path=get_path("COS221 Assignment 1 2025.pdf"), tags=[]),
        File(name="cpp_api.md", original_path=get_path("cpp_api.md"), tags=[]),
        File(name="DeeBee.png", original_path=get_path("DeeBee.png"), tags=[]),
        File(name="Importing the Database.md", original_path=get_path("Importing the Database.md"), tags=[]),
        File(name="L01_Ch01a(1).pdf", original_path=get_path("L01_Ch01a(1).pdf"), tags=[]),
        File(name="L05_Ch02c.pdf", original_path=get_path("L05_Ch02c.pdf"), tags=[]),
        File(name="login_wireframe.png", original_path=get_path("login_wireframe.png"), tags=[]),
        File(name="MP Progress report.txt", original_path=get_path("MP Progress report.txt"), tags=[]),
        File(name="mp11_design_specification.md", original_path=get_path("mp11_design_specification.md"), tags=[]),
        File(name="mp11_requirement_spec.md", original_path=get_path("mp11_requirement_spec.md"), tags=[]),
        File(name="MPChecklist.txt", original_path=get_path("MPChecklist.txt"), tags=[]),
        File(name="Prac1Triggers.txt", original_path=get_path("Prac1Triggers.txt"), tags=[]),
        File(name="Screenshot_2025-02-26_at_15.36.48.png", original_path=get_path("Screenshot_2025-02-26_at_15.36.48.png"), tags=[]),
        File(name="statistics_page_wireframe.png", original_path=get_path("statistics_page_wireframe.png"), tags=[]),
        File(name="TODO mar30 Meeting.txt", original_path=get_path("TODO mar30 Meeting.txt"), tags=[]),
        File(name="Tututorial_2.pdf", original_path=get_path("Tututorial_2.pdf"), tags=[]),
        File(name="UseCase.png", original_path=get_path("UseCase.png"), tags=[]),
        File(name="~$ecutive summary", original_path=get_path("~$ecutive summary.docx"), tags=[]),
        File(name="~WRL0005.tmp", original_path=get_path("~WRL0005.tmp"), tags=[]),
        File(name="~WRL1847.tmp", original_path=get_path("~WRL1847.tmp"), tags=[]),
        File(name="3.6.4 Survey data to be analysed and visualised", original_path=get_path("3.6.4 Survey data to be analysed and visualised for project report mine.xlsx"), tags=[]),
        File(name="Document[1]", original_path=get_path("Document[1].pdf"), tags=[]),
        File(name="ENjoyment", original_path=get_path("ENjoyment.png"), tags=[]),
        File(name="Gantt chart", original_path=get_path("Gantt chart.png"), tags=[]),
        File(name="Gauteng", original_path=get_path("Gauteng.png"), tags=[]),
        File(name="most challanging", original_path=get_path("most challanging.png"), tags=[]),
        File(name="Most rewarding", original_path=get_path("Most rewarding.png"), tags=[]),
        File(name="Picture1", original_path=get_path("Picture1.png"), tags=[]),
        File(name="Picture2", original_path=get_path("Picture2.png"), tags=[]),
        File(name="Presentation speech", original_path=get_path("Presentation speech.docx"), tags=[]),
        File(name="Project Budget Form 2024", original_path=get_path("Project Budget Form 2024.pdf"), tags=[]),
        File(name="Taiichi ohno", original_path=get_path("Taiichi ohno.jpeg"), tags=[]),
        File(name="Week 3_Tutorial_2024_with Answers", original_path=get_path("Week 3_Tutorial_2024_with Answers.pdf"), tags=[]),
        File(name="Week 4_Tutorial_with answers", original_path=get_path("Week 4_Tutorial_with answers.pdf"), tags=[]),
        File(name="Week 5_Tutorial_2024_with answers", original_path=get_path("Week 5_Tutorial_2024_with answers.pdf"), tags=[]),
        File(name="Form", original_path=get_path("Main.form"), tags=[]),
        File(name="Main", original_path=get_path("Main.java"), tags=[]),
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
        File(name="math_5", original_path=get_path("math_5.pdf"), tags=[])
    ]

    # Some tags for the files
    tag_1 = Tag(name="COS122")
    tag_2 = Tag(name="COS122")
    files1[3].tags.append(tag_1)
    files1[12].tags.append(tag_2)

    root_dir = Directory(
        name = "test_files_3",
        path = get_path("test_files_3"),
        files = files1
    )    


    req = DirectoryRequest(root=root_dir, requestType="CLUSTERING", serverSecret=os.environ["SFM_SERVER_SECRET"], preferredCase = "KEBAB")

    yield req



def responseSizeChecker(curDir : Directory) -> int:

    count = len(curDir.files)


    # Recurisve call
    for subdir in curDir.directories:
        count += responseSizeChecker(subdir)

    return count

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
    assert responseSizeChecker(response.root) == 66

    # Check if response is well formed
    assert response.response_code == 200
    assert response.response_msg != "No file could be opened"

def test_invalid_credential_req(grpc_test_server):
    req = DirectoryRequest(root=None, requestType = "CLUSTERING", serverSecret = "wrongSecret")
    response = grpc_test_server.SendDirectoryStructure(req)
    assert response.response_code == 401
    assert response.response_msg == "Unauthorized: Incorrect server secret"
