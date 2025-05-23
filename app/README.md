# Smart File Manager(SFM) - Flutter App

SFM is a cross-platform desktop application built with Flutter, supporting macOS, Windows, and Linux.

## Prerequisites

Before you can run this project, you'll need to have Flutter installed and configured for desktop development.

### 1. Install Flutter SDK

Download and install Flutter from the official website:
- **Visit**: https://flutter.dev/docs/get-started/install
- **Follow** the installation guide for your operating system (use the Vscode instructions, it is the easiest)

### 2. Enable Desktop Support

After installing Flutter, enable desktop support by running these commands in your terminal:

```bash
flutter config --enable-windows-desktop
flutter config --enable-macos-desktop  
flutter config --enable-linux-desktop
```

## Getting Started

### 1. Install Dependencies

```bash
flutter pub get
```

### 2. Verify Setup

Check that everything is properly configured:

```bash
flutter doctor
```

Resolve any issues that appear in the output before proceeding.

### 3. Run the Application

#### Run on your current platform:
```bash
flutter run
```

#### Run on a specific platform:
```bash
# Windows
flutter run -d windows

# macOS
flutter run -d macos

# Linux
flutter run -d linux
```