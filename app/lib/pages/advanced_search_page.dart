import 'package:app/api.dart';
import 'package:flutter/material.dart';
import 'package:app/custom_widgets/file_details_panel.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/search_sub_page/folder_view_search.dart';
import 'package:app/custom_widgets/custom_search_bar.dart';
import 'package:app/custom_widgets/custom_dropdown_menu.dart';

class AdvancedSearchPage extends StatefulWidget {
  final List<String> managerNames;
  final String selectedManager;

  const AdvancedSearchPage({
    super.key,
    this.managerNames = const [],
    this.selectedManager = "",
  });

  @override
  State<AdvancedSearchPage> createState() => _AdvancedSearchPageState();
}

class _AdvancedSearchPageState extends State<AdvancedSearchPage> {
  bool _isDetailsVisible = false;
  final bool _disposed = false;
  FileTreeNode? _treeData;
  FileTreeNode? _selectedFile;
  bool _isLoading = false;
  String managername = "";
  bool _searchHappened = false;

  @override
  void initState() {
    super.initState();
    if (widget.selectedManager != "") {
      _updateSelectedManager(widget.selectedManager ?? "");
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
    Navigator.pop(context, {'action': 'navigate', 'path': folderPath});
  }

  void _handleDetailPanelClose() {
    if (!_disposed && mounted) {
      setState(() {
        _isDetailsVisible = false;
        _selectedFile = null;
      });
    }
  }

  void _callGoSearch(String query) async {
    if (query.trim().isEmpty) return;

    setState(() {
      _isLoading = true;
      _searchHappened = true;
    });

    try {
      FileTreeNode response = await Api.searchGo(managername, query);
      if (mounted) {
        setState(() {
          _treeData = response;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _treeData = null;
          _isLoading = false;
        });
      }
    }
  }

  void _updateSelectedManager(String selectedManager) {
    setState(() {
      managername = selectedManager;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body: Column(
        children: [
          _buildSearchBar(),
          Expanded(child: _buildFolderViewLayout()),
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
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Text(
                          'Advanced Search',
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
              ],
            ),
          ),
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
                  CustomDropdownMenu<String>(
                    hint: "Select Manager",
                    items:
                        widget.managerNames
                            .map(
                              (name) => DropdownMenuItem(
                                value: name,
                                child: Text(name),
                              ),
                            )
                            .toList(),
                    onChanged: (v) => _updateSelectedManager(v ?? ""),
                    value:
                        widget.selectedManager == ""
                            ? null
                            : widget.selectedManager,
                  ),
                  const SizedBox(width: 12),
                  CustomSearchBar(
                    icon: Icons.search_rounded,
                    hint: "Search for files inside manager",
                    isActive: managername.isNotEmpty,
                    onChanged: (s) => _callGoSearch(s),
                  ),
                ],
              ),
            ),
          ],
        ),
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

    if (!_searchHappened) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Search results will appear here',
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
    return FolderViewSearch(
      managerName: managername,
      treeData:
          _treeData ??
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
      showGoToFolder: false,
      currentBreadcrumbs: [],
      managerPath: "",
    );
  }

  Widget _buildDetailsPanel() {
    return FileDetailsPanel(
      managerName: managername,
      selectedFile: _selectedFile,
      isVisible: _isDetailsVisible,
      onClose: _handleDetailPanelClose,
    );
  }
}
