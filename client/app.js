import { Field, Ship } from "./field.js";

const fieldElement = document.getElementById("field");
const shipElement = document.getElementById("ship");
const shipElement2 = document.getElementById("ship2");

const field = new Field(fieldElement, {x: 10, y: 10})
const ship = new Ship(shipElement, {x: 2, y: 2}, {x: 1, y: 0}, 3)
const ship2 = new Ship(shipElement2, {x: 5, y: 5}, {x: 0, y: 1}, 4)
field.addShip(ship)
field.addShip(ship2)