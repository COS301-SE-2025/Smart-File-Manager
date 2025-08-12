import 'dart:convert';
import 'package:app/custom_widgets/custom_dropdown_menu.dart';
import 'package:app/models/file_tree_node.dart';
import 'package:flutter/material.dart';
import 'package:app/constants.dart';
import 'package:app/api.dart';
import 'package:app/models/file_model.dart';

class BulkDialog extends StatefulWidget {
  final String name;
  final String type;
  final bool umbrella;
  List<FileModel>? files;
  final Function(String, FileTreeNode) updateOnDelete;

  BulkDialog({
    required this.name,
    required this.type,
    required this.umbrella,
    required this.updateOnDelete,
    this.files,
    super.key,
  });

  @override
  State<BulkDialog> createState() => _BulkDialogState();
}

class _BulkDialogState extends State<BulkDialog> {
  final TextEditingController _bulkController = TextEditingController();
  bool _isLoading = true;
  final List<FileModel> _currentSelectedFiles = [];
  bool _selectAll = false;

  String _selectedCategory = "Documents";
  String? _selectedFileType;

  String _selectedBulkOperation = "Bulk Delete";

  final List<String> _bulkOperations = [
    "Bulk Delete",
    "Bulk Add Tag",
    "Bulk Remove Tag",
  ];

