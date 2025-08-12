class FileModel {
  final String name;
  final String filePath;

  FileModel({required this.name, required this.filePath});

  factory FileModel.fromJson(Map<String, dynamic> json) {
    return FileModel(name: json['file_name'], filePath: json['file_path']);
  }

  Map<String, dynamic> toJsonForDeleteDuplicate() {
    return {'name': name, 'file_path': filePath};
  }
}
