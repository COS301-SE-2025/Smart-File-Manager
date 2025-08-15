import 'package:app/constants.dart';
import 'package:app/custom_widgets/hoverable_button.dart';
import 'package:app/custom_widgets/sort_preview_dialog.dart';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/api.dart';

class SmartManagersPage extends StatefulWidget {
  Map<String, FileTreeNode> managerTreeData = {};
  List<String> managerNames;
  Map<String, bool> pendingSorts = {};
  Map<String, FileTreeNode> sortResults = {};
  final Function(String, FileTreeNode) onManagerSort;
  final Function(String, FileTreeNode) onSortApprove;
  final Function(String) onSortDecline;
  final Function(String)? onManagerDelete;

  SmartManagersPage({
    super.key,
    required this.managerTreeData,
    required this.managerNames,
    required this.pendingSorts,
    required this.sortResults,
    required this.onManagerSort,
    required this.onSortApprove,
    required this.onSortDecline,
    this.onManagerDelete,
  });

  @override
  State<SmartManagersPage> createState() => _SmartManagersPageState();
}

class _SmartManagersPageState extends State<SmartManagersPage> {
  List<FileTreeNode> _currentItems = [];

  @override
  void initState() {
    super.initState();
    _updateCurrentItems();
  }

  @override
  void didUpdateWidget(SmartManagersPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.managerTreeData != oldWidget.managerTreeData ||
        widget.pendingSorts != oldWidget.pendingSorts ||
        widget.sortResults != oldWidget.sortResults) {
      _updateCurrentItems();
    }
  }

  void _handleSortManager(String managerName) {
    // Get the tree data for this manager (we don't actually use it for the API call)
    final treeData = widget.managerTreeData[managerName];
    if (treeData != null) {
      widget.onManagerSort(managerName, treeData);
    }
  }

  void _showSortPreview(String managerName) {
    final sortResult = widget.sortResults[managerName];
    if (sortResult != null) {
      showDialog(
        context: context,
        barrierDismissible: false,
        builder:
            (context) => SortPreviewDialog(
              managerName: managerName,
              sortedData: sortResult,
              onApprove: widget.onSortApprove,
              onDecline: widget.onSortDecline,
            ),
      );
    }
  }

  void _handleDeleteManager(String managerName) {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            title: const Text('Delete Smart Manager'),
            content: Text(
              'Are you sure you want to delete "$managerName"? This action cannot be undone.',
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('Cancel'),
              ),
              TextButton(
                onPressed: () async {
                  Navigator.of(context).pop();
                  await _deleteManager(managerName);
                },
                style: TextButton.styleFrom(foregroundColor: Colors.red),
                child: const Text('Delete'),
              ),
            ],
          ),
    );
  }

  Future<void> _deleteManager(String managerName) async {
    try {
      final success = await Api.deleteSmartManager(managerName);

      if (success) {
        widget.onManagerDelete?.call(managerName);

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Smart Manager "$managerName" deleted successfully'),
            backgroundColor: kYellowText,
            duration: Duration(seconds: 2),
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to delete Smart Manager "$managerName"'),
            backgroundColor: Colors.redAccent,
            duration: Duration(seconds: 3),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error deleting Smart Manager "$managerName": $e'),
          backgroundColor: Colors.redAccent,
          duration: Duration(seconds: 3),
        ),
      );
    }
  }

  void _updateCurrentItems() {
    setState(() {
      _currentItems = List.from(widget.managerTreeData.values);
    });
  }

  @override
  Widget build(BuildContext context) {
    return widget.managerTreeData.isEmpty
        ? _buildEmptyState()
        : _buildFileGrid();
  }

  Widget _buildFileGrid() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: SizedBox(
        child: ListView.builder(
          padding: EdgeInsets.zero,
          itemCount: _currentItems.length,
          itemBuilder: (context, index) {
            final item = _currentItems[index];
            return LayoutBuilder(
              builder: (context, constraints) {
                return Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16.0,
                    vertical: 8.0,
                  ),
                  child: Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xff242424),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: const Color(0xff3D3D3D),
                        width: 1,
                      ),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            const SizedBox(width: 8),
                            Text(item.name, style: kTitle2),
                            const Spacer(),
                            Text(
                              "Managing 234 files across 24 directories.",
                              style: TextStyle(color: const Color(0xff9CA3AF)),
                            ),
                          ],
                        ),
                        Divider(color: const Color(0xff3D3D3D), thickness: 1),
                        Column(
                          //crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            _buildSortButton(item.name),
                            SizedBox(height: 8),
                            HoverableButton(
                              name: "Delete Manager",
                              icon: Icons.delete_forever_rounded,
                              expanded: true,
                              onTap: () => _handleDeleteManager(item.name),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.folder_open, size: 64, color: Color(0xff6B7280)),
          SizedBox(height: 16),
          Text(
            'No managers created',
            style: TextStyle(color: Color(0xff9CA3AF), fontSize: 16),
          ),
        ],
      ),
    );
  }

  Widget _buildSortButton(String managerName) {
    final isPending = widget.pendingSorts[managerName] ?? false;
    final hasSortResult = widget.sortResults.containsKey(managerName);

    if (hasSortResult) {
      // Show "View Sorted" button when results are ready
      return HoverableButton(
        onTap: () => _showSortPreview(managerName),
        name: "View Sorted",
        icon: Icons.preview_rounded,
        expanded: true,
      );
    } else if (isPending) {
      // Show "Sorting..." when in progress
      return HoverableButton(
        onTap: null,
        name: "Sorting...",
        icon: Icons.filter_tilt_shift_rounded,
        expanded: true,
      );
    } else {
      // Show "Sort Manager" when idle
      return HoverableButton(
        onTap: () => _handleSortManager(managerName),
        name: "Sort Manager",
        icon: Icons.filter_tilt_shift_rounded,
        expanded: true,
      );
    }
  }
}
