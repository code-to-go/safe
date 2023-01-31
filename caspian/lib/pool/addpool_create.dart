import 'package:caspian/safepool/safepool.dart';
import 'package:caspian/safepool/safepool_def.dart';
import 'package:flutter/material.dart';
import 'package:basic_utils/basic_utils.dart';

class CreatePool extends StatefulWidget {
  const CreatePool({super.key});

  @override
  State<CreatePool> createState() => _CreatePoolState();
}

const urlHint = "Enter a supported URL and click +";
const validSchemas = ["s3", "sftp", "file"];
const availableApps = ["chat", "library", "gallery"];

class _CreatePoolState extends State<CreatePool> {
  final _formKey = GlobalKey<FormState>();

  final _config = PoolConfig();
  final _appsMap = {for (var e in availableApps) e: true};
  String? _publicExchange;
  String? _privateExchange;

  String? _validateExchangeUrl(String? text) {
    if (text == null || text == "") {
      return null;
    }
    var uri = Uri.parse(text);
    return uri.host != "" && validSchemas.contains(uri.scheme)
        ? null
        : "Enter a valid url and click +";
  }

  bool _validExchangeUrl(String? text) {
    return (text != null && text != "" && _validateExchangeUrl(text) == null);
  }

  bool _validConfig() {
    var apps = [];
    _appsMap.forEach((key, value) {
      if (value) apps.add(key);
    });

    return _config.name != "" && _config.public_.isNotEmpty && apps.isNotEmpty;
  }

  final _publicController = TextEditingController();
  final _privateController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    var appsList = availableApps
        .map((e) => CheckboxListTile(
            title: Text(StringUtils.capitalize(e)),
            value: _appsMap[e],
            onChanged: (a) {
              setState(() {
                _appsMap[e] = !(_appsMap[e] ?? false);
              });
            }))
        .toList();

    return Scaffold(
      appBar: AppBar(
        title: const Text("Create Pool"),
      ),
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 16.0),
        child: Builder(
          builder: (context) => Form(
            key: _formKey,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const Text(
                    "Enter a name and at least a public Exchange, i.e. a URI to sftp or s3"),
                TextFormField(
                  decoration: const InputDecoration(labelText: 'Name'),
                  validator: (value) {
                    return null;
                  },
                  onChanged: (val) => setState(() => _config.name = val),
                ),
                TextFormField(
                  maxLines: 3,
                  controller: _publicController,
                  decoration: InputDecoration(
                    labelText: 'Public Exchanges',
                    suffixIcon: IconButton(
                      onPressed: _validExchangeUrl(_publicExchange)
                          ? () {
                              setState(() {
                                _config.public_.add(_publicExchange ?? "");
                                _publicController.clear();
                                _publicExchange = null;
                              });
                            }
                          : null,
                      icon: const Icon(Icons.add),
                    ),
                  ),
                  validator: _validateExchangeUrl,
                  onChanged: (val) => setState(() => _publicExchange = val),
                ),
                ListView.builder(
                  scrollDirection: Axis.vertical,
                  shrinkWrap: true,
                  itemCount: _config.public_.length,
                  itemBuilder: (context, index) => ListTile(
                    trailing: IconButton(
                        icon: const Icon(Icons.delete),
                        onPressed: () {
                          setState(() {
                            _config.public_.removeAt(index);
                          });
                        }),
                    title: Text(_config.public_[index]),
                  ),
                ),
                TextFormField(
                  maxLines: 3,
                  controller: _privateController,
                  decoration: InputDecoration(
                    labelText: 'Private Exchanges',
                    suffixIcon: IconButton(
                      onPressed: _validExchangeUrl(_privateExchange)
                          ? () {
                              setState(() {
                                _config.private_.add(_privateExchange ?? "");
                                _publicController.clear();
                                _privateExchange = null;
                              });
                            }
                          : null,
                      icon: const Icon(Icons.add),
                    ),
                  ),
                  validator: _validateExchangeUrl,
                  onChanged: (val) => setState(() => _publicExchange = val),
                ),
                ListView.builder(
                  scrollDirection: Axis.vertical,
                  shrinkWrap: true,
                  itemCount: _config.private_.length,
                  itemBuilder: (context, index) => ListTile(
                    trailing: IconButton(
                        icon: const Icon(Icons.delete),
                        onPressed: () {
                          setState(() {
                            _config.private_.removeAt(index);
                          });
                        }),
                    title: Text(_config.private_[index]),
                  ),
                ),
                Column(children: appsList),
                Container(
                  padding: const EdgeInsets.symmetric(
                      vertical: 16.0, horizontal: 16.0),
                  child: ElevatedButton(
                    onPressed: _validConfig()
                        ? () {
                            List<String> apps = [];
                            _appsMap.forEach((key, value) {
                              if (value) apps.add(key);
                            });
                            createPool(_config, apps);
                          }
                        : null,
                    child: const Text('Create'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
