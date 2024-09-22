package main

import "fmt"
import "flag"

import "voda/game"
import "player/connector"

func main() {
  // --portでポート番号を設定
  var port *int = flag.Int("port", 8000, "port number");
  // --playerでプレイヤー名を設定
  /*
    random: randomPlayer --- ランダム
    g0F: g0F             --- 乱択アルゴリズム
  */
  var player *string = flag.String("player", "random", "player");

  flag.Parse();

  // プレイヤー関数を設定
  var player_func func(game.PlayerParam) game.PlayerRet;
  switch *player {
  case "random":
    player_func = randomPlayer;
  case "g0F":
    player_func = g0F;
  default:
    fmt.Println(fmt.Sprintf("Unknown Player Name `%s`", player));
    return;
  }

  // ゲームに接続
  connector.Play(player_func, uint(*port));
}
