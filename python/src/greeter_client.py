#This is Go's job

from __future__ import print_function

import logging

import grpc
import message_structure_pb2_grpc

from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry

def run():
    tag1 = Tag(name="ImFixed")
    meta1 = MetadataEntry(key="author", value="johnny")

    file1 = File(
        name="gopdoc.pdf",
        original_path="/usr/go/bin/gopdoc.pdf",
        new_path="/usr/trash/gopdoc.pdf",
        tags=[tag1],
        metadata=[meta1]
    )

    dir1 = Directory(
        name="useless_files",
        path="/usr/trash",
        files=[file1],
        directories=[]
    )
    req = DirectoryRequest(root=dir1)        
    with grpc.insecure_channel('localhost:50051') as channel:        
        stub = message_structure_pb2_grpc.DirectoryServiceStub(channel)
        response = stub.SendDirectoryStructure(req)
        print("Greeter client received: ")
        print(f"{response.root}")


if __name__ == "__main__":
    logging.basicConfig()
    run()
