import 'package:flutter/material.dart';
import 'package:toggle_switch/toggle_switch.dart';

class ManagerPage extends StatelessWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: PreferredSize(
        preferredSize: Size.fromHeight(50),
        child: AppBar(
          backgroundColor: const Color(0xff2E2E2E),
          automaticallyImplyLeading: false,
          shape: Border(bottom: BorderSide(color: Color(0xff3D3D3D), width: 1)),
          flexibleSpace: Center(
            child: SafeArea(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: Row(
                  children: [
                    Expanded(
                      child: Row(
                        children: [
                          Text(
                            '$name Manager',
                            style: const TextStyle(
                              fontWeight: FontWeight.bold,
                              fontSize: 20,
                              color: Colors.white,
                            ),
                          ),
                          const SizedBox(width: 12),
                          Container(
                            width: 100,
                            height: 20,
                            padding: const EdgeInsets.symmetric(
                              vertical: 2,
                              horizontal: 8,
                            ),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(10),
                              color: const Color(0xff094E3A),
                            ),
                            child: const Center(
                              child: Text(
                                '75% organized',
                                style: TextStyle(
                                  fontSize: 10,
                                  color: Color(0xff6EE79B),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    MouseRegion(
                      cursor: SystemMouseCursors.click,
                      child: ToggleSwitch(
                        initialLabelIndex: 0,
                        minWidth: 90,
                        minHeight: 30,
                        fontSize: 12,
                        inactiveBgColor: Color(0xff242424),
                        inactiveFgColor: Color(0xff9CA3AF),
                        borderColor: [Color(0xff3D3D3D)],
                        cornerRadius: 4,
                        borderWidth: 1,
                        totalSwitches: 2,
                        labels: const ['Folder View', 'Graph View'],
                        onToggle: (index) {},
                      ),
                    ),
                    const SizedBox(width: 8),
                    SizedBox(
                      height: 33,
                      child: TextButton(
                        onPressed: () {},
                        style: TextButton.styleFrom(
                          backgroundColor: const Color(0xff242424),
                          side: const BorderSide(
                            color: Color(0xff3D3D3D),
                            width: 1,
                          ),

                          foregroundColor: const Color(0xff9CA3AF),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(4),
                          ),
                        ),

                        child: const Text(
                          "Sort",
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.normal,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
      body: Scaffold(
        appBar: PreferredSize(
          preferredSize: const Size.fromHeight(50),
          child: AppBar(
            backgroundColor: const Color(0xff2E2E2E),
            automaticallyImplyLeading: false,
            shape: const Border(
              bottom: BorderSide(color: Color(0xff3D3D3D), width: 1),
            ),
            titleSpacing: 0,
            title: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Row(
                children: [
                  // Search bar
                  Container(
                    width: 200,
                    height: 33,
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    decoration: BoxDecoration(
                      color: const Color(0xff242424),
                      borderRadius: BorderRadius.circular(4),
                      border: Border.all(
                        color: const Color(0xff3D3D3D),
                        width: 1,
                      ),
                    ),
                    child: Center(
                      child: const TextField(
                        cursorColor: Color(0xffFFB400),
                        style: TextStyle(
                          fontSize: 12,
                          color: Color(0xff9CA3AF),
                        ),
                        decoration: InputDecoration(
                          hintText: 'Search...',
                          hintStyle: TextStyle(
                            color: Color(0xff9CA3AF),
                            fontSize: 12,
                          ),
                          border: InputBorder.none,
                          isCollapsed: true,
                        ),
                      ),
                    ),
                  ),

                  const SizedBox(width: 12),

                  Container(
                    height: 33,
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    decoration: BoxDecoration(
                      color: const Color(0xff242424),
                      borderRadius: BorderRadius.circular(4),
                      border: Border.all(
                        color: const Color(0xff3D3D3D),
                        width: 1,
                      ),
                    ),
                    child: DropdownButtonHideUnderline(
                      child: DropdownButton<String>(
                        value: 'Name',
                        dropdownColor: const Color(0xff2E2E2E),
                        iconEnabledColor: const Color(0xff9CA3AF),
                        style: const TextStyle(
                          fontSize: 12,
                          color: Color(0xff9CA3AF),
                        ),
                        items: const [
                          DropdownMenuItem(
                            value: 'Name',
                            child: Text('Sort by Name'),
                          ),

                          DropdownMenuItem(
                            value: 'Size',
                            child: Text('Sort by Size'),
                          ),
                        ],
                        onChanged: (value) {},
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        body: const Placeholder(),
      ),
    );
  }
}
