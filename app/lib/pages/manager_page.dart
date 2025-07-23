import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/manager_page_sub/folder_view_page.dart';
import 'package:app/pages/manager_page_sub/graph_view_page.dart';
import 'package:app/custom_widgets/file_details_panel.dart';
import 'package:app/api.dart';
import 'package:app/custom_widgets/hoverable_button.dart';

class ManagerPage extends StatefulWidget {
  final String name;
  final FileTreeNode? treeData;
  final Function(String, FileTreeNode)? onTreeDataUpdate;
  const ManagerPage({
    required this.name,
    this.treeData,
    this.onTreeDataUpdate,
    super.key,
  });
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
  bool _isSorting = false;
  bool _disposed = false;

  @override
  void dispose() {
    _disposed = true;
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    if (widget.treeData != null) {
      setState(() {
        _treeData = widget.treeData;
        _isLoading = false;
      });
    } else {
      getTree(); // Fallback method to get the data if not exists
    }
  }

  @override
  void didUpdateWidget(ManagerPage oldWidget) {
    super.didUpdateWidget(oldWidget);

    // If receive new tree data, update state
    if (widget.treeData != null && widget.treeData != oldWidget.treeData) {
      if (!_disposed && mounted) {
        setState(() {
          _treeData = widget.treeData;
          _isLoading = false;
        });
      }
    }

    // If manager name changed and do not have tree data, load it
    if (widget.name != oldWidget.name && widget.treeData == null) {
      getTree();
    }
  }

  Future<void> getTree() async {
    if (!_disposed && mounted) {
      setState(() {
        _isLoading = true;
      });
    }

    try {
      FileTreeNode response = await Api.loadTreeData(widget.name);

      if (!_disposed && mounted) {
        setState(() {
          _treeData = response;
          _isLoading = false;
        });

        // Update the parent loaded tree data
        widget.onTreeDataUpdate?.call(widget.name, response);
      }
    } catch (e) {
      if (!_disposed && mounted) {
        setState(() {
          _isLoading = false;
        });
      }
      print('Error loading tree data: $e');
    }
  }

  Future<void> _handleSortManager() async {
    if (!_disposed && mounted) {
      setState(() {
        _isLoading = true;
        _isSorting = true;
      });
    }

    try {
      FileTreeNode response = await Api.sortManager(widget.name);

      if (!_disposed && mounted) {
        setState(() {
          _treeData = response;
          _isLoading = false;
          _isSorting = false;
        });

        widget.onTreeDataUpdate?.call(widget.name, response);
      }
    } catch (e) {
      if (!_disposed && mounted) {
        setState(() {
          _isLoading = false;
          _isSorting = false;
        });
      }
      print('Error sorting manager: $e');
    }
  }

  void _handleViewChange(int index) {
    if (!_disposed && mounted) {
      setState(() {
        _currentView = index;
        _isDetailsVisible = false;
        _selectedFile = null;
      });
    }
  }

  void _handleFileSelect(FileTreeNode file) {
    if (!_disposed && mounted) {
      setState(() {
        _selectedFile = file;
        _isDetailsVisible = true;
      });
    }
  }

  void _handleNavigation(List<String> newPath) {
    if (!_disposed && mounted) {
      setState(() {
        _currentPath = newPath;
      });
    }
  }

  void _handleDetailPanelClose() {
    if (!_disposed && mounted) {
      setState(() {
        _isDetailsVisible = false;
        _selectedFile = null;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body: Column(
        children: [
          _buildSearchBar(),
          Expanded(
            child:
                _currentView == 0
                    ? _buildFolderViewLayout()
                    : _buildGraphViewLayout(),
          ),
        ],
      ),
    );
  }

  Widget _buildFolderViewLayout() {
    return Row(
      children: [
        Expanded(flex: _isDetailsVisible ? 3 : 1, child: _buildMainContent()),
        if (_isDetailsVisible)
          SizedBox(width: 200, child: _buildDetailsPanel()),
      ],
    );
  }

  Widget _buildGraphViewLayout() {
    return Stack(
      children: [
        _buildMainContent(),

        if (_isDetailsVisible)
          Positioned(
            top: 0,
            right: 0,
            bottom: 0,
            width: 200,
            child: Container(
              decoration: BoxDecoration(
                color: const Color(0xff2E2E2E),
                border: const Border(
                  left: BorderSide(color: Color(0xff3D3D3D), width: 1),
                ),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.3),
                    blurRadius: 10,
                    offset: const Offset(-2, 0),
                  ),
                ],
              ),
              child: _buildDetailsPanel(),
            ),
          ),
      ],
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
                        const SizedBox(width: 12),
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
            Expanded(
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
                        style: const TextStyle(
                          fontSize: 12,
                          color: Colors.white,
                        ),
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
                          style: TextStyle(
                            fontSize: 12,
                            color: Color(0xff9CA3AF),
                          ),
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
            HoverableButton(
              onTap: _isSorting ? null : _handleSortManager,
              name: _isSorting ? "Sorting..." : "Sort Manager",
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMainContent() {
    if (_isLoading) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const CircularProgressIndicator(color: Color(0xffFFB400)),
            const SizedBox(height: 16),
            Text(
              _isSorting ? 'Sorting files...' : 'Loading files...',
              style: const TextStyle(color: Color(0xff9CA3AF)),
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
        return GraphViewPage(
          treeData: _treeData!,
          currentPath: _currentPath,
          onFileSelected: _handleFileSelect,
          onNavigate: _handleNavigation,
        );
      default:
        return const Placeholder();
    }
  }

  Widget _buildDetailsPanel() {
    return FileDetailsPanel(
      managerName: widget.name,
      selectedFile: _selectedFile,
      isVisible: _isDetailsVisible,
      onClose: _handleDetailPanelClose,
    );
  }
}
