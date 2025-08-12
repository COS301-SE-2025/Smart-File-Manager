import 'package:app/constants.dart';
import 'package:app/custom_widgets/custom_dropdown_menu.dart';
import 'package:flutter/material.dart';
import 'package:app/custom_widgets/overview_widget.dart';
import 'package:app/custom_widgets/manager_item_widget.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body: Column(children: [Expanded(child: _buildMainSection())]),
    );
  }

  PreferredSize _buildTopBar() {
    return PreferredSize(
      preferredSize: const Size.fromHeight(50),
      child: AppBar(
        scrolledUnderElevation: 0,
        backgroundColor: const Color(0xff2E2E2E),
        automaticallyImplyLeading: false,
        shape: const Border(
          bottom: BorderSide(color: Color(0xff3D3D3D), width: 1),
        ),
        flexibleSpace: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Row(
              children: [
                Expanded(
                  child: Center(
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Text(
                          'Dashboard',
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 20,
                            color: Colors.white,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildFileList() {
    List<Map<String, String>> sampleFiles = [
      {'name': 'Document1.pdf', 'size': '2.5 MB', 'date': '2024-01-15'},
      {'name': 'Spreadsheet.xlsx', 'size': '1.2 MB', 'date': '2024-01-14'},
      {'name': 'Presentation.pptx', 'size': '5.8 MB', 'date': '2024-01-13'},
      {'name': 'Report.docx', 'size': '3.1 MB', 'date': '2024-01-12'},
      {'name': 'Image.jpg', 'size': '4.2 MB', 'date': '2024-01-11'},
    ];

    return SizedBox(
      height: 300,
      child: Column(
        children:
            sampleFiles.map((file) {
              return Container(
                margin: EdgeInsets.only(bottom: 10),
                padding: EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Color(0xff242424),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Color(0xff3D3D3D), width: 1),
                ),
                child: Row(
                  children: [
                    SizedBox(width: 8),
                    Expanded(
                      flex: 2,
                      child: Text(
                        file['name']!,
                        style: TextStyle(
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
                        file['size']!,
                        style: TextStyle(
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
                        file['date']!,
                        style: TextStyle(
                          color: Color(0xff6b7280),
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                        ),
                        textAlign: TextAlign.right,
                      ),
                    ),
                    SizedBox(width: 10),
                  ],
                ),
              );
            }).toList(),
      ),
    );
  }

  Widget _buildFileTypeStats() {
    List<Map<String, dynamic>> fileTypes = [
      {'type': 'PDF', 'count': 100, 'total': 230, 'color': Colors.red},
      {'type': 'XLSX', 'count': 45, 'total': 230, 'color': Colors.green},
      {'type': 'DOCX', 'count': 30, 'total': 230, 'color': Colors.blue},
      {'type': 'PPTX', 'count': 25, 'total': 230, 'color': Colors.orange},
      {'type': 'Images', 'count': 20, 'total': 230, 'color': Colors.purple},
      {'type': 'Other', 'count': 10, 'total': 230, 'color': Color(0xff9CA3AF)},
    ];

    return Column(
      children:
          fileTypes.map((fileType) {
            double percentage = fileType['count'] / fileType['total'];

            return Container(
              margin: EdgeInsets.only(bottom: 16),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        fileType['type'],
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      Text(
                        '${fileType['count']}/${fileType['total']}',
                        style: TextStyle(
                          color: Color(0xff9CA3AF),
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 6),
                  Container(
                    width: double.infinity,
                    height: 6,
                    decoration: BoxDecoration(
                      color: Color(0xff2E2E2E),
                      borderRadius: BorderRadius.circular(3),
                    ),
                    child: FractionallySizedBox(
                      alignment: Alignment.centerLeft,
                      widthFactor: percentage,
                      child: Container(
                        decoration: BoxDecoration(
                          color: fileType['color'],
                          borderRadius: BorderRadius.circular(3),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            );
          }).toList(),
    );
  }

  Widget _buildMainSection() {
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text("Overview", style: kTitle1),
            Divider(color: Color(0xff3D3D3D)),
            LayoutBuilder(
              builder: (context, constraints) {
                double itemWidth =
                    (constraints.maxWidth - 36) /
                    2; // 2 items per row with spacing
                return Wrap(
                  spacing: 12,
                  runSpacing: 12,
                  children: [
                    SizedBox(
                      width:
                          constraints.maxWidth > 800
                              ? (constraints.maxWidth - 36) / 4
                              : itemWidth,
                      child: OverviewWidget(
                        title: 'Total Managers',
                        value: '4',
                        icon: Icons.manage_accounts_rounded,
                      ),
                    ),
                    SizedBox(
                      width:
                          constraints.maxWidth > 800
                              ? (constraints.maxWidth - 36) / 4
                              : itemWidth,
                      child: OverviewWidget(
                        title: 'Total Files',
                        value: '12 483',
                        icon: Icons.file_present_rounded,
                      ),
                    ),
                    SizedBox(
                      width:
                          constraints.maxWidth > 800
                              ? (constraints.maxWidth - 36) / 4
                              : itemWidth,
                      child: OverviewWidget(
                        title: 'Total Folders',
                        value: '2 482',
                        icon: Icons.folder,
                      ),
                    ),
                    SizedBox(
                      width:
                          constraints.maxWidth > 800
                              ? (constraints.maxWidth - 36) / 4
                              : itemWidth,
                      child: OverviewWidget(
                        title: 'Total Storage',
                        value: '126 Gb',
                        icon: Icons.storage_rounded,
                      ),
                    ),
                  ],
                );
              },
            ),
            SizedBox(height: 20),
            Text("Managers", style: kTitle1),
            Divider(color: Color(0xff3D3D3D)),
            ManagerItemWidget(
              managerName: "Manager Name",
              managerFiles: "1 249",
              managerFolders: "456",
              managerSize: "12Gb",
            ),
            SizedBox(height: 10),
            ManagerItemWidget(
              managerName: "Manager Name",
              managerFiles: "1 249",
              managerFolders: "456",
              managerSize: "12Gb",
            ),
            SizedBox(height: 10),
            ManagerItemWidget(
              managerName: "Manager Name",
              managerFiles: "1 249",
              managerFolders: "456",
              managerSize: "12Gb",
            ),

            SizedBox(height: 20),
            Text("Files", style: kTitle1),
            Divider(color: Color(0xff3D3D3D)),
            Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Expanded(
                  flex: 1,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            "Show files by:",
                            style: TextStyle(color: Color(0xff9CA3AF)),
                          ),
                          SizedBox(width: 10),
                          CustomDropdownMenu(
                            items: [
                              DropdownMenuItem(child: Text("Recents")),
                              DropdownMenuItem(child: Text("Oldest")),
                              DropdownMenuItem(child: Text("Largest")),
                            ],
                            hint: "select",
                          ),
                        ],
                      ),
                      SizedBox(height: 20),
                      _buildFileList(),
                    ],
                  ),
                ),
                SizedBox(width: 20),
                Container(width: 1, height: 300, color: Color(0xff3D3D3D)),
                SizedBox(width: 20),
                Expanded(
                  flex: 1,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        "File Types Managed",
                        style: TextStyle(
                          color: Color(0xff9CA3AF),
                          fontSize: 16,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      SizedBox(height: 20),
                      _buildFileTypeStats(),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
