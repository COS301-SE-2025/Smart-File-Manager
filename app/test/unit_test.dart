import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mockito/annotations.dart';

import 'package:app/models/file_tree_node.dart';
import 'package:app/custom_widgets/create_manager.dart';
import 'package:app/navigation/main_navigation.dart';

//flutter packages pub run build_runner build

@GenerateMocks([http.Client])
void main() {
  group('FileTreeNode Tests', () {
    test('should create FileTreeNode from JSON correctly', () {
      final json = {
        'name': 'test_file.txt',
        'path': '/home/user/test_file.txt',
        'isFolder': false,
        'tags': ['important', 'work'],
        'metadata': {'size': '1024', 'mimeType': 'text/plain'},
        'children': null,
      };

      final node = FileTreeNode.fromJson(json);

      expect(node.name, equals('test_file.txt'));
      expect(node.path, equals('/home/user/test_file.txt'));
      expect(node.isFolder, isFalse);
      expect(node.tags, equals(['important', 'work']));
      expect(node.metadata?['size'], equals('1024'));
      expect(node.metadata?['mimeType'], equals('text/plain'));
      expect(node.children, isNull);
    });

    test('should create folder FileTreeNode with children', () {
      final json = {
        'name': 'Documents',
        'path': '/home/user/Documents',
        'isFolder': true,
        'children': [
          {
            'name': 'file1.txt',
            'path': '/home/user/Documents/file1.txt',
            'isFolder': false,
            'tags': [],
            'metadata': {},
          },
        ],
        'tags': [],
        'metadata': {},
      };

      final node = FileTreeNode.fromJson(json);

      expect(node.name, equals('Documents'));
      expect(node.isFolder, isTrue);
      expect(node.children, isNotNull);
      expect(node.children!.length, equals(1));
      expect(node.children!.first.name, equals('file1.txt'));
      expect(node.children!.first.isFolder, isFalse);
    });

    test('should convert FileTreeNode to JSON correctly', () {
      final node = FileTreeNode(
        name: 'test.dart',
        path: '/project/test.dart',
        isFolder: false,
        tags: ['code', 'dart'],
        metadata: {'size': '2048', 'mimeType': 'text/x-dart'},
      );

      final json = node.toJson();

      expect(json['name'], equals('test.dart'));
      expect(json['path'], equals('/project/test.dart'));
      expect(json['isFolder'], isFalse);
      expect(json['tags'], equals(['code', 'dart']));
      expect(json['metadata']['size'], equals('2048'));
    });

    test('should handle empty or null values gracefully', () {
      final json = {'name': '', 'isFolder': false};

      final node = FileTreeNode.fromJson(json);

      expect(node.name, equals(''));
      expect(node.path, equals(''));
      expect(node.isFolder, isFalse);
      expect(node.children, isNull);
      expect(node.tags, equals([]));
      expect(node.metadata, equals({}));
    });

    test('toString should return correct format', () {
      final node = FileTreeNode(name: 'test_folder', isFolder: true);

      final result = node.toString();

      expect(result, equals('FileTreeNode(name: test_folder, isFolder: true)'));
    });
  });

  group('SmartManagerInfo Tests', () {
    test('should create SmartManagerInfo with name and directory', () {
      final info = SmartManagerInfo('My Manager', '/home/user/documents');

      expect(info.name, equals('My Manager'));
      expect(info.directory, equals('/home/user/documents'));
    });

    test('should allow modification of name and directory', () {
      final info = SmartManagerInfo('Original', '/original/path');

      info.name = 'Modified';
      info.directory = '/modified/path';

      expect(info.name, equals('Modified'));
      expect(info.directory, equals('/modified/path'));
    });
  });

  group('NavigationItem Tests', () {
    test('should create NavigationItem correctly', () {
      final item = NavigationItem(icon: Icons.dashboard, label: 'Dashboard');

      expect(item.icon, equals(Icons.dashboard));
      expect(item.label, equals('Dashboard'));
    });

    test('should create ManagerNavigationItem with additional properties', () {
      final item = ManagerNavigationItem(
        icon: Icons.folder,
        label: 'Project Manager',
        directory: '/home/user/projects',
        isLoading: true,
      );

      expect(item.icon, equals(Icons.folder));
      expect(item.label, equals('Project Manager'));
      expect(item.directory, equals('/home/user/projects'));
      expect(item.isLoading, isTrue);
    });

    test('should have default isLoading value as false', () {
      final item = ManagerNavigationItem(
        icon: Icons.folder,
        label: 'Test Manager',
        directory: '/test/path',
      );

      expect(item.isLoading, isFalse);
    });
  });

  group('API Tests', () {
    test('should construct correct URLs for API calls', () {
      const expectedBaseUri = "http://localhost:51000";
      const managerName = "test_manager";
      const path = "/test/path";
      const tag = "important";

      expect(
        "$expectedBaseUri/loadTreeData?name=$managerName",
        equals("http://localhost:51000/loadTreeData?name=test_manager"),
      );

      expect(
        "$expectedBaseUri/addDirectory?name=$managerName&path=$path",
        equals(
          "http://localhost:51000/addDirectory?name=test_manager&path=/test/path",
        ),
      );

      expect(
        "$expectedBaseUri/deleteDirectory?name=$managerName",
        equals("http://localhost:51000/deleteDirectory?name=test_manager"),
      );

      expect(
        "$expectedBaseUri/addTag?name=$managerName&path=$path&tag=$tag",
        equals(
          "http://localhost:51000/addTag?name=test_manager&path=/test/path&tag=important",
        ),
      );

      expect(
        "$expectedBaseUri/deleteTag?name=$managerName&path=$path&tag=$tag",
        equals(
          "http://localhost:51000/deleteTag?name=test_manager&path=/test/path&tag=important",
        ),
      );
    });
  });

  group('File Path Utilities Tests', () {
    test('should convert WSL path to Windows path correctly', () {
      String convertWSLPath(String wslPath) {
        final match = RegExp(r"^/mnt/([a-zA-Z])/").firstMatch(wslPath);
        if (match != null) {
          final driveLetter = match.group(1)!.toUpperCase();
          final windowsPath = wslPath
              .replaceFirst(RegExp(r"^/mnt/[a-zA-Z]/"), "$driveLetter:/")
              .replaceAll('/', r'\');
          return windowsPath;
        }
        return wslPath;
      }

      expect(
        convertWSLPath('/mnt/c/Users/John/Documents/file.txt'),
        equals(r'C:\Users\John\Documents\file.txt'),
      );

      expect(
        convertWSLPath('/mnt/d/Projects/app/main.dart'),
        equals(r'D:\Projects\app\main.dart'),
      );

      expect(
        convertWSLPath('/home/user/file.txt'),
        equals('/home/user/file.txt'),
      );
    });

    test('should handle file size formatting correctly', () {
      String formatFileSize(String sizeStr) {
        try {
          int sizeInBytes = int.parse(sizeStr);
          if (sizeInBytes == 0) return '0 bytes';

          const List<String> units = ['bytes', 'KB', 'MB', 'GB', 'TB'];
          int unitIndex = 0;
          double size = sizeInBytes.toDouble();

          while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
          }

          if (unitIndex == 0) {
            return '${size.toInt()} ${units[unitIndex]}';
          } else {
            return '${size.toStringAsFixed(1)} ${units[unitIndex]}';
          }
        } catch (e) {
          return sizeStr;
        }
      }

      expect(formatFileSize('0'), equals('0 bytes'));
      expect(formatFileSize('512'), equals('512 bytes'));
      expect(formatFileSize('1024'), equals('1.0 KB'));
      expect(formatFileSize('1536'), equals('1.5 KB'));
      expect(formatFileSize('1048576'), equals('1.0 MB'));
      expect(formatFileSize('1073741824'), equals('1.0 GB'));
      expect(formatFileSize('invalid'), equals('invalid'));
    });

    test('should format dates correctly', () {
      String formatDate(String dateStr) {
        try {
          DateTime dateTime = DateTime.parse(dateStr);
          return '${dateTime.day}/${dateTime.month}/${dateTime.year} ${dateTime.hour.toString().padLeft(2, '0')}:${dateTime.minute.toString().padLeft(2, '0')}';
        } catch (e) {
          return dateStr;
        }
      }

      expect(formatDate('2024-03-15T14:30:00.000Z'), equals('15/3/2024 14:30'));
      expect(formatDate('2024-12-01T09:05:00.000Z'), equals('1/12/2024 09:05'));
      expect(formatDate('invalid-date'), equals('invalid-date'));
    });
  });

  group('File Type Detection Tests', () {
    test('should return correct file icons based on extension', () {
      IconData getFileIcon(String fileName) {
        final extension = fileName.split('.').last.toLowerCase();

        switch (extension) {
          case 'dart':
            return Icons.code;
          case 'json':
            return Icons.data_object;
          case 'yaml':
          case 'yml':
            return Icons.settings;
          case 'md':
            return Icons.article;
          case 'png':
          case 'jpg':
          case 'jpeg':
          case 'gif':
            return Icons.image;
          case 'pdf':
            return Icons.picture_as_pdf;
          default:
            return Icons.description;
        }
      }

      expect(getFileIcon('main.dart'), equals(Icons.code));
      expect(getFileIcon('package.json'), equals(Icons.data_object));
      expect(getFileIcon('config.yaml'), equals(Icons.settings));
      expect(getFileIcon('config.yml'), equals(Icons.settings));
      expect(getFileIcon('README.md'), equals(Icons.article));
      expect(getFileIcon('image.png'), equals(Icons.image));
      expect(getFileIcon('photo.JPG'), equals(Icons.image));
      expect(getFileIcon('document.pdf'), equals(Icons.picture_as_pdf));
      expect(getFileIcon('unknown.xyz'), equals(Icons.description));
      expect(getFileIcon('no_extension'), equals(Icons.description));
    });
  });

  group('Graph Depth Calculation Tests', () {
    test('should calculate node depths correctly with max depth limit', () {
      const int maxDepth = 3;

      int findNodeDepthFromRoot(
        FileTreeNode current,
        FileTreeNode target,
        int currentDepth,
      ) {
        if (current == target) {
          return currentDepth;
        }

        if (current.children != null && currentDepth < maxDepth) {
          for (FileTreeNode child in current.children!) {
            int result = findNodeDepthFromRoot(child, target, currentDepth + 1);
            if (result != -1) {
              return result;
            }
          }
        }

        return -1;
      }

      final leafNode = FileTreeNode(name: 'leaf.txt', isFolder: false);
      final childNode = FileTreeNode(
        name: 'child',
        isFolder: true,
        children: [leafNode],
      );
      final rootNode = FileTreeNode(
        name: 'root',
        isFolder: true,
        children: [childNode],
      );

      expect(findNodeDepthFromRoot(rootNode, rootNode, 0), equals(0));
      expect(findNodeDepthFromRoot(rootNode, childNode, 0), equals(1));
      expect(findNodeDepthFromRoot(rootNode, leafNode, 0), equals(2));
    });
  });

  group('System Directory Validation Tests', () {
    test('should identify system directories correctly', () {
      bool isSystemDirectory(String path) {
        final normalizedPath = path.toLowerCase().replaceAll('\\', '/');

        final windowsSystemPaths = [
          'c:/windows',
          'c:/program files',
          'c:/program files (x86)',
        ];

        final linuxSystemPaths = ['/bin', '/usr', '/etc', '/var'];

        for (String systemPath in windowsSystemPaths) {
          if (normalizedPath.startsWith(systemPath)) {
            return true;
          }
        }

        for (String systemPath in linuxSystemPaths) {
          if (normalizedPath.startsWith(systemPath)) {
            return true;
          }
        }

        return false;
      }

      // Windows system paths
      expect(isSystemDirectory('C:/Windows/System32'), isTrue);
      expect(isSystemDirectory('C:/Program Files/App'), isTrue);
      expect(isSystemDirectory('C:/Program Files (x86)/App'), isTrue);

      // Linux system paths
      expect(isSystemDirectory('/bin/bash'), isTrue);
      expect(isSystemDirectory('/usr/local/bin'), isTrue);
      expect(isSystemDirectory('/etc/config'), isTrue);
      expect(isSystemDirectory('/var/log'), isTrue);

      // User directories (should be allowed)
      expect(isSystemDirectory('C:/Users/John/Documents'), isFalse);
      expect(isSystemDirectory('/home/user/projects'), isFalse);
      expect(isSystemDirectory('/Users/john/Desktop'), isFalse);
    });
  });

  group('Color Utility Tests', () {
    test('should generate level colors correctly', () {
      const List<Color> levelColors = [
        Color(0xFFFFB400), // kprimaryColor
        Color(0xFF3498DB),
        Color(0xFF2ECC71),
        Color(0xFFE74C3C),
        Color(0xFF9B59B6),
      ];

      Color getLevelColor(int depth) {
        return levelColors[depth % levelColors.length];
      }

      expect(getLevelColor(0), equals(Color(0xFFFFB400)));
      expect(getLevelColor(1), equals(Color(0xFF3498DB)));
      expect(getLevelColor(5), equals(Color(0xFFFFB400))); // Wraps around
      expect(getLevelColor(7), equals(Color(0xFF2ECC71))); // 7 % 5 = 2
    });
  });
}

