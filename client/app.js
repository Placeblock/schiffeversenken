

const field = document.getElementById("field");
const ship = document.getElementById("ship");

for (let y = 0; y < 10; y++) {
    const row = document.createElement("tr")  
    for (let x = 0; x < 10; x++) {
        const cell = document.createElement("td")
        cell.classList.add("cell")
        cell.ondragover = onDragOver;
        cell.ondrop = onDrop;
        cell.id = getCellID(x, y)
        cell.setAttribute("cell-x", x)
        cell.setAttribute("cell-y", y)
        row.appendChild(cell)
    }
    field.appendChild(row)
}

function getCellID(x, y) {
    return "cell-"+x+"-"+y
}

function getCell(x, y) {
    return document.getElementById(getCellID(x, y))
}

ship.ondragstart = onDragStart;
function onDragStart(ev) {
    const dx = Math.floor(ev.layerX / 40);
    const dy = Math.floor(ev.layerY / 40);
    ev.dataTransfer.setData("drag-delta-x", dx)
    ev.dataTransfer.setData("drag-delta-y", dy)
    ev.dataTransfer.dropEffect = "move";
    ev.dataTransfer.setData("text/plain", ev.target.id);
}

function onDragOver(ev) {
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "move";
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
    targetElement.appendChild(document.getElementById(data));
}