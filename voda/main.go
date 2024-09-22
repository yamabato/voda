package main

import "flag"

import "voda/game"

/*
#Voda - Вода(Water)
コネクトフォー

board --- 盤面
  board.go --- 盤面の操作

game --- ゲーム
  game.go      --- ゲームの管理
  connector.go --- ソケット通信の管理
  game_data.go --- ゲーム管理のための構造体等
  game_browser.go --- ブラウザ上でのゲーム実行
*/

func main() {
  // 実行時引数の取得
  var port *int = flag.Int("port", 8080, "port number");
  var black_port *int = flag.Int("port1", 8000, "port number for black");
  var white_port *int = flag.Int("port2", 8001, "port number for white");

  var show_board *bool = flag.Bool("board", true, "output the board or not");
  var show_result *bool = flag.Bool("result", true, "output the result or not")

  var cli *bool = flag.Bool("cli", false, "cli");
  flag.Parse();

  var g game.Game;

  if (*cli) {
    g.StartCLI(uint(*black_port), uint(*white_port), *show_board, *show_result);
  }else {
    g.StartBrowser(uint(*port), uint(*black_port), uint(*white_port), *show_board, *show_result);
  }
}
