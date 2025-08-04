import 'package:flutter/material.dart';
import '../models/file_tree_node.dart';

class SearchItemWidget extends StatefulWidget {
  final FileTreeNode item;
  final Function(FileTreeNode) onTap;
  final Function(FileTreeNode) onDoubleTap;
  final String managerPath;

  const SearchItemWidget({
    required this.item,
    required this.onTap,
    required this.onDoubleTap,
    required this.managerPath,
    super.key,
  });

  @override
  State<SearchItemWidget> createState() => _SearchItemWidgetState();
}

class _SearchItemWidgetState extends State<SearchItemWidget> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return GestureDetector(
            onTap: () => widget.onTap(widget.item),
            onDoubleTap: () => widget.onDoubleTap(widget.item),
            child: SizedBox(
              width: constraints.maxWidth,
              height: constraints.maxHeight,
              child: Positioned.fill(
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color:
                        _isHovered
                            ? const Color(0xff374151)
                            : const Color(0xff242424),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color:
                          _isHovered
                              ? const Color(0xffFFB400)
                              : const Color(0xff3D3D3D),
                      width: _isHovered ? 2 : 1,
                    ),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 32,
                        height: 32,
                        decoration: BoxDecoration(
                          color: _getFileColor(),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: Icon(
                          _getFileIcon(),
                          color: Colors.white,
                          size: 18,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        flex: 1,
                        child: Text(
                          widget.item.name,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        flex: 3,
                        child: Text(
                          _getTruncatedPath(
                            widget.item.path ?? "",
                            widget.managerPath,
                          ),
                          style: const TextStyle(
                            color: Color(0xff6b7280),
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                          textAlign: TextAlign.right,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  String _getTruncatedPath(
    String fullPath,
    String managerPath, {
    int maxVisibleSegments = 4,
  }) {
    if (fullPath.isEmpty) return "";

    // Remove the managerPath
    String relative = fullPath;
    if (managerPath.isNotEmpty && fullPath.startsWith(managerPath)) {
      relative = fullPath.substring(managerPath.length);
      relative = relative.replaceFirst(RegExp(r'^[/\\]+'), '');
    }
    // Split into segments
    final segments =
        relative.split(RegExp(r'[/\\]')).where((s) => s.isNotEmpty).toList();

    // Prepend a fixed root label
    segments.insert(0, "Root");

    if (segments.length <= maxVisibleSegments) {
      return ".../${["Root", ...segments].skip(1).join('/')}";
    } else {
      final tailCount = maxVisibleSegments - 1;
      final lastParts = segments.sublist(segments.length - tailCount);
      return ".../${["Root", ...lastParts].skip(1).join('/')}";
    }
  }

  Color _getFileColor() {
    if (widget.item.isFolder) return const Color(0xffFFB400);
    return const Color(0xff2563EB);
  }

  IconData _getFileIcon() {
    if (widget.item.isFolder) return Icons.folder;
    return Icons.description;
  }
}
