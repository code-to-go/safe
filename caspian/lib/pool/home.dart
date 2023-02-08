import 'package:caspian/safepool/safepool.dart';
import 'package:flutter/material.dart';

import '../common/main_navigation_bar.dart';

class Home extends StatelessWidget {
  const Home({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    var poolList = getPoolList();
    var widgets = poolList
        .map(
          (e) => Card(
            child: ListTile(
              title: Text(e),
              leading: const Icon(Icons.waves),
              onTap: () {
                Navigator.pushNamed(context, "/pool", arguments: e);
              },
            ),
          ),
        )
        .toList();

    return Scaffold(
      appBar: AppBar(
        title: const Text("Your pools"),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 20.0),
            child: GestureDetector(
                onTap: () {
                  Navigator.pushNamed(context, "/addPool");
                },
                child: const Icon(
                  Icons.add,
                )),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(8),
        children: widgets,
      ),
      bottomNavigationBar: const MainNavigatorBar(),
    );
  }
}
