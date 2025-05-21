import 'package:flutter/material.dart';

void main() {
  runApp(const MaterialApp(home: SafeArea(child: AppBody())));
}

class AppBody extends StatefulWidget {
  const AppBody({super.key});

  @override
  State<AppBody> createState() => _AppBodyState();
}

class _AppBodyState extends State<AppBody> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Color(0xff1E1E1E),
      appBar: AppBar(
        leading: Padding(
          padding: const EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Image.asset("images/logo.png"),
        ),
        backgroundColor: Color(0xff2E2E2E),
        title: Text("SMART FILE MANAGER"),
        titleTextStyle: TextStyle(
          color: Color(0xffFFB400),
          fontWeight: FontWeight.bold,
          fontSize: 20,
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 10, 0),
            child: FilledButton.icon(
              onPressed: () {},
              label: Text("Login"),
              style: FilledButton.styleFrom(backgroundColor: Color(0xff242424)),
              icon: Icon(Icons.account_circle),
            ),
          ),
        ],
      ),
    );
  }
}
