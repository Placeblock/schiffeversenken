
export class Cell {
    constructor(element) {
        this.element = element
        this.ship = false
        this.shot = false
        this.canPlaceShip = true
    }
}

export class Ship {
    constructor(element, position, direction, length) {
        this.element = element;
        this.position = position;
        this.direction = direction;
        this.length = length;
        this.updateElementSize();
    }

    sink() {
        this.element.classList.add("sunk")
    }

    getOccupiedCells(position=this.position, direction=this.direction) {
        const cells = []
        for (let l = 0; l < this.length; l++) {
            const x = position.x + l*direction.x
            const y = position.y + l*direction.y
            cells.push({x, y})
        }
        return cells
    }

    getOppositeDirection() {
        if (this.direction.x == 1) {
            return {x: 0, y: 1}
        } else {
            return {x: 1, y: 0}
        }
    }

    rotate() {
        this.direction = this.getOppositeDirection()
        this.updateElementSize();
    }

    updateElementSize() {
        this.element.style.width = Math.max(this.direction.x*this.length, 1)*31-1 + "px"
        this.element.style.height = Math.max(this.direction.y*this.length, 1)*31-1 + "px"
    }
}

let i = 0
export function newShip(position, direction, length) {
    const shipElement = document.createElement("div")
    shipElement.classList.add("ship")
    shipElement.id = "ship-"+i
    const ship = new Ship(shipElement, position, direction, length)
    i++
    return ship
}

export class Field {
    constructor(element, size) {
        this.element = element;
        this.element.innerHTML = ""
        this.cells = [];
        this.ships = [];
        this.size = size;
        this.createFieldElements()
        this.updateShipBoxes()
    }

    addShip(ship, draggable=true) {
        this.ships.push(ship)
        this.updateShipBoxes()
        if (draggable) {
            this.addDragHandlers(ship)
        }
        this.getCell(ship.position.x, ship.position.y).element.appendChild(ship.element)
    }

    hideShips() {
        for (let ship of this.ships) {
            ship.element.style.display = "none"
        }
    }

    showShips() {
        for (let ship of this.ships) {
            ship.element.style.display = "block"
        }
    }

    disableOccupied() {
        this.element.classList.remove("visualized")
    }

    addDragHandlers(ship) {
        ship.element.ondragstart = this.dragStartHandler.bind(this)
        ship.element.ondragend = this.dragEndHandler.bind(this)
        ship.element.onclick = this.clickHandler.bind(this)
        ship.element.draggable = true;
    }

    removeDragHandlers() {
        for (let ship of this.ships) {
            ship.element.ondragstart = null;
            ship.element.ondragend = null;
            ship.element.onclick = null;
            ship.element.draggable = false;
        }
    }

    createFieldElements() {
        for (let y = 0; y < this.size.y; y++) {
            const rowElement = document.createElement("tr")
            const row = []
            for (let x = 0; x < this.size.x; x++) {
                const cellElement = document.createElement("td")
                cellElement.classList.add("cell")
                cellElement.setAttribute("cell-x", x)
                cellElement.setAttribute("cell-y", y)
                cellElement.ondragover = this.dragOverHandler.bind(this)
                cellElement.ondrop = this.dropHandler.bind(this)
                rowElement.appendChild(cellElement)
                const cell = new Cell(cellElement);
                row.push(cell);
            }
            this.element.appendChild(rowElement)
            this.cells.push(row)
        }
    }
    
    getCellID(x, y) {
        return "cell-"+x+"-"+y
    }

    getCell(x, y) {
        if (this.cells.length <= y || y < 0 ||
            this.cells[0].length <= x || x < 0) return undefined
        return this.cells[y][x]
    }

    shoot(x, y) {
        const cell = this.getCell(x, y);
        cell.shot = true
        cell.element.innerText = "X"
    }

