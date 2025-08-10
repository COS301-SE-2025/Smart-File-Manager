import 'package:flutter/material.dart';
import 'package:app/constants.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/pages/manager_page_sub/folder_view_page.dart';
import 'package:app/pages/manager_page_sub/graph_view_page.dart';

class SortPreviewDialog extends StatefulWidget {
  final String managerName;
  final FileTreeNode sortedData;
  final Function(String, FileTreeNode) onApprove;
  final Function(String) onDecline;

  const SortPreviewDialog({
    required this.managerName,
    required this.sortedData,
    required this.onApprove,
    required this.onDecline,
    super.key,
  });

  @override
  State<SortPreviewDialog> createState() => _SortPreviewDialogState();
}

class _SortPreviewDialogState extends State<SortPreviewDialog> {
  int _currentView = 0; // 0 = folder, 1 = graph
  List<String> _currentPath = [];
  FileTreeNode? _selectedFile;
  bool _isDetailsVisible = false;

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

  void _handleApprove() {
    widget.onApprove(widget.managerName, widget.sortedData);
    Navigator.of(context).pop();
  }

  void _handleDecline() {
    widget.onDecline(widget.managerName);
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: kScaffoldColor,
      child: GestureDetector(
        onSecondaryTap: () {
          // Disable right-click context menu for entire dialog
        },
        child: SizedBox(
          width: MediaQuery.of(context).size.width * 0.9,
          height: MediaQuery.of(context).size.height * 0.9,
          child: Column(
            children: [
              _buildHeader(),
              Expanded(child: _buildContent()),
              _buildFooter(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: Color(0xff3D3D3D))),
      ),
      child: Row(
        children: [
          Text(
            'Sort Preview - ${widget.managerName}',
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const Spacer(),
          _buildViewToggle(),
          const SizedBox(width: 16),
          IconButton(
            onPressed: _handleDecline,
            icon: const Icon(Icons.close, color: Colors.white),
          ),
        ],
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
          color: isSelected ? kYellowText : Colors.transparent,
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

  Widget _buildContent() {
    return Container(
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: Color(0xff3D3D3D))),
      ),
      child: _currentView == 0 ? _buildFolderView() : _buildGraphView(),
    );
  }

  Widget _buildFolderView() {
    return FolderViewPage(
      treeData: widget.sortedData,
      currentPath: _currentPath,
      onFileSelected: _handleFileSelect,
      onNavigate: _handleNavigation,
      managerName: widget.managerName,
      onTagChanged: () {
        // Tags can't be changed in preview mode
      },
    );
  }

  Widget _buildGraphView() {
    return GraphViewPage(
      treeData: widget.sortedData,
      currentPath: _currentPath,
      onFileSelected: _handleFileSelect,
      onNavigate: _handleNavigation,
      managerName: widget.managerName,
      onTagChanged: () {
        // Tags can't be changed in preview mode
      },
    );
  }

  Widget _buildFooter() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          TextButton(
            onPressed: _handleDecline,
            style: TextButton.styleFrom(
              foregroundColor: Colors.grey,
              side: const BorderSide(color: Colors.grey),
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            ),
            child: const Text('Decline'),
          ),
          const SizedBox(width: 12),
          ElevatedButton(
            onPressed: _handleApprove,
            style: ElevatedButton.styleFrom(
              backgroundColor: kYellowText,
              foregroundColor: Colors.black,
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            ),
            child: const Text('Approve & Apply'),
          ),
        ],
      ),
    );
  }
}
