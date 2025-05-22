import 'package:flutter/material.dart';

class NavigationItem {
  final IconData icon;
  final String label;

  NavigationItem({required this.icon, required this.label});
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
  List<NavigationItem> _managers = [];

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 250,
      padding: EdgeInsets.fromLTRB(0, 0, 0, 50),
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
                  NavigationItem item = entry.value;

                  return HoverableNavigationTile(
                    icon: item.icon,
                    label: item.label,
                    selected: widget.selectedManager == item.label,
                    onTap: () => widget.onManagerTap?.call(item.label),
                  );
                }),
              ],
            ),
          ),
          TextButton(
            onPressed: () {
              setState(() {
                _managers.add(
                  NavigationItem(
                    icon: Icons.circle,
                    label: 'Manager ${_managers.length + 1}',
                  ),
                );
              });
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
  final bool selected;
  final VoidCallback onTap;

  const HoverableNavigationTile({
    super.key,
    required this.icon,
    required this.label,
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

    return MouseRegion(
      onEnter: (_) => setState(() => _hovering = true),
      onExit: (_) => setState(() => _hovering = false),
      child: Container(
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
      ),
    );
  }
}
