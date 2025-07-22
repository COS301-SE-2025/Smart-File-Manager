class StartupResponse {
  final String responseMessage;
  final List<String> managerNames;

  StartupResponse({
    required this.responseMessage,
    required this.managerNames,
  });

  factory StartupResponse.fromJson(Map<String, dynamic> json) {
    return StartupResponse(
      responseMessage: json['responseMessage'] as String,
      managerNames: List<String>.from(json['managerNames'] as List),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'responseMessage': responseMessage,
      'managerNames': managerNames,
    };
  }
}