import 'package:http/http.dart' as http;
import 'dart:convert';
import 'models/file_tree_node.dart';
import 'models/startup_response.dart';

const uri = "http://localhost:51000";

class Api {
  //Call to load tree data
  static Future<FileTreeNode> loadTreeData(String name) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/loadTreeData?name=$name"),
      );

      if (response.statusCode == 200) {
        print(jsonDecode(response.body) as Map<String, dynamic>);

        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to load data: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error loading tree data from loadTreeData: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call to initialize app and get existing smart managers
  static Future<StartupResponse> startUp() async {
    try {
      final response = await http.get(Uri.parse("$uri/startUp"));

      if (response.statusCode == 200) {
        return StartupResponse.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to load data: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error loading initial data from startup: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call To Sort Tree structure
  static Future<FileTreeNode> sortManager(String name) async {
    try {
      final response = await http.get(Uri.parse("$uri/sortTree?name=$name"));

      if (response.statusCode == 200) {
        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to load data: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error sorting tree structure from sortManager: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call to add Mangager to backend
  static Future<bool> addSmartManager(String name, String path) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/addDirectory?name=$name&path=$path"),
      );
      print("$uri/addDirectory?name=$name&path=$path");
      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception(
          'Failed to add SmartManager: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error creating smart manager at addDirectory: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call to delete Mangager to backend
  static Future<bool> deleteSmartManager(String name) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/deleteDirectory?name=$name"),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception(
          'Failed to delete SmartManager: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error deleting smart manager at deleteDirectory: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call to add tag to File
  static Future<bool> addTagToFile(String name, String path, String tag) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/addTag?name=$name&path=$path&tag=$tag"),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception('Failed to add Tag: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error adding tag at addTag deleteDirectory: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Call to delete tag from File
  static Future<bool> deleteTagFromFile(
    String name,
    String path,
    String tag,
  ) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/removeTag?name=$name&path=$path&tag=$tag"),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception('Failed to delete Tag: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error deleting tag at deleteTag: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //lock files/folders
  static Future<bool> locking(String path) async {
    try {
      final response = await http.post(Uri.parse("$uri/lock?path=$path"));
      print(response.body);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception(
          'Failed to lock File/Folder: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error locking folder/file: $e');
      print(stackTrace);
      rethrow;
    }
  }

  static Future<bool> unlocking(String path) async {
    try {
      final response = await http.post(Uri.parse("$uri/unlock?path=$path"));
      print(response.body);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return true;
      } else {
        throw Exception(
          'Failed to unlock File/Folder: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error unlocking folder/file: $e');
      print(stackTrace);
      rethrow;
    }
  }
}
