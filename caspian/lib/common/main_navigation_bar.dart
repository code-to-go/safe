import 'package:flutter/material.dart';

class MainNavigatorBar extends StatelessWidget {
  const MainNavigatorBar({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
      currentIndex: 0,
      onTap: (idx) {
        switch (idx) {
          case 0:
            Navigator.pop(context);
            break;
          case 1:
            Navigator.pushNamed(context, '/settings');
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Pools'),
        BottomNavigationBarItem(icon: Icon(Icons.settings), label: "Settings")
      ],
    );
  }
}
