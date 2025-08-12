import 'package:flutter/material.dart';

class HoverableButton extends StatefulWidget {
  final VoidCallback? onTap;
  final String name;
  final IconData icon;
  final bool expanded;

  const HoverableButton({
    super.key,
    this.onTap,
    required this.name,
    required this.icon,
    this.expanded = false,
  });

  @override
  State<HoverableButton> createState() => _HoverableButtonState();
}

class _HoverableButtonState extends State<HoverableButton> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final child = Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: _isHovered ? const Color(0xffFFB400) : const Color(0xff242424),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: const Color(0xff3D3D3D)),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        mainAxisSize: MainAxisSize.max,
        children: [
          Icon(
            widget.icon,
            size: 16,
            color: _isHovered ? Colors.black : const Color(0xff9CA3AF),
          ),
          const SizedBox(width: 4),
          Text(
            widget.name,
            style: TextStyle(
              fontSize: 12,
              color: _isHovered ? Colors.black : const Color(0xff9CA3AF),
            ),
          ),
        ],
      ),
    );

    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: widget.onTap ?? () {},
        child:
            widget.expanded
                ? SizedBox(width: double.infinity, child: child)
                : child,
      ),
    );
  }
}
