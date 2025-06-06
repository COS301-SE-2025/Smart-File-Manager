import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/folder_view_page.dart';
import 'package:app/pages/graph_view_page.dart';
import 'package:app/custom_widgets/file_details_panel.dart';
import 'package:http/http.dart' as http;
import 'dart:math';

class ManagerPage extends StatefulWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  State<ManagerPage> createState() => _ManagerPageState();
}

class _ManagerPageState extends State<ManagerPage> {
  int _currentView = 0; // 0 = folder, 1 = graph
  List<String> _currentPath = [];
  FileTreeNode? _treeData;
  FileTreeNode? _selectedFile;
  bool _isDetailsVisible = false;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadTreeData();
  }

  Future<void> _loadTreeData() async {
    final response1 = await http.get(
      Uri.parse('https://run.mocky.io/v3/b3097f03-5576-4e45-ab9e-54e12fa12d87'),
    );

    final response2 = await http.get(
      Uri.parse('https://run.mocky.io/v3/a809ac12-e410-4a79-95b3-604837f22e59'),
    );

    final randomChoice = Random().nextBool();
    final selectedResponse = randomChoice ? response1 : response2;

    if (selectedResponse.statusCode == 200) {
      setState(() {
        _treeData = FileTreeNode.fromJson(
          jsonDecode(selectedResponse.body) as Map<String, dynamic>,
        );
        _isLoading = false;
      });
    } else {
      throw Exception('Failed to load data');
    }
  }

  void _handleViewChange(int index) {
    setState(() {
      _currentView = index;
      _isDetailsVisible = false;
      _selectedFile = null;
    });
  }

  void _handleFileSelect(FileTreeNode file) {
    setState(() {
      _selectedFile = file;
      _isDetailsVisible = true;
    });
  }

  void _handleNavigation(List<String> newPath) {
    setState(() {
      _currentPath = newPath;
    });
  }

  void _handleDetailPanelClose() {
    setState(() {
      _isDetailsVisible = false;
      _selectedFile = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body: Column(
        children: [
          _buildSearchBar(),
          Expanded(
            child: Row(
              children: [
                Expanded(
                  flex: _isDetailsVisible ? 3 : 1,
                  child: _buildMainContent(),
                ),
                if (_isDetailsVisible)
                  SizedBox(width: 200, child: _buildDetailsPanel()),
              ],
            ),
          ),
        ],
      ),
    );
  }

  PreferredSize _buildTopBar() {
    return PreferredSize(
      preferredSize: const Size.fromHeight(50),
      child: AppBar(
        backgroundColor: const Color(0xff2E2E2E),
        automaticallyImplyLeading: false,
        shape: const Border(
          bottom: BorderSide(color: Color(0xff3D3D3D), width: 1),
        ),
        flexibleSpace: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Row(
              children: [
                Expanded(
                  child: Center(
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
                          padding: const EdgeInsets.symmetric(
                            vertical: 4,
                            horizontal: 8,
                          ),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(12),
                            color: const Color(0xff094E3A),
                          ),
                          child: const Text(
                            '75% organized',
                            style: TextStyle(
                              fontSize: 10,
                              color: Color(0xff6EE79B),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                _buildViewToggle(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildViewToggle() {
    return Container(
      decoration: BoxDecoration(
        color: const Color(0xff242424),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: const Color(0xff3D3D3D)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          _buildToggleButton('Folder', 0, Icons.folder),
          _buildToggleButton('Graph', 1, Icons.account_tree),
        ],
      ),
    );
  }

  Widget _buildToggleButton(String label, int index, IconData icon) {
    final isSelected = _currentView == index;
    return GestureDetector(
      onTap: () => _handleViewChange(index),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isSelected ? const Color(0xffFFB400) : Colors.transparent,
          borderRadius: BorderRadius.circular(4),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 16,
              color: isSelected ? Colors.black : const Color(0xff9CA3AF),
            ),
            const SizedBox(width: 4),
            Text(
              label,
              style: TextStyle(
                fontSize: 12,
                color: isSelected ? Colors.black : const Color(0xff9CA3AF),
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchBar() {
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
              width: 250,
              height: 32,
              padding: const EdgeInsets.symmetric(horizontal: 12),
              decoration: BoxDecoration(
                color: const Color(0xff242424),
                borderRadius: BorderRadius.circular(6),
                border: Border.all(color: const Color(0xff3D3D3D)),
              ),
              child: Center(
                child: TextField(
                  onChanged: null,
                  cursorColor: const Color(0xffFFB400),
                  style: const TextStyle(fontSize: 12, color: Colors.white),
                  decoration: const InputDecoration(
                    hintText: 'Search files and folders...',
                    hintStyle: TextStyle(
                      color: Color(0xff9CA3AF),
                      fontSize: 12,
                    ),
                    border: InputBorder.none,
                    isCollapsed: true,
                    prefixIcon: Icon(
                      Icons.search,
                      color: Color(0xff9CA3AF),
                      size: 16,
                    ),
                    prefixIconConstraints: BoxConstraints(
                      minWidth: 20,
                      minHeight: 16,
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Container(
              height: 32,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: BoxDecoration(
                color: const Color(0xff242424),
                borderRadius: BorderRadius.circular(6),
                border: Border.all(color: const Color(0xff3D3D3D)),
              ),
              child: DropdownButtonHideUnderline(
                child: DropdownButton<String>(
                  value: null,
                  hint: const Text(
                    'Sort by',
                    style: TextStyle(fontSize: 12, color: Color(0xff9CA3AF)),
                  ),
                  dropdownColor: const Color(0xff2E2E2E),
                  iconEnabledColor: const Color(0xff9CA3AF),
                  style: const TextStyle(
                    fontSize: 12,
                    color: Color(0xff9CA3AF),
                  ),
                  items: const [
                    DropdownMenuItem(value: 'name', child: Text('Name')),
                    DropdownMenuItem(value: 'size', child: Text('Size')),
                  ],
                  onChanged: null,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMainContent() {
    if (_isLoading) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(color: Color(0xffFFB400)),
            SizedBox(height: 16),
            Text(
              'Loading files...',
              style: TextStyle(color: Color(0xff9CA3AF)),
            ),
          ],
        ),
      );
    }

    if (_treeData == null) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 64, color: Color(0xffDC2626)),
            SizedBox(height: 16),
            Text(
              'Failed to load data',
              style: TextStyle(color: Colors.white, fontSize: 16),
            ),
            Text(
              'Please try again later',
              style: TextStyle(color: Color(0xff9CA3AF), fontSize: 14),
            ),
          ],
        ),
      );
    }

    switch (_currentView) {
      case 0:
        return FolderViewPage(
          treeData: _treeData!,
          currentPath: _currentPath,
          onFileSelected: _handleFileSelect,
          onNavigate: _handleNavigation,
        );
      case 1:
        return GraphViewPage();
      default:
        return Placeholder();
    }
  }

  Widget _buildDetailsPanel() {
    return FileDetailsPanel(
      selectedFile: _selectedFile,
      isVisible: _isDetailsVisible,
      onClose: _handleDetailPanelClose,
    );
  }
}
