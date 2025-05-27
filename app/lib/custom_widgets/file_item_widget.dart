import 'package:flutter/material.dart';
import '../models/file_tree_node.dart';

class FileItemWidget extends StatefulWidget {
  final FileTreeNode item;
  final Function(FileTreeNode) onTap;

  const FileItemWidget({required this.item, required this.onTap, super.key});

  @override
  State<FileItemWidget> createState() => _FileItemWidgetState();
}

class _FileItemWidgetState extends State<FileItemWidget> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: () => widget.onTap(widget.item),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color:
                _isHovered ? const Color(0xff374151) : const Color(0xff242424),
            borderRadius: BorderRadius.circular(8),
            border: Border.all(
              color:
                  _isHovered
                      ? const Color(0xffFFB400)
                      : const Color(0xff3D3D3D),
              width: _isHovered ? 2 : 1,
            ),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: _getFileColor(),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Icon(_getFileIcon(), color: Colors.white, size: 18),
              ),
              const SizedBox(height: 8),
              Text(
                widget.item.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                ),
                textAlign: TextAlign.center,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),
            ],
          ),
        ),
      ),
    );
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
