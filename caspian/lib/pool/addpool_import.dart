import 'dart:io';

import 'package:caspian/common/main_navigation_bar.dart';
import 'package:caspian/safepool/safepool.dart';
import 'package:caspian/safepool/safepool_def.dart';
import 'package:flutter/material.dart';
import 'package:share_plus/share_plus.dart';
import 'package:flutter/services.dart';

class ImportPool extends StatefulWidget {
  const ImportPool({super.key});

  @override
  State<ImportPool> createState() => _ImportPoolState();
}

class _ImportPoolState extends State<ImportPool> {
  final _formKey = GlobalKey<FormState>();

  String _token = "";
  @override
  Widget build(BuildContext context) {
    var selfId = getSelfId();

    Widget shareButton;
    if (Platform.isAndroid || Platform.isIOS) {
      shareButton = ElevatedButton.icon(
          icon: const Icon(Icons.share),
          label: const Text("Share the id"),
          onPressed: () {
            final box = context.findRenderObject() as RenderBox?;
            Share.share(selfId,
                subject: "Can you add me to your pool?",
                sharePositionOrigin:
                    box!.localToGlobal(Offset.zero) & box.size);
          });
    } else {
      shareButton = ElevatedButton.icon(
          icon: const Icon(Icons.copy),
          label: const Text("Copy id to the clipboard"),
          onPressed: () {
            Clipboard.setData(ClipboardData(text: selfId)).then((_) {
              ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text("Id copied to clipboard")));
            });
          });
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text("Import Pool"),
      ),
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 16.0),
        child: Column(
          children: [
            const Text("The below text is your global id. Share the id with a "
                "peer in the pool to get an invite token."),
            const Text("Your peer will provide you an invite. Add it in the "
                "input below and click Import"),
            const SizedBox(height: 20),
            Text(selfId,
                style: TextStyle(
                  color: Colors.grey[800],
                  fontWeight: FontWeight.bold,
                )),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                shareButton,
                ElevatedButton(
                    child: const Text("Continue"),
                    onPressed: () {
                      Share.share(selfId,
                          subject: "Can you add me to your pool?");
                    }),
              ],
            ),
            const SizedBox(height: 20),
            Builder(
              builder: (context) => Form(
                key: _formKey,
                autovalidateMode: AutovalidateMode.onUserInteraction,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    TextFormField(
                      maxLines: 6,
                      decoration:
                          const InputDecoration(labelText: 'Invite Token'),
                      validator: (value) {
                        return null;
                      },
                      onChanged: (val) => setState(() => _token = val),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(
                          vertical: 16.0, horizontal: 16.0),
                      child: ElevatedButton(
                        onPressed: _token != ""
                            ? () {
                                try {
                                  var config = addPool(_token);
                                  ScaffoldMessenger.of(context)
                                      .showSnackBar(SnackBar(
                                          content: Text(
                                    "Congrats! You successfully joined ${config.name}",
                                  )));
                                } catch (e) {
                                  ScaffoldMessenger.of(context)
                                      .showSnackBar(SnackBar(
                                          backgroundColor: Colors.red,
                                          content: Text(
                                            "Invalid token: $e",
                                          )));
                                }
                              }
                            : null,
                        child: const Text('Create'),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