    updateShipBoxes() {
        for (let y = 0; y < this.cells.length; y++) {
            const row = this.cells[y]
            for (let x = 0; x < row.length; x++) {
                const cell = this.getCell(x, y)
                if (cell == undefined) continue
                //cell.ship = false
                cell.canPlaceShip = true
                cell.element.classList.remove("ship-blocked")
            }
        }
        for (let ship of this.ships) {
            const occupiedCells = ship.getOccupiedCells();
            for (let occupiedCell of occupiedCells) {
                for (let dx = -1; dx <= 1; dx++) {
                    for (let dy = -1; dy <= 1; dy++) {
                        const ox = occupiedCell.x+dx
                        const oy = occupiedCell.y+dy
                        const cell = this.getCell(ox, oy)
                        if (cell == undefined) continue
                        cell.canPlaceShip = false
                        if (occupiedCells.some(c => c.x == ox && c.y == oy)) continue
                        cell.element.classList.add("ship-blocked")
                    }
                }
            }
        }
    }

    updateCells() {
        for (let y = 0; y < this.cells.length; y++) {
            for (let x = 0; x < this.cells[y].length; x++) {
                const cell = this.getCell(x, y)
                if (cell.shot) {
                    let shotMarker = cell.element.querySelector("#shot-marker")
                    if (shotMarker == null) {
                        shotMarker = document.createElement("span")
                        shotMarker.id = "shot-marker"
                        cell.element.appendChild(shotMarker)
                    }
                    if (cell.ship) {
                        shotMarker.innerText = "X"
                        cell.element.classList.add("ship-cell")
                    } else {
                        shotMarker.innerText = "O"
                        cell.element.classList.remove("ship-cell")
                    }
                }
            }
        }
    }

    getShip(element) {
        for (let ship of this.ships) {
            if (ship.element == element) {
                return ship
            }
        }
        return undefined
    }

    dragStartHandler(ev) {
        const ship = this.getShip(ev.target)
        if (ship == undefined) {
            ev.setCancelled(true)
            return
        }
        console.log(ev);
        console.log(ev.layerX);
        console.log(ev.layerY);
        const dx = Math.floor(ev.layerX / 30);
        const dy = Math.floor(ev.layerY / 30);
        console.log(dx);
        console.log(dy);
        const data = {id: ev.target.id, dx, dy}
        ev.dataTransfer.dropEffect = "move";
        ev.dataTransfer.setData("text/plain", JSON.stringify(data));
        setTimeout(() => {
            ev.target.style.display = "none"
        }, 0)
    }

    dragEndHandler(ev) {
        const ship = this.getShip(ev.target)
        if (ship == undefined) return
        ev.target.style.display = "block"
    }

    dragOverHandler(ev) {
        ev.preventDefault();
        ev.dataTransfer.dropEffect = "move";
    }

    dropHandler(ev) {
        ev.preventDefault();
        console.log(ev)
        const data = JSON.parse(ev.dataTransfer.getData("text/plain"));
        const src = document.getElementById(data.id);
        const ship = this.getShip(src);
        if (ship == undefined) return

        const dx = data.dx
        const dy = data.dy
        console.log(dx);
        console.log(dy);
        const x = parseInt(ev.target.getAttribute("cell-x"))
        const y = parseInt(ev.target.getAttribute("cell-y"))
        console.log(x);
        console.log(y);
        const targetX = x-dx;
        const targetY = y-dy;
        const targetCell = this.getCell(targetX, targetY);
        const newPos = {x: targetX, y: targetY}
        if (!this.canMoveShip(ship, newPos)) return;
        ship.position = newPos
        targetCell.element.appendChild(src);
        this.updateShipBoxes();
    }

    canMoveShip(ship, newPosition, newDirection) {
        if (newDirection == undefined) {
            newDirection = ship.direction
        }
        const cells = ship.getOccupiedCells(newPosition, newDirection);
        this.ships.splice(this.ships.indexOf(ship), 1)
        this.updateShipBoxes()
        for (let occupiedCell of cells) {
            const cell = this.getCell(occupiedCell.x, occupiedCell.y);
            if (cell == undefined || !cell.canPlaceShip) {
                this.addShip(ship);
                return false;
            }
        }
        this.addShip(ship);
        return true
    }

    clickHandler(ev) {
        const ship = this.getShip(ev.target);
        console.log(ship)
        if (ship == undefined) return
        if (this.canMoveShip(ship, ship.position, ship.getOppositeDirection())) {
            ship.rotate()
            this.updateShipBoxes()
        }   
    }
}