import pytest
import sys
import os
import time
from typing import List

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

from src.message_structure_pb2 import Directory, File, DirectoryRequest


def create_real_files(tmp_path, num_files: int) -> List[File]:
    """Create real dummy text files with content for keyword extraction."""
    files: List[File] = []
    for i in range(num_files):
        file_path = tmp_path / f"file_{i}.txt"
        file_path.write_text(
            f"This is file {i}. It contains some test text for keyword extraction. "
            f"Keywords like project, data, clustering, testing, analysis.\n"
        )
        files.append(File(name=file_path.name, original_path=str(file_path)))
    return files


def make_request(request_type: str, files: List[File], test_dir: str) -> DirectoryRequest:
    """Helper to build a DirectoryRequest with real files."""
    root = Directory(name="scalability_test", path=test_dir, files=files)
    return DirectoryRequest(
        root=root,
        requestType=request_type,
        serverSecret=os.environ["SFM_SERVER_SECRET"]
    )


@pytest.mark.parametrize("num_files", [10, 50, 100, 200])
def test_scalability_keywords(grpc_test_server, num_files, tmp_path):
    """Benchmark keywords request with increasing number of real files."""
    files = create_real_files(tmp_path, num_files)
    req = make_request("KEYWORDS", files, str(tmp_path))

    start = time.perf_counter()
    response = grpc_test_server.SendDirectoryStructure(req)
    duration = time.perf_counter() - start

    print(f"[KEYWORDS] {num_files} files -> {duration:.4f} seconds")

    assert response.response_code == 200
    assert response.response_msg != "No file could be opened"
    assert len(response.root.files) == num_files
