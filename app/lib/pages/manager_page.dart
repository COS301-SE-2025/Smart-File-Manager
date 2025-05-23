import 'package:flutter/material.dart';
import 'package:toggle_switch/toggle_switch.dart';
import 'dart:convert'; //JSON

import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/folder_view_page.dart';
import 'package:app/pages/graph_view_page.dart';
import 'package:app/custom_widgets/file_details_panel.dart';

class ManagerPage extends StatefulWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  State<ManagerPage> createState() => _ManagerPageState();
}

class _ManagerPageState extends State<ManagerPage> {
  int _currentView = 0; //folder or graph
  List<String> _currentPath = []; //breadcrums
  FileTreeNode? _treeData; //full tree structure
  FileTreeNode? _selectedFile; //for details of file
  bool _isDetailsVisible = false;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadTreeData();
  }

  //API - Replace with api call (use mocking data now)
  Future<void> _loadTreeData() async {
    //netword delay test
    await Future.delayed(const Duration(seconds: 3));

    //Mock JSON
    final mockJsonData = {
      "name": "root",
      "isFolder": true,
      "id": "root",
      "children": [
        {"name": "file1.txt", "isFolder": false, "id": "file1"},
        {"name": "file2.docx", "isFolder": false, "id": "file2"},
        {
          "name": "Documents",
          "isFolder": true,
          "id": "documents",
          "children": [
            {"name": "resume.pdf", "isFolder": false, "id": "resume"},
            {
              "name": "Projects",
              "isFolder": true,
              "id": "projects",
              "children": [
                {"name": "project1.docx", "isFolder": false, "id": "project1"},
                {"name": "project2.xlsx", "isFolder": false, "id": "project2"},
              ],
            },
          ],
        },
        {
          "name": "Pictures",
          "isFolder": true,
          "id": "pictures",
          "children": [
            {"name": "vacation.jpg", "isFolder": false, "id": "vacation"},
            {
              "name": "Family",
              "isFolder": true,
              "id": "family",
              "children": [
                {"name": "photo1.png", "isFolder": false, "id": "photo1"},
                {"name": "photo2.png", "isFolder": false, "id": "photo2"},
              ],
            },
          ],
        },
      ],
    };

    //parse jjson to class
    setState(() {
      _treeData = FileTreeNode.fromJson(mockJsonData);
      _isLoading = false;
    });
  }

  //changing view
  void _handleViewChange(int? index) {
    setState(() {
      _currentView = index ?? 0;
      _isDetailsVisible = false;
      _selectedFile = null;
    });
  }

  //file selection
  void _handleFileSelect(FileTreeNode file) {
    setState(() {
      _selectedFile = file;
      _isDetailsVisible = true;
    });
  }

  //directory navigation
  void _handleNavigation(List<String> newPath) {
    setState(() {
      _currentPath = newPath;
    });
  }

  //detailsPanelClose
  void _handleDetailPanelClose() {
    setState(() {
      _isDetailsVisible = false;
      _selectedFile = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _topBar(),
      body: Stack(
        children: [
          Column(children: [_searchBar(), Expanded(child: _mainContent())]),
          _detailsPanel(),
        ],
      ),
    );
  }

  PreferredSize _topBar() {
    return PreferredSize(
      preferredSize: const Size.fromHeight(50),
      child: AppBar(
        backgroundColor: const Color(0xff2E2E2E),
        automaticallyImplyLeading: false,
        shape: const Border(
          bottom: BorderSide(color: Color(0xff3D3D3D), width: 1),
        ),
        flexibleSpace: Center(
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Row(
                children: [
                  Expanded(
                    child: Row(
                      children: [
                        Text(
                          '${widget.name} Manager',
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 20,
                            color: Colors.white,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Container(
                          width: 100,
                          height: 20,
                          padding: const EdgeInsets.symmetric(
                            vertical: 2,
                            horizontal: 8,
                          ),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(10),
                            color: const Color(0xff094E3A),
                          ),
                          child: const Center(
                            child: Text(
                              '75% organized',
                              style: TextStyle(
                                fontSize: 10,
                                color: Color(0xff6EE79B),
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  MouseRegion(
                    cursor: SystemMouseCursors.click,
                    child: ToggleSwitch(
                      initialLabelIndex: _currentView,
                      minWidth: 90,
                      minHeight: 30,
                      fontSize: 12,
                      inactiveBgColor: const Color(0xff242424),
                      inactiveFgColor: const Color(0xff9CA3AF),
                      borderColor: const [Color(0xff3D3D3D)],
                      cornerRadius: 4,
                      borderWidth: 1,
                      totalSwitches: 2,
                      labels: const ['Folder View', 'Graph View'],
                      onToggle: _handleViewChange,
                    ),
                  ),
                  const SizedBox(width: 8),
                  SizedBox(
                    height: 33,
                    child: TextButton(
                      onPressed: () {},
                      style: TextButton.styleFrom(
                        backgroundColor: const Color(0xff242424),
                        side: const BorderSide(
                          color: Color(0xff3D3D3D),
                          width: 1,
                        ),
                        foregroundColor: const Color(0xff9CA3AF),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(4),
                        ),
                      ),
                      child: const Text(
                        "Sort",
                        style: TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.normal,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _searchBar() {
    return Container(
      height: 50,
      decoration: const BoxDecoration(
        color: Color(0xff2E2E2E),
        border: Border(bottom: BorderSide(color: Color(0xff3D3D3D), width: 1)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16.0),
        child: Row(
          children: [
            Container(
              width: 200,
              height: 33,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: BoxDecoration(
                color: const Color(0xff242424),
                borderRadius: BorderRadius.circular(4),
                border: Border.all(color: const Color(0xff3D3D3D), width: 1),
              ),
              child: Center(
                child: TextField(
                  cursorColor: const Color(0xffFFB400),
                  style: const TextStyle(
                    fontSize: 12,
                    color: Color(0xff9CA3AF),
                  ),
                  decoration: const InputDecoration(
                    hintText: 'Search...',
                    hintStyle: TextStyle(
                      color: Color(0xff9CA3AF),
                      fontSize: 12,
                    ),
                    border: InputBorder.none,
                    isCollapsed: true,
                  ),
                  //onChanged:  Need code
                ),
              ),
            ),
            const SizedBox(width: 12),
            Container(
              height: 33,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: BoxDecoration(
                color: const Color(0xff242424),
                borderRadius: BorderRadius.circular(4),
                border: Border.all(color: const Color(0xff3D3D3D), width: 1),
              ),
              child: DropdownButtonHideUnderline(
                child: DropdownButton<String>(
                  //value: need code
                  dropdownColor: const Color(0xff2E2E2E),
                  iconEnabledColor: const Color(0xff9CA3AF),
                  style: const TextStyle(
                    fontSize: 12,
                    color: Color(0xff9CA3AF),
                  ),
                  items: const [
                    DropdownMenuItem(
                      value: 'Name',
                      child: Text('Sort by Name'),
                    ),
                    DropdownMenuItem(
                      value: 'Size',
                      child: Text('Sort by Size'),
                    ),
                    DropdownMenuItem(
                      value: 'Date Modified',
                      child: Text('Sort by Date Modified'),
                    ),
                    DropdownMenuItem(
                      value: 'Date Created',
                      child: Text('Sort by Date Created'),
                    ),
                  ],
                  onChanged: null, //Code needed
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _mainContent() {
    if (_isLoading) {
      return const Center(
        child: CircularProgressIndicator(color: Color(0xffFFB400)),
      );
    }

    if (_treeData == null) {
      return const Center(
        child: Text(
          'Failed to load data',
          style: TextStyle(color: Colors.white),
        ),
      );
    }

    // Step 6: View Switching Logic
    switch (_currentView) {
      case 0:
        return FolderViewPage(
          treeData: _treeData!,
          currentPath: _currentPath,
          onFileSelected: _handleFileSelect,
          onNavigate: _handleNavigation,
        );
      case 1:
        return GraphViewPage(
          treeData: _treeData!,
          onFileSelected: _handleFileSelect,
        );
      default:
        return Container();
    }
  }

  Widget _detailsPanel() {
    return FileDetailsPanel(
      selectedFile: _selectedFile,
      isVisible: _isDetailsVisible,
      onClose: _handleDetailPanelClose,
    );
  }
}
