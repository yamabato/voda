var move_count = 0; // 手数
var board = 0; // 盤面

var is_black_human = false; // 先手が人間か
var is_white_human = false; // 後手が人間か

var col = 0; // 石を落とす列

// ゲームを進める
async function game() {
	/*
	result
		0: 先手勝ち
		1: 後手勝ち
		2: 引き分け
		3: ゲーム中
		255: 異常終了
	*/
	result = 3;
	for (i=0; i<42; i++) {
		// ゲームが終了
		if (result != 3) {
			break;
		}
		result = await getNextMove();
	}

	// i==42
	if (result == 3) { result=2; }

	showResult(result);
}

// ゲームを初期化する
async function startGame() {
	// プレイヤー名を待機中にする
	newRotatingStr("Waiting...", "black-player-lbl");
	newRotatingStr("Waiting...", "white-player-lbl");

	// 盤の表示をリセット
	clearBoard();
	// 結果をリセット
	showResult(255);
	// 手数、盤面をリセット
	move_count = 1;
	board = 0;
	showTurn();

	// 開始を通知
	res = await sendRequest({
		command: "start"
	});

	// 各プレイヤーが人間か否か設定する
	if (res["BlackPort"] == 0) {
		is_black_human = true;
	}
	if (res["WhitePort"] == 0) {
		is_white_human = true;
	}

	// プレイヤーの名称
	// 人間なら"Human"
	var black_name = res["BlackName"];
	var white_name = res["WhiteName"];

	if (black_name == "") { black_name = "Human"; }
	if (white_name == "") { white_name = "Human"; }

	document.querySelector("#black-player-lbl").innerHTML = `<span>${black_name}</span>`;
	document.querySelector("#white-player-lbl").innerHTML = `<span>${white_name}</span>`;

	// ゲーム開始
	game();
}

// ゲーム終了
async function quitGame() {
	res = await sendRequest({
		command: "quit"
	});

	document.querySelector("#cover").style.display = "flex";
}


// 手数、手番を表示
function showTurn() {
	document.querySelector("#move-count").innerText = move_count;

	let turn_lbl = document.querySelector("#turn-lbl");
	if (move_count%2 == 1) {
		turn_lbl.innerText = "Black";
	} else {
		turn_lbl.innerText = "White";
	}
}

// 次の手が入力されるのを待つ
function waitCol() {
	return new Promise(function(resolve, reject){
		wait_interval_id = setInterval(()=>{
			if(col!=-1) { clearInterval(wait_interval_id); resolve(); }
			else { return }
		}, 500);
	})
}

// 次の手に進める
async function getNextMove() {
	var pos;
	let is_human = move_count%2==1? is_black_human : is_white_human;

	if (is_human) {
		col = -1;
		await waitCol();
		res = await sendRequest({
			command: "drop",
			col: col,
		});
	} else {
		res = await sendRequest({
			command: "move",
		});
	}

	pos = res["Pos"];
	drop(pos, (move_count+1)%2)
	move_count++;
	showTurn();

	return res["Result"];
}

// 結果表示
function showResult(result) {
	let lbl = document.querySelector("#result-lbl");

	if (result == 0) {
		lbl.innerText = "Black Wins";
	} else if (result == 1) {
		lbl.innerText = "White Wins";
	} else if (result == 2) {
		lbl.innerText = "Draw";
	} else {
		lbl.innerText = "---";
	}
}

async function sendRequest(body) {
	return await fetch("/game", {
		method: "POST",
		body: JSON.stringify(body),
	})
		.then((result) => {
		if (result.ok) {
			return result.json();
		};
	}).then((data) => {
		return data;
	});
}

// 回転する文字列をidの要素に追加する
function newRotatingStr(str, id) {
	let elm = document.querySelector("#" + id);
	elm.innerHTML = "";

	for (i=0; i<str.length; i++) {
		char = document.createElement("span");
		char.classList.add("rotating-char");
		char.innerText = str[i];
		char.style.animationDelay = `${i*0.2}s`;

		elm.appendChild(char);
	}
}

// 石を落とす
function drop(pos, turn) {
  let cell = document.querySelector("#board-" + ("0"+String(pos)).slice(-2));
  cell.innerHTML = `<div class="stone" id="stone-${("0"+String(pos)).slice(-2)}"></div>`
  
  let stone = document.querySelector("#stone-" + ("0"+String(pos)).slice(-2));

	if (turn == 0) {
		stone.classList.add("black");
	} else {
		stone.classList.add("white");
	}

  y = pos % 6;
  stone.animate(
		[
			{ top: `-${(6-y)*10}vmin` },
			{ top: "1vmin" },
		],
		{
			duration: Math.sqrt((6-y))*100,
			easing: "ease-in",
		}
	);
}

// 盤面の表示をリセット
function clearBoard() {
	for (i=0; i<42; i++) {
		let cell = document.querySelector("#board-" + ("0"+String(i)).slice(-2));
		cell.innerHTML = "";
	}
}

// 盤面のクリックを検知し、列を取得
document.querySelector("#board").addEventListener("click", function(event) {
	let board = document.querySelector("#board")
	let boardX = board.getBoundingClientRect().x;

	let clickX = event.pageX - boardX;

	let cell = document.querySelector("#board-00");
	let cellSize = cell.getBoundingClientRect().width;

	col = Math.floor(clickX / cellSize);
});
