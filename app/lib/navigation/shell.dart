import 'package:flutter/material.dart';
import '../pages/dashboard_page.dart';
import '../pages/smart_managers_page.dart';
import '../pages/settings_page.dart';
import '../pages/advanced_search_page.dart';
import 'main_navigation.dart';
import '../pages/manager_page.dart';
import 'package:app/constants.dart';

class Shell extends StatefulWidget {
  const Shell({super.key});

  @override
  State<Shell> createState() => _ShellState();
}

class _ShellState extends State<Shell> {
  int _selectedIndex = 0; //index selected form the main menu (0 to 3)
  String? _selectedManager; //name of the Manager selected

  //pages are stored in a Widget List to send the correct widget(content) given the selected Index
  final List<Widget> _pages = [
    const DashboardPage(),
    const SmartManagersPage(),
    const AdvancedSearchPage(),
    const SettingsPage(),
  ];

  //list of the navigation items for selection states
  final List<NavigationItem> _navigationItems = [
    NavigationItem(icon: Icons.dashboard_rounded, label: 'Dashboard'),
    NavigationItem(icon: Icons.web_stories_outlined, label: 'Smart Managers'),
    NavigationItem(icon: Icons.search_outlined, label: 'Advanced Search'),
    NavigationItem(icon: Icons.settings, label: 'Settings'),
  ];

  //when main navigation item is tapped, set its index and unselect manager if one is selected
  void _onNavigationTap(int index) {
    setState(() {
      _selectedIndex = index;
      _selectedManager = null;
    });
  }

  //select the manager by passed in name and deselect main navigation if selected
  void _onManagerTap(String managerName) {
    setState(() {
      _selectedManager = managerName;
      _selectedIndex = -1;
    });
  }

  //find the active page and return its widget
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
      backgroundColor: kScaffoldColor,
      //Main Appbar with app title and login button
      appBar: AppBar(
        leading: Padding(
          padding: const EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Image.asset("images/logo.png"),
        ),
        backgroundColor: kAppBarColor,
        title: const Text("SMART FILE MANAGER"),
        centerTitle: false,
        titleTextStyle: kTitle1,
        actions: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 10, 0),
            child: FilledButton.icon(
              onPressed: () {}, //TODO: need to add login functionality
              label: const Text("Login"),
              style: FilledButton.styleFrom(backgroundColor: kScaffoldColor),
              icon: const Icon(Icons.account_circle),
            ),
          ),
        ],
        shape: Border(bottom: BorderSide(color: Color(0xff3D3D3D), width: 1)),
      ),
      body: Row(
        children: [
          //Main Navigation Widget with parmeters used to navigate
          MainNavigation(
            items: _navigationItems,
            selectedIndex: _selectedIndex,
            selectedManager: _selectedManager,
            onTap: _onNavigationTap,
            onManagerTap: _onManagerTap,
          ),
          //Page that needs to be rendered depending on navigation index
          Expanded(child: _getCurrentPage()),
        ],
      ),
    );
  }
}
