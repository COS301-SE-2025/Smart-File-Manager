import request_handler

def main():
    # Start the gRPC server
    handler = request_handler.RequestHandler()
    handler.serve()

if __name__ == "__main__":
    main()
