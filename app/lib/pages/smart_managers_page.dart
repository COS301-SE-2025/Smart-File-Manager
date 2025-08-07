import 'package:app/constants.dart';
import 'package:app/custom_widgets/hoverable_button.dart';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';

class SmartManagersPage extends StatefulWidget {
  Map<String, FileTreeNode> managerTreeData = {};
  List<String> managerNames;
  final Function(String) onManagerDelete;
  final Function(String, FileTreeNode) onManagerSort;

  SmartManagersPage({
    super.key,
    required this.managerTreeData,
    required this.managerNames,
    required this.onManagerDelete,
    required this.onManagerSort,
  });

  @override
  State<SmartManagersPage> createState() => _SmartManagersPageState();
}

class _SmartManagersPageState extends State<SmartManagersPage> {
  List<FileTreeNode> _currentItems = [];
  final Map<String, bool> _isSorting = {};

  @override
  void initState() {
    super.initState();
    _updateCurrentItems();
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
          itemCount: widget.managerTreeData.length,
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
                            HoverableButton(
                              onTap: _isSorting ? null : _handleSortManager,
                              name: _isSorting ? "Sorting..." : "Sort Manager",
                              icon: Icons.filter_tilt_shift_rounded,
                              expanded: true,
                            ),
                            SizedBox(height: 8),
                            HoverableButton(
                              name: "Delete Manager",
                              icon: Icons.delete_forever_rounded,
                              expanded: true,
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
}
