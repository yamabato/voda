import random

from board import *

# 乱択アルゴリズム
def g0F(stones, opp_stones, valid_moves, moves):
    result_tbl = {}

    # 試行回数
    TIMES = 500
    move_count = len(moves) + 1;

    for move in valid_moves:
        result_tbl[move] = [0, 0, 0]

        # 最初の一手を反映した盤面を生成
        _stones = make_move(stones, opp_stones, move)
        if (move_count%2 == 1):
            black_stones = _stones
            white_stones = opp_stones
        else:
            black_stones = opp_stones
            white_stones = _stones

        # 指定回数試行
        for i in range(TIMES):
            result = play_out(black_stones, white_stones, move_count)
            result_tbl[move][result] += 1;

    if (move_count%2 == 1):
        valid_moves.sort(key=lambda x: result_tbl[x][0])
    else:
        valid_moves.sort(key=lambda x: result_tbl[x][1])

    return valid_moves[-1]

# ゲームを最後までプレイ
def play_out(black_stones, white_stones, move_count):
    stones = [black_stones, white_stones]
    side = 0

    for i in range(move_count, 42):
        side = i%2

        if check_alignment(stones[not side]):
            return int(not side)

        valid_moves = gen_valid_moves(*stones)
        move = random.choice(valid_moves)
        stones[side] = make_move(stones[side], stones[not side], move)

    return 2

