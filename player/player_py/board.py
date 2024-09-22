"""
コネクトフォーの盤面を操作する関数群
各関数の動作についてはvoda/board/board.goを参照のこと
"""

# 着手可能であるか検証
def can_move(black_stones, white_stones, col):
    board = black_stones | white_stones
    top_mask = 1 << (col*7 + 5)
    return (board & top_mask) == 0

# 石を落とした盤面を生成
def make_move(stones, opp_stones, col):
    board = stones | opp_stones
    col_mask = 63 << (col*7);
    return stones^((((board&col_mask)>>(col*7))+1)<<(col*7))

# 石を取り除く
def remove_stone(stones, opp_stones, col):
    board = stones | opp_stones
    col_mask = 63 << (col*7)
    return stones^((((board&col_mask)+1)<<(col*7))>>1)

# 合法手を生成
def gen_valid_moves(black_stones, white_stones):
    board = black_stones | white_stones;
    top_stones = board & 141845657554976;

    moves = [];
    for i in range(7):
        if (top_stones&(1<<(i*7+5))==0): moves.append(i)

    return moves

# 石が揃ったか検証
def check_alignment(stones):
    if (stones&(stones>>6)&(stones>>12)&(stones>>18)!=0): return True
    if (stones&(stones>>8)&(stones>>16)&(stones>>24)!=0): return True
    if (stones&(stones>>7)&(stones>>14)&(stones>>21)!=0): return True
    if (stones&(stones>>1)&(stones>>2)&(stones>>3)!=0): return True

    return False

# 盤面を表示
def print_board(black_stones, white_stones):
    BLACK_SYMBOL = "o";
    WHITE_SYMBOL = "x";
    EMPTY_SYMBOL = "-";

    for y in range(5, -1, -1):
        for x in range(7):
            bit_mask = (1<<(y+x*7))
            if (black_stones & bit_mask != 0):
                print(BLACK_SYMBOL, end="")
            elif (white_stones & bit_mask != 0):
                print(WHITE_SYMBOL, end="")
            else:
                print(EMPTY_SYMBOL, end="")
        print()
