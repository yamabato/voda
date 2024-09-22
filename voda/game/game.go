package game

import "fmt"
import "sync"

import "voda/board"

/*
#game
コネクトフォーのゲームを管理
・盤面操作にvoda/boardを利用
・盤面、履歴をGameDataに記録
・プレイヤー関数のエラーやタイムアウトは考慮しない
*/

/*
#プレイヤー
・プレイヤーは一つの関数の形で指定
・プレイヤーへの要求はPlayerParamを通して伝達
・プレイヤーからの応答はPlayerRetを通して受付
*/

/*
#Start
ゲームを初期化し、開始する

*引数
black_port string: 先手のポート
white_port string: 後手のポート

show_board bool : 盤面の出力
show_result bool: 結果の出力

*返り値
uint8: 結果
  0: 先手勝
  1: 後手勝
  2: 引き分け
  255: 異常終了
*/
func (g *Game) StartCLI(
  black_port uint, white_port uint,
  show_board bool, show_result bool,
) uint8 {
  // ゲームの初期化
  g.initializeGame(black_port, white_port);

  // プレイヤーに開始コマンドを送信
  var start bool = g.sendStartCommand();
  // いずれかが開始に失敗した場合、異常終了
  if (!start) {
    g.WaitGroup.Done();
    return 255;
  }

  // ゲームを進める
  var valid bool; // 手が合法か
  var stones uint64; // 打った側の石
  var result uint8; // 結果
  for i:=0; i<42; i++ { // 最大42手(7*6)
    // 次の手に進む
    // 正当な手であった(valid)かを返す
    valid, _ = g.inquireNextMove();

    // 盤面表示
    if (show_board) {
      fmt.Println(g.Board.Counter); // 手数
      board.PrintBoard(g.Board.BlackStones, g.Board.WhiteStones);
      fmt.Println();
    }

    // 非合法手が選択された場合
    if (!valid) {
      // 相手の勝ちとしてゲームを終了
      g.endGame(g.Board.Counter%2);
      break;
    }

    if (g.Board.Counter % 2 == 0) {
      stones = g.Board.WhiteStones;
    } else {
      stones = g.Board.BlackStones;
    }

    // いずれかが勝利した場合
    if (board.CheckAlignment(stones)) {
      g.endGame((g.Board.Counter+1)%2);
      result = (g.Board.Counter+1) % 2;
      break;
    }
  }

  if (g.Board.Counter == 42) {
    // すべて埋まった場合は引き分け
    g.endGame(2);
    result = 2;
  }

  // 結果表示
  if (show_result) {
    switch result {
    case 0:
      fmt.Println("Win: Black, Lose: White");
    case 1:
      fmt.Println("Win: White, Lose: Black");
    case 2:
      fmt.Println("Draw");
    }
  }

  // プレイヤーを終了させる
  g.quitPlayer();

  // 並行処理のWaitGroupを待つ
  g.WaitGroup.Wait();

  return result;
}

/*
#InitializeGame
ゲームを初期化する
・盤面の初期化
・ゲーム情報の初期化
・プレイヤーとの接続
・プレイヤー情報の取得

*引数
black_port uint: 先手のポート番号
white_port uint: 後手のポート番号
*/
func (g *Game) initializeGame(black_port uint, white_port uint) {
  // 盤面を初期化
  g.initializeBoard();

  // プレイヤーと接続
  g.initializeConnection(black_port, white_port);

  // プレイヤーの名前を設定
  g.inquirePlayerName();
}

/*
#initializeBoard
盤面の初期化
*/
func (g *Game) initializeBoard() {
  // 盤面の初期化
  (*g).Board = BoardData {
    0,          // BlackStones
    0,          // WhiteStones
    []uint8{},  // Moves
    0,          // Counter
  };
}

/*
#initializeConnection
プレイヤーと接続を確立
・通信は各々のプレイヤーに対し、並行で実行する

*引数
black_port uint:  先手のポート番号
white_port uint: 後手のポート番号
*/
func (g *Game) initializeConnection(black_port uint, white_port uint) {
  // プレイヤーとの通信は並行処理で実行するため、それらのWaitGroupを準備
  var wg *sync.WaitGroup = new(sync.WaitGroup);
  // 両プレイヤーに対し通信を行うため、プロセス2つ分待つ
  wg.Add(2);

  // WaitGroup, portを保持
  (*g).WaitGroup = wg;
  (*g).BlackPort = black_port;
  (*g).WhitePort = white_port;

  // 通信を行うためのチャネルを用意
  // param_channel: プレイヤーへの送信
  // ret_channel  : プレイヤーからの受信

  if (black_port != 0) {
    var black_param_channel chan PlayerParam = make(chan PlayerParam);
    var black_ret_channel chan PlayerRet = make(chan PlayerRet);
    (*g).BlackParamChannel = black_param_channel;
    (*g).BlackRetChannel = black_ret_channel;
    // 通信を確立する
    go connectToPlayer(black_param_channel, black_ret_channel, black_port, wg);
  }

  if (white_port != 0) {
    var white_param_channel chan PlayerParam = make(chan PlayerParam);
    var white_ret_channel chan PlayerRet = make(chan PlayerRet);
    (*g).WhiteParamChannel = white_param_channel;
    (*g).WhiteRetChannel = white_ret_channel;
    // 通信を確立する
    go connectToPlayer(white_param_channel, white_ret_channel, white_port, wg);
  }
}

