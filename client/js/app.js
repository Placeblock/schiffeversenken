import { Field, Ship } from "./field.js";
import "./ws.js"

const fieldElement = document.getElementById("field");

createField({x: 10, y: 10}, {"5": 1,"4": 2,"3": 3,"2": 4})
function createField(size, settings) {
    const field = new Field(fieldElement, {x: 10, y: 10})
    let y = 0;
    let x = 0;
    let n = 0
    for (let key in settings) {
        const length = Number(key)
        const amount = settings[key]
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