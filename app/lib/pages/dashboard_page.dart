import 'package:app/constants.dart';
import 'package:app/custom_widgets/hoverable_button.dart';
import 'package:flutter/material.dart';
import 'package:app/custom_widgets/overview_widget.dart';
import 'package:app/custom_widgets/manager_item_widget.dart';
import 'package:app/api.dart';
import 'package:app/models/stats_model.dart';

class DashboardPage extends StatefulWidget {
  List<String> managerNames;
  DashboardPage({required this.managerNames, super.key});

  @override
  State<DashboardPage> createState() => DashboardPageState();
}

class DashboardPageState extends State<DashboardPage> {
  ManagersStatsResponse? _managerStats;
  final double _gigInBytes = 1e+9;
  bool _gig = false;
  bool _refreshing = false;
  bool _noManagers = true;
  String _selectedFileSort = 'Recents';

  @override
  void initState() {
    super.initState();
    _onPageload();
  }

  void _onPageload() async {
    if (widget.managerNames.isNotEmpty) {
      var response = await Api.loadStatsData();
      setState(() {
        _noManagers = false;
        _managerStats = response;
      });
    } else {
      _noManagers = true;
    }
  }

  void loadStatsData() async {
    var response = await Api.loadStatsData();
    setState(() {
      _managerStats = response;
    });
  }

  void refreshData() async {
    setState(() {
      _refreshing = true;
    });
    var response = await Api.loadStatsData();
    setState(() {
      _managerStats = response;
      _refreshing = false;
    });
  }

  int _sumTotalFiles() {
    int totalfiles = 0;
    for (StatsModel manager in _managerStats!.managers!) {
      totalfiles = totalfiles + manager.files!;
    }
    return totalfiles;
  }

  int _sumTotalFolders() {
    int totalfolders = 0;
    for (StatsModel manager in _managerStats!.managers!) {
      totalfolders = totalfolders + manager.folders!;
    }
    return totalfolders;
  }

  double _sumTotalSize() {
    double totalsize = 0;
    for (StatsModel manager in _managerStats!.managers!) {
      totalsize = totalsize + manager.size!;
    }
    if (totalsize < _gigInBytes) {
      _gig = false;
      return (totalsize / 1e+6).roundToDouble();
    } else {
      _gig = true;
      return (totalsize / 8e+9).roundToDouble();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildTopBar(),
      body:
          _refreshing == true
              ? _refreshState()
              : _noManagers
              ? _loadingState()
              : _managerStats == null
              ? _loadingState()
              : _buildMainSection(),
    );
  }

