import { Field } from "./field.js";
import { State, setState } from "./state.js";

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
}

export function startPlaying(size, ships) {
    createField(size, ships);
    createOpponentField(size);
    setState(State.Playing);
}