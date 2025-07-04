class FileTreeNode {
  final String name;
  final bool isFolder;
  final List<FileTreeNode>? children;
  final List<String>? tags;
  final String? path;
  final Map<String, String>? metadata;

  FileTreeNode({
    required this.name,
    required this.isFolder,
    this.children,
    this.tags,
    this.path,
    this.metadata,
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
    };
  }

  @override
  String toString() => 'FileTreeNode(name: $name, isFolder: $isFolder)';
}
