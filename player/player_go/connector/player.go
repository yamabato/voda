package connector

import "sync"

import "voda/game"

/*
#Play
通信プロセスとプレイヤー関数とのやり取りを担う
・プレイヤー関数において、ロジック部分のみ記述すれば良い

*引数
player func(PlayerParam) PlayerRet: プレイヤー関数
port string                       : 通信に用いるポート番号
*/
func Play(player func(game.PlayerParam) game.PlayerRet, port uint) {
  // 通信は別プロセスで起動
  var wg *sync.WaitGroup = new(sync.WaitGroup)
  wg.Add(1)

  // 通信プロセスとやり取りするためのチャネル
  var msg_channel chan game.PlayerParam = make(chan game.PlayerParam);
  var ret_channel chan game.PlayerRet = make(chan game.PlayerRet);

  // ゲームとの通信のためのプロセスを起動
  go ConnectToGame(msg_channel, ret_channel, port, wg);

  var param game.PlayerParam;
  for {
    // 通信プロセスからメッセージを受け取り
    param = <-msg_channel;

    // 終了コマンドを受信した場合、即座に終了する
    if param.Command == "quit" { 
      break;
    }

    // プレイヤー関数の応答を通信プロセスへ送信
    var ret game.PlayerRet = player(param);
    ret_channel <- ret;
  }

  wg.Wait();
}

