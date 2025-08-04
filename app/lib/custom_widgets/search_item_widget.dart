import 'package:flutter/material.dart';
import '../models/file_tree_node.dart';

class SearchItemWidget extends StatefulWidget {
  final FileTreeNode item;
  final Function(FileTreeNode) onTap;
  final Function(FileTreeNode) onDoubleTap;

  const SearchItemWidget({
    required this.item,
    required this.onTap,
    required this.onDoubleTap,
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
                          _getTruncatedPath(widget.item.path ?? ""),
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

  String _getTruncatedPath(String path) {
    if (path.isEmpty) return "";

    // Split path
    final parts = path.split(RegExp(r'[/\\]'));

    // If path is short ,retrun
    if (parts.length <= 4) return path;

    // Show only the last 3
    final lastParts = parts.sublist(parts.length - 4);
    return ".../${lastParts.join('/')}";
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
