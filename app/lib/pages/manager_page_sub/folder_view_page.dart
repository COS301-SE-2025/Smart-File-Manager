import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/custom_widgets/file_item_widget.dart';
import 'package:app/custom_widgets/breadcrumb_widget.dart';
import 'package:open_file/open_file.dart';
import 'dart:io';

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
        BreadcrumbWidget(
          currentPath: widget.currentPath,
          onNavigate: widget.onNavigate,
        ),
        Expanded(
          child: _currentItems.isEmpty ? _buildEmptyState() : _buildFileGrid(),
        ),
      ],
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
                onDoubleTap: _handleNodeDoubleTap,
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

  void _handleNodeDoubleTap(FileTreeNode item) {
    if (!item.isFolder) {
      _openDocument(item.path ?? '');
    }
  }

  String _convertWSLPath(String wslPath) {
    if (Platform.isWindows) {
      final match = RegExp(r"^/mnt/([a-zA-Z])/").firstMatch(wslPath);
      if (match != null) {
        final driveLetter = match.group(1)!.toUpperCase();
        final windowsPath = wslPath
            .replaceFirst(RegExp(r"^/mnt/[a-zA-Z]/"), "$driveLetter:/")
            .replaceAll('/', r'\');
        return windowsPath;
      }
      return wslPath;
    } else {
      return wslPath;
    }
  }

  void _openDocument(String originalWSLPath) async {
    final nativePath = _convertWSLPath(originalWSLPath);

    final file = File(nativePath);
    if (await file.exists()) {
      await OpenFile.open(nativePath);
    } else {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('File not found: $nativePath'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }
}
