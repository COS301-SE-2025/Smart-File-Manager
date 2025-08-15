import 'package:flutter/material.dart';
import '../pages/dashboard_page.dart';
import '../pages/smart_managers_page.dart';
import '../pages/settings_page.dart';
import '../pages/advanced_search_page.dart';
import 'main_navigation.dart';
import '../pages/manager_page.dart';
import 'package:app/constants.dart';
import 'package:app/api.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:app/models/file_tree_node.dart';

GlobalKey<DashboardPageState> globalKey = GlobalKey();

class Shell extends StatefulWidget {
  const Shell({super.key});

  @override
  State<Shell> createState() => _ShellState();
}

class _ShellState extends State<Shell> {
  int _selectedIndex = 0; //index selected form the main menu (0 to 3)
  String? _selectedManager; //name of the Manager selected
  List<String> _managerNames = []; //list of manager names from startup
  String _selectedManagerForSearch = "";

  final Map<String, FileTreeNode> _managerTreeData = {};
  final Map<String, bool> _pendingSorts = {}; // Track active sort operations
  final Map<String, FileTreeNode> _sortResults =
      {}; // Store sort results for approval

  final Uri _url = Uri.parse(
    'https://cos301-se-2025.github.io/Smart-File-Manager/',
  );

  //pages are created dynamically to pass manager names to AdvancedSearchPage
  List<Widget> get _pages => [
    DashboardPage(managerNames: _managerNames, key: globalKey),
    SmartManagersPage(
      managerTreeData: _managerTreeData,
      managerNames: _managerNames,
      pendingSorts: _pendingSorts,
      sortResults: _sortResults,
      onManagerSort: _onManagerSort,
      onSortApprove: _onSortApprove,
      onSortDecline: _onSortDecline,
      onManagerDelete: _onManagerDelete,
    ),
    AdvancedSearchPage(
      managerNames: _managerNames,
      selectedManager: _selectedManagerForSearch,
    ),
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


  //update stats
  void _updateStats() {
    if (_managerNames.isNotEmpty) {
      globalKey.currentState?.loadStatsData();
    }
  }

  //when manager is deleted updated values here:
  void _onManagerDelete(String managerName) {
    setState(() {
      _managerNames.remove(managerName);
      _managerTreeData.remove(managerName);
      _pendingSorts.remove(managerName);
      _sortResults.remove(managerName);
      
      // If the deleted manager was selected, deselect it
      if (_selectedManager == managerName) {
        _selectedManager = null;
        _selectedIndex = 0; // Go to dashboard
      }
      
      // If the deleted manager was selected for search, clear it
      if (_selectedManagerForSearch == managerName) {
        _selectedManagerForSearch = "";
      }
    });

    // Update stats to reflect the deletion
    _updateStats();
  }

  //when manager is sorted(move directory is called, update treedata for manager)
  void _onManagerSort(String managerName, FileTreeNode managerData) async {
    // Start background sorting
    setState(() {
      _pendingSorts[managerName] = true;
      _sortResults.remove(managerName); // Clear any previous results
    });

    try {
      final sortedData = await Api.sortManager(managerName);

      setState(() {
        _pendingSorts[managerName] = false;
        _sortResults[managerName] = sortedData;
      });

      // Show notification if user is not on SmartManagersPage
      if (_selectedIndex != 1) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              'Sort completed for "$managerName" - View results to approve',
            ),
            backgroundColor: kYellowText,
            duration: Duration(seconds: 4),
            action: SnackBarAction(
              label: 'View',
              onPressed: () {
                setState(() {
                  _selectedIndex = 1;
                  _selectedManager = null;
                });
              },
            ),
          ),
        );
      }
    } catch (e) {
      setState(() {
        _pendingSorts[managerName] = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Sort failed for "$managerName": $e'),
          backgroundColor: Colors.redAccent,
          duration: Duration(seconds: 3),
        ),
      );
    }
  }

  void _onSortApprove(String managerName, FileTreeNode sortedData) {
    setState(() {
      _managerTreeData[managerName] = sortedData;
      _sortResults.remove(managerName);
    });

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Sort applied for "$managerName"'),
        backgroundColor: kYellowText,
        duration: Duration(seconds: 2),
      ),
    );
  }

  void _onSortDecline(String managerName) {
    setState(() {
      _sortResults.remove(managerName);
    });

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Sort declined for "$managerName"'),
        backgroundColor: Colors.grey,
        duration: Duration(seconds: 2),
      ),
    );
  }

  Future<void> _launchUrl() async {
    if (!await launchUrl(_url, mode: LaunchMode.externalApplication)) {
      throw Exception('Could not launch $_url');
    }
  }

  //select the manager by passed in name and deselect main navigation if selected
  void _onManagerTap(String managerName, FileTreeNode? treeData) {
    setState(() {
      _selectedManager = managerName;
      _selectedIndex = -1;
    });
  }

  void _onManagerTreeDataUpdate(String managerName, FileTreeNode treeData) {
    setState(() {
      _managerTreeData[managerName] = treeData;
    });
  }

  void _onManagerNamesUpdate(List<String> managerNames) {
    setState(() {
      _managerNames = managerNames;
    });
  }

  void _updateManagerTreeData(String managerName, FileTreeNode treeData) {
    setState(() {
      _managerTreeData[managerName] = treeData;
    });
  }

  void _goToAdvancedSearch(String managerName) {
    setState(() {
      _selectedManager = null;
      _selectedIndex = 2;
      _selectedManagerForSearch = managerName;
    });
  }

  void _onManagerAdded(String managerName) {
    setState(() {
      if (!_managerNames.contains(managerName)) {
        _managerNames.add(managerName);
      }
    });
  }

  //find the active page and return its widget
  Widget _getCurrentPage() {
    if (_selectedManager != null) {
      return ManagerPage(
        key: ValueKey(_selectedManager),
        name: _selectedManager!,
        treeData: _managerTreeData[_selectedManager!],
        onTreeDataUpdate: _updateManagerTreeData,
        onGoToAdvancedSearch: _goToAdvancedSearch,
      );
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
        scrolledUnderElevation: 0,
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
            child: IconButton(
              onPressed: _launchUrl,
              icon: Icon(Icons.help_outline_rounded, color: kprimaryColor),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 10, 0),
            child: FilledButton.icon(
              onPressed: () {},
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
            updateStats: _updateStats,
            onManagerTap: _onManagerTap,
            onManagerTreeDataUpdate: _onManagerTreeDataUpdate,
            onManagerNamesUpdate: _onManagerNamesUpdate,
            onManagerAdded: _onManagerAdded,
            onManagerDelete: _onManagerDelete,
          ),
          //Page that needs to be rendered depending on navigation index
          Expanded(child: _getCurrentPage()),
        ],
      ),
    );
  }
}
