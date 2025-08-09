import 'package:flutter/material.dart';
import 'package:app/constants.dart';
import 'package:app/models/duplicate_model.dart';
import 'package:app/api.dart';

class DuplicateDialog extends StatefulWidget {
  List<DuplicateModel>? duplicates;
  final String name;

  DuplicateDialog({required this.name, super.key, this.duplicates});

  @override
  State<DuplicateDialog> createState() => _DuplicateDialogState();
}

class _DuplicateDialogState extends State<DuplicateDialog> {
  final TextEditingController _duplicateController = TextEditingController();
  bool _isLoading = true;

  @override
  void dispose() {
    _duplicateController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _loadDuplicateData(widget.name);
  }

  Future<void> _loadDuplicateData(String name) async {
    setState(() {
      _isLoading = true;
    });

    final duplicates = await Api.loadDuplicates(name);

    if (mounted) {
      setState(() {
        widget.duplicates = duplicates;
        _isLoading = false;
      });
    }
  }

  void _deleteDuplicate(String managerName, String filePath) async {
    bool response = await Api.deleteSingleFile(managerName, filePath);
    if (response == true) {
      setState(() {
        widget.duplicates?.removeWhere(
          (item) => item.duplicatePath == filePath,
        );
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Deleted duplicate successfully'),
          backgroundColor: kYellowText,
          duration: Duration(seconds: 2),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Could not delete duplicate '),
          backgroundColor: Colors.redAccent,
          duration: Duration(seconds: 2),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: kScaffoldColor,
      title: const Text('Duplicates', style: kTitle1),
      content: SizedBox(
        width: double.maxFinite,
        height: double.maxFinite,
        child:
            _isLoading
                ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      CircularProgressIndicator(color: Color(0xffFFB400)),
                      SizedBox(height: 16),
                      Text(
                        'Loading duplicates...',
                        style: TextStyle(color: Color(0xff9CA3AF)),
                      ),
                    ],
                  ),
                )
                : widget.duplicates == null || widget.duplicates!.isEmpty
                ? const Center(
                  child: Text(
                    'No duplicates found.',
                    style: TextStyle(color: Color(0xff9CA3AF), fontSize: 16),
                  ),
                )
                : Column(
                  children: [
                    Expanded(
                      child: ListView(
                        shrinkWrap: true,
                        children:
                            widget.duplicates!.map((object) {
                              return Card(
                                color: Color(0xff242424),
                                elevation: 0,
                                margin: const EdgeInsets.all(8.0),
                                child: Padding(
                                  padding: const EdgeInsets.all(16.0),
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Text(
                                        object.name,
                                        style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 20,
                                        ),
                                      ),
                                      const SizedBox(height: 4),
                                      Text(
                                        'Original Path: ${object.originalPath}',
                                        style: TextStyle(
                                          color: Color(0xff9CA3AF),
                                        ),
                                      ),
                                      Text(
                                        'Duplicate Path: ${object.duplicatePath}',
                                        style: TextStyle(
                                          color: Color(0xff9CA3AF),
                                        ),
                                      ),
                                      SizedBox(height: 20),
                                      TextButton(
                                        onPressed:
                                            () => _deleteDuplicate(
                                              widget.name,
                                              object.duplicatePath,
                                            ),
                                        style: TextButton.styleFrom(
                                          foregroundColor: Colors.grey,
                                          side: const BorderSide(
                                            color: Colors.grey,
                                          ),
                                          padding: const EdgeInsets.symmetric(
                                            horizontal: 24,
                                            vertical: 12,
                                          ),
                                        ),
                                        child: const Text('Delete Duplicate'),
                                      ),
                                    ],
                                  ),
                                ),
                              );
                            }).toList(),
                      ),
                    ),
                  ],
                ),
      ),
      actions: [
        if (widget.duplicates == null || widget.duplicates!.isEmpty)
          TextButton(
            onPressed: () {
              if (mounted) {
                Navigator.pop(context);
              }
            },
            style: TextButton.styleFrom(
              foregroundColor: Colors.grey,
              side: const BorderSide(color: Colors.grey),
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            ),
            child: const Text('Close'),
          )
        else
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              TextButton(
                onPressed: () {
                  if (mounted) {
                    Navigator.pop(context);
                  }
                },
                style: TextButton.styleFrom(
                  foregroundColor: Colors.grey,
                  side: const BorderSide(color: Colors.grey),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 24,
                    vertical: 12,
                  ),
                ),
                child: const Text('Close'),
              ),
              SizedBox(width: 20),
              ElevatedButton(
                onPressed: () {},
                style: ElevatedButton.styleFrom(
                  backgroundColor: kYellowText,
                  foregroundColor: Colors.black,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 24,
                    vertical: 12,
                  ),
                ),
                child: const Text('Delete All Duplicates'),
              ),
            ],
          ),
      ],
    );
  }
}
