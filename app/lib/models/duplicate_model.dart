class DuplicateModel {
  final String name;
  final String originalPath;
  final String duplicatePath;

  DuplicateModel({
    required this.name,
    required this.originalPath,
    required this.duplicatePath,
  });

  factory DuplicateModel.fromJson(Map<String, dynamic> json) {
    return DuplicateModel(
      name: json['name'],
      originalPath: json['original'],
      duplicatePath: json['duplicate'],
    );
  }

  Map<String, dynamic> toJsonForDeleteDuplicate() {
    return {'filepath': duplicatePath};
  }
}
