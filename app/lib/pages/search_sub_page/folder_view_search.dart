import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/custom_widgets/file_item_widget.dart';
import 'package:app/custom_widgets/tag_dialog.dart';
import 'package:app/api.dart';
import 'package:open_file/open_file.dart';
import 'package:flutter_context_menu/flutter_context_menu.dart';
import 'dart:io';

class FolderViewSearch extends StatefulWidget {
  final String managerName;
  final FileTreeNode treeData;
  final Function(FileTreeNode) onFileSelected;
  final VoidCallback? onTagChanged;

  const FolderViewSearch({
    required this.managerName,
    required this.treeData,
    required this.onFileSelected,
    this.onTagChanged,
    super.key,
  });

  @override
  State<FolderViewSearch> createState() => _FolderViewSearchState();
}

class _FolderViewSearchState extends State<FolderViewSearch> {
  List<FileTreeNode> _currentItems = [];

  @override
  void initState() {
    super.initState();
    _updateCurrentItems();
  }


  void _updateCurrentItems() {
    setState(() {
      _currentItems = List.from(widget.treeData.children ?? []);
    });
  }

  @override
  Widget build(BuildContext context) {
    return _currentItems.isEmpty ? _buildEmptyState() : _buildFileGrid();
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
              return GestureDetector(
                onSecondaryTapDown:
                    (details) => _handleItemRightTap(
                      widget.managerName ?? "",
                      _currentItems[index],
                      details.globalPosition,
                    ),
                child: FileItemWidget(
                  item: _currentItems[index],
                  onTap: _handleItemTap,
                  onDoubleTap: _handleNodeDoubleTap,
                ),
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
    if (!item.isFolder) {
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

  void _handleItemRightTap(
    String managerName,
    FileTreeNode node,
    Offset globalPosition,
  ) {
    if (node.isFolder) {
      final entries = <ContextMenuEntry>[
        MenuItem(
          label: 'Lock',
          icon: Icons.lock,
          onSelected: () async {
            bool response = await Api.locking(managerName, node.path ?? '');
            if (response == true) {
              _lockNode(node);
            }
          },
        ),
        MenuItem(
          label: 'Unlock',
          icon: Icons.lock_open,
          onSelected: () async {
            bool response = await Api.unlocking(managerName, node.path ?? '');
            if (response == true) {
              _unlockNode(node);
            }
          },
        ),
      ];

      final menu = ContextMenu(
        entries: entries,
        position: globalPosition,
        padding: const EdgeInsets.all(8.0),
      );

      showContextMenu(context, contextMenu: menu);
    } else {
      final entries = <ContextMenuEntry>[
        MenuItem(
          label: 'Details',
          icon: Icons.info_outline,
          onSelected: () => widget.onFileSelected(node),
        ),
        MenuItem(
          label: 'Add Tag',
          icon: Icons.label,
          onSelected: () => _showAddTagDialog(node),
        ),
        MenuItem(
          label: 'Lock',
          icon: Icons.lock,
          onSelected: () async {
            bool response = await Api.locking(managerName, node.path ?? '');
            if (response == true) {
              _lockNode(node);
            }
          },
        ),
        MenuItem(
          label: 'Unlock',
          icon: Icons.lock_open,
          onSelected: () async {
            bool response = await Api.unlocking(managerName, node.path ?? '');
            if (response == true) {
              _unlockNode(node);
            }
          },
        ),
      ];

      final menu = ContextMenu(
        entries: entries,
        position: globalPosition,
        padding: const EdgeInsets.all(8.0),
      );

      showContextMenu(context, contextMenu: menu);
    }
  }

  void _showAddTagDialog(FileTreeNode node) async {
    final addedTag = await showDialog<String>(
      context: context,
      builder:
          (context) => TagDialog(node: node, managerName: widget.managerName),
    );

    if (addedTag != null && mounted) {
      setState(() {
        node.tags?.add(addedTag);
      });
      // Notify parent that tags have changed
      widget.onTagChanged?.call();
    }
  }

  void _lockNode(FileTreeNode node) {
    setState(() {
      widget.treeData.lockItem(node.path ?? '');
    });
  }

  void _unlockNode(FileTreeNode node) {
    setState(() {
      widget.treeData.unlockItem(node.path ?? '');
    });
  }
}
