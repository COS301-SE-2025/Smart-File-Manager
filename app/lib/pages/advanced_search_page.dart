import 'package:flutter/material.dart';

class AdvancedSearchPage extends StatelessWidget {
  const AdvancedSearchPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Advanced Search',
            style: TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: 20),
          Text(
            'Welcome to your file manager search!',
            style: TextStyle(color: Colors.grey),
          ),
        ],
      ),
    );
  }
}
