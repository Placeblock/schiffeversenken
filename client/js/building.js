import { State, setState } from "./state.js"
import { messageTarget, sendMessage } from "./ws.js"
import { Field, Ship } from "./field.js"
import { startPlaying } from "./playing.js";

const fieldElement = document.getElementById("building-field");

let field;
let fieldSize;

messageTarget.addEventListener("FIELD", (data) => {
    const ships = data.detail.settings.ships
    createField(fieldElement, data.detail.size, ships)
})

messageTarget.addEventListener("STATE", (data) => {
    if (data.detail == "building") {
        setState(State.Building)
    }
    if (data.detail == "playing") {
        startPlaying(fieldSize, field.ships)
    }
})

function createField(element, size, ships) {
    field = new Field(element, size)
    fieldSize = size;
    let y = 0;
    let x = 0;
    let n = 0
    for (let key in ships) {
        const length = Number(key)
        const amount = ships[key]
        for (let i = 0; i < amount; i++) {
            if (x+length >= size.x) {
                y +=2;
                x = 0;
            }
            const shipElement = document.createElement("div")
            shipElement.classList.add("ship")
            shipElement.id = "ship-"+n
            const ship = new Ship(shipElement, {x: x, y: y}, {x: 1, y: 0}, length)
            field.addShip(ship)
            x+=length+1
            n++;
        }
    }
}

const submitBtn = document.getElementById("submit-field")
const waitingText = document.getElementById("building-waiting")

submitBtn.onclick = () => {
    for (let ship of field.ships) {
        sendMessage("SHIP", ship)
    }
    submitBtn.style.display = "none"
    waitingText.style.display = "block"
}