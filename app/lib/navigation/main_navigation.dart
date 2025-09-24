import 'package:app/models/file_tree_node.dart';
import 'package:flutter/material.dart';
import 'package:app/custom_widgets/create_manager.dart';
import 'package:app/constants.dart';
import 'package:app/api.dart';

//class to keep track of main navigation items icon and labels(easy to add more in the future)
class NavigationItem {
  final IconData icon;
  final String label;

  NavigationItem({required this.icon, required this.label});
}

//child class that adds directory field for managers
class ManagerNavigationItem extends NavigationItem {
  final String directory;
  final bool isLoading;
  final FileTreeNode? treeData;

  ManagerNavigationItem({
    required super.icon,
    required super.label,
    required this.directory,
    this.isLoading = false,
    this.treeData,
  });
}

class MainNavigation extends StatefulWidget {
  final List<NavigationItem> items;
  final int selectedIndex;
  final Function(int) onTap;
  final Function() updateStats;
  final Function(String, FileTreeNode?)? onManagerTap;
  final Function(String, FileTreeNode)? onManagerTreeDataUpdate;
  final Function(List<String>)? onManagerNamesUpdate;
  final Function(String)? onManagerAdded;
  final Function(String)? onManagerDelete;
  final String? selectedManager;

  const MainNavigation({
    super.key,
    required this.items,
    required this.selectedIndex,
    required this.onTap,
    required this.updateStats,
    this.onManagerTap,
    this.onManagerTreeDataUpdate,
    this.onManagerNamesUpdate,
    this.onManagerAdded,
    this.onManagerDelete,
    this.selectedManager,
  });

  @override
  State<MainNavigation> createState() => MainNavigationState();
}

class MainNavigationState extends State<MainNavigation> {
  //has a list of managers that are created
  final List<ManagerNavigationItem> _managers = [];
  bool _isInitialized = false;
  bool _disposed = false;
  bool _isInitialLoading = true;
  int _loadingManagersCount = 0;

  @override
  void dispose() {
    _disposed = true;
    super.dispose();
  }

  //if manager exist, ignore otherwise proceed
  bool _managerNameExists(String name) {
    for (ManagerNavigationItem item in _managers) {
      if (name.toLowerCase() == item.label.toLowerCase()) {
        return false;
      }
    }
    return true;
  }

