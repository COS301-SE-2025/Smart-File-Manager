import 'package:app/custom_widgets/duplicate_dialog.dart';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/manager_page_sub/folder_view_page.dart';
import 'package:app/pages/manager_page_sub/graph_view_page.dart';
import 'package:app/custom_widgets/file_details_panel.dart';
import 'package:app/api.dart';
import 'package:app/custom_widgets/hoverable_button.dart';
import 'package:app/custom_widgets/custom_search_bar.dart';
import 'package:app/pages/search_sub_page/folder_view_search.dart';
import 'package:app/custom_widgets/bulk_dialog.dart';

class ManagerPage extends StatefulWidget {
  final String name;
  final FileTreeNode? treeData;
  final Function(String, FileTreeNode)? onTreeDataUpdate;
  final Function(String)? onGoToAdvancedSearch;
  const ManagerPage({
    required this.name,
    this.treeData,
    this.onTreeDataUpdate,
    this.onGoToAdvancedSearch,
    super.key,
  });
  @override
  State<ManagerPage> createState() => _ManagerPageState();
}

class _ManagerPageState extends State<ManagerPage> {
  int _currentView = 0; // 0 = folder, 1 = graph
  List<String> _currentPath = [];
  FileTreeNode? _treeData;
  FileTreeNode? _searchTreeData;
  bool _searchHappened = false;
  FileTreeNode? _selectedFile;
  bool _isDetailsVisible = false;
  bool _isLoading = true;
  bool _disposed = false;
  late final ScrollController _scrollController;
  late final TextEditingController _searchController;

  @override
  void dispose() {
    _scrollController.dispose();
    _searchController.dispose();
    _disposed = true;
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _scrollController = ScrollController();
    _searchController = TextEditingController();
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

  void _handleGoToFolder(List<String> folderPath) {
    if (!_disposed && mounted) {
      setState(() {
        _searchHappened = false;
        _searchTreeData = null;
        _searchController.clear();
        _currentPath = folderPath;
        _currentView = 0;
        _selectedFile = null;
        _isDetailsVisible = false;
      });
    }
  }

  void _handleNavigation(List<String> newPath) {
    if (!_disposed && mounted) {
      setState(() {
        _currentPath = newPath;
        print(_currentPath);
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

  void _showDuplicateDialog(String name) async {
    showDialog<String>(
      context: context,
      builder:
          (context) => DuplicateDialog(
            name: name,
            updateOnDuplicateDelete: _updateOnDuplicateDelete,
          ),
    );
  }

  void _showBulkDialog(String name) async {
    showDialog<String>(
      context: context,
      builder:
          (context) => BulkDialog(
            name: name,
            type: "Documents",
            umbrella: true,
            updateOnDelete: _updateOnDuplicateDelete,
          ),
    );
  }

  void _updateOnDuplicateDelete(String managerName, FileTreeNode treeData) {
    if (!_disposed && mounted) {
      setState(() {
        _treeData = treeData;
      });
    }
    widget.onTreeDataUpdate?.call(managerName, treeData);
  }

  void _callGoSearch(String query) async {
    if (query.trim().isEmpty) {
      setState(() {
        _isLoading = false;
        _searchHappened = false;
        _searchTreeData = null;
      });
      return;
    }

    setState(() {
      _isLoading = true;
      _searchHappened = true;
    });

    try {
      FileTreeNode response = await Api.searchGo(widget.name, query);
      if (mounted) {
        setState(() {
          _searchTreeData = response;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _searchTreeData = null;
          _isLoading = false;
        });
      }
    }
  }

  Widget mainContent() {
    if (_searchHappened == true) {
      return _buildFolderViewSearch();
    } else {
      return _currentView == 0
          ? _buildFolderViewLayout()
          : _buildGraphViewLayout();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body: Column(
        children: [_buildSearchBar(), Expanded(child: mainContent())],
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

  Widget _buildFolderViewSearch() {
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
        scrolledUnderElevation: 0,
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
      padding: const EdgeInsets.symmetric(horizontal: 16.0),
      child: Row(
        children: [
          CustomSearchBar(
            icon: Icons.search_rounded,
            hint: "Search for files inside manager",
            isActive: true,
            controller: _searchController,
            onChanged: (s) => _callGoSearch(s),
          ),
          const SizedBox(width: 12),
          HoverableButton(
            onTap: () {
              widget.onGoToAdvancedSearch?.call(widget.name);
            },
            name: "Advanced Search",
            icon: Icons.manage_search_rounded,
          ),
          const VerticalDivider(color: Color(0xff3D3D3D)),
          Expanded(
            child: Scrollbar(
              thickness: 2,
              thumbVisibility: true,
              interactive: true,
              controller: _scrollController,
              child: SingleChildScrollView(
                controller: _scrollController,
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    HoverableButton(
                      onTap: () {
                        _showDuplicateDialog(widget.name);
                      },
                      name: "Find Duplicates",
                      icon: Icons.filter_none_rounded,
                    ),
                    const SizedBox(width: 12),
                    HoverableButton(
                      onTap: () {
                        _showBulkDialog(widget.name);
                      },
                      name: "Bulk Operations",
                      icon: Icons.factory_rounded,
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
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
              'Loading files...',
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

    if (_searchHappened == true) {
      return FolderViewSearch(
        managerName: widget.name,
        treeData:
            _searchTreeData ??
            FileTreeNode(
              name: '',
              path: '',
              isFolder: false,
              locked: false,
              children: [],
            ),
        onFileSelected: _handleFileSelect,
        onTagChanged: () {
          // Trigger rebuild of details panel when tags change
          if (mounted) setState(() {});
        },
        onGoToFolder: _handleGoToFolder,
        showGoToFolder: true,
        currentBreadcrumbs: _currentPath,
        managerPath: widget.treeData!.rootPath ?? "",
      );
    }

    switch (_currentView) {
      case 0:
        return FolderViewPage(
          treeData: _treeData!,
          currentPath: _currentPath,
          onFileSelected: _handleFileSelect,
          onNavigate: _handleNavigation,
          managerName: widget.name,
          onTagChanged: () {
            // Trigger rebuild of details panel when tags change
            if (mounted) setState(() {});
          },
        );
      case 1:
        return GraphViewPage(
          treeData: _treeData!,
          currentPath: _currentPath,
          onFileSelected: _handleFileSelect,
          onNavigate: _handleNavigation,
          managerName: widget.name,
          onTagChanged: () {
            // Trigger rebuild of details panel when tags change
            if (mounted) setState(() {});
          },
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
