import socket

# ゲームとプレイヤー関数のやり取り
def play(player_func, name, port):
    # 指定ポート番号に接続
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(("localhost", port))

    while True:
        # ゲームからのメッセージを受信
        msg = client.recv(1024).decode("utf-8")

        command = msg.split()[0]
        # 送信されたコマンドごとの処理
        if (command == "name"):
            res = ret_player_name(name)
        elif (command == "go"):
            res = build_next_move_response(msg, player_func)
        elif (command == "start"):
            res = "ready"
        elif (command == "end"):
            res = "bye"
        elif (command == "quit"):
            break

        # レスポンスを送信
        client.sendall(res.encode("utf-8"))

    # 通信終了
    client.close()

# プレイヤー名を設定するためのレスポンスを設定
def ret_player_name(name):
    return f"setname {name}"

# 次の手を返すためのレスポンスを構成
def build_next_move_response(msg, player_func):
    params = msg.split()[1:]

    stones = int(params[0]) # こちら側の石
    opp_stones = int(params[1]) # 相手方の石
    valid_moves = list(map(int, list(params[2]))) # 合法手

    # これまでの履歴
    if (len(params) == 3):
        moves = []
    else:
        moves = list(map(int, list(params[3])))

    # 次の手を生成
    next_move = player_func(stones, opp_stones, valid_moves, moves)

    return f"move {next_move}"
