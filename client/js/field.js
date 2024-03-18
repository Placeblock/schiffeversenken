
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

export class Field {
    constructor(element, size) {
        this.element = element;
        this.cells = [];
        this.ships = [];
        this.size = size;
        this.createFieldElements()
        this.updateShips()
    }

    addShip(ship, draggable=true) {
        this.ships.push(ship)
        this.updateShips()
        if (draggable) {
            this.addDragHandlers(ship)
            ship.element.draggable = true;
        }
        this.getCell(ship.position.x, ship.position.y).element.appendChild(ship.element)
    }

    disableOccupied() {
        this.element.classList.remove("visualized")
    }

    addDragHandlers(ship) {
        ship.element.ondragstart = this.dragStartHandler.bind(this)
        ship.element.ondragend = this.dragEndHandler.bind(this)
        ship.element.onclick = this.clickHandler.bind(this)
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

    updateShips() {
        for (let y = 0; y < this.cells.length; y++) {
            const row = this.cells[y]
            for (let x = 0; x < row.length; x++) {
                const cell = this.getCell(x, y)
                if (cell == undefined) continue
                cell.ship = false
                cell.canPlaceShip = true
                cell.element.classList.remove("ship-blocked")
            }
        }
        for (let ship of this.ships) {
            const occupiedCells = ship.getOccupiedCells();
            for (let occupiedCell of occupiedCells) {
                this.getCell(occupiedCell.x, occupiedCell.y).ship = true
                for (let dx = -1; dx <= 1; dx++) {
                    for (let dy = -1; dy <= 1; dy++) {
                        const ox = occupiedCell.x+dx
                        const oy = occupiedCell.y+dy
                        const cell = this.getCell(ox, oy)
                        if (cell == undefined) continue
                        cell.canPlaceShip = false
                        cell.element.classList.add("ship-blocked")
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
        const dx = Math.floor(ev.layerX / 30);
        const dy = Math.floor(ev.layerY / 30);
        ev.dataTransfer.setData("drag-delta-x", dx)
        ev.dataTransfer.setData("drag-delta-y", dy)
        ev.dataTransfer.dropEffect = "move";
        ev.dataTransfer.setData("text/plain", ev.target.id);
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
        const data = ev.dataTransfer.getData("text/plain");
        if (data == "") return
        const src = document.getElementById(data);
        const ship = this.getShip(src);
        if (ship == undefined) return

        const dx = parseInt(ev.dataTransfer.getData("drag-delta-x"))
        const dy = parseInt(ev.dataTransfer.getData("drag-delta-y"))
        const x = parseInt(ev.target.getAttribute("cell-x"))
        const y = parseInt(ev.target.getAttribute("cell-y"))
        const targetX = x-dx;
        const targetY = y-dy;
        const targetCell = this.getCell(targetX, targetY);
        const newPos = {x: targetX, y: targetY}
        if (!this.canMoveShip(ship, newPos)) return;
        ship.position = newPos
        targetCell.element.appendChild(src);
        this.updateShips();
    }

    canMoveShip(ship, newPosition, newDirection) {
        if (newDirection == undefined) {
            newDirection = ship.direction
        }
        const cells = ship.getOccupiedCells(newPosition, newDirection);
        this.ships.splice(this.ships.indexOf(ship), 1)
        this.updateShips()
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
            this.updateShips()
        }   
    }
}