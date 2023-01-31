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
