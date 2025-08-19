class FileModel {
  final String name;
  final String filePath;
  final List<String>? fileTags;

  FileModel({
    required this.name,
    required this.filePath,
    required this.fileTags,
  });

  factory FileModel.fromJson(Map<String, dynamic> json) {
    return FileModel(
      name: json['file_name'],
      filePath: json['file_path'],
      fileTags:
          json['file_tags'] != null ? List<String>.from(json['file_tags']) : [],
    );
  }

  Map<String, dynamic> toJsonForDeleteDuplicate() {
    return {'name': name, 'file_path': filePath, 'file_tags': fileTags};
  }
}
