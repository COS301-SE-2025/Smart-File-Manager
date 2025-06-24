# Use Python base image
FROM python:3.11-slim-bookworm

# Install system packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic1 \
    file \
    build-essential \
    gcc \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Set working directory inside container
WORKDIR /app

# Copy Python code
COPY python/src/ ./src/
COPY python/testing/ ./testing/

# Copy shared protos
COPY protos/ ./protos/

# Use app as PATH
ENV PYTHONPATH=/app

# Install Python dependencies
RUN pip install --no-cache-dir \
    pytest \
    python-magic \
    mutagen \
    pypdf \
    python-docx \
    Pillow \
    pymediainfo \
    yake \
    scikit-learn \
    sentence-transformers \
    grpcio \
    grpcio-tools

# Generate gRPC code
RUN python -m grpc_tools.protoc \
    -I./protos \
    --python_out=./src \
    --pyi_out=./src \
    --grpc_python_out=./src \
    ./protos/message_structure.proto


# Expose gRPC server port
EXPOSE 50051

# Start the gRPC server`
CMD ["python", "src/main.py"]

# Old command for running tests directly in docker (taking this out for now as the workflow will do this)
# CMD ["pytest", "-v", "-s", "--color=yes", "--tb=short", "testing/"]



