FROM ghcr.io/cirruslabs/flutter:3.24.4

WORKDIR /app

COPY pubspec.yaml pubspec.lock ./

# dependencies
RUN flutter pub get

COPY . .

RUN flutter packages pub run build_runner build --delete-conflicting-outputs

# Run tests
# CMD ["flutter", "test", "--reporter=expanded"]