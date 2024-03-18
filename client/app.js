import { Field, Ship } from "./field.js";

const fieldElement = document.getElementById("field");
const shipElement = document.getElementById("ship");
const shipElement2 = document.getElementById("ship2");

const field = new Field(fieldElement, {x: 10, y: 10})
const ship = new Ship(shipElement, {x: 2, y: 2}, {x: 1, y: 0}, 3)
const ship2 = new Ship(shipElement2, {x: 5, y: 5}, {x: 0, y: 1}, 4)
field.addShip(ship)
field.addShip(ship2)
/*
shipElement.ondragstart = onDragStart;
function onDragStart(ev) {
    const dx = Math.floor(ev.layerX / 40);
    const dy = Math.floor(ev.layerY / 40);
    ev.dataTransfer.setData("drag-delta-x", dx)
    ev.dataTransfer.setData("drag-delta-y", dy)
    ev.dataTransfer.dropEffect = "move";
    ev.dataTransfer.setData("text/plain", ev.target.id);
    requestAnimationFrame(() => {
        ev.target.style.display = "none"
    })
}
shipElement.ondragend = onDragEnd;
function onDragEnd(ev) {
    ev.target.style.display = "block"
}

function onDrop(ev) {
    ev.preventDefault();
    const dx = parseInt(ev.dataTransfer.getData("drag-delta-x"))
    const dy = parseInt(ev.dataTransfer.getData("drag-delta-y"))
    const x = parseInt(ev.target.getAttribute("cell-x"))
    const y = parseInt(ev.target.getAttribute("cell-y"))
    const targetX = x-dx;
    const targetY = y-dy;
    const targetElement = getCell(targetX, targetY);
    const data = ev.dataTransfer.getData("text/plain");
    const src = document.getElementById(data);
    targetElement.appendChild(src);
}*/