  void _showApiError(String message) {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            title: const Text("Error"),
            content: Text(message),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text("OK"),
              ),
            ],
          ),
    );
  }

  void _handleDeleteManager(String managerName) async {
    try {
      final success = await Api.deleteSmartManager(managerName);

      if (success) {
        setState(() {
          _managers.removeWhere((m) => m.label == managerName);
        });

        // Notify parent about the deletion
        widget.onManagerDelete?.call(managerName);

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Smart Manager "$managerName" deleted successfully'),
            backgroundColor: kYellowText,
            duration: Duration(seconds: 2),
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to delete Smart Manager "$managerName"'),
            backgroundColor: Colors.redAccent,
            duration: Duration(seconds: 3),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error deleting Smart Manager "$managerName": $e'),
          backgroundColor: Colors.redAccent,
          duration: Duration(seconds: 3),
        ),
      );
    }
  }

  void removeManagerFromNavigation(String managerName) {
    if (mounted) {
      setState(() {
        _managers.removeWhere((m) => m.label == managerName);
      });
    }
  }

  Future<void> loadTreeDataForManager(String managerName) async {
    final index = _managers.indexWhere((m) => m.label == managerName);
    if (index == -1) return;

    // If already loaded or loading, return
    if (_managers[index].treeData != null || _managers[index].isLoading) {
      return;
    }

    // Check if widget is still mounted before updating state
    if (!_disposed && mounted) {
      setState(() {
        _managers[index] = ManagerNavigationItem(
          icon: _managers[index].icon,
          label: _managers[index].label,
          directory: _managers[index].directory,
          isLoading: true,
          treeData: _managers[index].treeData,
        );
      });
    }

    try {
      final treeData = await Api.loadTreeData(managerName);

      if (!_disposed && mounted) {
        setState(() {
          _managers[index] = ManagerNavigationItem(
            icon: Icons.folder,
            label: managerName,
            directory: _managers[index].directory,
            isLoading: false,
            treeData: treeData,
          );
        });

        // Notify parent
        widget.onManagerTreeDataUpdate?.call(managerName, treeData);
      }
    } catch (e) {
      print('Error loading tree data for $managerName: $e');

      if (!_disposed && mounted) {
        setState(() {
          _managers[index] = ManagerNavigationItem(
            icon: Icons.folder,
            label: managerName,
            directory: _managers[index].directory,
            isLoading: false,
            treeData: null,
          );
        });
      }
    }
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _initializeApp();
    });
  }

  void _initializeApp() {
    if (_isInitialLoading) {
      _showBootupScreen();
    }
    _loadExistingManagers();
  }

  void _showBootupScreen() {
    showDialog<void>(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return PopScope(
          canPop: false,
          child: Container(
            width: MediaQuery.of(context).size.width,
            height: MediaQuery.of(context).size.height,
            color: kScaffoldColor,
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Image.asset('images/logo.png', width: 120, height: 120),
                  const SizedBox(height: 40),
                  Text(
                    'SMART FILE MANAGER',
                    style: TextStyle(
                      color: kprimaryColor,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 60),
                  CircularProgressIndicator(color: kYellowText, strokeWidth: 3),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  void _hideBootupScreen() {
    if (!_disposed && mounted && _isInitialLoading) {
      Navigator.of(context).pop();
      setState(() {
        _isInitialLoading = false;
      });
      widget.updateStats.call();
    }
  }

  void _loadExistingManagers() async {
    if (_isInitialized) return;

    try {
      final startupResponse = await Api.startUp();

      if (!_disposed && mounted) {
        setState(() {
          _managers.clear();
          for (String managerName in startupResponse.managerNames) {
            _managers.add(
              ManagerNavigationItem(
                icon: Icons.folder,
                label: managerName,
                directory: '',
                isLoading: true,
              ),
            );
          }
          _isInitialized = true;
          _loadingManagersCount = startupResponse.managerNames.length;
        });

        // Notify parent with manager names
        widget.onManagerNamesUpdate?.call(startupResponse.managerNames);

        // If no managers, hide bootup screen immediately
        if (startupResponse.managerNames.isEmpty) {
          _hideBootupScreen();
        } else {
          // Load tree data for each manager in background
          for (String managerName in startupResponse.managerNames) {
            _loadTreeDataInBackground(managerName);
          }
        }
      }
    } catch (e) {
      print('Error loading existing managers: $e');
      if (!_disposed && mounted) {
        setState(() {
          _isInitialized = true;
        });
        _hideBootupScreen();
      }
    }
  }

  void _loadTreeDataInBackground(String managerName) async {
    try {
      final treeData = await Api.loadTreeData(managerName);

      if (!_disposed && mounted) {
        setState(() {
          final index = _managers.indexWhere((m) => m.label == managerName);
          if (index != -1) {
            _managers[index] = ManagerNavigationItem(
              icon: Icons.folder,
              label: managerName,
              directory: _managers[index].directory,
              isLoading: false,
              treeData: treeData,
            );
          }
          _loadingManagersCount--;
        });

        widget.onManagerTreeDataUpdate?.call(managerName, treeData);

        // Hide bootup screen when all managers are loaded
        if (_loadingManagersCount <= 0) {
          _hideBootupScreen();
        }
      }
    } catch (e) {
      print('Error loading tree data for $managerName: $e');

      if (!_disposed && mounted) {
        setState(() {
          final index = _managers.indexWhere(
            (m) => m.label == managerName && m.isLoading,
          );
          if (index != -1) {
            _managers[index] = ManagerNavigationItem(
              icon: Icons.folder,
              label: managerName,
              directory: _managers[index].directory,
              isLoading: false,
              treeData: null,
            );
          }
          _loadingManagersCount--;
        });

        // Hide bootup screen even if loading failed
        if (_loadingManagersCount <= 0) {
          _hideBootupScreen();
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 250,
      padding: EdgeInsets.only(bottom: 50),
      decoration: BoxDecoration(
        border: Border(right: BorderSide(color: kOutlineBorder)),
        color: kScaffoldColor,
      ),
      child: Column(
        children: [
          //add all tab
          const SizedBox(height: 20),
          ...widget.items.asMap().entries.map((entry) {
            int index = entry.key;
            NavigationItem item = entry.value;

            return HoverableNavigationTile(
              icon: item.icon,
              label: item.label,
              selected:
                  index == widget.selectedIndex &&
                  widget.selectedManager == null,
              onTap: () => widget.onTap(index),
            );
          }),
          //start smart manager section
          Align(
            alignment: Alignment.centerLeft,
            widthFactor: 1.5,
            heightFactor: 2.0,
            child: Text(
              "SMART MANAGERS",
              style: TextStyle(color: Color(0xff9CA3AF)),
            ),
          ),
          Expanded(
            //add all managers to the below scrollable section with ... operator
            child: ListView(
              children: [
                ..._managers.asMap().entries.map((entry) {
                  ManagerNavigationItem item = entry.value;

                  return HoverableNavigationTile(
                    icon: item.icon,
                    label: item.label,
                    selected:
                        widget.selectedManager == item.label && !item.isLoading,
                    isLoading: item.isLoading,
                    onTap:
                        item.isLoading
                            ? null
                            : () {
                              // Load tree data if not available
                              if (item.treeData == null) {
                                loadTreeDataForManager(item.label);
                              }
                              widget.onManagerTap?.call(
                                item.label,
                                item.treeData,
                              );
                            },
                  );
                }),
              ],
            ),
          ),
          //button to create a new manager
          TextButton(
            onPressed: () async {
              final result = await createManager(context);

              if (result != null) {
                bool isUnique = _managerNameExists(result.name);

                if (isUnique) {
                  // Add manager with loading state
                  setState(() {
                    _managers.add(
                      ManagerNavigationItem(
                        icon: Icons.folder,
                        label: result.name,
                        directory: result.directory,
                        isLoading: true,
                      ),
                    );
                  });

                  try {
                    // Attempt to create the manager via API
                    final success = await Api.addSmartManager(
                      result.name,
                      result.directory,
                    );

                    if (success) {
                      // Update manager to remove loading state and load treedata
                      setState(() {
                        final index = _managers.indexWhere(
                          (m) => m.label == result.name,
                        );
                        if (index != -1) {
                          _managers[index] = ManagerNavigationItem(
                            icon: Icons.folder,
                            label: result.name,
                            directory: result.directory,
                            isLoading: false,
                          );
                        }
                      });

                      //load data
                      loadTreeDataForManager(result.name);

                      // Notify parent that a new manager was added
                      widget.onManagerAdded?.call(result.name);

                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text(
                            'Smart Manager "${result.name}" created successfully',
                          ),
                          backgroundColor: kYellowText,
                          duration: Duration(seconds: 2),
                        ),
                      );
                    } else {
                      // Remove manager API call failed - directory already in use
                      setState(() {
                        _managers.removeWhere((m) => m.label == result.name);
                      });

                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text(
                            'Cannot create Smart Manager "${result.name}": Directory is already contained in another Smart Manager',
                          ),
                          backgroundColor: Colors.redAccent,
                          duration: Duration(seconds: 3),
                        ),
                      );
                    }
                  } catch (e) {
                    // Remove manager  API call threw  exception
                    setState(() {
                      _managers.removeWhere((m) => m.label == result.name);
                    });

                    _showApiError("Error occurred: $e");
                  }
                } else {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        'Smart Manager with that name already exists.',
                      ),
                      backgroundColor: Colors.redAccent,
                      duration: Duration(seconds: 2),
                    ),
                  );
                }
              }
            },
            child: Text(
              "+ Create Smart Manager",
              style: TextStyle(color: kYellowText),
            ),
          ),
        ],
      ),
    );
  }
}

