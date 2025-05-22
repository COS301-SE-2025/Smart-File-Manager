import 'package:flutter/material.dart';

class ManagerPage extends StatelessWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            name,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 20),
          Text(
            'Welcome to your $name dashboard!',
            style: const TextStyle(color: Colors.grey),
          ),
          const SizedBox(height: 20),
          const Text(
            'This is where you can manage files and configurations for this smart manager.',
            style: TextStyle(color: Colors.grey),
          ),
        ],
      ),
    );
  }
}
