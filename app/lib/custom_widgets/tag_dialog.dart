import 'package:flutter/material.dart';
import 'package:app/constants.dart';
import 'package:app/api.dart';
import 'package:app/models/file_tree_node.dart';

class TagDialog extends StatefulWidget {
  final FileTreeNode node;
  final String? managerName;

  const TagDialog({required this.node, required this.managerName, super.key});

  @override
  State<TagDialog> createState() => _TagDialogState();
}

class _TagDialogState extends State<TagDialog> {
  final TextEditingController _tagController = TextEditingController();

  @override
  void dispose() {
    _tagController.dispose();
    super.dispose();
  }

  Future<void> _addTag() async {
    if (_tagController.text.isEmpty) return;

    bool response = await Api.addTagToFile(
      widget.managerName ?? '',
      widget.node.path ?? '',
      _tagController.text,
    );

    if (mounted) {
      Navigator.pop(context, response ? _tagController.text : null);
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
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
        onSubmitted: (_) => _addTag(),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context, false),
          child: const Text(
            'Cancel',
            style: TextStyle(color: Color(0xff9CA3AF)),
          ),
        ),
        TextButton(
          onPressed: _addTag,
          child: const Text('Add', style: TextStyle(color: Color(0xffFFB400))),
        ),
      ],
    );
  }
}
