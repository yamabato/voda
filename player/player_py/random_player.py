import random

# ランダムに手を選択する
def random_player(stones, opp_stones, valid_moves, moves):
    # ゲームより与えられた合法手一覧からランダムに選択
    return random.choice(valid_moves)
