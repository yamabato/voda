package connector

import "fmt"
import "net"
import "sync"
import "strconv"
import "strings"

import "voda/game"

/*
#ConnectToGame
ゲームに接続

*引数
msg_channel chan string         : 受信したメッセージ送信用のチャネル
ret_channel chan game.PlayerRet : 送信するメッセージ受信用のチャネル
port string                     : 通信に用いるポート番号
wg *sync.WaitGroup              : 並行処理管理のためのWaitGroup
*/
func ConnectToGame(msg_channel chan game.PlayerParam, ret_channel chan game.PlayerRet, port uint, wg *sync.WaitGroup) {
  // 処理終了時にWaitGroupのカウンタを1減ず
  defer wg.Done();

  // 指定ポート番号での接続
  conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
  if err != nil {
  }

  for {
    // ゲームから送信されたメッセージの受信
    buf := make([]byte, 1024)
    count, err := conn.Read(buf)
    if err != nil {
      wg.Done()
    }

    // 受信したメッセージをプレイヤーのプロセスへ送信
    msg_channel <- decodeMessage(string(buf[:count]));
    // quitの場合、処理を終了する
    if string(buf[:count]) == "quit" {
      break
    }

    // メッセージを構成し送信
    count, err = conn.Write([]byte(buildRetMsg(<-ret_channel)))
    if err != nil {
      break;
    }
  }
}

/*
#DecodeMessage
受け取ったメッセージからPlayerParamを構成

*引数
msg string: ゲームから受け取ったメッセージ

*返り値
PlayerParam: メッセージから生成したPlayerParam構造体
*/
func decodeMessage(msg string) game.PlayerParam {
  // メッセージを分割
  var msg_words []string = strings.Split(msg, " ");

  // メッセージの最初の単語をコマンド、残りを引数
  var command string;
  command = msg_words[0];
  var param_words []string = msg_words[1:];

  var param game.PlayerParam;
  param.Command = command;

  switch command {
  case "name": // nameコマンドは引数なし
    break;
  case "start":
    // 先後を受け取る
    var side string = param_words[0];
    if (side == "black") {
      param.Turn = true;
    } else {
      param.Turn = false;
    }
  case "go":
    return buildGoPlayerParam(param_words);
  case "end":
    // 終了コードを設定
    var result string = param_words[0];
    switch result {
    case "win": param.Result = 0;
    case "lose": param.Result = 1;
    case "draw": param.Result = 2;
    default: param.Result = 255;
    }
  }

  return param;
}

/*
#buildRetMsg
PlayerRetからゲームに送信するメッセージを生成

*引数
ret PlayerRet: メッセージ生成元

*返り値
string: 生成したメッセージ
*/
func buildRetMsg(ret game.PlayerRet) string {
  var command string = ret.Command;

  switch command {
  case "setname":
    return fmt.Sprintf("setname %s", ret.Name);
  case "ready":
    return "ready";
  case "move":
    return fmt.Sprintf("move %d", ret.Move);
  case "bye":
    return fmt.Sprintf("bye");
  default:
  }
  return ""
}

/*
#buildGoPlayerParam
goコマンドに対するPlayerParamを構成

*引数
param_words []string: goコマンドの引数

*返り値
PlayerParam: goコマンドに対するPlayerParam構造体
*/
func buildGoPlayerParam(param_words []string) game.PlayerParam {
  var param game.PlayerParam;
  param.Command = "go";

  var stones uint64;
  var opp_stones uint64;
  var valid_moves []uint8;
  var moves []uint8;

  // 石の配置を取得
  stones_int, err := strconv.Atoi(param_words[0]);
  if (err==nil) { stones = uint64(stones_int); }
  opp_stones_int, err := strconv.Atoi(param_words[1]);
  if (err==nil) { opp_stones = uint64(opp_stones_int); }

  // 着手可能な手の一覧を取得
  for _, m:=range strings.Split(param_words[2], "") {
    move_int, err := strconv.Atoi(m);
    if (err==nil){
      valid_moves = append(valid_moves, uint8(move_int));
    }
  }

  // 移動履歴を取得
  for _, m:=range strings.Split(param_words[3], "") {
    move_int, err := strconv.Atoi(m);
    if (err == nil) {
      moves = append(moves, uint8(move_int));
    }
  }

  param.Stones = stones;
  param.OppStones = opp_stones;
  param.ValidMoves = valid_moves;
  param.Moves = moves;

  return param;
}
