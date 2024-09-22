package main

import "fmt"
import "math/rand"

import "voda/game"

/*
#randomPlayer
完全にランダムに手を選択するプレイヤー

*引数
param game.PlayerParam: ゲームから受信したパラメータ

*返り値
game.PlayerRet: ゲームへの応答
*/
func randomPlayer(param game.PlayerParam) game.PlayerRet {
  var command string = param.Command;
  switch (command) {
  case "name":
    return retPlayerName("RandomPlayer-Go");
  case "start":
    return game.PlayerRet{ Command: "ready" };
  case "go":
    return randomChoiceMove(param);
  case "end":
    return game.PlayerRet{ Command: "bye" };
  default:
    fmt.Println(fmt.Sprintf("Unknown Command: `%s`", command));
  }
  return game.PlayerRet{};
}

/*
#retPlayerName
プレイヤー名問い合わせに対する応答を生成

*引数
player_name string: プレイヤー名

*返り値
PlayerRet: プレイヤー名の指定
*/
func retPlayerName(player_name string) game.PlayerRet {
  var ret game.PlayerRet;
  ret.Command = "setname";
  ret.Name = player_name;
  return ret;
}

/*
#randomChoiceMove
ランダムに手を選択する

*引数
param game.PlayerParam: ゲームから受信したパラメータ

*返り値
game.PlayerRet: ゲームへの応答
*/
func randomChoiceMove(param game.PlayerParam) game.PlayerRet {
  var valid_moves = param.ValidMoves;
  var move = valid_moves[rand.Intn(len(valid_moves))];

  return game.PlayerRet {
    Command: "move",
    Move: move,
  };
}
