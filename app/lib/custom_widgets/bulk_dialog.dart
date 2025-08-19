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

  String _selectedCategory = "All Files";
  String? _selectedFileType;

  String _selectedBulkOperation = "Bulk Delete";
  String? _selectedTagFilter;

  final List<String> _bulkOperations = [
    "Bulk Delete",
    "Bulk Add Tag",
    "Bulk Remove Tag",
  ];

  final Map<String, List<String>> _fileTypeMap = {
    "All Files": [],
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
    if (_selectedBulkOperation == "Bulk Remove Tag") {
      _loadBulkData(widget.name, "TAGS", widget.umbrella);
    } else {
      String apiType =
          _selectedCategory == "All Files" ? "ALL" : _selectedCategory;
      _loadBulkData(widget.name, apiType, widget.umbrella);
    }
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
      String apiType =
          _selectedCategory == "All Files" ? "ALL" : _selectedCategory;
      _loadBulkData(widget.name, apiType, true);
    }
  }

  void _onBulkOperationChanged(String? newBulkOperation) {
    if (newBulkOperation != null &&
        newBulkOperation != _selectedBulkOperation) {
      setState(() {
        _selectedCategory = "All Files";
        _selectedFileType = null;
        _selectedBulkOperation = newBulkOperation;
        _selectedTagFilter = null;
      });
      if (newBulkOperation == "Bulk Remove Tag") {
        _loadBulkData(widget.name, "TAGS", true);
      } else {
        String apiType =
            _selectedCategory == "All Files" ? "ALL" : _selectedCategory;
        _loadBulkData(widget.name, apiType, true);
      }
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
        String apiType =
            _selectedCategory == "All Files" ? "ALL" : _selectedCategory;
        _loadBulkData(widget.name, apiType, true);
      }
    }
  }

  void _onTagFilterChanged(String? newTag) {
    setState(() {
      _selectedTagFilter = newTag;
    });
  }

  List<String> _getAllAvailableTags() {
    Set<String> allTags = {};
    if (widget.files != null) {
      for (FileModel file in widget.files!) {
        if (file.fileTags != null) {
          allTags.addAll(file.fileTags!);
        }
      }
    }
    return allTags.toList()..sort();
  }

  List<String> _getTagsFromSelectedFiles() {
    Set<String> selectedFileTags = {};
    for (FileModel file in _currentSelectedFiles) {
      if (file.fileTags != null) {
        selectedFileTags.addAll(file.fileTags!);
      }
    }
    return selectedFileTags.toList()..sort();
  }

  void _toggleSelectAll() {
    setState(() {
      _selectAll = !_selectAll;
      _currentSelectedFiles.clear();
      if (_selectAll && widget.files != null) {
        // Filter files based on bulk operation and tag filter before adding to selection
        List<FileModel> filesToAdd =
            widget.files!.where((file) {
              if (_selectedBulkOperation == "Bulk Remove Tag") {
                if (file.fileTags == null || file.fileTags!.isEmpty) {
                  return false;
                }
                if (_selectedTagFilter != null) {
                  return file.fileTags!.contains(_selectedTagFilter);
                }
                return true;
              }
              return true;
            }).toList();
        _currentSelectedFiles.addAll(filesToAdd);
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
      // Calculate the total number of available files based on filter
      int totalAvailableFiles = 0;
      if (widget.files != null) {
        totalAvailableFiles =
            widget.files!.where((file) {
              if (_selectedBulkOperation == "Bulk Remove Tag") {
                if (file.fileTags == null || file.fileTags!.isEmpty) {
                  return false;
                }
                if (_selectedTagFilter != null) {
                  return file.fileTags!.contains(_selectedTagFilter);
                }
                return true;
              }
              return true;
            }).length;
      }
      _selectAll = _currentSelectedFiles.length == totalAvailableFiles;
    });
  }

  String _convertToJsonDelete() {
    List<Map<String, String>> fileList =
        _currentSelectedFiles
            .map((file) => {"file_path": file.filePath})
            .toList();
    return jsonEncode(fileList);
  }

  void _deleteMultipleFiles(String managerName, List<FileModel> files) async {
    String jsonPaths = _convertToJsonDelete();
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

  String _convertToJsonAddTag(String tag) {
    List<Map<String, dynamic>> fileList =
        _currentSelectedFiles
            .map(
              (file) => {
                "file_path": file.filePath,
                "tags": [tag],
              },
            )
            .toList();
    return jsonEncode(fileList);
  }

  String _convertToJsonRemoveTag(String? tag) {
    List<Map<String, dynamic>> fileList = [];

    for (FileModel file in _currentSelectedFiles) {
      if (tag == null) {
        // Remove all tags from this file
        fileList.add({"file_path": file.filePath, "tags": file.fileTags ?? []});
      } else {
        // Only include this file if it has the specific tag to remove
        if (file.fileTags != null && file.fileTags!.contains(tag)) {
          fileList.add({
            "file_path": file.filePath,
            "tags": [tag],
          });
        }
      }
    }

    return jsonEncode(fileList);
  }

  void _showAddTagDialog() {
    final TextEditingController tagController = TextEditingController();

    showDialog(
      context: context,
      builder: (BuildContext dialogContext) {
        return AlertDialog(
          backgroundColor: kScaffoldColor,
          title: const Text('Add Tag', style: kTitle1),
          content: TextField(
            controller: tagController,
            style: const TextStyle(color: Colors.white),
            decoration: const InputDecoration(
              hintText: 'Enter tag name',
              hintStyle: TextStyle(color: Color(0xff9CA3AF)),
              enabledBorder: UnderlineInputBorder(
                borderSide: BorderSide(color: Colors.grey),
              ),
              focusedBorder: UnderlineInputBorder(
                borderSide: BorderSide(color: kYellowText),
              ),
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.of(dialogContext).pop();
                tagController.dispose();
              },
              child: const Text('Cancel', style: TextStyle(color: Colors.grey)),
            ),
            ElevatedButton(
              onPressed: () {
                final tag = tagController.text.trim();
                if (tag.isNotEmpty) {
                  _tagMultipleFiles(widget.name, tag);
                }
                tagController.dispose();
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: kYellowText,
                foregroundColor: Colors.black,
              ),
              child: const Text('Add Tag'),
            ),
          ],
        );
      },
    );
  }

  void _showRemoveTagDialog() {
    List<String> availableTags = _getTagsFromSelectedFiles();
    String? selectedTag;

    // Check if any selected files have tags
    if (availableTags.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('No tags found on selected files'),
          backgroundColor: Colors.orange,
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    showDialog(
      context: context,
      builder: (BuildContext dialogContext) {
        return StatefulBuilder(
          builder: (context, setState) {
            return AlertDialog(
              backgroundColor: kScaffoldColor,
              title: const Text('Remove Tag', style: kTitle1),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Select which tags to remove from selected files:',
                    style: TextStyle(color: Colors.white, fontSize: 16),
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.maxFinite,
                    child: CustomDropdownMenu<String>(
                      items: [
                        const DropdownMenuItem<String>(
                          value: "ALL_TAGS",
                          child: Text('All Tags'),
                        ),
                        ...availableTags.map((String tag) {
                          return DropdownMenuItem<String>(
                            value: tag,
                            child: Text(tag),
                          );
                        }),
                      ],
                      value: selectedTag,
                      onChanged: (String? newValue) {
                        setState(() {
                          selectedTag = newValue;
                        });
                      },
                      hint: "Select tag to remove",
                      minWidth: 200,
                      maxWidth: 300,
                    ),
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: () {
                    Navigator.of(dialogContext).pop();
                  },
                  child: const Text(
                    'Cancel',
                    style: TextStyle(color: Colors.grey),
                  ),
                ),
                ElevatedButton(
                  onPressed:
                      selectedTag != null
                          ? () {
                            Navigator.of(dialogContext).pop();
                            _removeTagsFromFiles(widget.name, selectedTag);
                          }
                          : null,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: kYellowText,
                    foregroundColor: Colors.black,
                  ),
                  child: const Text('Remove Tag'),
                ),
              ],
            );
          },
        );
      },
    );
  }

  void _tagMultipleFiles(String managerName, String tag) async {
    String jsonPaths = _convertToJsonAddTag(tag);
    FileTreeNode response = await Api.bulkAddTag(managerName, jsonPaths);
    if (response.name == managerName) {
      setState(() {
        widget.files?.clear();
        _currentSelectedFiles.clear();
      });
      widget.updateOnDelete.call(managerName, response);
      if (mounted) {
        Navigator.pop(context);
      }
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Successfully added tag "$tag" to selected files'),
          backgroundColor: kYellowText,
          duration: Duration(seconds: 2),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Could not add tag to files'),
          backgroundColor: Colors.redAccent,
          duration: Duration(seconds: 2),
        ),
      );
    }
  }

  void _removeTagsFromFiles(String managerName, String? selectedTag) async {
    String? tagToRemove = selectedTag == "ALL_TAGS" ? null : selectedTag;

    // Count how many files will be affected
    int affectedFilesCount = 0;
    for (FileModel file in _currentSelectedFiles) {
      if (tagToRemove == null) {
        // Count files that have any tags
        if (file.fileTags != null && file.fileTags!.isNotEmpty) {
          affectedFilesCount++;
        }
      } else {
        // Count files that have the specific tag
        if (file.fileTags != null && file.fileTags!.contains(tagToRemove)) {
          affectedFilesCount++;
        }
      }
    }

    if (affectedFilesCount == 0) {
      String message =
          selectedTag == "ALL_TAGS"
              ? 'No tags found on selected files'
              : 'Tag "$selectedTag" not found on any selected files';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(message),
          backgroundColor: Colors.orange,
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    String jsonPaths = _convertToJsonRemoveTag(tagToRemove);
    FileTreeNode response = await Api.bulkRemoveTag(managerName, jsonPaths);
    if (response.name == managerName) {
      setState(() {
        widget.files?.clear();
        _currentSelectedFiles.clear();
      });
      widget.updateOnDelete.call(managerName, response);
      if (mounted) {
        Navigator.pop(context);
      }
      String message =
          selectedTag == "ALL_TAGS"
              ? 'Successfully removed all tags from $affectedFilesCount file(s)'
              : 'Successfully removed tag "$selectedTag" from $affectedFilesCount file(s)';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(message),
          backgroundColor: kYellowText,
          duration: Duration(seconds: 2),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Could not remove tags from files'),
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
                if (_selectedBulkOperation != "Bulk Remove Tag") ...[
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
                      hint: "All Files",
                      minWidth: 120,
                      maxWidth: 180,
                    ),
                  ),
                  if (_selectedCategory != "All Files") ...[
                    SizedBox(width: 10),
                    Expanded(
                      child: CustomDropdownMenu<String>(
                        items:
                            _fileTypeMap[_selectedCategory]!.map((
                              String fileType,
                            ) {
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
                ],
              ],
            ),
            // Show tag filter dropdown only for Bulk Remove Tag operation
            if (_selectedBulkOperation == "Bulk Remove Tag") ...[
              const SizedBox(height: 10),
              Row(
                children: [
                  const Text(
                    'Filter by tag: ',
                    style: TextStyle(color: Colors.white, fontSize: 14),
                  ),
                  Expanded(
                    child: CustomDropdownMenu<String>(
                      items: [
                        const DropdownMenuItem<String>(
                          value: null,
                          child: Text('All tags'),
                        ),
                        ..._getAllAvailableTags().map((String tag) {
                          return DropdownMenuItem<String>(
                            value: tag,
                            child: Text(tag),
                          );
                        }),
                      ],
                      value: _selectedTagFilter,
                      onChanged: _onTagFilterChanged,
                      hint: "Select tag to filter",
                      minWidth: 120,
                      maxWidth: 250,
                    ),
                  ),
                ],
              ),
            ],
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
                      'Select All (${_currentSelectedFiles.length}/${widget.files!.where((file) {
                        if (_selectedBulkOperation == "Bulk Remove Tag") {
                          if (file.fileTags == null || file.fileTags!.isEmpty) {
                            return false;
                          }
                          if (_selectedTagFilter != null) {
                            return file.fileTags!.contains(_selectedTagFilter);
                          }
                          return true;
                        }
                        return true;
                      }).length})',
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
                            widget.files!
                                .where((object) {
                                  // Filter files based on bulk operation and tag filter
                                  if (_selectedBulkOperation ==
                                      "Bulk Remove Tag") {
                                    if (object.fileTags == null ||
                                        object.fileTags!.isEmpty) {
                                      return false;
                                    }
                                    if (_selectedTagFilter != null) {
                                      return object.fileTags!.contains(
                                        _selectedTagFilter,
                                      );
                                    }
                                    return true;
                                  }
                                  return true;
                                })
                                .map((object) {
                                  bool isSelected = _currentSelectedFiles
                                      .contains(object);
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
                                                    _toggleFileSelection(
                                                      object,
                                                    ),
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
                                                  overflow:
                                                      TextOverflow.ellipsis,
                                                ),
                                                const SizedBox(height: 6),
                                                _selectedBulkOperation ==
                                                        "Bulk Delete"
                                                    ? Text(
                                                      'Path: ${object.filePath}',
                                                      style: const TextStyle(
                                                        color: Color(
                                                          0xff9CA3AF,
                                                        ),
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
                                                        color: Color(
                                                          0xff9CA3AF,
                                                        ),
                                                        fontSize: 12,
                                                      ),
                                                    )
                                                    : Wrap(
                                                      spacing: 8,
                                                      runSpacing: 8,
                                                      children:
                                                          object.fileTags!
                                                              .map(
                                                                (
                                                                  tag,
                                                                ) => Container(
                                                                  padding:
                                                                      const EdgeInsets.symmetric(
                                                                        horizontal:
                                                                            12,
                                                                        vertical:
                                                                            6,
                                                                      ),
                                                                  decoration: BoxDecoration(
                                                                    color: const Color(
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
                                                                        width:
                                                                            4,
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
                                })
                                .toList(),
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
                        _showAddTagDialog();
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
                        _showRemoveTagDialog();
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
