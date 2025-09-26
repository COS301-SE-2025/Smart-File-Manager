import 'package:app/models/file_model.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:io';
import 'models/file_tree_node.dart';
import 'models/startup_response.dart';
import 'models/duplicate_model.dart';
import 'models/stats_model.dart';

const uri = "http://localhost:51000";

class Api {
  static String get _goApiSecret =>
      Platform.environment['SFM_API_SECRET'] ?? '';

  static Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'apiSecret': _goApiSecret,
  };
  //Call to load tree data
  static Future<FileTreeNode> loadTreeData(String name) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/loadTreeData?name=$name"),
        headers: _headers,
      );

      if (response.statusCode == 200) {
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
      final response = await http.get(
        Uri.parse("$uri/startUp"),
        headers: _headers,
      );

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
  static Future<FileTreeNode> sortManager(
    String name, {
    String? caseType,
  }) async {
    try {
      String url = "$uri/sortTree?name=$name";
      if (caseType != null && caseType.isNotEmpty) {
        url += "&case=$caseType";
      }

      final response = await http.get(Uri.parse(url), headers: _headers);

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
        headers: _headers,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        final body = response.body.trim();
        return body == "true";
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
        Uri.parse("$uri/deleteManager?name=$name"),
        headers: _headers,
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
        headers: _headers,
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
        headers: _headers,
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
  static Future<bool> locking(String managerName, String path) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/lock?name=$managerName&path=$path"),
        headers: _headers,
      );
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

  static Future<bool> unlocking(String managerName, String path) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/unlock?name=$managerName&path=$path"),
        headers: _headers,
      );
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

  //Call to load duplicate data
  static Future<List<DuplicateModel>> loadDuplicates(String name) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/findDuplicateFiles?name=$name"),
        headers: _headers,
      );

      if (response.statusCode == 200) {
        final List<dynamic> jsonList = jsonDecode(response.body);
        return jsonList
            .map(
              (item) => DuplicateModel.fromJson(item as Map<String, dynamic>),
            )
            .toList();
      } else {
        throw Exception(
          'Failed to load duplicate data: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error loading duplicates: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Endpoint to use GO search(Simple file name search)
  static Future<FileTreeNode> searchGo(String name, String searchString) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/search?compositeName=$name&searchText=$searchString"),
        headers: _headers,
      );

      if (response.statusCode == 200) {
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

  //Endpoint to use Advanced search (keywords)
  static Future<FileTreeNode> searchAdvanced(
    String name,
    String searchString,
  ) async {
    try {
      final response = await http.get(
        Uri.parse(
          "$uri/keywordSearch?compositeName=$name&searchText=$searchString",
        ),
        headers: _headers,
      );

      if (response.statusCode == 200) {
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

  //return true/false depending if advanced search is ready
  static Future<bool> searchAdvancedReady(String name) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/isKeywordSearchReady?compositeName=$name"),
        headers: _headers,
      );
      if (response.statusCode == 200) {
        final body = response.body.trim().toLowerCase();
        return body == "true" || body == "1";
      } else {
        throw Exception('Failed to load data: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error loading tree data from loadTreeData: $e');
      print(stackTrace);
      rethrow;
    }
  }

  static Future<FileTreeNode> deleteSingleFile(
    String managerName,
    String path,
  ) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/deleteFile?name=$managerName&path=$path"),
        headers: _headers,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to delete File: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error deleting file: $e');
      print(stackTrace);
      rethrow;
    }
  }

  static Future<FileTreeNode> bulkDeleteFiles(
    String managerName,
    String jsonPaths,
  ) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/bulkDeleteFiles?name=$managerName"),
        headers: _headers,
        body: jsonPaths,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to delete files: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error deleting files: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Bulk add Tag
  static Future<FileTreeNode> bulkAddTag(
    String managerName,
    String jsonPaths,
  ) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/bulkAddTag?name=$managerName"),
        headers: _headers,
        body: jsonPaths,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to delete files: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error deleting files: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Bulk remove Tag
  static Future<FileTreeNode> bulkRemoveTag(
    String managerName,
    String jsonPaths,
  ) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/bulkRemoveTag?name=$managerName"),
        headers: _headers,
        body: jsonPaths,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        return FileTreeNode.fromJson(
          jsonDecode(response.body) as Map<String, dynamic>,
        );
      } else {
        throw Exception('Failed to remove tags: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error removing tags: $e');
      print(stackTrace);
      rethrow;
    }
  }

  static Future<List<FileModel>> bulkOperation(
    String name,
    String type,
    bool umbrella,
  ) async {
    try {
      final response = await http.get(
        Uri.parse("$uri/returnType?name=$name&type=$type&umbrella=$umbrella"),
        headers: _headers,
      );

      if (response.statusCode == 200) {
        if (response.body == "" || response.body == "null") {
          return <FileModel>[];
        }
        try {
          final dynamic jsonData = jsonDecode(response.body);
          print(jsonData);

          if (jsonData == null) {
            return <FileModel>[];
          }

          if (jsonData is! List) {
            return <FileModel>[];
          }

          final List<dynamic> jsonList = jsonData;
          return jsonList
              .map((item) => FileModel.fromJson(item as Map<String, dynamic>))
              .toList();
        } catch (e) {
          return <FileModel>[];
        }
      } else {
        throw Exception(
          'Failed to load duplicate data: HTTP ${response.statusCode}',
        );
      }
    } catch (e, stackTrace) {
      print('Error loading duplicates: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //loadStats
  static Future<ManagersStatsResponse> loadStatsData() async {
    try {
      final response = await http.get(
        Uri.parse("$uri/returnStats"),
        headers: _headers,
      );

      if (response.statusCode == 200) {
        return ManagersStatsResponse.fromJson(
          jsonDecode(response.body) as List<dynamic>,
        );
      } else {
        throw Exception('Failed to load stats: HTTP ${response.statusCode}');
      }
    } catch (e, stackTrace) {
      print('Error loading stats data: $e');
      print(stackTrace);
      rethrow;
    }
  }

  //Move directory
  static Future<bool> moveDirectory(String managerName) async {
    try {
      final response = await http.post(
        Uri.parse("$uri/moveDirectory?name=$managerName"),
        headers: _headers,
      );
      if (response.statusCode == 200 || response.statusCode == 201) {
        if (response.body == "true") {
          return true;
        } else {
          return false;
        }
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
