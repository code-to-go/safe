import 'package:flutter/material.dart';

import '../common/main_navigation_bar.dart';

class Settings extends StatelessWidget {
  const Settings({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Settings"),
        automaticallyImplyLeading: false,
      ),
      // body: ListView(
      //   padding: const EdgeInsets.all(8),
      //   children: widgets,
      // ),
      bottomNavigationBar: const MainNavigatorBar(),
    );
  }
}
