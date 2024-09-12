const playerColors = ["red", "blue"];

function getColor(userID) {
    const players = sessionStorage.getItem("players").split(",");
    const userIndex = players.indexOf(userID);
    return userIndex !== -1 ? playerColors[userIndex] : null;
}

async function getBoardData() {
    const response = await fetch(`/game/${sessionStorage.getItem("game_id")}`, {
        method: 'GET',
        headers: {
            "Authorization": `Bearer ${sessionStorage.getItem("token")}`
        }
    });
    if (response.ok) {
        return response.json();
    } else {
        throw new Error(`Failed to fetch board data: ${response.status}`);
    }
}

function setupBoard(boardData, stateData) {
    const board = document.querySelector(".board");
    const groups = groupBoardData(boardData);
    createGroupElements(board, groups);

    for (const [controlGroup, elements] of Object.entries(groups)) {
        const groupElement = document.querySelector(`.group.g${controlGroup}`);
        elements.forEach(data => {
            const stackElement = createStack(data);
            groupElement.appendChild(stackElement);
            groupElement.appendChild(createGroupLabel(controlGroup));
        });
    }

    applyStateToBoard(stateData);
}

function groupBoardData(boardData) {
    return boardData.reduce((groups, item) => {
        if (!groups[item.control_group]) {
            groups[item.control_group] = [];
        }
        groups[item.control_group].push(item);
        return groups;
    }, {});
}

function createGroupElements(board, groups) {
    Object.keys(groups).sort((a, b) => a - b).forEach(controlGroup => {
        const groupElement = document.createElement("div");
        groupElement.className = `group g${controlGroup}`;
        board.appendChild(groupElement);
    });
}

function createGroupLabel(controlGroup) {
    const groupLabel = document.createElement("div");
    groupLabel.className = 'group_label';
    groupLabel.textContent = controlGroup;
    return groupLabel;
}

function createStack(data) {
    const stackContainer = document.createElement("div");
    stackContainer.className = "stack_container";

    const stackLabel = document.createElement("div");
    stackLabel.className = 'stack_label';
    stackLabel.textContent = data.points;
    stackContainer.appendChild(stackLabel);

    const stack = document.createElement("div");
    stack.className = 'stack';
    stack.setAttribute("data", data._id);

    for (let i = 0; i < data.number_of_segments; i++) {
        const segment = document.createElement("div");
        segment.className = `segment ${i}`;
        segment.setAttribute("data", `${data._id}_${i}`);
        stack.appendChild(segment);
    }
    stackContainer.appendChild(stack);
    return stackContainer;
}

function applyStateToBoard(stateData) {
    applyPlayerState(stateData.player_0, "red");
    applyPlayerState(stateData.player_1, "blue");
}

function applyPlayerState(playerState, color) {
    if (!playerState) return;
    playerState.forEach(segmentId => {
        const segment = document.querySelector(`.segment[data="${segmentId}"]`);
        chooseSegment(segment, color);
    });
}

function handleBoardClick(event) {
    const el = event.target;
    if (el.classList.contains('segment') && !el.classList.contains("red") && !el.classList.contains("blue") && document.querySelector(".board").classList.contains("select")) {
        const segmentId = el.getAttribute("data");
        const gameId = sessionStorage.getItem("game_id");

        fetch(`http://${window.location.hostname}:8080/game/${gameId}/occupy/${segmentId}`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${sessionStorage.getItem("token")}`
            }
        });
    }
}

function chooseSegment(el, color) {
    el.classList.add(color);
    const segmentCount = el.parentElement.querySelectorAll(".segment").length;
    const selectedSegmentCount = el.parentElement.querySelectorAll(`.segment.${color}`).length;
    const canOwnIt = selectedSegmentCount > Math.floor(segmentCount / 2);

    if (canOwnIt && !el.parentElement.parentElement.querySelector(".stack_label").classList.contains(color)) {
        el.parentElement.parentElement.querySelector(".stack_label").classList.add(color);
        audioObjects['got_stack.wav'].play();
    } else {
        audioObjects['choose_segment.wav'].play();
    }
}

function chooseGroups(groups) {
    groups.forEach(group => {
        const el = document.querySelector(`.group.g${group}`);
        if (el) {
            el.classList.add("selected");
        }
    });
}

function deSelectGroups(groups) {
    document.querySelectorAll(".group.selected").forEach(el => {
        el.classList.remove("selected");
    });
}

window.addEventListener("load", async function () {
    const userColor = getColor(sessionStorage.getItem("user_id"));
    const board = document.querySelector(".board");
    board.classList.add(userColor);

    try {
        const boardData = await getBoardData();
        setupBoard(boardData.board, boardData.state);

        board.addEventListener("click", handleBoardClick);
        document.getElementById('accountName').innerText = sessionStorage.getItem('user_id');
        game_timeout_in_secs = boardData.state.timeout_in_secs;
        let players = boardData.state.players;
        if (players[boardData.state.currentPlayer] === sessionStorage.getItem('user_id')) {
            document.getElementById("activeplayer").classList.add("select");
        }

    } catch (error) {
        console.error("Error fetching board data:", error);
    }
});

