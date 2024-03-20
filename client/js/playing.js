import { Field, newShip } from "./field.js";
import { State, setState } from "./state.js";
import { messageTarget, sendMessage } from "./ws.js";

const fieldElement = document.getElementById("field")
const opponentFieldElement = document.getElementById("field-opponent")

let field;
function createField(size, ships) {
    field = new Field(fieldElement, size)

    for (let ship of ships) {
        field.addShip(ship)
    }
    field.removeDragHandlers()
}

let opponentField;
function createOpponentField(size) {
    opponentField = new Field(opponentFieldElement, size)
    for (let y = 0; y < opponentField.cells.length; y++) {
        for (let x = 0; x < opponentField.cells[y].length; x++) {
            const cell = opponentField.cells[y][x];
            cell.element.onclick = () => shoot(x, y)
        }
    }
}

export function startPlaying(size, ships) {
    createField(size, ships);
    createOpponentField(size);
    setState(State.Playing);
}

let playersTurn = false;
function shoot(x, y) {
    if (!playersTurn) return;
    sendMessage("SHOOT", {cell: {x, y}});
}

const turnEndInfo = document.getElementById("turn-end-info")
messageTarget.addEventListener("TURN_START", () => {
    opponentFieldElement.style.opacity = 1;
    opponentFieldElement.classList.add("shootable")
    turnEndInfo.style.display = "none";
    playersTurn = true;
})

messageTarget.addEventListener("TURN_END", () => {
    opponentFieldElement.style.opacity = 0.5;
    opponentFieldElement.classList.remove("shootable")
    turnEndInfo.style.display = "block";
    playersTurn = false;
})

messageTarget.addEventListener("HIT_OTHER", (e) => {
    const {x, y} = e.detail
    const cell = opponentField.getCell(x, y)
    cell.shot = true
    cell.ship = true
    opponentField.updateCells(true, true)
})

messageTarget.addEventListener("HIT_SELF", (e) => {
    const {x, y} = e.detail
    const cell = field.getCell(x, y)
    cell.shot = true
    cell.ship = true
    field.updateCells(true, false)
})

messageTarget.addEventListener("NO_HIT_OTHER", (e) => {
    const {x, y} = e.detail
    const cell = opponentField.getCell(x, y)
    cell.shot = true
    opponentField.updateCells(true, true)
})

messageTarget.addEventListener("NO_HIT_SELF", (e) => {
    const {x, y} = e.detail
    const cell = field.getCell(x, y)
    cell.shot = true
    field.updateCells(true, false)
})

messageTarget.addEventListener("SUNK_OTHER", (e) => {
    const {position, direction, length} = e.detail
    const ship = newShip(position, direction, length)
    opponentField.addShip(ship)
})

const wonTitle = document.getElementById("won-title")
const lostTitle = document.getElementById("lost-title")
messageTarget.addEventListener("WON", (e) => {
    wonTitle.style.display = "block"
    lostTitle.style.display = "none"
    setState(State.Ended)
})
messageTarget.addEventListener("LOST", (e) => {
    wonTitle.style.display = "none"
    lostTitle.style.display = "block"
    setState(State.Ended)
})

const playAgainBtn = document.getElementById("play-again-btn")
playAgainBtn.onclick = () => {
    sendMessage("POOL", null)
}