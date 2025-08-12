import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_force_directed_graph/flutter_force_directed_graph.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:app/constants.dart';
import 'dart:math' as math;
import 'package:app/custom_widgets/breadcrumb_widget.dart';
import 'package:open_file/open_file.dart';
import 'package:app/api.dart';
import 'package:flutter_context_menu/flutter_context_menu.dart';
import 'dart:io';
import 'package:app/custom_widgets/tag_dialog.dart';

class GraphViewPage extends StatefulWidget {
  final FileTreeNode treeData;
  final List<String> currentPath;
  final Function(FileTreeNode) onFileSelected;
  final Function(List<String>) onNavigate;
  final String? managerName;
  final VoidCallback? onTagChanged;

  const GraphViewPage({
    required this.treeData,
    required this.currentPath,
    required this.onFileSelected,
    required this.onNavigate,
    this.managerName,
    this.onTagChanged,
    super.key,
  });

  @override
  State<GraphViewPage> createState() => _GraphViewPageState();
}

class _GraphViewPageState extends State<GraphViewPage> {
  late ForceDirectedGraphController<FileTreeNode> controller;
  Map<FileTreeNode, int> nodeDepths = {};
  FileTreeNode? draggedNode;
  Set<FileTreeNode> highlightedNodes = {};

  //max depth
  static const int maxDepth = 3;
  FileTreeNode? currentRootNode;

  static const List<Color> levelColors = [
    kprimaryColor,
    Color(0xFF3498DB),
    Color(0xFF2ECC71),
    Color(0xFFE74C3C),
    Color(0xFF9B59B6),
  ];

  @override
  void initState() {
    super.initState();
    _initializeController();
    _buildFileTreeGraph();
  }

