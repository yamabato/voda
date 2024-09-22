package board

import "fmt"

/*
#board
コネクトフォーの盤面
・盤面をuint64として表現する
・操作履歴の保持は行わない
*/

/*
#ルール
・縦6、横7の盤面を使い、先手・後手が順に石を置く
・先に石を縦、横、斜めいずれかで4つ揃えた側の勝利
・石を置けるのは各列最上段の石の上のみ、石がない場合は最下段
*/

/*
#盤面の表現

black: 先手
white: 後手

black_stones uint64: 先手の盤面
white_stones uint64: 後手の盤面

// 各マス表現
    06 13 20 27 34 41 48
    --------------------
    05 12 19 26 33 40 47
    04 11 18 25 32 39 46
    03 10 17 24 31 38 45
    02 09 16 23 30 37 44
    01 08 15 22 29 36 43
    00 07 14 21 28 35 42
col 00 01 02 03 04 05 06

06, 13, 20, 27, 34, 41, 48は常に0
*/

/*
#CanMove
列に置けるか確認する

*引数
black_stones uint64: 先手の盤面
white_stones uint64: 後手の盤面
col uint8          : 列 0~6

*返り値
bool: その列に置けるか否か

board uint64     : 盤面
top_mask uint64  : col列の最上段用のマスク
*/
func CanMove(black_stones uint64 , white_stones uint64, col uint8) bool {
  // 先手、後手の石を合成し、盤面ビット列を生成する
  var board uint64 = black_stones | white_stones;

  // 最上段用のマスクを生成する
  // col*7の左シフトでcol列分ずらす
  // 5の左シフトで最上段までずらす
  var top_mask uint64 = 1 << (col*7 + 5);

  // 最上段に石がなければ置くことができる
  return (board&top_mask)==0;
}

/*
#MakeMove
与えられた盤面boardに対し、col列に石を置いた状態を返す
その列に置けるか否か、又盤面の正当性は検証しない

*引数
stones uint64    : 変更を加える盤面
opp_stones uint64: 相手方の盤面
col uint8        : 列 0~6

*返り値
uint64: 石を置いた後の盤面

board uint64   : 盤面
col_mask uint64: 列を抜き出すマスク
*/
func MakeMove(stones uint64 , opp_stones uint64, col uint8) uint64 {
  // 先手、後手の石を合成し、盤面ビット列を生成する
  var board uint64 = stones | opp_stones;

  // 列を抜き出すマスク
  // 63(=0b111111)を(col*7)左シフトすることでcol列目のみ111111となる
  var col_mask uint64 = 63 << (col*7);

  /*
  #石を置く手順
  1. (board&col_mask)で、boardのcol列を抜き出す
  2. (board&col_mask)を(col*7)だけ右シフトし、col列だけを残す
  3. boardのcol列に於いて、石を置くマスまでは1で埋まっている(はず)
  4. そのため、(((board&col_mask)>>(col*7))+1)は石を置くマスのみ1が立つ
  5. (((board&col_mask)>>(col*7))+1)を(col*7)だけ左シフトすることで、目的の列まで移動させる
  6. 4までで得られたマスクと盤面の排他的論理和を取ると盤面を更新できる

  排他的論理和を使うと、同じマスクをもう一度適用することで操作取り消しができる
  */
  return stones^((((board&col_mask)>>(col*7))+1)<<(col*7));
}

/*
#RemoveStone
指定列最上段の石を除く

*引数
stones uint64: 変更を加える盤面
col uint8    : 列 0~6

*返り値
uint64: 石を置いた後の盤面

board uint64   : 盤面
col_mask uint64: 列を抜き出すマスク
*/
func RemoveStone(stones uint64 , opp_stones uint64, col uint8) uint64 {
  // 先手、後手の石を合成し、盤面ビット列を生成する
  var board uint64 = stones | opp_stones;

  // 列を抜き出すマスク
  var col_mask uint64 = 63 << (col*7);

  // MakeMoveと同じ手順で次に石を置くためのマスクを生成
  // そのマスクを1ビット右シフトすることで、除く石の位置を取得
  return stones^((((board&col_mask)+1)<<(col*7))>>1);
}

/*
#GenValidMoves
石を置ける列のリストを返す

*引数
black_stones uint64: 先手の盤面
white_stones uint64: 後手の盤面

*返り値
[]uint8: 石を置ける列の配列

moves []uint8    : 石を置ける列の配列
board uint64     : 盤面
top_stones uint64: 最上段の石
*/
func GenValidMoves(black_stones uint64 , white_stones uint64) []uint8 {
  var moves []uint8;

  // 先手、後手の石を合成し、盤面ビット列を生成する
  var board uint64 = black_stones | white_stones;

  // 最上段の石だけを抜き出す
  // (1<<5)|(1<<12)|(1<<19)|(1<<26)|(1<<33)|(1<<40)|(1<<47) = 141845657554976
  var top_stones uint64 = board & 141845657554976;

  var i uint8;
  for i=0; i<7; i++ {
    // 最上段が0の列のみ抜き出し
    if (top_stones&(1<<(i*7+5))==0) { moves=append(moves,i) }
  }

  return moves;
}

/*
#CheckAlignment
石が揃ったか検出する

*引数
stones uint64: 盤面

*返り値
bool: 石が揃ったか
*/
func CheckAlignment(stones uint64) bool {
  /*
  #連続の検出

  4つ連続で石が並んでいれば、3シフトしても一つは重なる
  111|1|
   11|1|1
    1|1|11
     |1|111
  */

  // 右下がり
  if (stones&(stones>>6)&(stones>>12)&(stones>>18)!=0) { return true; }

  // 右上がり
  if (stones&(stones>>8)&(stones>>16)&(stones>>24)!=0) { return true; }

  // 横
  if (stones&(stones>>7)&(stones>>14)&(stones>>21)!=0) { return true; }

  // 縦
  if (stones&(stones>>1)&(stones>>2)&(stones>>3)!=0) { return true; }

  return false;
}

/*
#PrintBoard
盤面を標準出力に出力

*引数
black_stones uint64: 先手の盤面
white_stones uint64: 後手の盤面
*/
func PrintBoard(black_stones uint64 , white_stones uint64) {
  // 各コマの記号
  const BLACK_SYMBOL string = "o";
  const WHITE_SYMBOL string = "x";
  const EMPTY_SYMBOL string = "-";

  // 特定のマスを抜き出すためのマスク
  var bit_mask uint64;

  for y:=5; y>=0; y-- {
    for x:=0; x<7; x++ {
      bit_mask = (1<<(y+x*7));

      if (black_stones & bit_mask != 0) {
        fmt.Print(BLACK_SYMBOL);
      } else if (white_stones & bit_mask != 0) {
        fmt.Print(WHITE_SYMBOL);
      }else {
        fmt.Print(EMPTY_SYMBOL);
      }
    }
    fmt.Print("\n");
  }
}