class HoverableNavigationTile extends StatefulWidget {
  final IconData icon;
  final String label;
  final String? tooltip;
  final bool selected;
  final bool isLoading;
  final VoidCallback? onTap;

  const HoverableNavigationTile({
    super.key,
    required this.icon,
    required this.label,
    this.tooltip,
    required this.selected,
    this.isLoading = false,
    required this.onTap,
  });

  @override
  State<HoverableNavigationTile> createState() =>
      _HoverableNavigationTileState();
}

class _HoverableNavigationTileState extends State<HoverableNavigationTile> {
  bool _hovering = false;

  @override
  Widget build(BuildContext context) {
    final isSelected = widget.selected;
    final isLoading = widget.isLoading;

    Color bgColor =
        isSelected
            ? kYellowText
            : _hovering && !isLoading
            ? kAppBarColor
            : Colors.transparent;

    Color iconTextColor = isSelected ? Colors.black : const Color(0xffF5F5F5);

    Widget tile = Container(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      padding: const EdgeInsets.symmetric(horizontal: 5),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(5),
      ),
      child: ListTile(
        leading:
            isLoading
                ? SizedBox(
                  width: 24,
                  height: 24,
                  child: CircularProgressIndicator(
                    color: Color(0xffFFB400),
                    strokeWidth: 2,
                  ),
                )
                : Icon(widget.icon, color: iconTextColor),
        title: Text(
          widget.label,
          style: TextStyle(color: isLoading ? Colors.grey : iconTextColor),
        ),
        onTap: widget.onTap,
      ),
    );

    Widget wrappedTile = MouseRegion(
      onEnter: isLoading ? null : (_) => setState(() => _hovering = true),
      onExit: isLoading ? null : (_) => setState(() => _hovering = false),
      child: tile,
    );

    // Add tooltip
    if (widget.tooltip != null && !isLoading) {
      return Tooltip(message: widget.tooltip!, child: wrappedTile);
    }

    return wrappedTile;
  }
}