  @override
  void didUpdateWidget(GraphViewPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.currentPath != widget.currentPath) {
      _buildFileTreeGraph();
    }
  }

  void _initializeController() {
    controller = ForceDirectedGraphController<FileTreeNode>(
      graph: ForceDirectedGraph.generateNNodes(
        nodeCount: 0,
        generator:
            () => FileTreeNode(name: "data", isFolder: false, locked: false),
        config: GraphConfig(scaling: 0.1, elasticity: 0.8, repulsionRange: 200),
      ),
    );
  }

  void _buildFileTreeGraph() {
    // Clear existing graph
    _initializeController();
    nodeDepths.clear();
    // current root based on currentPath
    currentRootNode = _getCurrentRoot();

    List<FileTreeNode> allNodes = [];
    List<(FileTreeNode, FileTreeNode)> edges = [];

    _traverseTreeWithDepthLimit(currentRootNode!, allNodes, edges, 0);

    // Populate depths
    for (var node in allNodes) {
      _calculateDepth(node);
    }

    // Add nodes
    for (var node in allNodes) {
      controller.addNode(node);
    }

    // Add edges
    for (var edge in edges) {
      try {
        controller.addEdgeByData(edge.$1, edge.$2);
      } catch (e) {
        print("Error adding edge: $e");
      }
    }

    WidgetsBinding.instance.addPostFrameCallback((_) {
      controller.needUpdate();
    });
  }

  FileTreeNode _getCurrentRoot() {
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

    return currentFolder;
  }

  // depth limit treversing
  void _traverseTreeWithDepthLimit(
    FileTreeNode node,
    List<FileTreeNode> allNodes,
    List<(FileTreeNode, FileTreeNode)> edges,
    int currentDepth,
  ) {
    allNodes.add(node);

    if (node.children != null && currentDepth < maxDepth) {
      for (FileTreeNode child in node.children!) {
        edges.add((node, child));
        _traverseTreeWithDepthLimit(child, allNodes, edges, currentDepth + 1);
      }
    }
  }

  int _calculateDepth(FileTreeNode targetNode) {
    if (nodeDepths.containsKey(targetNode)) {
      return nodeDepths[targetNode]!;
    }

    int depth = _findNodeDepthFromRoot(currentRootNode!, targetNode, 0);
    nodeDepths[targetNode] = depth;
    return depth;
  }

  int _findNodeDepthFromRoot(
    FileTreeNode current,
    FileTreeNode target,
    int currentDepth,
  ) {
    if (current == target) {
      return currentDepth;
    }

    if (current.children != null && currentDepth < maxDepth) {
      for (FileTreeNode child in current.children!) {
        int result = _findNodeDepthFromRoot(child, target, currentDepth + 1);
        if (result != -1) {
          return result;
        }
      }
    }

    return -1;
  }

  Set<FileTreeNode> _getConnectedNodes(FileTreeNode targetNode) {
    Set<FileTreeNode> connected = {};
    if (targetNode.children != null) {
      connected.addAll(targetNode.children!);
    }

    return connected;
  }

  void _handleNodeDoubleTap(FileTreeNode node) {
    if (node.isFolder) {
      if (node == currentRootNode) {
        // Don't navigate if its current directory
        return;
      }

      List<String> pathToNode = _getPathToNode(node);
      widget.onNavigate(pathToNode);
    } else {
      openDocument(node.path ?? '');
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

  void _handleNodeRightTap(
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

  List<String> _getPathToNode(FileTreeNode targetNode) {
    List<String> path = [];
    _buildPathToNode(widget.treeData, targetNode, [], path);
    return path;
  }

  bool _buildPathToNode(
    FileTreeNode current,
    FileTreeNode target,
    List<String> currentPath,
    List<String> resultPath,
  ) {
    if (current == target) {
      resultPath.addAll(currentPath);
      return true;
    }

    if (current.children != null) {
      for (FileTreeNode child in current.children!) {
        List<String> newPath = [...currentPath];
        if (child.isFolder || child == target) {
          newPath.add(child.name);
        }

        if (_buildPathToNode(child, target, newPath, resultPath)) {
          return true;
        }
      }
    }

    return false;
  }

  Color _getLevelColor(int depth) {
    return levelColors[depth % levelColors.length];
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        BreadcrumbWidget(
          currentPath: widget.currentPath,
          onNavigate: widget.onNavigate,
        ),
        Expanded(child: _buildGraph()),
      ],
    );
  }

  Widget _buildGraph() {
    return Listener(
      onPointerSignal: (event) {
        setState(() {
          if (event is PointerScrollEvent) {
            double yScroll = event.scrollDelta.dy;
            if (yScroll < 0 && controller.scale < 2) {
              controller.scale += 0.1;
            } else if (yScroll > 0 && controller.scale > 0.1) {
              controller.scale -= 0.1;
            }
          }
        });
      },
      child: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [kScaffoldColor, kScaffoldColor.withBlue(30)],
          ),
        ),
        child: ForceDirectedGraphWidget(
          controller: controller,
          onDraggingStart: (data) {
            setState(() {
              draggedNode = data;
              highlightedNodes = _getConnectedNodes(data);
            });
          },
          onDraggingEnd: (data) {
            setState(() {
              draggedNode = null;
              highlightedNodes.clear();
            });
          },
          onDraggingUpdate: (data) {},
          nodesBuilder: (context, data) {
            final isFolder = data.isFolder;
            final hasChildren = data.children?.isNotEmpty ?? false;
            final depth = _calculateDepth(data);
            final isBeingDragged = draggedNode == data;
            final isHighlighted = highlightedNodes.contains(data);
            final levelColor = _getLevelColor(depth);

            final baseSize = isFolder ? 50.0 : 35.0;
            final nodeSize = math.max(20.0, baseSize - (depth * 6.0));

            return GestureDetector(
              onDoubleTap: () => _handleNodeDoubleTap(data),
              onSecondaryTapUp:
                  (details) => _handleNodeRightTap(
                    widget.managerName ?? "",
                    data,
                    details.globalPosition,
                  ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: nodeSize,
                    height: nodeSize,
                    decoration: BoxDecoration(
                      color: _getNodeColor(
                        isFolder,
                        depth,
                        isBeingDragged,
                        isHighlighted,
                        levelColor,
                      ),
                      border: Border.all(
                        color: _getBorderColor(
                          isFolder,
                          isBeingDragged,
                          isHighlighted,
                          levelColor,
                        ),
                        width: isBeingDragged ? 3.0 : 2.0,
                      ),
                      borderRadius: BorderRadius.circular(
                        isFolder ? 8 : nodeSize / 2,
                      ),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withValues(alpha: 0.2),
                          blurRadius: isBeingDragged ? 6.0 : 3.0,
                          offset: Offset(0, isBeingDragged ? 3.0 : 1.5),
                        ),
                      ],
                    ),
                    child: Stack(
                      children: [
                        Center(
                          child: Icon(
                            isFolder
                                ? (hasChildren
                                    ? Icons.folder
                                    : Icons.folder_outlined)
                                : _getFileIcon(data.name),
                            color: _getIconColor(
                              isFolder,
                              isBeingDragged,
                              isHighlighted,
                              levelColor,
                            ),
                            size: math.max(
                              10.0,
                              (isFolder ? 26.0 : 20.0) - (depth * 3.0),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),

                  // Label below node
                  const SizedBox(height: 4),
                  Container(
                    constraints: const BoxConstraints(maxWidth: 80),
                    padding: const EdgeInsets.symmetric(
                      horizontal: 6,
                      vertical: 2,
                    ),
                    decoration: BoxDecoration(
                      color: kScaffoldColor.withValues(alpha: 0.9),
                      borderRadius: BorderRadius.circular(4),
                      border: Border.all(
                        color: isHighlighted ? levelColor : kOutlineBorder,
                        width: isHighlighted ? 1.0 : 0.5,
                      ),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        (data.locked == true)
                            ? Icon(Icons.lock, size: 15, color: Colors.white70)
                            : Icon(
                              Icons.lock_open,
                              size: 15,
                              color: Colors.white70,
                            ),
                        const SizedBox(width: 2),
                        Flexible(
                          child: Text(
                            data.name,
                            style: TextStyle(
                              color:
                                  isHighlighted ? levelColor : Colors.white70,
                              fontSize: math.max(9.0, 11.0 - (depth * 0.5)),
                              fontWeight:
                                  isHighlighted
                                      ? FontWeight.w600
                                      : FontWeight.w500,
                            ),
                            textAlign: TextAlign.center,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            );
          },
          edgesBuilder: (context, a, b, distance) {
            final depthA = _calculateDepth(a);
            final depthB = _calculateDepth(b);
            final parentDepth = math.min(depthA, depthB);
            final levelColor = _getLevelColor(parentDepth);

            // Modified edge highlighting - only highlight edges going to children
            final isHighlighted =
                (draggedNode == a && highlightedNodes.contains(b));

            return Container(
              width: distance,
              height: isHighlighted ? 2.5 : 1.5,
              decoration: BoxDecoration(
                color:
                    isHighlighted
                        ? levelColor.withValues(alpha: 0.8)
                        : levelColor.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(1),
              ),
            );
          },
        ),
      ),
    );
  }
}

IconData _getFileIcon(String fileName) {
  final extension = fileName.split('.').last.toLowerCase();

  switch (extension) {
    case 'dart':
      return Icons.code;
    case 'json':
      return Icons.data_object;
    case 'yaml':
    case 'yml':
      return Icons.settings;
    case 'md':
      return Icons.article;
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
      return Icons.image;
    case 'pdf':
      return Icons.picture_as_pdf;
    default:
      return Icons.description;
  }
}

Color _getNodeColor(
  bool isFolder,
  int depth,
  bool isBeingDragged,
  bool isHighlighted,
  Color levelColor,
) {
  if (isBeingDragged) {
    return levelColor.withValues(alpha: 0.7);
  }
  if (isHighlighted) {
    return levelColor.withValues(alpha: 0.3);
  }

  return isFolder
      ? levelColor.withValues(alpha: 0.2)
      : kScaffoldColor.withValues(alpha: 0.8);
}

Color _getBorderColor(
  bool isFolder,
  bool isBeingDragged,
  bool isHighlighted,
  Color levelColor,
) {
  if (isBeingDragged || isHighlighted) {
    return levelColor;
  }
  return isFolder ? levelColor.withValues(alpha: 0.6) : kOutlineBorder;
}

Color _getIconColor(
  bool isFolder,
  bool isBeingDragged,
  bool isHighlighted,
  Color levelColor,
) {
  if (isBeingDragged) {
    return Colors.white;
  }
  if (isHighlighted) {
    return levelColor;
  }
  return isFolder ? levelColor : Colors.white70;
}

String convertWSLPath(String wslPath) {
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

void openDocument(String originalWSLPath) async {
  final nativePath = convertWSLPath(originalWSLPath);

  final file = File(nativePath);
  if (await file.exists()) {
    OpenFile.open(nativePath);
  }
}
