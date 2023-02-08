import 'dart:math';
import 'dart:typed_data';

import 'package:caspian/safepool/safepool.dart' as sp;
import 'package:flutter/material.dart';
import 'package:chatview/chatview.dart' as cv;

import '../common/main_navigation_bar.dart';

class Chat extends StatefulWidget {
  const Chat({Key? key}) : super(key: key);

  @override
  State<Chat> createState() => _ChatState();
}

class _ChatState extends State<Chat> {
  int _idStart = 0;
  int _idEnd = 1 << 63 - 1;
  String _poolName = "";
  late cv.ChatUser _currentUser;
  late cv.ChatController _chatController;

  void loadMessages(String poolName) {
    var messages = sp.getMessages(poolName, _idStart, _idEnd, 100).map((e) {
      var id = 0;
      _idStart = _idStart == 0 ? id : min(_idStart, id);
      _idEnd = _idEnd == 1 << 63 ? id : max(_idEnd, id);
      switch (e.contentType) {
        case 'text/html':
          return cv.Message(
            id: e.id,
            message: e.text,
            createdAt: e.time,
            sendBy: e.author,
          );
        default:
          return cv.Message(
            id: e.id,
            message: "unsupported content ${e.contentType}",
            createdAt: e.time,
            sendBy: e.author,
            messageType: cv.MessageType.text,
          );
      }
    }).toList();
    _chatController = cv.ChatController(
      initialMessageList: messages,
      scrollController: ScrollController(),
    );
  }

  void setUser(String poolName) {
    var p = sp.getPool(poolName);
    _currentUser = cv.ChatUser(id: sp.getSelfId(), name: p.self.nick);
  }

  @override
  Widget build(BuildContext context) {
    final poolName = ModalRoute.of(context)!.settings.arguments as String;

    if (_poolName.isEmpty) {
      setUser(poolName);
      loadMessages(poolName);
      _poolName = poolName;
    }

    return Scaffold(
      appBar: AppBar(
        title: Text("Chat $poolName"),
      ),
      body: cv.ChatView(
          sender: _currentUser,
          chatController: _chatController,
          onSendTap: _onSendTap,
          receiver: _currentUser,
          sendMessageConfig: cv.SendMessageConfiguration(
              imagePickerIconsConfig: cv.ImagePickerIconsConfiguration(
            onImageSelected: (imagePath, error) {},
            cameraIconColor: Colors.black,
            galleryIconColor: Colors.black,
          ))),
      bottomNavigationBar: const MainNavigatorBar(),
    );
  }

  void _onSendTap(
    String message,
    cv.ReplyMessage replyMessage,
  ) {
    var id = sp.postMessage(_poolName, "text/html", message, Uint8List(0));

    _chatController.addMessage(
      cv.Message(
        id: id.toString(),
        createdAt: DateTime.now(),
        message: message,
        sendBy: _currentUser.id,
        replyMessage: replyMessage,
        messageType: cv.MessageType.text,
      ),
    );
  }
}
