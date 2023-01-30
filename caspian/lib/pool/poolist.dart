import 'package:caspian/safepool/safepool.dart';
import 'package:flutter/material.dart';

import '../common/main_navigation_bar.dart';

class PoolList extends StatelessWidget {
  const PoolList({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    var poolList = getPoolList();
    var widgets = poolList
        .map(
          (e) => Card(
            child: ListTile(
              title: Text(e),
              leading: const Icon(Icons.waves),
            ),
          ),
        )
        .toList();

    return Scaffold(
      appBar: AppBar(
        title: const Text("Pools"),
        actions: [
          Padding(
            padding: EdgeInsets.only(right: 20.0),
            child: GestureDetector(
                onTap: () {},
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