  final Map<String, List<String>> _fileTypeMap = {
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

  @override
  void dispose() {
    _bulkController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _loadBulkData(widget.name, _selectedCategory, widget.umbrella);
  }

  Future<void> _loadBulkData(String name, String type, bool umbrella) async {
    setState(() {
      _isLoading = true;
    });

    final files = await Api.bulkOperation(name, type, umbrella);

    if (mounted) {
      setState(() {
        widget.files = files ?? [];
        _currentSelectedFiles.clear();
        _selectAll = false;
        _isLoading = false;
      });
    }
  }

  void _onCategoryChanged(String? newCategory) {
    if (newCategory != null && newCategory != _selectedCategory) {
      setState(() {
        _selectedCategory = newCategory;
        _selectedFileType = null;
      });
      _loadBulkData(widget.name, _selectedCategory, true);
    }
  }

  void _onBulkOperationChanged(String? newBulkOperation) {
    if (newBulkOperation != null &&
        newBulkOperation != _selectedBulkOperation) {
      setState(() {
        _selectedCategory = "Documents";
        _selectedFileType = null;
        _selectedBulkOperation = newBulkOperation;
      });
      _loadBulkData(widget.name, _selectedCategory, true);
    }
  }

  void _onFileTypeChanged(String? newFileType) {
    if (newFileType != _selectedFileType) {
      setState(() {
        _selectedFileType = newFileType;
      });
      if (newFileType != null) {
        _loadBulkData(widget.name, newFileType, false);
      } else {
        _loadBulkData(widget.name, _selectedCategory, true);
      }
    }
  }

  void _toggleSelectAll() {
    setState(() {
      _selectAll = !_selectAll;
      _currentSelectedFiles.clear();
      if (_selectAll && widget.files != null) {
        _currentSelectedFiles.addAll(widget.files!);
      }
    });
  }

  void _toggleFileSelection(FileModel file) {
    setState(() {
      if (_currentSelectedFiles.contains(file)) {
        _currentSelectedFiles.remove(file);
      } else {
        _currentSelectedFiles.add(file);
      }
      _selectAll =
          widget.files != null &&
          _currentSelectedFiles.length == widget.files!.length;
    });
  }

  String _convertToJsonDuplicates() {
    List<Map<String, String>> fileList =
        _currentSelectedFiles
            .map((file) => {"file_path": file.filePath})
            .toList();
    return jsonEncode(fileList);
  }

  void _deleteMultipleFiles(String managerName, List<FileModel> files) async {
    String jsonPaths = _convertToJsonDuplicates();
    FileTreeNode response = await Api.bulkDeleteFiles(managerName, jsonPaths);
    if (response.name == managerName) {
      setState(() {
        widget.files?.clear();
        _currentSelectedFiles.clear();
      });
      widget.updateOnDelete.call(managerName, response);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Deleted all files successfully'),
          backgroundColor: kYellowText,
          duration: Duration(seconds: 2),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Could not delete files'),
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
      title: const Text('Bulk Operations', style: kTitle1),
      content: SizedBox(
        width: double.maxFinite,
        height: MediaQuery.of(context).size.height * 0.9,
        child: Column(
          children: [
            Row(
              children: [
                Expanded(
                  child: CustomDropdownMenu<String>(
                    items:
                        _bulkOperations.map((String operation) {
                          return DropdownMenuItem<String>(
                            value: operation,
                            child: Text(operation),
                          );
                        }).toList(),
                    value: _selectedBulkOperation,
                    onChanged: _onBulkOperationChanged,
                    hint: "Select Bulk Operation",
                    minWidth: 120,
                    maxWidth: 180,
                  ),
                ),
                SizedBox(width: 10),
                Expanded(
                  child: CustomDropdownMenu<String>(
                    items:
                        _fileTypeMap.keys.map((String category) {
                          return DropdownMenuItem<String>(
                            value: category,
                            child: Text(category),
                          );
                        }).toList(),
                    value: _selectedCategory,
                    onChanged: _onCategoryChanged,
                    hint: "Documents",
                    minWidth: 120,
                    maxWidth: 180,
                  ),
                ),
                SizedBox(width: 10),
                Expanded(
                  child: CustomDropdownMenu<String>(
                    items:
                        _fileTypeMap[_selectedCategory]!.map((String fileType) {
                          return DropdownMenuItem<String>(
                            value: fileType,
                            child: Text(fileType.toUpperCase()),
                          );
                        }).toList(),
                    value: _selectedFileType,
                    onChanged: _onFileTypeChanged,
                    hint: "Select File type",
                    minWidth: 120,
                    maxWidth: 180,
                  ),
                ),
              ],
            ),
            const Divider(color: Color(0xff3D3D3D)),
            if (widget.files != null && widget.files!.isNotEmpty)
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 8.0,
                  vertical: 4.0,
                ),
                child: Row(
                  children: [
                    Checkbox(
                      value: _selectAll,
                      onChanged: (bool? value) => _toggleSelectAll(),
                      fillColor: WidgetStateProperty.resolveWith<Color?>((
                        states,
                      ) {
                        if (states.contains(WidgetState.selected)) {
                          return kYellowText;
                        }
                        return Colors.transparent;
                      }),
                      side: BorderSide(color: Colors.grey, width: 1.5),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      'Select All (${_currentSelectedFiles.length}/${widget.files!.length})',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),
            Expanded(
              child:
                  _isLoading
                      ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: const [
                            CircularProgressIndicator(color: Color(0xffFFB400)),
                            SizedBox(height: 16),
                            Text(
                              'Loading files...',
                              style: TextStyle(color: Color(0xff9CA3AF)),
                            ),
                          ],
                        ),
                      )
                      : widget.files == null || widget.files!.isEmpty
                      ? const Center(
                        child: Text(
                          'No Files found.',
                          style: TextStyle(
                            color: Color(0xff9CA3AF),
                            fontSize: 16,
                          ),
                        ),
                      )
                      : ListView(
                        shrinkWrap: true,
                        children:
                            widget.files!.map((object) {
                              bool isSelected = _currentSelectedFiles.contains(
                                object,
                              );
                              return Container(
                                margin: const EdgeInsets.symmetric(
                                  vertical: 4.0,
                                  horizontal: 8.0,
                                ),
                                decoration: BoxDecoration(
                                  color:
                                      isSelected
                                          ? const Color(0xff374151)
                                          : const Color(0xff242424),
                                  borderRadius: BorderRadius.circular(8),
                                  border: Border.all(
                                    color:
                                        isSelected
                                            ? const Color(0xffFFB400)
                                            : const Color(0xff3D3D3D),
                                    width: 1,
                                  ),
                                ),
                                child: Padding(
                                  padding: const EdgeInsets.all(12.0),
                                  child: Row(
                                    children: [
                                      Checkbox(
                                        value: isSelected,
                                        onChanged:
                                            (bool? value) =>
                                                _toggleFileSelection(object),
                                        fillColor:
                                            WidgetStateProperty.resolveWith<
                                              Color?
                                            >((states) {
                                              if (states.contains(
                                                WidgetState.selected,
                                              )) {
                                                return kYellowText;
                                              }
                                              return Colors.transparent;
                                            }),
                                        side: BorderSide(
                                          color: Colors.grey,
                                          width: 1.5,
                                        ),
                                      ),
                                      const SizedBox(width: 12),
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          mainAxisSize: MainAxisSize.min,
                                          children: [
                                            Text(
                                              object.name,
                                              style: const TextStyle(
                                                color: Colors.white,
                                                fontSize: 16,
                                                fontWeight: FontWeight.w500,
                                              ),
                                              maxLines: 2,
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                            const SizedBox(height: 6),
                                            _selectedBulkOperation ==
                                                    "Bulk Delete"
                                                ? Text(
                                                  'Path: ${object.filePath}',
                                                  style: const TextStyle(
                                                    color: Color(0xff9CA3AF),
                                                    fontSize: 12,
                                                  ),
                                                  maxLines: 2,
                                                  overflow:
                                                      TextOverflow.ellipsis,
                                                )
                                                : object.fileTags!.isEmpty
                                                ? Text(
                                                  "no tags",
                                                  style: const TextStyle(
                                                    color: Color(0xff9CA3AF),
                                                    fontSize: 12,
                                                  ),
                                                )
                                                : Wrap(
                                                  spacing: 8,
                                                  runSpacing: 8,
                                                  children:
                                                      object.fileTags!
                                                          .map(
                                                            (tag) => Container(
                                                              padding:
                                                                  const EdgeInsets.symmetric(
                                                                    horizontal:
                                                                        12,
                                                                    vertical: 6,
                                                                  ),
                                                              decoration: BoxDecoration(
                                                                color:
                                                                    const Color(
                                                                      0xff374151,
                                                                    ),
                                                                borderRadius:
                                                                    BorderRadius.circular(
                                                                      16,
                                                                    ),
                                                                border: Border.all(
                                                                  color: const Color(
                                                                    0xff4B5563,
                                                                  ),
                                                                ),
                                                              ),
                                                              child: Row(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .min,
                                                                children: [
                                                                  Text(
                                                                    tag,
                                                                    style: const TextStyle(
                                                                      color: Color(
                                                                        0xffE5E7EB,
                                                                      ),
                                                                      fontSize:
                                                                          12,
                                                                    ),
                                                                  ),
                                                                  const SizedBox(
                                                                    width: 4,
                                                                  ),
                                                                ],
                                                              ),
                                                            ),
                                                          )
                                                          .toList(),
                                                ),
                                          ],
                                        ),
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
        if (widget.files == null || widget.files!.isEmpty)
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
              if (_currentSelectedFiles.isNotEmpty) ...[
                SizedBox(width: 20),
                _selectedBulkOperation == "Bulk Delete"
                    ? ElevatedButton(
                      onPressed: () {
                        _deleteMultipleFiles(
                          widget.name,
                          _currentSelectedFiles,
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: kYellowText,
                        foregroundColor: Colors.black,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 24,
                          vertical: 12,
                        ),
                      ),
                      child: const Text('Delete Selected Files'),
                    )
                    : _selectedBulkOperation == "Bulk Add Tag"
                    ? ElevatedButton(
                      onPressed: () {
                        _deleteMultipleFiles(
                          widget.name,
                          _currentSelectedFiles,
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: kYellowText,
                        foregroundColor: Colors.black,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 24,
                          vertical: 12,
                        ),
                      ),
                      child: const Text('Tag Selected Files'),
                    )
                    : ElevatedButton(
                      onPressed: () {
                        _deleteMultipleFiles(
                          widget.name,
                          _currentSelectedFiles,
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: kYellowText,
                        foregroundColor: Colors.black,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 24,
                          vertical: 12,
                        ),
                      ),
                      child: const Text('Remove Tags from Selected Files'),
                    ),
              ],
              SizedBox(width: 20),
            ],
          ),
      ],
    );
  }
}
