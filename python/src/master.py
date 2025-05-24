from message_structure_pb2 import DirectoryResponse, Directory

def process(request):
    response_directory = request.root
    response = DirectoryResponse(root=response_directory)
    return response