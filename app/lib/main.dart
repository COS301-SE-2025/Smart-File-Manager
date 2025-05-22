import 'package:flutter/material.dart';
import 'package:app/shell.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Smart File Manager',
      theme: ThemeData(
        primaryColor: const Color(0xffFFB400),
        scaffoldBackgroundColor: const Color(0xff1E1E1E),
      ),
      home: const Shell(),
    );
  }
}
