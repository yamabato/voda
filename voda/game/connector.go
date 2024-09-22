package game

import "net"
import "fmt"
import "sync"
import "strconv"
import "strings"

/*
#コマンド
name
  *プレイヤー名を要求
  Param: -
  Msg: name

start
  *ゲーム開始
  Param:
    Turn bool: 先後
      先手: true, 後手: false
  Msg: start (black, white)

go
  *次の手を要求
  Param:
    Stones uint64      : 自分の石の配置
    OppStones uint64   : 相手の石の配置
    Moves []uint8      : 操作履歴
    ValidMoves []uint8 : 合法手のリスト
  Msg: go (Stones) (OppStones) (Moves) (ValidMoves)

end
  *ゲーム終了
  Param:
    Result string: 結果(win, lose, draw)
  Msg: end (win, lose, draw)

quit
  *プレイヤー終了
  Param: -
  Msg: -
*/

/*
#connectToPlayer
プレイヤーとソケット通信を行う

*引数
param_channel chan PlayerParam: プレイヤーに送る情報を受け取るチャネル
ret_channel chan PlayerRet    : プレイヤーから受け受け取った情報を送信するチャネル
port uint                     : 通信先のポート番号
wg *sync.WaitGroup            : 通信終了のためWaitGroupが必要
*/
func connectToPlayer(
  param_channel chan PlayerParam,
  ret_channel chan PlayerRet,
  port uint, wg *sync.WaitGroup,
) {
  // 通信を開始
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port));
  if err != nil {
  }

  // 通信に用いるTCPConn構造体を取得
  conn, err := ln.Accept();
  // コネクションは終了時に必ず切断
  defer conn.Close();
  // WaitGroupのカウンタを1減ずる
  defer wg.Done();
  if err != nil {
    return;
  }

  for {
    // プレイヤーに送信するパラメータを受け取る
    param, ok := <- param_channel;
    if (!ok) {
      break
    }
    // プレイヤーに送信するメッセージを組み立て
    ok, msg := build_message(param);
    fmt.Println(fmt.Sprintf("msg(%d)", port), msg)
    if (!ok) {
    }

    // プレイヤーにメッセージを送信
    _, err = conn.Write([]byte(msg))
    if err != nil {
    }

    // 送信したメッセージがquitだった場合、通信を終了
    if param.Command == "quit" {
      break;
    }

    // プレイヤーからメッセージを受信
    var buf []byte = make([]byte, 1024);
    n, err := conn.Read(buf);
    if err != nil {
    }
    fmt.Println(fmt.Sprintf("rsv(%d)", port), string(buf[:n]))

    // プレイヤーからの応答をPlayerRet構造体に変換し、ゲームに通知
    var ret PlayerRet = buildPlayerRet(string(buf[:n]));
    ret_channel <- ret;
  }
}

/*
#sendMessage
プレイヤーにメッセージを送信する

*引数
param PlayerParam             : 送信するパラメータ
param_channel chan PlayerParam: パラメータを送信するチャネル
ret_channel chan PlayerRet    : プレイヤーの応答を受け取るチャネル
*/
func sendMessage(param PlayerParam, param_channel chan PlayerParam, ret_channel chan PlayerRet) PlayerRet{
  // パラメータをconnectToPlayerのプロセスへ送信
  param_channel <- param;

  // connectToPlayerのプロセスからプレイヤーの応答を受け取り
  var ret PlayerRet = <-ret_channel;
  return ret;
}

/*
#build_message
プレイヤーに送信するメッセージ文字列を構成

*引数
param PlayerParam: コマンドのパラメータ

*返り値
bool  : エラーの有無
string: メッセージ文字列
*/
func build_message(param PlayerParam) (bool, string) {
  // メッセージ
  var msg string;

  var command string = param.Command;
  switch command {
  case "name": // 名称の問い合せ
    msg = "name";
  case "start": // ゲーム開始の通知
    // 引数としてターンを通知
    msg = fmt.Sprintf("start %s", make_turn_str(param));
  case "go": // 次の手の要求
    msg = fmt.Sprintf("go %s", build_go_msg_param(param));
  case "end": // ゲーム終了の通知
    msg = fmt.Sprintf("end %s", makeResultStr(param));
  case "quit": // プレイヤーの終了
    msg = fmt.Sprintf("quit");
  default:
    return false, "";
  }
  return true, msg;
}

/*
#make_turn_str
手番を示す文字列を生成

*引数
param PlayerParam: パラメータ

*返り値
string: 手番の文字列(black, white)
*/
func make_turn_str(param PlayerParam) string {
  // trueが先手、falseが後手
  if (param.Turn) {
    return "black";
  } else {
    return "white";
  }
}

/*
#build_go_msg_param
goコマンドのパラメータを組み立て

*引数
param PlayerParam: パラメータ

*返り値
string: パラメータ文字列

stones_str string     : 石配置(10進表記)
opp_stones_str string : 相手方石配置(10進表記)
moves_str string      : 操作履歴
valid_moves string    : 有効手
*/
func build_go_msg_param(param PlayerParam) string {
  // 石の配置を10進表記で文字列化
  var stones_str string = fmt.Sprintf("%d", param.Stones);
  var opp_stones_str string = fmt.Sprintf("%d", param.OppStones);

  // 操作履歴と合法手を文字列化
  /*
  1. fmt.Sprintでデフォルト形式によって配列を文字列化(ex. [1 0 0])
  2. strings.Replaceで空白を削除(ex. [100])
  3. strings.Trimで両端の[]を削除(ex. 100)
  */
  var moves_str string = strings.Trim(strings.Replace(fmt.Sprint(param.Moves), " ", "", -1), "[]");
  var valid_moves_str string = strings.Trim(strings.Replace(fmt.Sprint(param.ValidMoves), " ", "", -1), "[]");

  // パラメータを連結する
  return fmt.Sprintf("%s %s %s %s", stones_str, opp_stones_str, valid_moves_str, moves_str);
}

/*
#makeResultStr
結果を通知する文字列を生成

*引数
param PlayerParam

*返り値
string: 結果
*/
func makeResultStr(param PlayerParam) string {
  switch param.Result {
  case 0: return "win";
  case 1: return "lose";
  case 2: return "draw";
  default: return "";
  }
}

/*
#buildPlayerRet
プレイヤーからの応答を文字列からPlayerRet構造体へ変換

*引数
ret_msg string: プレイヤーから受信した文字列

*返り値
PlayerRet: プレイヤーからのメッセージより生成したPlayerRet構造体
*/
func buildPlayerRet(ret_msg string) (PlayerRet){
  // プレイヤーからの応答を分解
  var ret_words []string = strings.Split(ret_msg, " ");

  // 応答の最初の単語はコマンド、残りはその引数として解釈
  var command string = ret_words[0];
  var ret_param_list []string = ret_words[1:];

  var ret PlayerRet;
  ret.Command = command;

  switch command {
  case "setname": // 名称の設定
    ret.Name = ret_param_list[0];
  case "ready": // 準備完了
    ret.Ready = true;
  case "move": // 次の手
    move_int, err := strconv.Atoi(ret_param_list[0]);
    if (err==nil) {
      ret.Move = uint8(move_int);
    }
  case "bye": // 終了
  default:
    fmt.Println(fmt.Sprintf("Unknown Commnad `%s`", command));
  }

  return ret;
}