/*
#inquirePlayerName
プレイヤーの名称を取得

Command: name
Param: -

Ret:
  Name string: プレイヤー名
*/
func (g *Game) inquirePlayerName() {
  var black_name string;
  var white_name string;

  // 名称の要求を送信
  if (g.BlackPort != 0) {
    black_name = sendMessage(PlayerParam{ Command: "name" }, g.BlackParamChannel, g.BlackRetChannel).Name;
  }

  if (g.WhitePort != 0) {
    white_name = sendMessage(PlayerParam{ Command: "name" }, g.WhiteParamChannel, g.WhiteRetChannel).Name;
  }

  // プレイヤー名を設定
  (*g).BlackName = black_name;
  (*g).WhiteName = white_name;
}

/*
#SendStartCommand
プレイヤーに開始コマンドを送信

*返り値
bool: ゲーム開始の可否

Command: start
Param:
  Turn bool: 先後
    先手: true, 後手: false

Ret:
  Ok bool: 開始の確認
*/
func (g Game) sendStartCommand() bool {
  // 盤面のリセット
  g.initializeBoard();

  // 開始を通知する
  // 準備ができた場合、Readyがtrue
  var black_ok bool = true;
  var white_ok bool = true;

  if (g.BlackPort != 0) {
    black_ok = sendMessage(PlayerParam{ Command: "start", Turn: true }, g.BlackParamChannel, g.BlackRetChannel).Ready;
  }
  if (g.WhitePort != 0) {
    white_ok = sendMessage(PlayerParam{ Command: "start", Turn: false }, g.WhiteParamChannel, g.WhiteRetChannel).Ready;
  }

  // 両方準備できていた場合、ゲームを開始する
  return black_ok && white_ok;
}

/*
#inquireNextMove
次の手を取得、適用する

*返り値
bool: 返却された手が合法手であるか
uint8: 返却された手

Command: go
Param:
  Stones uint64      : 自分の石の配置
  OppStones uint64   : 相手の石の配置
  Moves []uint8      : 操作履歴
  ValidMoves []uint8 : 合法手のリスト

Ret:
  Move uint8: 操作
*/
func (g *Game) inquireNextMove() (bool, uint8) {
  var black bool = g.Board.Counter%2==0; // 先後

  var stones uint64;
  var opp_stones uint64;
  var param_channel chan PlayerParam;
  var ret_channel chan PlayerRet;

  // 先後に応じ、プレイヤー関数、石の配置を設定
  if (black) {
    stones = g.Board.BlackStones;
    opp_stones = g.Board.WhiteStones;
    param_channel = g.BlackParamChannel;
    ret_channel = g.BlackRetChannel;
  } else {
    stones = g.Board.WhiteStones;
    opp_stones = g.Board.BlackStones;
    param_channel = g.WhiteParamChannel;
    ret_channel = g.WhiteRetChannel;
  }

  // 次の操作を要求する
  var next_move uint8 = sendMessage(PlayerParam { 
    Command: "go",
    Stones: stones,
    OppStones: opp_stones,
    Moves: g.Board.Moves,
    ValidMoves: board.GenValidMoves(stones, opp_stones),
  }, param_channel, ret_channel).Move;

  return g.dropStone(next_move);
}

func (g *Game) dropStone(move uint8) (bool, uint8) {
  var black bool = g.Board.Counter%2==0; // 先後

  // 手数のカウンタを進める
  (*g).Board.Counter++;
  (*g).Board.Moves = append(g.Board.Moves, move);

  // 非合法手が返された
  if (!board.CanMove(g.Board.BlackStones, g.Board.WhiteStones, move)) { return false, move; }

  // 盤面の情報を更新
  if (black) {
    (*g).Board.BlackStones = board.MakeMove(
      g.Board.BlackStones, g.Board.WhiteStones, move,
    );
  } else {
    (*g).Board.WhiteStones = board.MakeMove(
      g.Board.WhiteStones, g.Board.BlackStones, move,
    );
  }

  return true, move;
}

/*
#endGame
ゲームを終了させる

*引数
result uint8: 結果
  0: 先手勝
  1: 後手勝
  2: 引き分け

Command: end
Param:
  Result uint8: 結果(0:win, 1:lose, 2:draw)

Ret: -
*/
func (g *Game) endGame(result uint8) {
  var black_result uint8 = 2; // draw
  var white_result uint8 = 2; // draw

  if (result == 0) {
    // 先手勝利の場合
    black_result = 0; // win
    white_result = 1; // lose
  } else {
    // 後手勝利の場合
    black_result = 1; // lose
    white_result = 0; // win
  }

  // 終了メッセージを送信
  if (g.BlackPort != 0) {
    sendMessage(PlayerParam{ Command: "end", Result: black_result }, g.BlackParamChannel, g.BlackRetChannel);
  }
  if (g.WhitePort != 0) {
    sendMessage(PlayerParam{ Command: "end", Result: white_result }, g.WhiteParamChannel, g.WhiteRetChannel);
  }
}

/*
#quitPlayer
プレイヤーを終了させる
*/
func (g *Game) quitPlayer() {
  // 終了命令を送信
  if (g.BlackPort != 0) {
    go sendMessage(PlayerParam{ Command: "quit" }, g.BlackParamChannel, g.BlackRetChannel);
  }
  if (g.WhitePort != 0) {
    go sendMessage(PlayerParam{ Command: "quit" }, g.WhiteParamChannel, g.WhiteRetChannel);
  }
}

