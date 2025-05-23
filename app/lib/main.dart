import 'package:flutter/material.dart';
import 'package:app/navigation/shell.dart';
import 'package:window_manager/window_manager.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await windowManager.ensureInitialized();

  WindowOptions windowOptions = WindowOptions(
    minimumSize: Size(900, 600),
    size: Size(900, 600),
    center: true,
    backgroundColor: Colors.transparent,
  );
  windowManager.waitUntilReadyToShow(windowOptions, () async {
    await windowManager.show();
    await windowManager.focus();
  });

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
