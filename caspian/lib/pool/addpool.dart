import 'package:flutter/material.dart';

import '../common/main_navigation_bar.dart';

class AddPool extends StatelessWidget {
  const AddPool({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Add Pool"),
      ),
      body: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          children: [
            Row(children: [
              Expanded(
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.pushNamed(context, "/addPool/create");
                  },
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: const [
                      SizedBox(height: 20),
                      Icon(Icons.create),
                      SizedBox(height: 10),
                      Text("Create"),
                      SizedBox(height: 20),
                    ],
                  ),
                ),
              ),
            ]),
            const SizedBox(height: 10),
            Row(children: [
              Expanded(
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.pushNamed(context, "import");
                  },
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: const [
                      SizedBox(height: 20),
                      Icon(Icons.upload),
                      SizedBox(height: 10),
                      Text("Import"),
                      SizedBox(height: 20),
                    ],
                  ),
                ),
              ),
            ]),
          ],
        ),
      ),
      bottomNavigationBar: const MainNavigatorBar(),
    );
  }
}
