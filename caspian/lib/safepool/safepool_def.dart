import 'dart:convert';
import 'dart:typed_data';

class PoolConfig {
  String name = "";
  List<String> public_ = [];
  List<String> private_ = [];
  PoolConfig();

  PoolConfig.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        public_ = json['public'],
        private_ = json['private'];

  Map<String, dynamic> toJson() => {
        'name': name,
        'public': public_,
        'private': private_,
      };
}

class Key {
  String public_;
  String private_;

  Key()
      : private_ = "",
        public_ = "";

  Key.fromJson(Map<String, dynamic> json)
      : public_ = json['pu'],
        private_ = json['pr'];
}

class Identity {
  String nick = "";
  String email = "";
  Key signatureKey = Key();
  Key encryptionKey = Key();

  Identity.fromJson(Map<String, dynamic> json)
      : nick = json['n'],
        email = json['m'] ?? "",
        signatureKey = Key.fromJson(json["s"]),
        encryptionKey = Key.fromJson(json["e"]);
}

class Pool {
  String name = "";
  Identity self;
  List<String> apps = [];
  bool trusted = false;

  Pool.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        apps = json['apps'] ?? ['chat', 'library'],
        self = Identity.fromJson(json['self']),
        trusted = json['trusted'];

  Map<String, dynamic> toJson() => {
        'name': name,
        'apps': apps,
      };
}
// type Message struct {
// 	Id          uint64    `json:"id,string"`
// 	Author      string    `json:"author"`
// 	Time        time.Time `json:"time"`
// 	Content     string    `json:"content"`
// 	ContentType string    `json:"contentType"`
// 	Attachments [][]byte  `json:"attachments"`
// 	Signature   []byte    `json:"signature"`
// }

class Message {
  String id;
  String author;
  DateTime time;
  String contentType;
  String text;
  Uint8List binary;

  Message()
      : id = "0",
        author = "",
        time = DateTime.now(),
        contentType = "",
        text = "",
        binary = Uint8List(0);

  Message.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        author = json['author'],
        time = DateTime.parse(json['time']),
        contentType = json['contentType'],
        text = json['text'],
        binary = base64Decode(json['binary'] ?? "");
}
