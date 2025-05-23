class FileTreeNode {
  final String name;
  final bool isFolder;
  final List<FileTreeNode>? children;
  final String? id;
  final List<String>? tags;

  FileTreeNode({
    required this.name,
    required this.isFolder,
    this.children,
    this.id,
    this.tags,
  });

  factory FileTreeNode.fromJson(Map<String, dynamic> json) {
    return FileTreeNode(
      name: json['name'] ?? '',
      isFolder: json['isFolder'] ?? false,
      children:
          json['children'] != null
              ? (json['children'] as List)
                  .map((child) => FileTreeNode.fromJson(child))
                  .toList()
              : null,
      id: json['id'],
      tags: json['tags'] != null ? List<String>.from(json['tags']) : [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'isFolder': isFolder,
      'children': children?.map((child) => child.toJson()).toList(),
      'id': id,
      'tags': tags,
    };
  }

  @override
  String toString() => 'FileTreeNode(name: $name, isFolder: $isFolder)';
}
