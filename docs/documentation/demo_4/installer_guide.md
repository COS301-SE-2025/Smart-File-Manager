# Complete Guide: Flutter App to Windows EXE with Servers

## Prerequisites

- Flutter SDK installed and configured
- Go installed and configured
- Python with PyInstaller installed (`pip install pyinstaller`)
- Inno Setup Compiler installed
- Bat to EXE converter (like Bat To Exe Converter or Advanced BAT to EXE Converter)

## Step 1: Prepare Build Environment

1. **Create build directory structure**

   ```
   project_root/
   ├── build/ (add this folder)
   │   └── storage/ (add this folder)
   ├── golang/ (your Go server code)
   ├── python/ (your Python server code)
   ├── app/ (your Flutter app)
   └── installer.iss
   ```

## Step 2: Update Go Server for Windows Deployment

Before compiling the Go server, you need to update the path handling code to work correctly on Windows instead of WSL paths.

1. **Update the imports in directoryToObject.go**

   Add `"runtime"` to your imports section:

   ```go
   import (
       "encoding/json"
       "fmt"
       "os"
       "path/filepath"
       "runtime"  // Add this import
       "strings"
       "time"
       // ... other imports
   )
   ```

2. **Update the ConvertToWSLPath function**

   Replace the existing `ConvertToWSLPath` function with this Windows-compatible version:

   ```go
   func ConvertToWSLPath(winPath string) string {
       winPath = strings.TrimSpace(winPath)

       // Check if running on Windows
       if runtime.GOOS == "windows" {
           // For Windows deployment: normalize path separators and clean up double slashes
           winPath = strings.ReplaceAll(winPath, "\\", "/")
           // Remove any double slashes that might have been created
           for strings.Contains(winPath, "//") {
               winPath = strings.ReplaceAll(winPath, "//", "/")
           }
           return winPath
       }

       // Original WSL conversion logic for non-Windows systems
       winPath = strings.ReplaceAll(winPath, "\\", "/")
       if len(winPath) > 2 && winPath[1] == ':' {
           drive := strings.ToLower(string(winPath[0]))
           rest := winPath[2:]
           return "/mnt/" + drive + rest
       }

       return winPath
   }
   ```

   This change ensures that:

   - On Windows: Paths are normalized with forward slashes and double slashes are removed

## Step 3: Compile Go Server

Navigate to your Go server directory (cd golang) and run:

```bash
go build -ldflags="-H=windowsgui" -o ../build/go_server.exe .
```

**Note**: The `-H=windowsgui` flag prevents console windows from appearing.

## Step 4: Compile Python Server

### 4.1 Download and Setup SentenceTransformers Model

First, download the required model to your project (cd python/src):

```python
# download_model.py
from sentence_transformers import SentenceTransformer
import os

# Create models directory in your project
models_dir = './models'
os.makedirs(models_dir, exist_ok=True)

print("Downloading all-MiniLM-L6-v2 model...")
# Download directly to your project folder for offline use
model = SentenceTransformer('all-MiniLM-L6-v2', cache_folder=models_dir)

print("✅ Model downloaded successfully!")

# Test the model
embeddings = model.encode(["Hello world", "This is a test"])
print(f"✅ Model working! Embedding shape: {embeddings.shape}")

# Show what was downloaded
print(f"\nModel saved to: {models_dir}")
for root, dirs, files in os.walk(models_dir):
    level = root.replace(models_dir, '').count(os.sep)
    if level < 3:  # Don't go too deep
        indent = ' ' * 2 * level
        print(f"{indent}{os.path.basename(root)}/")
```

Run this script:

```bash
python download_model.py
```

### 4.2 Add file runtime_hook_transformers.py

This file validates the YAKE stopwords:

```python
# runtime_hook_transformers.py
import os
import yake

yake_path = os.path.dirname(yake.__file__)
stopwords_path = os.path.join(yake_path, 'core', 'StopwordsList')

print(f"YAKE path: {yake_path}")
print(f"Stopwords path: {stopwords_path}")
print(f"Exists: {os.path.exists(stopwords_path)}")

# Test the specific file that was missing
missing_file = os.path.join(stopwords_path, 'stopwords_noLang.txt')
print(f"Missing file exists: {os.path.exists(missing_file)}")

if os.path.exists(missing_file):
    print("✅ Ready to build!")
else:
    print("❌ File still missing")
```

Keep this file in the python/src directory.

### 4.3 Create python_server.spec in python/src

Your spec file needs to include both the SentenceTransformers model and YAKE stopwords:

```python
# -*- mode: python ; coding: utf-8 -*-
import os
import yake

# Get the yake stopwords path
yake_path = os.path.dirname(yake.__file__)
stopwords_path = os.path.join(yake_path, 'core', 'StopwordsList')

a = Analysis(
    ['main.py'],
    pathex=[],
    binaries=[],
    datas=[
        ('models', 'models'),  # Include the downloaded SentenceTransformers model
        (stopwords_path, 'yake/core/StopwordsList'),  # Include YAKE stopwords files
    ],
    hiddenimports=[
        'magic', 'mutagen', 'docx', 'PIL', 'pymediainfo',
        'grpc', 'grpc_tools', 'yake', 'yake.core', 'yake.core.yake',
        'pandas', 'sklearn', 'sentence_transformers',
        'sentence_transformers.models', 'sentence_transformers.SentenceTransformer',
        'nltk', 'pyinstrument', 'fitz', 'fitz._fitz',
        'torch', 'tokenizers', 'transformers',
        'transformers.models', 'transformers.models.auto', 'numpy',
    ],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=['runtime_hook_transformers.py'],
    excludes=[],
    noarchive=False,
    optimize=0,
)

pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.datas,
    [],
    name='python_server',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=False,  # This prevents console window
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)
```

### 4.4 Update your Python code for offline model loading

Update your RequestHandler class to properly load the model in both development and PyInstaller environments:

```python
# request_handler.py
from concurrent import futures
from sentence_transformers import SentenceTransformer
import grpc
import message_structure_pb2_grpc
import os
import sys

import master

class RequestHandler(message_structure_pb2_grpc.DirectoryServiceServicer):

    def __init__(self):
        # Configure for offline mode
        os.environ["TRANSFORMERS_OFFLINE"] = "1"

        # Determine the correct path for the model
        if hasattr(sys, '_MEIPASS'):
            # When running as PyInstaller bundle
            cache_folder = os.path.join(sys._MEIPASS, 'models')
        else:
            # When running in development
            cache_folder = './models'

        print(f"Loading model with cache_folder: {cache_folder}")

        # Load SentenceTransformers model with local files only
        transformer = SentenceTransformer('all-MiniLM-L6-v2',
                                         cache_folder=cache_folder,
                                         local_files_only=True)
        self.master = master.Master(10, transformer)

    def SendDirectoryStructure(self, request, context):
        response = self.master.submit_task(request).result()
        return response

    def serve(self):
        port = "50051"
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        message_structure_pb2_grpc.add_DirectoryServiceServicer_to_server(self, server)
        server.add_insecure_port("[::]:" + port)
        server.start()
        print("Server started, listening on " + port)
        server.wait_for_termination()
```

### 4.5 Compile Python server

```bash
python -m PyInstaller --clean python_server.spec
```

Move the generated `python_server.exe` in the dist folder to your `build/` folder.

## Step 5: Build Flutter Application

1. **Build Flutter for Windows**

   ```bash
   cd flutter_app
   flutter build windows --release
   ```

2. **Copy and rename**
   ```bash
   cp -r build/windows/runner/Release ../build/FlutterApp
   ```

## Step 6: Create Enhanced Launcher with Progress Bar

### 6.1 Create launcher.bat inside the build folder

