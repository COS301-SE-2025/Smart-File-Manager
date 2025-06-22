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

  ManagerNavigationItem({
    required super.icon,
    required super.label,
    required this.directory,
    this.isLoading = false,
  });
}

class MainNavigation extends StatefulWidget {
  final List<NavigationItem> items;
  final int selectedIndex;
  final Function(int) onTap;
  final Function(String)? onManagerTap;
  final String? selectedManager;

  const MainNavigation({
    super.key,
    required this.items,
    required this.selectedIndex,
    required this.onTap,
    this.onManagerTap,
    this.selectedManager,
  });

  @override
  State<MainNavigation> createState() => _MainNavigationState();
}

class _MainNavigationState extends State<MainNavigation> {
  //has a list of managers that are created
  final List<ManagerNavigationItem> _managers = [];

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
                            : () => widget.onManagerTap?.call(item.label),
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
                      // Update manager to remove loading state
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
                      // Remove manager API call failed
                      setState(() {
                        _managers.removeWhere((m) => m.label == result.name);
                      });

                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text(
                            'Failed to create Smart Manager "${result.name}"',
                          ),
                          backgroundColor: Colors.redAccent,
                          duration: Duration(seconds: 2),
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
