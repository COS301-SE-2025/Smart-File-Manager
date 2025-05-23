import 'package:flutter/material.dart';
import 'package:toggle_switch/toggle_switch.dart';

class ManagerPage extends StatelessWidget {
  final String name;
  const ManagerPage({required this.name, super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: const Color(0xff2E2E2E),
        title: Text('$name Manager'),
        centerTitle: false,
        titleTextStyle: const TextStyle(
          fontWeight: FontWeight.bold,
          fontSize: 20,
        ),
        shape: Border(bottom: BorderSide(color: Color(0xff3D3D3D), width: 1)),
        actions: [
          Row(
            children: [
              Expanded(
                flex: 3,
                child: Container(
                  width: 100,
                  height: 20,
                  padding: EdgeInsets.symmetric(vertical: 2, horizontal: 8),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color: Color(0xff094E3A),
                  ),
                  child: Center(
                    child: Text(
                      '75% organized',
                      style: TextStyle(fontSize: 10, color: Color(0xff6EE79B)),
                    ),
                  ),
                ),
              ),
              Expanded(
                flex: 1,
                child: ToggleSwitch(
                  initialLabelIndex: 0,
                  totalSwitches: 2,
                  labels: ['Folder View', 'Graph View'],
                  onToggle: (index) {
                    print('switched to: $index');
                  },
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