class TestHelpers {
  static FileTreeNode createTestFileTree() {
    return FileTreeNode(
      name: 'root',
      isFolder: true,
      children: [
        FileTreeNode(
          name: 'documents',
          isFolder: true,
          children: [
            FileTreeNode(
              name: 'file1.txt',
              isFolder: false,
              tags: ['important'],
              metadata: {'size': '1024', 'mimeType': 'text/plain'},
            ),
            FileTreeNode(
              name: 'file2.pdf',
              isFolder: false,
              tags: ['work', 'document'],
              metadata: {'size': '2048', 'mimeType': 'application/pdf'},
            ),
          ],
        ),
        FileTreeNode(
          name: 'pictures',
          isFolder: true,
          children: [
            FileTreeNode(
              name: 'photo.jpg',
              isFolder: false,
              tags: ['vacation'],
              metadata: {'size': '5120', 'mimeType': 'image/jpeg'},
            ),
          ],
        ),
      ],
    );
  }

  static Map<String, dynamic> createTestApiResponse() {
    return {
      'name': 'test_root',
      'path': '/test',
      'isFolder': true,
      'children': [
        {
          'name': 'test_file.txt',
          'path': '/test/test_file.txt',
          'isFolder': false,
          'tags': ['test'],
          'metadata': {'size': '100', 'mimeType': 'text/plain'},
        },
      ],
      'tags': [],
      'metadata': {},
    };
  }
}
