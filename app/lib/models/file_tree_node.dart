class FileTreeNode {
  final String name;
  final bool isFolder;
  final List<FileTreeNode>? children;
  final List<String>? tags;
  final String? path;
  final Map<String, String>? metadata;
  bool locked;

  FileTreeNode({
    required this.name,
    required this.isFolder,
    this.children,
    this.tags,
    this.path,
    this.metadata,
    required this.locked,
  });

  factory FileTreeNode.fromJson(Map<String, dynamic> json) {
    return FileTreeNode(
      name: json['name'] ?? '',
      path: json['path'] ?? '',
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
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'path': path,
      'isFolder': isFolder,
      'children': children?.map((child) => child.toJson()).toList(),
      'tags': tags,
      'metadata': metadata,
      'locked': locked,
    };
  }

  @override
  String toString() =>
      'FileTreeNode(name: $name, isFolder: $isFolder, locked: $locked)';

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
}
