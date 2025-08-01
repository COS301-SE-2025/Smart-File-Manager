from concurrent import futures
import logging

import grpc
import message_structure_pb2
import message_structure_pb2_grpc

import master

# Class for handling requests to gRPC python server
# Assigns request to master for processing (currently assigns each request to a single master)
class RequestHandler(message_structure_pb2_grpc.DirectoryServiceServicer):

    def __init__(self):
        self.master = master.Master(10)

    def SendDirectoryStructure(self, request, context):
        response = self.master.submit_task(request).result()
        return response

    def serve(self):
        port = "50051"
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        message_structure_pb2_grpc.add_DirectoryServiceServicer_to_server(RequestHandler(), server)
        server.add_insecure_port("[::]:" + port)
        server.start()
        print("Server started, listening on " + port)
        server.wait_for_termination()
