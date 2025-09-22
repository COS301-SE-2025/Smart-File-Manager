<p align="center">
  <img src="assets/spark_logo_dark_background.png" alt="Company Logo" width="300" height=100%/>
</p>

# gRPC endpoint Documentation 

**Version:** 3.0.0.0  
**Prepared By:** Spark Industries  
**Prepared For:** Southern Cross Solutions & Personal Development team use  

## Content
* [Introduction](#introduction)
* [A note on security](#a-note-on-security)
* [Data_Sent_and_Received](#data-sent-and-received)
* [Protocol_Buffer_Files](#protocol-buffer-files-proto)

## Introduction
Our team makes use of gRPC for communication between Go and Python. Our team decided on using gRPC rather than traditional HTTP requests due to the following advantages:

* High Performance: gRPC uses protocol buffers and HTTP/2 which allows us to keep smart managing files as fast and efficient as possible.
* Cross Platform: Since our project combined both Go and Python gRPC is an ideal candidate for this project.
* Contract-First Development: Allows us as developers to agree on the exact format of data sent and received by request and responses promoting a well-structured codebase and concurrent development by the go and python team.
* Reduced network usage: Protocol buffers use binary serialization which reduces network traffic.

This document serves as the exact documentation for the services and messages sent and received. 

## A note on security
Our application is deployed as a standalone desktop application which does not rely on (or use) internet connection. All servers simply run on localhost and cannot be accessed by any malicious user, except possibly the user of the app themselves. Should the user (for some reason) decide to do so they can do no more damage to their system then what they could do using the standard file explorer. Like any application it is protected from unauthorized use by the standard system login functionality. It is for this reason that our endpoints are unprotected and does not require authentication via a secret key based solution or other similar methods. Please note that we have discussed this decision with Mr. Avinash Singh who is both, involved in both the COS301 lecturing team, and COS330 module coordinator. He has approved this decision from a security point of view, noting that we should perhaps add a login to the application itself, a feature we will be adding for demo4.

## Data sent and received
Go must extract and send the following information to python processor:
* Current directory tree structure
* Filenames contained in directory tree (includes the extension)
* Tags defined and added for specific file

After reorganising the folder structure based on output from the clustering algorithm the following information must be sent back to go for creating the new file strucutre.
* Newly organised directory tree structure (not necessarily containing the same folders as before)
* Filenames contained in the directory tree (Includes the old location of the file on disk)
* Tags defined for original files

The following JSON serves as an example of what information will be sent. Note that JSON is not over gRPC and this merely serves as an example of how the data that we're sending might be structured.

```
{
  "name": "root",
  "path": "/root",
  "files": [
    {
      "name": "readme.md",
      "path": "/root/readme.md",
      "tags": ["intro"]
    }
  ],
  "directories": [
    {
      "name": "docs",
      "path": "/root/docs",
      "files": [
        {
          "name": "report.pdf",
          "path": "/root/docs/report.pdf",
          "tags": ["finance", "q1"]
        },
        {
          "name": "notes.txt",
          "path": "/root/docs/notes.txt",
          "tags": ["meeting"]
        }
      ],
      "directories": []
    },
    {
      "name": "images",
      "path": "/root/images",
      "files": [
        {
          "name": "logo.png",
          "path": "/root/images/logo.png",
          "tags": ["brand"]
        },
        {
          "name": "graph.jpg",
          "path": "/root/images/graph.jpg",
          "tags": ["analytics"]
        }
      ],
      "directories": []
    }
  ]
}
```

## Protocol Buffer Files (.proto)
The following .proto files can be created to represent the information that needs to be communicated as explained [here](#data-sent-and-received):

```
syntax = "proto3";

// Simple tags defined by users
message Tag {
  string name = 1;
}

// Key-val metadata as extracted (don't use maps as this is more portable for go)
message MetadataEntry {
  string key = 1;
  string value = 2;
}

// Keyword extracted from file along with the score it was given by yake
message Keyword {
  string keyword = 1;
  float score = 2;
}

message File {
  string name = 1;
  string original_path = 2;
  string new_path = 3;
  repeated Tag tags = 4;
  repeated MetadataEntry metadata = 5;
  bool is_locked = 6;
  repeated Keyword keywords = 6;
}

message Directory {
  string name = 1;
  string path = 2;
  repeated File files = 3;
  repeated Directory directories = 4;
  bool is_locked = 5;
}

// For request/response messages
// Go to python 
message DirectoryRequest {
  Directory root = 1;
  string requestType = 2; // This can be CLUSTERING or METADATA or KEYWORDS
  string prefferedCase = 3; // This can be CAMEL or PASCAL or SNAKE or KEBAB
  string serverSecret = 4;
}

// Python to go
message DirectoryResponse {
  Directory root = 1;
  int32 response_code = 2;
  string response_msg = 3;
}

service DirectoryService {
  rpc SendDirectoryStructure(DirectoryRequest) returns (DirectoryResponse);
}

```

The .proto files may then be used generate the required Go and Python code via
```
protoc --go_out=. directory.proto // Go
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. directory.proto // Python

```


Which may then be used in implementation

```
// Go
package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    pb "your-generated-go-package" // Import your generated Go protobuf package
)

func main() {
    // Set up gRPC client connection
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    // Create gRPC client
    client := pb.NewDirectoryServiceClient(conn)

    // Prepare directory structure (omitted for brevity)

    // Send directory structure to Python server
    response, err := client.SendDirectoryStructure(context.Background(), &pb.DirectoryRequest{ /* directory structure */ })
    if err != nil {
        log.Fatalf("Error sending directory structure: %v", err)
    }

    // Handle response from Python server (updated directory structure)
    // process response
}

```

```
// Python
import grpc
from concurrent import futures
import directory_pb2
import directory_pb2_grpc

class DirectoryService(directory_pb2_grpc.DirectoryServiceServicer):
    def SendDirectoryStructure(self, request, context):
        # Process directory structure (apply clustering algorithms, etc.)
        # Construct updated directory structure
        response = directory_pb2.DirectoryResponse( /* updated directory structure */ )
        return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    directory_pb2_grpc.add_DirectoryServiceServicer_to_server(DirectoryService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()

```

**Important Note**  
In .proto3 all field are optional by default. If a field is not set then the value is simply set to the field's default type. Hence if a field is not needed, for example the new_path or MetadataEntry in a DirectoryRequest message it should simply be left empty.