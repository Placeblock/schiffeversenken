package data

type Vector struct {
	X int
	Y int
}

func (v *Vector) Add(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y}
}

func (v *Vector) Multiply(scalar int) Vector {
	return Vector{int(v.X) * scalar, int(v.Y) * scalar}
}

type Ship struct {
	Position      Vector
	Direction     Vector
	OccupiedCells []Vector
	Length        uint8
	Sunk          bool
}

type Cell struct {
	Ship         *Ship
	PossibleShip bool
	Shot         bool
}

type Field struct {
	Size  Vector
	Ships []*Ship
	Cells map[Vector]*Cell
}

// Initializes a new Ship and calculates all occupied fields
func NewShip(position Vector, direction Vector, length uint8) Ship {
	ship := Ship{Position: position, Direction: direction, Length: length, Sunk: false}

	occupiedCells := make([]Vector, 0)
	for i := uint8(0); i < length; i++ {
		occupiedCells = append(occupiedCells, position.Add(direction.Multiply(int(i))))
	}
	ship.OccupiedCells = occupiedCells

	return ship
}

// Initializes a new Field and creates the associated cells
func NewField(size Vector) Field {
	field := Field{Ships: make([]*Ship, 0), Cells: make(map[Vector]*Cell), Size: size}
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			cell := &Cell{PossibleShip: true, Shot: false}
			field.Cells[Vector{x, y}] = cell
		}
	}
	return field
}

// Adds a ship to the field and marks cells as not accepting new ships
func (f *Field) AddShip(ship *Ship) {
	for _, occupiedCell := range ship.OccupiedCells {
		cell, exists := f.Cells[occupiedCell]
		if !exists {
			continue
		}
		cell.Ship = ship
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				impossibleShipCell := occupiedCell.Add(Vector{dx, dy})
				cell, exists := f.Cells[impossibleShipCell]
				if !exists {
					continue
				}
				cell.PossibleShip = false
			}
		}
	}
	f.Ships = append(f.Ships, ship)
}

// Checks if a new ship can added at a specific location.
func (f *Field) CanAddShip(ship *Ship) bool {
	for _, occupiedCell := range ship.OccupiedCells {
		cell, exists := f.Cells[occupiedCell]
		if !exists {
			return false
		}
		if !cell.PossibleShip {
			return false
		}
	}
	return true
}

// Returns the pointer to a ship at a specific field or nil if not present
func (f *Field) GetShip(position Vector) *Ship {
	cell, exists := f.Cells[position]
	if !exists {
		return nil
	}
	return cell.Ship
}

// Checks if this cell was already shot
func (f *Field) CanShoot(position Vector) bool {
	cell, exists := f.Cells[position]
	if !exists {
		return false
	}
	return !cell.Shot
}

// Shoots at a specific field and returns whether there was a first-time-hit and whether the ship sunk
func (f *Field) Shoot(position Vector) (hit, sunk bool) {
	cell, exists := f.Cells[position]
	if !exists {
		return false, false
	}
	ship := cell.Ship
	if ship == nil || cell.Shot {
		return false, false
	}
	cell.Shot = true
	for _, occupiedCell := range ship.OccupiedCells {
		if !f.Cells[occupiedCell].Shot {
			return true, false
		}
	}
	ship.Sunk = true
	return true, true
}

func (f *Field) isDefeated() bool {
	for _, ship := range f.Ships {
		if !ship.Sunk {
			return false
		}
	}
	return true
}
