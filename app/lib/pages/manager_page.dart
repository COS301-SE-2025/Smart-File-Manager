import 'package:flutter/material.dart';

class ManagerPage extends StatelessWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: const Color(0xff2E2E2E),
        title: Text('$name Manager'),
        centerTitle: false,
        titleTextStyle: const TextStyle(
          fontWeight: FontWeight.bold,
          fontSize: 20,
        ),
        shape: Border(bottom: BorderSide(color: Color(0xff3D3D3D), width: 1)),
      ),
    );
  }
}