```batch
@echo off
setlocal EnableDelayedExpansion
title Smart File Manager Launcher
cd /d "%~dp0"

REM Create a simple progress display
echo.
echo  ====================================================
echo  ^|               Smart File Manager                ^|
echo  ^|                   Starting Up...                ^|
echo  ====================================================
echo.

echo [----------] 0%% - Initializing...
timeout /t 1 >nul

echo [##--------] 20%% - Starting Go server...
start /B "" "%~dp0go_server.exe"

echo [####------] 40%% - Starting Python server (this may take a while)...
start /B "" "%~dp0python_server.exe"

echo [######----] 60%% - Waiting for servers to be ready...

REM Wait for Go server (port 51000) - no timeout, wait indefinitely
:waitgo
timeout /t 2 >nul
netstat -an 2>nul | find "LISTEN" | find ":51000" >nul
if errorlevel 1 (
    echo Waiting for Go server...
    goto waitgo
)

echo [########--] 80%% - Go server ready!

REM Wait for Python server (port 50051) - no timeout, wait indefinitely
:waitpython
timeout /t 2 >nul
netstat -an 2>nul | find "LISTEN" | find ":50051" >nul
if errorlevel 1 (
    echo Waiting for Python server...
    goto waitpython
)

echo [##########] 100%% - All servers ready!
echo.
echo Launching Smart File Manager...
timeout /t 1 >nul

REM Launch Flutter app
if exist "%~dp0FlutterApp\app.exe" (
    start "" "%~dp0FlutterApp\app.exe"
    echo Flutter app launched successfully
) else (
    echo ERROR: app.exe not found in FlutterApp folder
    echo Expected location: %~dp0FlutterApp\app.exe
    pause
    goto cleanup
)

REM Wait for Flutter app to close
:checkapp
timeout /t 2 >nul
tasklist /fi "imagename eq app.exe" 2>nul | find /i "app.exe" >nul
if %errorlevel%==0 goto checkapp

:cleanup
echo.
echo App closed, cleaning up servers...
taskkill /F /IM "go_server.exe" 2>nul
taskkill /F /IM "python_server.exe" 2>nul
echo Cleanup complete.
timeout /t 2 >nul
exit /b 0
```

### 6.2 Convert BAT to EXE

Use a BAT to EXE converter with these **critical** settings:

- **Visibility**: Set to "Visible" (not hidden) so you can see the progress bar
- **Include icon**: Add your app icon (inside FlutterApp folder under assets)

Name the output `launcher.exe` and place it in the `build/` folder.

## Step 7: Create Installer with Inno Setup in root directory

### 7.1 installer.iss

