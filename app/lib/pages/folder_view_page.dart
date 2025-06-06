import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/custom_widgets/file_item_widget.dart';

class FolderViewPage extends StatefulWidget {
  final FileTreeNode treeData;
  final List<String> currentPath;
  final Function(FileTreeNode) onFileSelected;
  final Function(List<String>) onNavigate;

  const FolderViewPage({
    required this.treeData,
    required this.currentPath,
    required this.onFileSelected,
    required this.onNavigate,
    super.key,
  });

  @override
  State<FolderViewPage> createState() => _FolderViewPageState();
}

class _FolderViewPageState extends State<FolderViewPage> {
  List<FileTreeNode> _currentItems = [];

  @override
  void initState() {
    super.initState();
    _updateCurrentItems();
  }

  @override
  void didUpdateWidget(FolderViewPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.currentPath != widget.currentPath) {
      _updateCurrentItems();
    }
  }

  void _updateCurrentItems() {
    FileTreeNode currentFolder = widget.treeData;

    for (String pathSegment in widget.currentPath) {
      final foundFolder = currentFolder.children?.firstWhere(
        (child) => child.name == pathSegment && child.isFolder,
        orElse: () => currentFolder,
      );
      if (foundFolder != null && foundFolder != currentFolder) {
        currentFolder = foundFolder;
      }
    }

    setState(() {
      _currentItems = List.from(currentFolder.children ?? []);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _buildBreadcrumb(),
        Expanded(
          child: _currentItems.isEmpty ? _buildEmptyState() : _buildFileGrid(),
        ),
      ],
    );
  }

  Widget _buildBreadcrumb() {
    return Center(
      child: Container(
        height: 40,
        padding: const EdgeInsets.symmetric(horizontal: 16),
        decoration: const BoxDecoration(
          color: Color(0xff242424),
          border: Border(bottom: BorderSide(color: Color(0xff3D3D3D))),
        ),
        child: Row(
          children: [
            _buildBreadcrumbItem('Root', [], widget.currentPath.isEmpty),
            ...widget.currentPath.asMap().entries.map((entry) {
              final index = entry.key;
              final pathSegment = entry.value;
              final isLast = index == widget.currentPath.length - 1;
              final pathToHere = widget.currentPath.sublist(0, index + 1);

              return Row(
                children: [
                  const Text('/', style: TextStyle(color: Color(0xff6B7280))),
                  _buildBreadcrumbItem(pathSegment, pathToHere, isLast),
                ],
              );
            }),
          ],
        ),
      ),
    );
  }

  Widget _buildBreadcrumbItem(String name, List<String> path, bool isActive) {
    return GestureDetector(
      onTap: () {
        widget.onNavigate(path);
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        decoration: BoxDecoration(borderRadius: BorderRadius.circular(4)),
        child: Text(
          name,
          style: TextStyle(
            color: isActive ? Color(0xffFFB400) : const Color(0xff9CA3AF),
            fontSize: 12,
            fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ),
    );
  }

  Widget _buildFileGrid() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: LayoutBuilder(
        builder: (context, constraints) {
          final itemWidth = 100.0;
          final spacing = 12.0;
          final availableWidth =
              constraints.maxWidth - 32.0; // Account for padding
          final crossAxisCount = ((availableWidth + spacing) /
                  (itemWidth + spacing))
              .floor()
              .clamp(1, 10);

          return GridView.builder(
            gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: crossAxisCount,
              childAspectRatio: 1.0,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
            ),
            itemCount: _currentItems.length,
            itemBuilder: (context, index) {
              return FileItemWidget(
                item: _currentItems[index],
                onTap: _handleItemTap,
              );
            },
          );
        },
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
            'This folder is empty',
            style: TextStyle(color: Color(0xff9CA3AF), fontSize: 16),
          ),
        ],
      ),
    );
  }

  void _handleItemTap(FileTreeNode item) {
    if (item.isFolder) {
      widget.onNavigate([...widget.currentPath, item.name]);
    } else {
      widget.onFileSelected(item);
    }
  }
}
