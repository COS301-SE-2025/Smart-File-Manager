import 'package:flutter/material.dart';
import 'package:app/navigation/shell.dart';
import 'package:window_manager/window_manager.dart';
import 'package:tray_manager/tray_manager.dart';
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

  windowManager.setPreventClose(true);

  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> with WindowListener, TrayListener {
  @override
  void initState() {
    super.initState();
    _initTray();
    windowManager.addListener(this);
  }

  @override
  void onWindowClose() async {
    windowManager.hide();
  }

  @override
  void onTrayIconMouseDown() {
    _showWindow();
  }

  Future<void> _initTray() async {
    trayManager.addListener(this);
    await trayManager.setIcon('assets/logo.ico'); // Set tray icon
    await trayManager.setToolTip("Smart File Manger"); // Tooltip
    await trayManager.setContextMenu(
      Menu(
        items: [
          MenuItem(label: "Show App", onClick: (menuItem) => _showWindow()),
          MenuItem(label: "Exit", onClick: (menuItem) => _exitApp()),
        ],
      ),
    );
  }

  void _showWindow() {
    windowManager.show(); // Restore window
    windowManager.focus();
  }

  void _exitApp() {
    trayManager.destroy();
    windowManager.destroy();
  }

  @override
  void onTrayIconRightMouseDown() {
    trayManager.popUpContextMenu();
  }

  @override
  Widget build(BuildContext context) {
    //Basic structure, calls the Shell wich creates structure of top Appbar and Navigation
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Smart File Manager',
      theme: ThemeData(
        primaryColor: kprimaryColor,
        scaffoldBackgroundColor: kScaffoldColor,
      ),
      home: const Shell(),
    );
  }
}
