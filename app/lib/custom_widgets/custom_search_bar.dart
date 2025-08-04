import 'package:flutter/material.dart';

class CustomSearchBar extends StatefulWidget {
  final String? hint;
  final TextEditingController? controller;
  final void Function(String)? onChanged;
  final IconData? icon;
  final bool isActive;
  final ValueChanged<bool>? onActiveChanged;

  const CustomSearchBar({
    super.key,
    this.hint,
    this.controller,
    this.onChanged,
    this.icon,
    this.isActive = false,
    this.onActiveChanged,
  });

  @override
  State<CustomSearchBar> createState() => _CustomSearchBarState();
}

class _CustomSearchBarState extends State<CustomSearchBar> {
  bool _isHovered = false;
  bool _hasFocus = false;
  late FocusNode _focusNode;

  @override
  void initState() {
    super.initState();
    _focusNode =
        FocusNode()..addListener(() {
          setState(() {
            _hasFocus = _focusNode.hasFocus;
          });
        });
  }

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  InputBorder _buildBorder() => OutlineInputBorder(
    borderRadius: BorderRadius.circular(6),
    borderSide: const BorderSide(color: Color(0xff3D3D3D)),
  );

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        decoration: BoxDecoration(
          color:
              widget.isActive
                  ? const Color(0xff242424)
                  : const Color(0xff1a1a1a),
          borderRadius: BorderRadius.circular(6),
          border: Border.all(
            color:
                widget.isActive
                    ? const Color(0xff3D3D3D)
                    : const Color(0xff2a2a2a),
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            ConstrainedBox(
              constraints: const BoxConstraints(minWidth: 200, maxWidth: 400),
              child: IntrinsicWidth(
                child: Row(
                  children: [
                    if (widget.icon != null) ...[
                      Icon(
                        widget.icon,
                        size: 16,
                        color:
                            widget.isActive
                                ? const Color(0xff9CA3AF)
                                : const Color(0xff6b7280),
                      ),
                      const SizedBox(width: 8),
                    ],
                    Expanded(
                      child: TextField(
                        focusNode: _focusNode,
                        controller: widget.controller,
                        onChanged: widget.onChanged,
                        enabled: widget.isActive,
                        style: TextStyle(
                          fontSize: 12,
                          color:
                              widget.isActive
                                  ? const Color(0xff9CA3AF)
                                  : const Color(0xff6b7280),
                        ),
                        cursorColor:
                            widget.isActive
                                ? const Color(0xff9CA3AF)
                                : const Color(0xff6b7280),
                        decoration: InputDecoration(
                          isDense: true,
                          hintText: widget.hint,
                          hintStyle: TextStyle(
                            fontSize: 12,
                            color:
                                widget.isActive
                                    ? const Color(0xff9CA3AF)
                                    : const Color(0xff6b7280),
                          ),
                          border: InputBorder.none,
                          contentPadding: const EdgeInsets.symmetric(
                            vertical: 8,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            if (widget.onActiveChanged != null) ...[
              const SizedBox(width: 8),
              GestureDetector(
                onTap: () => widget.onActiveChanged?.call(!widget.isActive),
                child: Icon(
                  widget.isActive ? Icons.toggle_on : Icons.toggle_off,
                  size: 20,
                  color:
                      widget.isActive
                          ? const Color(0xff9CA3AF)
                          : const Color(0xff6b7280),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
