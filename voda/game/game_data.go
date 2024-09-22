package game

import "sync"

// ゲームの情報
// 盤面、プレイヤーを保持
type Game struct {
  Board BoardData         // 盤面情報

  BlackPort uint // 先手のポート
  WhitePort uint // 後手のポート

  BlackParamChannel chan PlayerParam  // 先手パラメータ送信用チャネル
  WhiteParamChannel chan PlayerParam  // 後手パラメータ送信用チャネル

  BlackRetChannel chan PlayerRet  // 先手返り値受信用チャネル
  WhiteRetChannel chan PlayerRet  // 後手返り値受信用チャネル

  WaitGroup *sync.WaitGroup

  BlackName string        // 先手プレイヤー名
  WhiteName string        // 後手プレイヤー名
}

// コネクトフォーのゲーム情報を保持
type BoardData struct {
  BlackStones uint64  // 先手の石
  WhiteStones uint64  // 後手の石
  Moves []uint8       // 操作履歴
  Counter uint8       // 手数
}

// プレイヤーに与える引数
type PlayerParam struct {
  Command string

  Turn bool // 先後

  Stones uint64      // 石の配置
  OppStones uint64   // 相手方の石の配置
  Moves []uint8      // 操作履歴
  ValidMoves []uint8 // 合法手リスト

  Result uint8// 結果(0:win, 1:lose, 2:draw)
}

// プレイヤーからの返り値
type PlayerRet struct {
  Command string 

  Name string // プレイヤー名

  Ready bool //開始の確認

  Move uint8 // 操作
}

// クライアントへのレスポンス
type Response struct {
  Start bool
  Quit bool

  BlackName string // 先手の名称
  WhiteName string // 後手の名称
  BlackPort uint// 先手の名称
  WhitePort uint// 後手の名称

  BlackStones uint64  // 先手の石
  WhiteStones uint64  // 後手の石
  Board uint64
  Counter uint8       // 手数
  Pos uint8
  Result uint8

  NextMove uint8
  Valid bool
}
