from concurrent import futures
import threading
from sentence_transformers import SentenceTransformer
import grpc
import message_structure_pb2_grpc

import master

# Class for handling requests to gRPC python server
# Assigns request to master for processing (currently assigns each request to a single master)

class RequestHandler(message_structure_pb2_grpc.DirectoryServiceServicer):

    def __init__(self):
        # Early initialize sentence transformer
        transformer = SentenceTransformer('all-MiniLM-L6-v2')
        # If weights are none then they use a default value
        self.master = master.Master(10, transformer, None)

    def SendDirectoryStructure(self, request, context):
        response = self.master.submit_task(request).result()
        return response

    def serve(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        message_structure_pb2_grpc.add_DirectoryServiceServicer_to_server(self, server)
        port = server.add_insecure_port("0.0.0.0:0")

        server.start()
        print(f"export PYTHON_SERVER={port}")
        with open("server.env", "w") as f:
            f.write(f"PYTHON_SERVER={port}\n")

        print("Server ready and listening...")

        # Background thread listens for shutdown 
        def wait_for_stop():
            while True:
                cmd = input().strip().upper()
                if cmd == "STOP":
                    print("Stopping server...")
                    server.stop(grace=5)
                    break

        threading.Thread(target=wait_for_stop, daemon=True).start()

        server.wait_for_termination()