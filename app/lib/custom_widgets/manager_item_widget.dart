import 'package:flutter/material.dart';

class ManagerItemWidget extends StatefulWidget {
  final String managerName;
  final String managerFolders;
  final String managerFiles;
  final String managerSize;

  const ManagerItemWidget({
    required this.managerName,
    required this.managerFiles,
    required this.managerFolders,
    required this.managerSize,
    super.key,
  });

  @override
  State<ManagerItemWidget> createState() => _ManagerItemWidgetState();
}

class _ManagerItemWidgetState extends State<ManagerItemWidget> {
  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        return SizedBox(
          width: constraints.maxWidth,
          child: Positioned.fill(
            child: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: const Color(0xff242424),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: const Color(0xff3D3D3D), width: 1),
              ),
              child: Row(
                children: [
                  const SizedBox(width: 8),
                  Expanded(
                    flex: 2,
                    child: Text(
                      widget.managerName,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 15,
                        fontWeight: FontWeight.w500,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  Expanded(
                    flex: 1,
                    child: Text(
                      "Files: ${widget.managerFiles}",
                      style: const TextStyle(
                        color: Color(0xff6b7280),
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                      textAlign: TextAlign.right,
                    ),
                  ),
                  Expanded(
                    flex: 1,
                    child: Text(
                      "Folders: ${widget.managerFolders}",
                      style: const TextStyle(
                        color: Color(0xff6b7280),
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                      textAlign: TextAlign.right,
                    ),
                  ),
                  Expanded(
                    flex: 1,
                    child: Text(
                      "Size: ${widget.managerSize}",
                      style: const TextStyle(
                        color: Color(0xff6b7280),
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                      textAlign: TextAlign.right,
                    ),
                  ),
                  const SizedBox(width: 10),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
