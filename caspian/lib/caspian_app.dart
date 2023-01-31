import 'package:caspian/pool/addpool.dart';
import 'package:caspian/pool/addpool_create.dart';
import 'package:caspian/pool/poolist.dart';
import 'package:caspian/settings/settings.dart';
import 'package:flutter/material.dart';

class CaspianApp extends StatelessWidget {
  const CaspianApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Caspian',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      initialRoute: "/",
      routes: {
        "/": (context) => const PoolList(),
        "/addPool": (context) => const AddPool(),
        "/addPool/create": (context) => const CreatePool(),
        "/settings": (context) => const Settings(),
      },
    );
  }
}

// class CaspianPage extends StatefulWidget {
//   const CaspianPage({super.key, required this.title});

//   // This widget is the home page of your application. It is stateful, meaning
//   // that it has a State object (defined below) that contains fields that affect
//   // how it looks.

//   // This class is the configuration for the state. It holds the values (in this
//   // case the title) provided by the parent (in this case the App widget) and
//   // used by the build method of the State. Fields in a Widget subclass are
//   // always marked "final".

//   final String title;

//   @override
//   State<CaspianPage> createState() => _CaspianPageState();
// }

// class _CaspianPageState extends State<CaspianPage> {
//   int _selectedScreenIndex = 0;
//   late List _screens;

//   _CaspianPageState() {
//     _screens = [
//       PoolList(
//         selectScreen: _selectScreen,
//       ),
//       const Settings(),
//     ];
//   }

//   void _selectScreen(int index) {
//     setState(() {
//       _selectedScreenIndex = index;
//     });
//   }

//   @override
//   Widget build(BuildContext context) {
//     return _screens[_selectedScreenIndex];
//   }
// }
