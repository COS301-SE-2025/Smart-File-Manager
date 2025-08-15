import 'package:flutter/material.dart';

class OverviewWidget extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;

  const OverviewWidget({
    required this.title,
    required this.value,
    required this.icon,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        // Scale factor based on widget width
        final scale = constraints.maxWidth / 200;

        return Container(
          padding: EdgeInsets.all(12 * scale.clamp(0.8, 1.5)),
          decoration: BoxDecoration(
            color: const Color(0xff242424),
            borderRadius: BorderRadius.circular(8 * scale.clamp(0.8, 1.5)),
            border: Border.all(color: const Color(0xff3D3D3D), width: 1),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      color: const Color(0xff9CA3AF),
                      fontSize: 10 * scale.clamp(0.8, 1.5),
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  SizedBox(height: 4 * scale.clamp(0.8, 1.5)),
                  Text(
                    value,
                    style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: 25 * scale.clamp(0.8, 1.5),
                    ),
                  ),
                ],
              ),
              Icon(
                icon,
                size: 24 * scale.clamp(0.8, 1.5),
                color: Color(0xffFFB400),
              ),
            ],
          ),
        );
      },
    );
  }
}
