import 'package:http/http.dart' as http;
import 'dart:convert';
import 'models/file_tree_node.dart';

class Api {
  //APi call for generating tree structure
  static Future<FileTreeNode> loadTreeData(String endpoint) async {
    final response = await http.get(Uri.parse(endpoint));

    if (response.statusCode == 200) {
      return FileTreeNode.fromJson(
        jsonDecode(response.body) as Map<String, dynamic>,
      );
    } else {
      throw Exception('Failed to load data');
    }
  }
}
