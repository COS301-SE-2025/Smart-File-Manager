import 'package:flutter/material.dart';
import 'pages/dashboard_page.dart';
import 'pages/smart_managers_page.dart';
import 'pages/settings_page.dart';
import 'pages/advanced_search_page.dart';
import 'main_navigation.dart';
import 'pages/manager_page.dart';

class Shell extends StatefulWidget {
  const Shell({super.key});

  @override
  State<Shell> createState() => _ShellState();
}

class _ShellState extends State<Shell> {
  int _selectedIndex = 0;
  String? _selectedManager;

  final List<Widget> _pages = [
    const DashboardPage(),
    const SmartManagersPage(),
    const AdvancedSearchPage(),
    const SettingsPage(),
  ];

  final List<NavigationItem> _navigationItems = [
    NavigationItem(icon: Icons.dashboard_rounded, label: 'Dashboard'),
    NavigationItem(icon: Icons.web_stories_outlined, label: 'Smart Managers'),
    NavigationItem(icon: Icons.search_outlined, label: 'Advanced Search'),
    NavigationItem(icon: Icons.settings, label: 'Settings'),
  ];

  void _onNavigationTap(int index) {
    setState(() {
      _selectedIndex = index;
      _selectedManager = null;
    });
  }

  void _onManagerTap(String managerName) {
    setState(() {
      _selectedManager = managerName;
      _selectedIndex = -1;
    });
  }

  Widget _getCurrentPage() {
    if (_selectedManager != null) {
      return ManagerPage(name: _selectedManager!);
    } else if (_selectedIndex >= 0 && _selectedIndex < _pages.length) {
      return _pages[_selectedIndex];
    } else {
      return _pages[0];
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xff1E1E1E),
      appBar: AppBar(
        leading: Padding(
          padding: const EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Image.asset("images/logo.png"),
        ),
        backgroundColor: const Color(0xff2E2E2E),
        title: const Text("SMART FILE MANAGER"),
        centerTitle: false,
        titleTextStyle: const TextStyle(
          color: Color(0xffFFB400),
          fontWeight: FontWeight.bold,
          fontSize: 20,
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 10, 0),
            child: FilledButton.icon(
              onPressed: () {},
              label: const Text("Login"),
              style: FilledButton.styleFrom(
                backgroundColor: const Color(0xff242424),
              ),
              icon: const Icon(Icons.account_circle),
            ),
          ),
        ],
      ),
      body: Row(
        children: [
          MainNavigation(
            items: _navigationItems,
            selectedIndex: _selectedIndex,
            selectedManager: _selectedManager,
            onTap: _onNavigationTap,
            onManagerTap: _onManagerTap,
          ),
          Expanded(child: _getCurrentPage()),
        ],
      ),
    );
  }
}
