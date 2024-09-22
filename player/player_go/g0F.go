package main

import "fmt"
import "sort"
import "math/rand"

import "voda/game"
import "voda/board"

/*
#g0F
乱択アルゴリズム
次で勝つ、次で負ける、打ってはならないなどの手を検出しない

*引数
param game.PlayerParam: ゲームから受信したパラメータ

*返り値
game.PlayerRet: ゲームへの応答
*/
func g0F(param game.PlayerParam) game.PlayerRet {
  var command string = param.Command;
  switch (command) {
  case "name":
    return retPlayerName("g0F-Go");
  case "start":
    return game.PlayerRet{ Command: "ready" };
  case "go":
    return g0FChoiceMove(param);
  case "end":
    return game.PlayerRet{ Command: "bye" };
  default:
    fmt.Println(fmt.Sprintf("Unknown Command: `%s`", command));
  }
  return game.PlayerRet{};
}

/*
#g0FChoiceMove
乱択アルゴリズムにより手を選択

*引数
param game.PlayerParam: ゲームから受信したパラメータ

*返り値
game.PlayerRet: ゲームへの応答
*/
func g0FChoiceMove(param game.PlayerParam) game.PlayerRet {
  var valid_moves = param.ValidMoves;
  // ランダムに最後までプレイした結果
  var result_tbl map[uint8]*[3]uint = make(map[uint8]*[3]uint, 7);
  var next_move uint8;

  var stones uint64 = param.Stones;
  var opp_stones uint64 = param.OppStones;
  var black_stones uint64;
  var white_stones uint64;
  var result uint8;

  // 各手に対するプレイの回数
  const TIMES int = 500;
  var move_count int = len(param.Moves) + 1;

  for _, move := range valid_moves {
    result_tbl[move] = &[3]uint{ 0, 0, 0 };

    stones = param.Stones;
    stones = board.MakeMove(stones, opp_stones, move);

    if (move_count%2 == 1) {
      black_stones = stones;
      white_stones = opp_stones;
    } else {
      black_stones = opp_stones;
      white_stones = stones;
    }

    for i:=0; i<TIMES; i++ {
      result = playOut(black_stones, white_stones, move_count);
      result_tbl[move][result]++;
    }
    fmt.Println(move, result_tbl[move]);
  }

  if (move_count%2 == 1) {
    sort.Slice(valid_moves, func(i int, j int) bool { return result_tbl[valid_moves[i]][0] > result_tbl[valid_moves[j]][0]; });
  } else {
    sort.Slice(valid_moves, func(i int, j int) bool { return result_tbl[valid_moves[i]][1] > result_tbl[valid_moves[j]][1]; });
  }
  next_move = valid_moves[0];

  fmt.Println(next_move, "\n");

  return game.PlayerRet {
    Command: "move",
    Move: next_move,
  };
}

/*
#playOut
ランダムに最後までプレイする
*/
func playOut(black_stones uint64, white_stones uint64, move_count int) uint8 {
  var result uint8 = 2;
  var valid_moves []uint8;
  var move uint8;

  var stones *uint64;
  var opp_stones *uint64;

  for i:=move_count; i<42; i++ {
    if (i%2 == 0) {
      stones = &black_stones;
      opp_stones = &white_stones;
    } else {
      stones = &white_stones;
      opp_stones = &black_stones;
    }

    if (board.CheckAlignment(*opp_stones)) {
      result = uint8((i+1)%2);
      break;
    }

    valid_moves = board.GenValidMoves(black_stones, white_stones);
    move = valid_moves[rand.Intn(len(valid_moves))];

    *stones = board.MakeMove(*stones, *opp_stones, move);
  }

  return result;
}
