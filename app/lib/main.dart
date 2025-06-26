import 'package:flutter/material.dart';
import 'package:app/navigation/shell.dart';
import 'package:window_manager/window_manager.dart';
import 'constants.dart';

void main() async {
  //Package used to set minimum screen size
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
    //Basic structure, calls the Shell wich creates structure of top Appbar and Navigation
    return MaterialApp(
      title: 'Smart File Manager',
      theme: ThemeData(
        primaryColor: kprimaryColor,
        scaffoldBackgroundColor: kScaffoldColor,
      ),
      home: const Shell(),
    );
  }
}
