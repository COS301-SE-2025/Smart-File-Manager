class FileTreeNode {
  final String name;
  final bool isFolder;
  final List<FileTreeNode>? children;
  final List<String>? tags;
  String? path;
  final String? rootPath;
  final Map<String, String>? metadata;
  bool locked;
  final String? newPath;

  FileTreeNode({
    required this.name,
    required this.isFolder,
    this.children,
    this.tags,
    this.path,
    this.rootPath,
    this.metadata,
    required this.locked,
    this.newPath,
  });

  factory FileTreeNode.fromJson(Map<String, dynamic> json) {
    return FileTreeNode(
      name: json['name'] ?? '',
      path: json['path'] ?? '',
      rootPath: json['rootPath'] ?? '',
      isFolder: json['isFolder'] ?? false,
      children:
          json['children'] != null
              ? (json['children'] as List)
                  .map((child) => FileTreeNode.fromJson(child))
                  .toList()
              : null,
      tags: json['tags'] != null ? List<String>.from(json['tags']) : [],
      metadata:
          json['metadata'] != null
              ? Map<String, String>.from(json['metadata'])
              : {},
      locked: json['locked'] ?? false,
      newPath: json['newPath'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'path': path,
      'rootPath': rootPath,
      'isFolder': isFolder,
      'children': children?.map((child) => child.toJson()).toList(),
      'tags': tags,
      'metadata': metadata,
      'locked': locked,
      'newPath': newPath,
    };
  }

  @override
  String toString() {
    return _stringify(this, 0);
  }

  String _stringify(FileTreeNode node, int indent) {
    final indentStr = '  ' * indent;

    final buffer = StringBuffer();
    buffer.writeln('$indentStr- FileTreeNode(');
    buffer.writeln('$indentStr  name: ${node.name}');
    buffer.writeln('$indentStr  path: ${node.path}');
    buffer.writeln('$indentStr  rootPath: ${node.rootPath}');
    buffer.writeln('$indentStr  isFolder: ${node.isFolder}');
    buffer.writeln('$indentStr  locked: ${node.locked}');
    buffer.writeln('$indentStr  tags: ${node.tags}');
    buffer.writeln('$indentStr  metadata: ${node.metadata}');
    buffer.writeln('$indentStr  newPath: ${node.newPath}');

    if (node.children != null && node.children!.isNotEmpty) {
      buffer.writeln('$indentStr  children: [');
      for (var child in node.children!) {
        buffer.write(_stringify(child, indent + 2));
      }
      buffer.writeln('$indentStr  ]');
    } else {
      buffer.writeln('$indentStr  children: []');
    }

    buffer.writeln('$indentStr)');
    return buffer.toString();
  }

  // Function used to set the locked state of the structure to true (dependent on api lock endpoint)
  void lockItem(String path) {
    _findAndLockItem(this, path);
  }

  bool _findAndLockItem(FileTreeNode node, String targetPath) {
    // Check if target node
    if (node.path == targetPath) {
      _lockNodeAndDescendants(node);
      return true;
    }

    if (node.children != null) {
      for (var child in node.children!) {
        if (_findAndLockItem(child, targetPath)) {
          return true;
        }
      }
    }

    return false;
  }

  void _lockNodeAndDescendants(FileTreeNode node) {
    node.locked = true;
    if (node.children != null) {
      for (var child in node.children!) {
        _lockNodeAndDescendants(child);
      }
    }
  }

  // Function used to set the unlocked state of the structure to true (dependent on api unlock endpoint)
  void unlockItem(String path) {
    _findAndUnlockItem(this, path);
  }

  bool _findAndUnlockItem(FileTreeNode node, String targetPath) {
    // Check if target node
    if (node.path == targetPath) {
      _unlockNodeAndDescendants(node);
      return true;
    }

    if (node.children != null) {
      for (var child in node.children!) {
        if (_findAndUnlockItem(child, targetPath)) {
          return true;
        }
      }
    }

    return false;
  }

  void _unlockNodeAndDescendants(FileTreeNode node) {
    node.locked = false;
    if (node.children != null) {
      for (var child in node.children!) {
        _unlockNodeAndDescendants(child);
      }
    }
  }

  void replaceOldPath() {
    if (rootPath != null && rootPath!.isNotEmpty) {
      _updateNodePathsAndDescendants(this, rootPath!);
    }
  }

  void _updateNodePathsAndDescendants(FileTreeNode node, String baseRoot) {
    if (node.rootPath == null || node.rootPath!.isEmpty) {
      int lastSlash = baseRoot.lastIndexOf("/");
      if (lastSlash != -1) {
        String parentDir = baseRoot.substring(0, lastSlash);
        node.path = "$parentDir/${node.newPath}";
      } else {
        // No slash found → treat baseRoot as root
        node.path = "${node.newPath}";
      }
    }

    if (node.children != null) {
      for (var child in node.children!) {
        _updateNodePathsAndDescendants(child, node.path ?? baseRoot);
      }
    }
  }
}
