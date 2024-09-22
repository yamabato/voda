package game

import "fmt"
import "math"
import "net/http"
import "encoding/json"

import "voda/board"

/*
#StartBrowser
http通信を介してクライアントと通信し、ゲームを開始する
*/
func (g *Game) StartBrowser(
  port uint, 
  black_port uint, white_port uint,
  show_board bool, show_result bool,
) {

  // 静的ファイル(html, css, javascript)のディレクトリを設定
  http.Handle("/", http.FileServer(http.Dir("game/asset/")));
  // /gameへのハンドラを設定
  http.HandleFunc("/game", g.gameHandler)
  fmt.Println("http://localhost:8080")

  // ゲームの初期化
  go g.initializeGame(black_port, white_port)

  // 指定ポート番号を用いてhttp通信
  http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

/*
#gameHandler
/gameに対するハンドラ
クライアントとの通信、ゲームとのやり取り
*/
func (g *Game) gameHandler(w http.ResponseWriter, r *http.Request) {
  // リクエストを受け取る構造体
  var request struct {
    Command string
    Col uint8
  };
  // リクエストをパース
  json.NewDecoder(r.Body).Decode(&request);

  var response Response;

  // リクエストのコマンドに従い、処理を行う
  switch request.Command {
  case "start": // ゲーム開始
    g.startGameBrowser(&response);

  case "drop": // 石を落とす
    valid, next_move := g.dropStone(request.Col);
    g.dropStoneBrowser(&response, next_move, valid);
    case "move": // プレイヤーから次の手を取得
    valid, next_move := g.inquireNextMove();
    g.dropStoneBrowser(&response, next_move, valid);

  case "quit": // プレイヤーを終了
    g.quitPlayer();
    response.Quit = true;

  default:
    fmt.Println("???")
  }

  // レスポンス設定
  w.Header().Set("Content-Type", "application/json; charset=UTF-8");
  json.NewEncoder(w).Encode(response);
}

/*
#startGameBrowser
ゲームの開始
*/
func (g *Game) startGameBrowser(response *Response) {
  g.initializeBoard();
  var start bool = g.sendStartCommand();
  response.Start = start;
  response.BlackName = g.BlackName;
  response.WhiteName = g.WhiteName;
  response.BlackPort = g.BlackPort;
  response.WhitePort = g.WhitePort;
}

/*
#dropStoneBrowser
石を盤に落とす
*/
func (g *Game) dropStoneBrowser(response *Response, next_move uint8, valid bool) {
  response.BlackStones = g.Board.BlackStones;
  response.WhiteStones = g.Board.WhiteStones;

  var board_stones uint64 = g.Board.BlackStones | g.Board.WhiteStones;
  response.Board = board_stones;
  response.Counter = g.Board.Counter;
  response.NextMove = next_move;
  response.Valid = valid;
  // 落とした石の位置を取得
  response.Pos = uint8(math.Log2(float64((board_stones & (63<<(next_move*7)))>>(next_move*7))+1))-1 + (next_move*6);

  response.Result = 3;
  // 非合法手が選択された場合
  if (!valid) {
    // 相手の勝ちとしてゲームを終了
    response.Result = g.Board.Counter%2;
    g.endGame(g.Board.Counter%2);
  }

  var stones uint64
  if (g.Board.Counter % 2 == 0) {
    stones = g.Board.WhiteStones;
  } else {
    stones = g.Board.BlackStones;
  }

  // いずれかが勝利した場合
  if (board.CheckAlignment(stones)) {
    response.Result = (g.Board.Counter+1) % 2;
    g.endGame(response.Result);
  }
}
