import 'package:flutter/material.dart';
import '../models/file_tree_node.dart';

class FileItemWidget extends StatefulWidget {
  final FileTreeNode item;
  final Function(FileTreeNode) onTap;
  final Function(FileTreeNode) onDoubleTap;

  const FileItemWidget({
    required this.item,
    required this.onTap,
    required this.onDoubleTap,
    super.key,
  });

  @override
  State<FileItemWidget> createState() => _FileItemWidgetState();
}

class _FileItemWidgetState extends State<FileItemWidget> {
  bool _isHovered = false;

  // File categories with extensions
  static final Map<String, List<String>> _categories = {
    "Documents": ["pdf", "doc", "docx", "rtf", "txt", "odt", "md", "csv"],
    "Images": [
      "jpg",
      "jpeg",
      "png",
      "gif",
      "bmp",
      "tiff",
      "tif",
      "webp",
      "svg",
    ],
    "Music": ["mp3", "wav", "flac", "aac", "m4a", "ogg", "wma"],
    "Presentations": ["ppt", "pptx", "odp", "key"],
    "Videos": ["mp4", "mkv", "avi", "mov", "wmv", "webm"],
    "Spreadsheets": ["xls", "xlsx", "ods", "tsv", "xlsm"],
    "Archives": ["zip", "rar", "7z", "tar", "gz", "iso"],
    "Other": [],
  };

  // Icons & colors per category
  static final Map<String, IconData> _categoryIcons = {
    "Documents": Icons.description,
    "Images": Icons.image,
    "Music": Icons.music_note,
    "Presentations": Icons.slideshow,
    "Videos": Icons.movie,
    "Spreadsheets": Icons.table_chart,
    "Archives": Icons.archive,
    "Other": Icons.insert_drive_file,
  };

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
              child: Stack(
                children: [
                  Positioned.fill(
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
                            child: Icon(
                              _getFileIcon(),
                              color: Colors.white,
                              size: 18,
                            ),
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
                  Positioned(
                    top: 6,
                    left: 6,
                    child:
                        widget.item.locked
                            ? const Icon(
                              Icons.lock,
                              size: 15,
                              color: Colors.white70,
                            )
                            : const Icon(
                              Icons.lock_open,
                              size: 15,
                              color: Colors.white70,
                            ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  String _getFileExtension() {
    final name = widget.item.name.toLowerCase();
    if (!name.contains('.')) return "";
    return name.split('.').last;
  }

  String _getCategory() {
    if (widget.item.isFolder) return "Folder";

    final ext = _getFileExtension();
    for (var entry in _categories.entries) {
      if (entry.value.contains(ext)) {
        return entry.key;
      }
    }
    return "Other";
  }

  Color _getFileColor() {
    if (widget.item.isFolder) return const Color(0xffFFB400);
    return const Color(0xff2563EB);
  }

  IconData _getFileIcon() {
    if (widget.item.isFolder) return Icons.folder;
    return _categoryIcons[_getCategory()] ?? Icons.insert_drive_file;
  }
}
