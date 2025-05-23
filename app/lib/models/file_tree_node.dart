class FileTreeNode {
  final String name;
  final bool isFolder;
  final List<FileTreeNode>? children;
  final String? id;

  FileTreeNode({
    required this.name,
    required this.isFolder,
    this.children,
    this.id,
  });

  //convert json structure to dart object
  factory FileTreeNode.fromJson(Map<String, dynamic> json) {
    return FileTreeNode(
      name: json['name'],
      isFolder: json['isFolder'],
      children:
          json['children'] != null
              ? (json['children'] as List)
                  .map((child) => FileTreeNode.fromJson(child))
                  .toList()
              : null,
      id: json['id'],
    );
  }

  //used to send data back to Go server(converting back to jsin)
  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'isFolder': isFolder,
      'children': children?.map((child) => child.toJson()).toList(),
      'id': id,
    };
  }

  @override
  String toString() => 'FileTreeNode(name: $name, isFolder: $isFolder)';
}
