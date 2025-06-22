import 'package:flutter/material.dart';
import '../models/file_tree_node.dart';
import 'package:app/constants.dart';
import 'package:app/api.dart';

class FileDetailsPanel extends StatefulWidget {
  final String? managerName;
  final FileTreeNode? selectedFile;
  final bool isVisible;
  final VoidCallback onClose;

  const FileDetailsPanel({
    required this.managerName,
    required this.selectedFile,
    required this.isVisible,
    required this.onClose,
    super.key,
  });

  @override
  State<FileDetailsPanel> createState() => _FileDetailsPanelState();
}

class _FileDetailsPanelState extends State<FileDetailsPanel>
    with SingleTickerProviderStateMixin {
  final TextEditingController _tagController = TextEditingController();

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    _tagController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.selectedFile == null) return const SizedBox.shrink();

    return Container(
      decoration: const BoxDecoration(
        color: Color(0xff2E2E2E),
        border: Border(left: BorderSide(color: Color(0xff3D3D3D))),
      ),
      child: _buildContent(),
    );
  }

  Widget _buildContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildHeader(),
        Expanded(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildFileInfo(),
                const SizedBox(height: 24),
                _buildProperties(),
                const SizedBox(height: 24),
                _buildTagsSection(),
                const SizedBox(height: 24),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildHeader() {
    return Container(
      height: 40,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: Color(0xff3D3D3D))),
      ),
      child: Row(
        children: [
          const Text(
            'File Details',
            style: TextStyle(color: Colors.white, fontSize: 12),
          ),
          const Spacer(),
          IconButton(
            onPressed: widget.onClose,
            icon: const Icon(Icons.close, color: Color(0xff9CA3AF), size: 20),
            splashRadius: 20,
          ),
        ],
      ),
    );
  }

  Widget _buildFileInfo() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Center(
          child: Container(
            width: 32,
            height: 32,
            decoration: BoxDecoration(
              color: Color(0xff2563EB),
              borderRadius: BorderRadius.circular(6),
            ),
            child: Icon(Icons.description, color: Colors.white, size: 20),
          ),
        ),

        const SizedBox(height: 16),

        Text(
          widget.selectedFile!.name,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 15,
            fontWeight: FontWeight.w600,
          ),
        ),

        const SizedBox(height: 4),
      ],
    );
  }

  Widget _buildProperties() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Properties',
          style: TextStyle(
            color: Colors.white,
            fontSize: 15,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 10),

        if (widget.selectedFile!.isFolder &&
            widget.selectedFile!.children != null) ...[
          _buildPropertyRow(
            'Items',
            '${widget.selectedFile!.children!.length}',
          ),
          const SizedBox(height: 8),
        ],

        if (!widget.selectedFile!.isFolder &&
            widget.selectedFile!.metadata != null) ...[
          _buildPropertyRow(
            'Size',
            _formatFileSize(widget.selectedFile!.metadata!['size'].toString()),
          ),

          const SizedBox(height: 8),

          _buildPropertyRow(
            'Type',
            widget.selectedFile!.metadata!['mimeType'].toString(),
          ),

          const SizedBox(height: 8),

          _buildPropertyRow(
            'Created',
            _formatDate(
              widget.selectedFile!.metadata!['dateCreated'].toString(),
            ),
          ),

          const SizedBox(height: 8),

          _buildPropertyRow(
            'Modified',
            _formatDate(
              widget.selectedFile!.metadata!['lastModified'].toString(),
            ),
          ),
        ],
      ],
    );
  }

  // Helper method to format file size
  String _formatFileSize(String sizeStr) {
    try {
      int sizeInBytes = int.parse(sizeStr);
      if (sizeInBytes == 0) return '0 bytes';

      const List<String> units = ['bytes', 'KB', 'MB', 'GB', 'TB'];
      int unitIndex = 0;
      double size = sizeInBytes.toDouble();

      while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
      }

      if (unitIndex == 0) {
        return '${size.toInt()} ${units[unitIndex]}';
      } else {
        return '${size.toStringAsFixed(1)} ${units[unitIndex]}';
      }
    } catch (e) {
      return sizeStr;
    }
  }

  // Helper method to format dates
  String _formatDate(String dateStr) {
    try {
      DateTime dateTime = DateTime.parse(dateStr);
      return '${dateTime.day}/${dateTime.month}/${dateTime.year} ${dateTime.hour.toString().padLeft(2, '0')}:${dateTime.minute.toString().padLeft(2, '0')}';
    } catch (e) {
      return dateStr;
    }
  }

  Widget _buildPropertyRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(color: Color(0xff9CA3AF), fontSize: 14),
        ),
        Flexible(
          child: Text(
            value,
            style: const TextStyle(color: Colors.white, fontSize: 14),
            textAlign: TextAlign.right,
          ),
        ),
      ],
    );
  }

  Widget _buildTagsSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Tags',
              style: TextStyle(
                color: Colors.white,
                fontSize: 15,
                fontWeight: FontWeight.w600,
              ),
            ),
            TextButton.icon(
              onPressed: _showAddTagDialog,
              icon: const Icon(Icons.add, color: Color(0xffFFB400), size: 16),
              label: const Text(
                'Add',
                style: TextStyle(color: Color(0xffFFB400), fontSize: 12),
              ),
            ),
          ],
        ),

        const SizedBox(height: 8),

        if (widget.selectedFile!.tags != null &&
            widget.selectedFile!.tags!.isNotEmpty)
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children:
                widget.selectedFile!.tags!
                    .map(
                      (tag) => Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 6,
                        ),
                        decoration: BoxDecoration(
                          color: const Color(0xff374151),
                          borderRadius: BorderRadius.circular(16),
                          border: Border.all(color: const Color(0xff4B5563)),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              tag,
                              style: const TextStyle(
                                color: Color(0xffE5E7EB),
                                fontSize: 12,
                              ),
                            ),
                            const SizedBox(width: 4),
                            GestureDetector(
                              onTap: () => _removeTag(tag),
                              child: const Icon(
                                Icons.close,
                                color: Color(0xff9CA3AF),
                                size: 14,
                              ),
                            ),
                          ],
                        ),
                      ),
                    )
                    .toList(),
          ),
      ],
    );
  }

  void _showAddTagDialog() {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            backgroundColor: kScaffoldColor,
            title: const Text('Add Tag', style: kTitle1),
            content: TextField(
              controller: _tagController,
              style: const TextStyle(color: Colors.white),
              decoration: const InputDecoration(
                hintText: 'Enter tag name',
                hintStyle: TextStyle(color: Color(0xff9CA3AF)),
                enabledBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Color(0xff6B7280)),
                ),
                focusedBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Color(0xffFFB400)),
                ),
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text(
                  'Cancel',
                  style: TextStyle(color: Color(0xff9CA3AF)),
                ),
              ),
              TextButton(
                onPressed: () {
                  if (_tagController.text.isNotEmpty) {
                    _addTag(_tagController.text);
                    _tagController.clear();
                    Navigator.pop(context);
                  }
                },
                child: const Text(
                  'Add',
                  style: TextStyle(color: Color(0xffFFB400)),
                ),
              ),
            ],
          ),
    );
  }

  void _addTag(String tag) async {
    bool response = await Api.addTagToFile(
      widget.managerName ?? '',
      widget.selectedFile!.path ?? '',
      tag,
    );

    setState(() {
      if (response) {
        widget.selectedFile!.tags?.add(tag);
      }
    });
  }

  void _removeTag(String tag) async {
    bool response = await Api.deleteTagFromFile(
      widget.managerName ?? '',
      widget.selectedFile!.path ?? '',
      tag,
    );

    setState(() {
      if (response) {
        widget.selectedFile!.tags?.remove(tag);
      }
    });
  }
}
