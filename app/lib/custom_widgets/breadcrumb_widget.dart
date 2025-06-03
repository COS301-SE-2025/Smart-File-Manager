import 'package:flutter/material.dart';
import 'package:app/constants.dart';

class BreadcrumbWidget extends StatelessWidget {
  final List<String> currentPath;
  final Function(List<String>) onNavigate;
  final double height;
  final Color? backgroundColor;
  final Color? borderColor;

  const BreadcrumbWidget({
    required this.currentPath,
    required this.onNavigate,
    this.height = 40,
    this.backgroundColor,
    this.borderColor,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: backgroundColor ?? const Color(0xff242424),
        border: Border(
          bottom: BorderSide(color: borderColor ?? const Color(0xff3D3D3D)),
        ),
      ),
      child: Row(
        children: [
          Expanded(
            child: Row(
              children: [
                _buildBreadcrumbItem('Root', [], currentPath.isEmpty),
                ...currentPath.asMap().entries.map((entry) {
                  final index = entry.key;
                  final pathSegment = entry.value;
                  final isLast = index == currentPath.length - 1;
                  final pathToHere = currentPath.sublist(0, index + 1);

                  return Row(
                    children: [
                      const Text(
                        '/',
                        style: TextStyle(color: Color(0xff6B7280)),
                      ),
                      _buildBreadcrumbItem(pathSegment, pathToHere, isLast),
                    ],
                  );
                }),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBreadcrumbItem(String name, List<String> path, bool isActive) {
    return GestureDetector(
      onTap: () {
        onNavigate(path);
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        decoration: BoxDecoration(borderRadius: BorderRadius.circular(4)),
        child: Text(
          name,
          style: TextStyle(
            color: isActive ? kYellowText : const Color(0xff9CA3AF),
            fontSize: 12,
            fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}
