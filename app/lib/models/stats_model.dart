class ManagersStatsResponse {
  List<StatsModel>? managers;

  ManagersStatsResponse({this.managers});

  ManagersStatsResponse.fromJson(List<dynamic> json) {
    managers = [];
    for (var v in json) {
      managers!.add(StatsModel.fromJson(v));
    }
  }

  List<dynamic> toJson() {
    return managers?.map((v) => v.toJson()).toList() ?? [];
  }

  @override
  String toString() {
    return 'ManagersStatsResponse(managers: $managers)';
  }
}

class StatsModel {
  String? managerName;
  int? size;
  int? folders;
  int? files;
  List<Files>? recent;
  List<Files>? largest;
  List<Files>? oldest;
  List<int>? umbrellaCounts;

  StatsModel({
    this.managerName,
    this.size,
    this.folders,
    this.files,
    this.recent,
    this.largest,
    this.oldest,
    this.umbrellaCounts,
  });

  StatsModel.fromJson(Map<String, dynamic> json) {
    managerName = json['manager_name'];
    size = json['size'];
    folders = json['folders'];
    files = json['files'];
    if (json['recent'] != null) {
      recent = <Files>[];
      json['recent'].forEach((v) {
        recent!.add(Files.fromJson(v));
      });
    }
    if (json['largest'] != null) {
      largest = <Files>[];
      json['largest'].forEach((v) {
        largest!.add(Files.fromJson(v));
      });
    }
    if (json['oldest'] != null) {
      oldest = <Files>[];
      json['oldest'].forEach((v) {
        oldest!.add(Files.fromJson(v));
      });
    }
    umbrellaCounts = json['umbrella_counts']?.cast<int>();
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['manager_name'] = managerName;
    data['size'] = size;
    data['folders'] = folders;
    data['files'] = files;
    if (recent != null) {
      data['recent'] = recent!.map((v) => v.toJson()).toList();
    }
    if (largest != null) {
      data['largest'] = largest!.map((v) => v.toJson()).toList();
    }
    if (oldest != null) {
      data['oldest'] = oldest!.map((v) => v.toJson()).toList();
    }
    data['umbrella_counts'] = umbrellaCounts;
    return data;
  }

  @override
  String toString() {
    return 'StatsModel(managerName: $managerName, size: $size, folders: $folders, files: $files, '
        'recent: $recent, largest: $largest, oldest: $oldest, umbrellaCounts: $umbrellaCounts)';
  }
}

class Files {
  String? filePath;
  String? fileName;

  Files({this.filePath, this.fileName});

  Files.fromJson(Map<String, dynamic> json) {
    filePath = json['file_path'];
    fileName = json['file_name'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['file_path'] = filePath;
    data['file_name'] = fileName;
    return data;
  }

  @override
  String toString() {
    return 'Files(filePath: $filePath, fileName: $fileName)';
  }
}