```ini
[Setup]
AppName=Smart File Manager
AppVersion=1.0
AppPublisher=Southern Cross Solutions
DefaultDirName={autopf}\SmartFileManager
DefaultGroupName=Smart File Manager
OutputBaseFilename=SmartFileManagerSetup
Compression=lzma2/ultra
SolidCompression=yes
SetupIconFile=icon.ico
WizardStyle=modern
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
; Ignore file access errors during installation
IgnoreTasksCheck=yes

[Files]
Source: "build\go_server.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "build\python_server.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "build\launcher.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "build\storage\*"; DestDir: "{app}\storage"; Flags: recursesubdirs createallsubdirs ignoreversion
Source: "build\FlutterApp\*"; DestDir: "{app}\FlutterApp"; Flags: recursesubdirs createallsubdirs ignoreversion

[Icons]
Name: "{group}\Smart File Manager"; Filename: "{app}\launcher.exe"; IconFilename: "{app}\FlutterApp\data\flutter_assets\assets\logo.ico"
Name: "{commondesktop}\Smart File Manager"; Filename: "{app}\launcher.exe"; IconFilename: "{app}\FlutterApp\data\flutter_assets\assets\logo.ico"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "Create a desktop shortcut"; GroupDescription: "Additional icons:"

[Run]
Filename: "{app}\launcher.exe"; Description: "Launch Smart File Manager"; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "taskkill"; Parameters: "/F /IM go_server.exe"; Flags: runhidden skipifdoesntexist
Filename: "taskkill"; Parameters: "/F /IM python_server.exe"; Flags: runhidden skipifdoesntexist
Filename: "taskkill"; Parameters: "/F /IM app.exe"; Flags: runhidden skipifdoesntexist

[Code]
function InitializeSetup(): Boolean;
var
  ResultCode: Integer;
begin
  // Kill any running instances before installation
  Exec('taskkill', '/F /IM go_server.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  Exec('taskkill', '/F /IM python_server.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  Exec('taskkill', '/F /IM app.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  Result := True;
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  ResultCode: Integer;
begin
  if CurStep = ssPostInstall then
  begin
    // Add Windows Firewall exceptions for the servers (ignore errors)
    try
      Exec('netsh', 'advfirewall firewall add rule name="Smart File Manager - Go Server" dir=in action=allow program="' + ExpandConstant('{app}') + '\go_server.exe" enable=yes', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
      Exec('netsh', 'advfirewall firewall add rule name="Smart File Manager - Python Server" dir=in action=allow program="' + ExpandConstant('{app}') + '\python_server.exe" enable=yes', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);

      // Also add outbound rules
      Exec('netsh', 'advfirewall firewall add rule name="Smart File Manager - Go Server Out" dir=out action=allow program="' + ExpandConstant('{app}') + '\go_server.exe" enable=yes', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
      Exec('netsh', 'advfirewall firewall add rule name="Smart File Manager - Python Server Out" dir=out action=allow program="' + ExpandConstant('{app}') + '\python_server.exe" enable=yes', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    except
      // Silently continue if firewall rules fail
    end;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  ResultCode: Integer;
begin
  if CurUninstallStep = usPostUninstall then
  begin
    // Remove firewall rules during uninstall (ignore errors)
    try
      Exec('netsh', 'advfirewall firewall delete rule name="Smart File Manager - Go Server"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
      Exec('netsh', 'advfirewall firewall delete rule name="Smart File Manager - Python Server"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
      Exec('netsh', 'advfirewall firewall delete rule name="Smart File Manager - Go Server Out"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
      Exec('netsh', 'advfirewall firewall delete rule name="Smart File Manager - Python Server Out"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    except
      // Silently continue if firewall rule deletion fails
    end;
  end;
end;
```

### 7.2 Compile installer

1. Open Inno Setup Compiler
2. Load your `installer.iss` file
3. Click **Build** → **Compile**
4. The installer will be created as `SmartFileManagerSetup.exe` in the Output folder

## Step 8: Final Directory Structure

Your final `build/` directory should look like:

```
build/
├── storage/
├── FlutterApp/
│   ├── app.exe
│   ├── data/
│   │   ├── app.so
│   │   ├── flutter_assets/
│   │   │   ├── AssetManifest.bin
│   │   │   ├── AssetManifest.json
│   │   │   ├── assets/
│   │   │   │   └── logo.ico
│   │   │   ├── FontManifest.json
│   │   │   ├── fonts/
│   │   │   │   └── MaterialIcons-Regular.otf
│   │   │   ├── images/
│   │   │   │   └── logo.png
│   │   │   ├── NativeAssetsManifest.json
│   │   │   ├── NOTICES.Z
│   │   │   ├── packages/
│   │   │   │   └── cupertino_icons/
│   │   │   │       └── assets/
│   │   │   │           └── CupertinoIcons.ttf
│   │   │   └── shaders/
│   │   │       └── ink_sparkle.frag
│   │   └── icudtl.dat
│   ├── flutter_windows.dll
│   ├── screen_retriever_windows_plugin.dll
│   ├── tray_manager_plugin.dll
│   ├── url_launcher_windows_plugin.dll
│   └── window_manager_plugin.dll
├── go_server.exe
├── python_server.exe
└── launcher.exe
```