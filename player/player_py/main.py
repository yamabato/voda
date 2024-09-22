import sys
import argparse

from g0F import g0F
from play import play
from random_player import random_player

if __name__ == "__main__":
    # コマンドライン引数の取得
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", default=8000, type=int)
    parser.add_argument("--player", default="random", type=str)

    args = parser.parse_args()

    # プレイヤー名、プレイヤー関数の設定
    if args.player == "random":
        player_func = random_player
        player_name = "RandomPlayer-Py"
    elif args.player == "g0F":
        player_func = g0F
        player_name = "g0F-Py"
    else:
        print(f"Unknown Player Name {args.player}")
        sys.exit(0)

    play(player_func, player_name, args.port)
