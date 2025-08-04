import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/custom_widgets/search_item_widget.dart';
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
  final Function(String)? onGoToFolder;
  final bool showGoToFolder;

  const FolderViewSearch({
    required this.managerName,
    required this.treeData,
    required this.onFileSelected,
    this.onTagChanged,
    this.onGoToFolder,
    this.showGoToFolder = true,
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
      child: SingleChildScrollView(
        child: Column(
          children:
              _currentItems.map((item) {
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12.0),
                  child: SizedBox(
                    width: double.infinity,
                    height: 60,
                    child: GestureDetector(
                      onSecondaryTapDown:
                          (details) => _handleItemRightTap(
                            widget.managerName ?? "",
                            item,
                            details.globalPosition,
                          ),
                      child: SearchItemWidget(
                        item: item,
                        onTap: _handleItemTap,
                        onDoubleTap: _handleNodeDoubleTap,
                      ),
                    ),
                  ),
                );
              }).toList(),
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
    final entries = <ContextMenuEntry>[
      if (widget.showGoToFolder)
        MenuItem(
          label: 'Go to Folder',
          icon: Icons.drive_file_move_rounded,
          onSelected: () => _goToFolder(node),
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

  void _goToFolder(FileTreeNode node) {
    if (widget.onGoToFolder != null) {
      final filePath = node.path ?? '';
      if (filePath.isNotEmpty) {
        // Extract parent directory path
        final parts = filePath.split(RegExp(r'[/\\]'));
        if (parts.length > 1) {
          // Remove the file name, keep the directory path
          parts.removeLast();
          final parentPath = parts.join('/');
          
          // Find the folder path within the tree structure
          final folderPath = _findFolderPathInTree(parentPath);
          widget.onGoToFolder!(folderPath);
        }
      }
    }
  }

  String _findFolderPathInTree(String targetPath) {
    // Get the manager root path from the first child's path
    String managerRootPath = _getManagerRootPath();
    
    if (managerRootPath.isEmpty) {
      return ''; // Go to root if we can't determine manager root
    }
    
    // Normalize paths
    String normalizedTarget = targetPath.replaceAll('\\', '/');
    String normalizedManagerRoot = managerRootPath.replaceAll('\\', '/');
    
    // Remove manager root from target path to get relative path
    if (normalizedTarget.startsWith(normalizedManagerRoot)) {
      String relativePath = normalizedTarget.substring(normalizedManagerRoot.length);
      
      // Remove leading slash
      if (relativePath.startsWith('/')) {
        relativePath = relativePath.substring(1);
      }
      
      // Convert file system path to folder names by traversing the tree
      return _convertPathToFolderNames(relativePath);
    }
    
    return ''; // Go to root if path doesn't match
  }

  String _getManagerRootPath() {
    // If tree has children, use first child's path to determine manager root
    if (widget.treeData.children != null && widget.treeData.children!.isNotEmpty) {
      String firstChildPath = widget.treeData.children!.first.path ?? '';
      if (firstChildPath.isNotEmpty) {
        // Get parent directory of first child - that's the manager root
        List<String> pathParts = firstChildPath.split(RegExp(r'[/\\]'));
        pathParts.removeLast(); // Remove the child name
        return pathParts.join('/');
      }
    }
    
    // If no children, use the tree data's own path as manager root
    return widget.treeData.path ?? '';
  }

  String _convertPathToFolderNames(String relativePath) {
    if (relativePath.isEmpty) {
      return '';
    }
    
    // Split the relative path into segments
    List<String> pathSegments = relativePath.split('/').where((segment) => segment.isNotEmpty).toList();
    
    // Find the corresponding folder names in the tree structure
    List<String> folderNames = [];
    FileTreeNode currentNode = widget.treeData;
    
    for (String segment in pathSegments) {
      if (currentNode.children != null) {
        // Find the child folder with this path segment
        FileTreeNode? matchingChild;
        for (FileTreeNode child in currentNode.children!) {
          if (child.isFolder && child.path != null) {
            String childPath = child.path!.replaceAll('\\', '/');
            if (childPath.endsWith('/$segment') || childPath.endsWith(segment)) {
              matchingChild = child;
              break;
            }
          }
        }
        
        if (matchingChild != null) {
          folderNames.add(matchingChild.name);
          currentNode = matchingChild;
        } else {
          // If we can't find matching child, return what we have so far
          break;
        }
      }
    }
    
    return folderNames.join('/');
  }
}