  Widget _loadingState() {
    if (_noManagers) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: const [
            Icon(
              Icons.manage_accounts_outlined,
              size: 64,
              color: Color(0xff9CA3AF),
            ),
            SizedBox(height: 16),
            Text(
              'No managers found',
              style: TextStyle(
                color: Color(0xff9CA3AF),
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
            ),
            SizedBox(height: 8),
            Text(
              'Create a manager first to view dashboard data',
              style: TextStyle(color: Color(0xff6B7280), fontSize: 14),
            ),
          ],
        ),
      );
    }

    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: const [
          CircularProgressIndicator(color: Color(0xffFFB400)),
          SizedBox(height: 16),
          Text(
            'Loading Dashboard...',
            style: TextStyle(color: Color(0xff9CA3AF)),
          ),
        ],
      ),
    );
  }

  Widget _refreshState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: const [
          CircularProgressIndicator(color: Color(0xffFFB400)),
          SizedBox(height: 16),
          Text(
            'Refreshing Dashboard...',
            style: TextStyle(color: Color(0xff9CA3AF)),
          ),
        ],
      ),
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
                        SizedBox(width: 20),
                        HoverableButton(
                          name: "Refresh",
                          icon: Icons.refresh_rounded,
                          onTap: () => refreshData(),
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
    List<Files> files = [];

    for (var manager in _managerStats!.managers!) {
      switch (_selectedFileSort) {
        case 'Recents':
          if (manager.recent != null) files.addAll(manager.recent!);
          break;
        case 'Largest':
          if (manager.largest != null) files.addAll(manager.largest!);
          break;
        case 'Oldest':
          if (manager.oldest != null) files.addAll(manager.oldest!);
          break;
      }
    }

    if (files.isEmpty) {
      return SizedBox(
        height: 300,
        child: Center(
          child: Text(
            'No files available',
            style: TextStyle(color: Color(0xff9CA3AF), fontSize: 16),
          ),
        ),
      );
    }

    return SizedBox(
      height: 300,
      child: SingleChildScrollView(
        child: Column(
          children:
              files.take(10).map((file) {
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
                          file.fileName ?? 'Unknown File',
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
                        flex: 2,
                        child: Text(
                          file.filePath ?? '',
                          style: TextStyle(
                            color: Color(0xff6b7280),
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                          textAlign: TextAlign.right,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      SizedBox(width: 10),
                    ],
                  ),
                );
              }).toList(),
        ),
      ),
    );
  }

  Widget _buildFileTypeStats() {
    List<int> totalCounts = List.filled(8, 0);

    for (var manager in _managerStats!.managers!) {
      if (manager.umbrellaCounts != null &&
          manager.umbrellaCounts!.length >= 8) {
        for (int i = 0; i < 8; i++) {
          totalCounts[i] += manager.umbrellaCounts![i];
        }
      }
    }

    int grandTotal = totalCounts.fold(0, (sum, count) => sum + count);

    if (grandTotal == 0) {
      return Center(
        child: Text(
          'No file type data available',
          style: TextStyle(color: Color(0xff9CA3AF), fontSize: 16),
        ),
      );
    }

    List<Map<String, dynamic>> fileTypes = [
      {'type': 'Documents', 'count': totalCounts[0], 'color': Colors.blue},
      {'type': 'Images', 'count': totalCounts[1], 'color': Colors.green},
      {'type': 'Music', 'count': totalCounts[2], 'color': Colors.red},
      {
        'type': 'Presentations',
        'count': totalCounts[3],
        'color': Colors.orange,
      },
      {'type': 'Videos', 'count': totalCounts[4], 'color': Colors.purple},
      {'type': 'Spreadsheets', 'count': totalCounts[5], 'color': Colors.pink},
      {'type': 'Archives', 'count': totalCounts[6], 'color': Colors.tealAccent},
      {'type': 'Other', 'count': totalCounts[7], 'color': Color(0xff9CA3AF)},
    ];

    return Column(
      children:
          fileTypes.where((fileType) => fileType['count'] > 0).map((fileType) {
            double percentage = fileType['count'] / grandTotal;

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
                        '${fileType['count']}',
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
                        value: _managerStats!.managers!.length.toString(),
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
                        value: _sumTotalFiles().toString(),
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
                        value: _sumTotalFolders().toString(),
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
                        value: '${_sumTotalSize()} ${_gig ? 'Gb' : 'Mb'}',
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
            ..._managerStats!.managers!.map((manager) {
              double size = manager.size!.toDouble();
              String sizeUnit = 'MB';
              if (size >= _gigInBytes) {
                size = size / _gigInBytes;
                sizeUnit = 'GB';
              } else {
                size = size / 1e+6;
              }

              return Column(
                children: [
                  ManagerItemWidget(
                    managerName: manager.managerName ?? "Unknown Manager",
                    managerFiles: manager.files?.toString() ?? "0",
                    managerFolders: manager.folders?.toString() ?? "0",
                    managerSize: "${size.toStringAsFixed(1)}$sizeUnit",
                  ),
                  SizedBox(height: 10),
                ],
              );
            }),

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
                          DropdownButton<String>(
                            value: _selectedFileSort,
                            dropdownColor: Color(0xff242424),
                            style: TextStyle(color: Colors.white),
                            underline: Container(
                              height: 1,
                              color: Color(0xff3D3D3D),
                            ),
                            items:
                                [
                                  "Recents",
                                  "Largest",
                                  "Oldest",
                                ].map<DropdownMenuItem<String>>((String value) {
                                  return DropdownMenuItem<String>(
                                    value: value,
                                    child: Text(
                                      value,
                                      style: TextStyle(color: Colors.white),
                                    ),
                                  );
                                }).toList(),
                            onChanged: (String? newValue) {
                              setState(() {
                                _selectedFileSort = newValue!;
                              });
                            },
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
