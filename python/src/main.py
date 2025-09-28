import request_handler
import config

def main():

    # Check server secret setup
    if config.SERVER_SECRET == None:
        print("Server secret not found. Cancelling Server Startup...")
        return

    # Start the gRPC server
    handler = request_handler.RequestHandler()
    handler.serve()

if __name__ == "__main__":
    main()
