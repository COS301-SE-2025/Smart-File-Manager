from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves):
        self.slaves = ThreadPoolExecutor(maxSlaves)

    # Takes gRPC request's root and sends it to be processed by a slave
    def submitTask(self, request):
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    # Input: Directory Request
    # Output: DirectoryResponse
    def process(self, request):

        # Stub, returns what was sent
        response_directory = request.root
        response = DirectoryResponse(root=response_directory)
        return response
    
        '''
        What is this method supposed to really do?
        * Extract metadata
        * Perform clustering
        * Generate new tree strucutre as gRPC Directory

        Preferably each of the above should be handled by its own class for seperation of concerns
        '''
    
