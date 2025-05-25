# Instructions on how to use the python directory

## Dependancies

### PIP
Ensure that the python package manager pip is installed for further dependancy management. For detailed instructions on install pip please see [this](https://pip.pypa.io/en/stable/installation/).

### Python
Ensure python3 is installed. On windows this can be done via the installer, alternatively for linux use 
```
sudo apt install python3
```
Ensure that once installed python3 is added to PATH

### Pytest
Ensure pytest is installed. Using the pip package manager use
```
pip install -U pytest
```
Once again ensure pytest is added to PATH

### Other packages that require installion
* Magic (Used for extracting MIME type) : ```pip install python-magic```
* Mutagen (Used for audio metadata): ```pip install mutagen```
* PyPDF: ```pip install pypdf```
* docx: ```pip install python-docx```
* Pillow: ```pip install Pillow```
* gRPC: ```pip install grpcio```
* grpc tools: ```pip install grpcio-tools```
* protoc: ```sudo apt install -y protobuf-compiler```
* protoc gen:```go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest```
* protoc buf: ```go install google.golang.org/protobuf/cmd/protoc-gen-go@latest```
* yake: ```pip install yake```
* pdf processing: ```pip install pypdf```
* docx processing: ```pip install python-docx```



## Where to place your files
* Add all code to the src directory
* Add all pytest code to the testing directory

## Adding testing files
Ensure all tests are placed in python/src/testing and start with the prefix test_ to allow pytest to automatically collect all tests.

## gRPC (all python)
* Ensure `grpcio` and `grpcio-tools` are installed.

### Usage
* Add function to proto. `python/src/protos/your-name.proto`
* Add function definition to server. `python/src/greeter_server.py` 
* Add function call in client. `python/src/greeter_client.py`

### Build and run
**Generate pb2 files from the proto file (when there are any changes to proto)** 
```
make proto_gen
```
**Run the python server**
```
make python_server
```
**Run the python client** 
```
make python_client
```

## Build and run instructions

**Run all source files**
```
make python
```

**Run all testing**
```
make python_test
```