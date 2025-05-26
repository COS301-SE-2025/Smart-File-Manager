import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'dart:io';

//class to get name and directory that user types in
class SmartManagerInfo {
  String name;
  String directory;

  SmartManagerInfo(this.name, this.directory);
}

Future<SmartManagerInfo?> createManager(BuildContext context) {
  return showDialog<SmartManagerInfo>(
    context: context,
    builder: (BuildContext ctx) {
      return SmartManagerDialog();
    },
  );
}

class SmartManagerDialog extends StatefulWidget {
  const SmartManagerDialog({super.key});

  @override
  _SmartManagerDialogState createState() => _SmartManagerDialogState();
}

class _SmartManagerDialogState extends State<SmartManagerDialog> {
  final TextEditingController _nameController = TextEditingController();
  String? _selectedDirectory;
  bool _isDirectorySelected = false;

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }

  Future<void> _pickDirectory() async {
    String? selectedDirectory = await FilePicker.platform.getDirectoryPath();

    if (selectedDirectory != null) {
      // Check if directory is valid writable
      bool isValid = await _validateDirectory(selectedDirectory);

      if (isValid) {
        setState(() {
          _selectedDirectory = selectedDirectory;
          _isDirectorySelected = true;
        });
      } else {
        // Show error dialog
        _showDirectoryError(selectedDirectory);
      }
    }
  }

  Future<bool> _validateDirectory(String path) async {
    try {
      final directory = Directory(path);

      //directory exists
      if (!await directory.exists()) {
        return false;
      }

      // system directory
      if (_isSystemDirectory(path)) {
        return false;
      }

      final testFile = File('$path/.temp_write_test');
      try {
        await testFile.writeAsString('test');
        await testFile.delete();
        return true;
      } catch (e) {
        return false;
      }
    } catch (e) {
      return false;
    }
  }

  bool _isSystemDirectory(String path) {
    final normalizedPath = path.toLowerCase().replaceAll('\\', '/');

    // Windows
    final windowsSystemPaths = [
      'c:/windows',
      'c:/program files',
      'c:/program files (x86)',
      'c:/programdata',
      'c:/system volume information',
      'c:/recovery',
      'c:/perflogs',
      'c:/msocache',
      'c:/intel',
      'c:/amd',
      'c:/nvidia',
    ];

    // macOS
    final macSystemPaths = [
      '/system',
      '/library',
      '/usr',
      '/bin',
      '/sbin',
      '/etc',
      '/var',
      '/tmp',
      '/private',
      '/cores',
      '/dev',
      '/opt',
      '/applications',
    ];

    // Linux
    final linuxSystemPaths = [
      '/bin',
      '/boot',
      '/dev',
      '/etc',
      '/lib',
      '/lib32',
      '/lib64',
      '/proc',
      '/root',
      '/run',
      '/sbin',
      '/sys',
      '/tmp',
      '/usr',
      '/var',
      '/opt',
      '/srv',
      '/mnt',
      '/media',
    ];

    // Windows paths
    for (String systemPath in windowsSystemPaths) {
      if (normalizedPath.startsWith(systemPath)) {
        return true;
      }
    }

    // macOS paths
    for (String systemPath in macSystemPaths) {
      if (normalizedPath.startsWith(systemPath)) {
        return true;
      }
    }

    // Linux paths
    for (String systemPath in linuxSystemPaths) {
      if (normalizedPath.startsWith(systemPath)) {
        return true;
      }
    }

    // common restricted patterns
    if (normalizedPath.contains('/system/') ||
        normalizedPath.contains('\\system\\') ||
        normalizedPath.contains('/windows/') ||
        normalizedPath.contains('\\windows\\') ||
        normalizedPath.contains('program files')) {
      return true;
    }

    return false;
  }

  void _showDirectoryError(String path) {
    String errorMessage;

    if (_isSystemDirectory(path)) {
      errorMessage =
          "The selected directory is a system directory and cannot be used:\n\n$path\n\n"
          "Please choose a user directory.";
    } else {
      errorMessage =
          "The selected directory cannot be used:\n\n$path\n\n"
          "Please choose a directory where you have permissions.";
    }

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text("Invalid Directory"),
          content: Text(errorMessage),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text("OK"),
            ),
          ],
        );
      },
    );
  }

  bool _canCreate() {
    return _nameController.text.trim().isNotEmpty && _isDirectorySelected;
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: Color(0xff1E1E1E),
      title: const Text("Create Smart Manager"),
      titleTextStyle: TextStyle(
        color: Color(0xffFFB400),
        fontSize: 24,
        fontWeight: FontWeight.bold,
      ),
      contentTextStyle: TextStyle(color: Colors.white),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text("Enter manager name and select directory:"),
          const SizedBox(height: 16),
          TextField(
            style: TextStyle(color: Colors.white),
            cursorColor: Color(0xffFFB400),
            controller: _nameController,
            decoration: const InputDecoration(
              labelText: "Manager Name",
              labelStyle: TextStyle(color: Colors.white),

              hintText: "Enter a name for your manager",
              border: OutlineInputBorder(),
              focusedBorder: OutlineInputBorder(
                borderSide: BorderSide(color: Color(0xffFFB400)),
              ),
              enabledBorder: OutlineInputBorder(
                borderSide: BorderSide(color: Color(0xffFFB400)),
              ),
            ),
            onChanged: (value) {
              setState(() {});
            },
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: Text(
                  _isDirectorySelected
                      ? "Directory: ${_selectedDirectory!.split('/').last}"
                      : "No directory selected",
                  style: TextStyle(
                    color: _isDirectorySelected ? Colors.green : Colors.grey,
                  ),
                ),
              ),
              ElevatedButton(
                onPressed: _pickDirectory,
                style: ElevatedButton.styleFrom(
                  backgroundColor: Color(0xffFFB400),
                  foregroundColor: Colors.black,
                ),
                child: const Text("Browse"),
              ),
            ],
          ),
          if (_isDirectorySelected) ...[
            const SizedBox(height: 8),
            Text(
              _selectedDirectory!,
              style: const TextStyle(fontSize: 12, color: Colors.grey),
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ],
      ),
      actions: <Widget>[
        TextButton(
          onPressed: () {
            Navigator.of(context).pop(null);
          },
          style: TextButton.styleFrom(foregroundColor: Colors.grey),
          child: const Text("Cancel"),
        ),
        TextButton(
          style: TextButton.styleFrom(foregroundColor: Colors.grey),
          onPressed:
              _canCreate()
                  ? () {
                    Navigator.of(context).pop(
                      SmartManagerInfo(
                        _nameController.text.trim(),
                        _selectedDirectory!,
                      ),
                    );
                  }
                  : null,
          child: const Text("Create"),
        ),
      ],
    );
  }
}
