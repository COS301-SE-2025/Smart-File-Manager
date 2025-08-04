import 'package:flutter/material.dart';

class CustomDropdownMenu<T> extends StatefulWidget {
  final List<DropdownMenuItem<T>> items;
  final T? value;
  final void Function(T?)? onChanged;
  final String hint;
  final double minWidth;
  final double maxWidth;

  const CustomDropdownMenu({
    super.key,
    required this.items,
    required this.hint,
    this.value,
    this.onChanged,
    this.minWidth = 100,
    this.maxWidth = 150,
  });

  @override
  State<CustomDropdownMenu<T>> createState() => _CustomDropdownMenuState<T>();
}

class _CustomDropdownMenuState<T> extends State<CustomDropdownMenu<T>> {
  bool _isHovered = false;
  bool _hasFocus = false;
  bool _isOpen = false;
  late FocusNode _focusNode;
  final LayerLink _layerLink = LayerLink();
  T? _selected;
  OverlayEntry? _entry;

  static const _bgColor = Color(0xff242424);
  static const _borderColor = Color(0xff3D3D3D);
  static const _textLight = Color(0xff9CA3AF);

  FontWeight get _fontWeight {
    if (_isHovered || _hasFocus || _isOpen) return FontWeight.w600;
    return FontWeight.normal;
  }

  @override
  void initState() {
    super.initState();
    _selected = widget.value;
    _focusNode =
        FocusNode()
          ..addListener(() => setState(() => _hasFocus = _focusNode.hasFocus));
  }

  @override
  void didUpdateWidget(covariant CustomDropdownMenu<T> old) {
    super.didUpdateWidget(old);
    if (widget.value != old.value) {
      _selected = widget.value;
    }
  }

  @override
  void dispose() {
    _focusNode.dispose();
    _removeOverlay();
    super.dispose();
  }

  void _toggleOverlay() {
    if (_isOpen) {
      _removeOverlay();
    } else {
      _showOverlay();
    }
    setState(() => _isOpen = !_isOpen);
  }

  void _showOverlay() {
    _entry = OverlayEntry(
      builder: (context) {
        return GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () {
            _removeOverlay();
            setState(() => _isOpen = false);
          },
          child: Stack(
            children: [
              // transparent layer to catch outside taps
              Positioned.fill(child: Container(color: Colors.transparent)),
              Positioned(
                width: widget.maxWidth,
                child: CompositedTransformFollower(
                  link: _layerLink,
                  showWhenUnlinked: false,
                  offset: const Offset(0, 38),
                  child: Material(
                    color: _bgColor,
                    elevation: 8,
                    borderRadius: BorderRadius.circular(6),
                    child: ConstrainedBox(
                      constraints: const BoxConstraints(maxHeight: 250),
                      child: ListView(
                        padding: EdgeInsets.zero,
                        shrinkWrap: true,
                        children:
                            widget.items.map((item) {
                              final bool selected = item.value == _selected;
                              return InkWell(
                                onTap: () {
                                  _select(item.value);
                                },
                                hoverColor: const Color(0xff3D3D3D),
                                child: Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 12,
                                    vertical: 10,
                                  ),
                                  decoration: BoxDecoration(
                                    color:
                                        selected
                                            ? const Color(0xff3D3D3D)
                                            : Colors.transparent,
                                  ),
                                  child: DefaultTextStyle(
                                    style: TextStyle(
                                      fontSize: 12,
                                      color:
                                          selected ? Colors.white : _textLight,
                                      fontWeight:
                                          selected
                                              ? FontWeight.w600
                                              : FontWeight.normal,
                                    ),
                                    child: item.child,
                                  ),
                                ),
                              );
                            }).toList(),
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
    Overlay.of(context).insert(_entry!);
  }

  void _removeOverlay() {
    _entry?.remove();
    _entry = null;
  }

  void _select(T? value) {
    setState(() {
      _selected = value;
      _isOpen = false;
    });
    widget.onChanged?.call(value);
    _removeOverlay();
  }

  String _displayStringForValue(T value) {
    final matched = widget.items.firstWhere(
      (it) => it.value == value,
      orElse: () => widget.items.first,
    );
    if (matched.child is Text) {
      return (matched.child as Text).data ?? value.toString();
    }
    return value.toString();
  }

  @override
  Widget build(BuildContext context) {
    final displayText =
        _selected != null ? _displayStringForValue(_selected as T) : null;

    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: _toggleOverlay,
        behavior: HitTestBehavior.translucent,
        child: CompositedTransformTarget(
          link: _layerLink,
          child: Focus(
            focusNode: _focusNode,
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 150),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              constraints: BoxConstraints(
                minWidth: widget.minWidth,
                maxWidth: widget.maxWidth,
              ),
              decoration: BoxDecoration(
                color: _bgColor,
                borderRadius: BorderRadius.circular(6),
                border: Border.all(color: _borderColor),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Expanded(
                    child: Text(
                      displayText ?? widget.hint,
                      style: TextStyle(
                        fontSize: 12,
                        color:
                            displayText != null
                                ? _textLight
                                : const Color(0xff9CA3AF),
                        fontWeight:
                            displayText != null
                                ? _fontWeight
                                : FontWeight.normal,
                      ),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  const SizedBox(width: 6),
                  Icon(
                    _isOpen
                        ? Icons.keyboard_arrow_up
                        : Icons.keyboard_arrow_down,
                    size: 16,
                    color: _textLight,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
