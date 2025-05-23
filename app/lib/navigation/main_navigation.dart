import 'package:flutter/material.dart';
import 'package:app/custom_widgets/create_manager.dart';

class NavigationItem {
  final IconData icon;
  final String label;

  NavigationItem({required this.icon, required this.label});
}

class ManagerNavigationItem extends NavigationItem {
  final String directory;

  ManagerNavigationItem({
    required super.icon,
    required super.label,
    required this.directory,
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
  final List<ManagerNavigationItem> _managers = [];

  bool _managerNameExists(String name) {
    for (ManagerNavigationItem item in _managers) {
      if (name.toLowerCase() == item.label.toLowerCase()) {
        return false;
      }
    }
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 250,
      padding: EdgeInsets.only(bottom: 50),
      decoration: BoxDecoration(
        border: Border(right: BorderSide(color: Color(0xff3D3D3D))),
        color: const Color(0xff1E1E1E),
      ),
      child: Column(
        children: [
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
            child: ListView(
              children: [
                ..._managers.asMap().entries.map((entry) {
                  ManagerNavigationItem item = entry.value;

                  return HoverableNavigationTile(
                    icon: item.icon,
                    label: item.label,
                    tooltip: "Directory: ${item.directory}",
                    selected: widget.selectedManager == item.label,
                    onTap: () => widget.onManagerTap?.call(item.label),
                  );
                }),
              ],
            ),
          ),
          TextButton(
            onPressed: () async {
              final result = await createManager(context);

              if (result != null) {
                bool isUnique = _managerNameExists(result.name);

                if (isUnique) {
                  setState(() {
                    _managers.add(
                      ManagerNavigationItem(
                        icon: Icons.folder,
                        label: result.name,
                        directory: result.directory,
                      ),
                    );
                  });
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        'Smart Manager "${result.name}" created successfully',
                      ),
                      backgroundColor: Color(0xffFFB400),
                      duration: Duration(seconds: 2),
                    ),
                  );
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
              style: TextStyle(color: Color(0xffFFB400)),
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
  final VoidCallback onTap;

  const HoverableNavigationTile({
    super.key,
    required this.icon,
    required this.label,
    this.tooltip,
    required this.selected,
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

    Color bgColor =
        isSelected
            ? const Color(0xffFFB400)
            : _hovering
            ? const Color(0xff2E2E2E)
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
        leading: Icon(widget.icon, color: iconTextColor),
        title: Text(widget.label, style: TextStyle(color: iconTextColor)),
        onTap: widget.onTap,
      ),
    );

    return MouseRegion(
      onEnter: (_) => setState(() => _hovering = true),
      onExit: (_) => setState(() => _hovering = false),
      child: tile,
    );
  }
}
