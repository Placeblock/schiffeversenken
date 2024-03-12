package data

import (
	"testing"
)

func TestFieldCreation(t *testing.T) {
	size := Vector{10, 10}
	field := NewField(size)
	if len(field.Cells) != 100 {
		t.Fatal("Invalid field size")
	}
	for _, cell := range field.Cells {
		if !cell.PossibleShip {
			t.Fatal("Invalid impossible ship")
		}
		if cell.Shot {
			t.Fatal("Invalid shot")
		}
		if cell.Ship != nil {
			t.Fatal("Invalid ship")
		}
	}
}

func TestShipCreation(t *testing.T) {
	ship := NewShip(Vector{0, 0}, Vector{0, 1}, 4)
	test := Vector{0, 0}
	if ship.OccupiedCells[0] != test {
		t.Fatal("Invalid occupied cell")
	}
	test1 := Vector{0, 1}
	if ship.OccupiedCells[1] != test1 {
		t.Fatal("Invalid occupied cell")
	}
	test2 := Vector{0, 2}
	if ship.OccupiedCells[2] != test2 {
		t.Fatal("Invalid occupied cell")
	}
	test3 := Vector{0, 3}
	if ship.OccupiedCells[3] != test3 {
		t.Fatal("Invalid occupied cell")
	}
}

func TestShipPlace(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 0}, Vector{0, 1}, 4)
	if !field.CanAddShip(&ship) {
		t.Fatal("Cannot add ship in valid position")
	}
	ship2 := NewShip(Vector{0, 1}, Vector{0, 1}, 4)
	if field.CanAddShip(&ship2) {
		t.Fatal("Can add ship in invalid position")
	}
	ship3 := NewShip(Vector{3, 0}, Vector{0, 1}, 4)
	if !field.CanAddShip(&ship3) {
		t.Fatal("Cannot add ship in valid position")
	}
}

func TestShipOverlap(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 1}, Vector{0, 1}, 3)
	ship2 := NewShip(Vector{1, 0}, Vector{1, 0}, 3)
	field.AddShip(&ship)
	if field.CanAddShip(&ship2) {
		t.Fatal("Can add ship in invalid position")
	}
	field.AddShip(&ship2)
}

func TestShipNoOverlap(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 1}, Vector{0, 1}, 3)
	ship2 := NewShip(Vector{2, 0}, Vector{1, 0}, 2)
	if !field.CanAddShip(&ship) {
		t.Fatal("Cannot add ship in valid position")
	}
	field.AddShip(&ship)
	if !field.CanAddShip(&ship2) {
		t.Fatal("Cannot add ship in valid position")
	}
	field.AddShip(&ship2)
}

func TestShoot(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 0}, Vector{0, 1}, 4)
	field.AddShip(&ship)
	if field.CanShoot(Vector{0, 4}) {
		t.Fatal("Can shoot in invalid position")
	}
	if !field.CanShoot(Vector{0, 3}) {
		t.Fatal("Cannot shoot in valid position")
	}
	hit, sunk := field.Shoot(Vector{0, 3})
	if !hit {
		t.Fatal("Hit ship without hit being true")
	}
	if sunk {
		t.Fatal("Sunk ship without shot all cells")
	}
}

func TestShipSink(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 0}, Vector{0, 1}, 4)
	field.AddShip(&ship)
	hit, sunk := field.Shoot(Vector{0, 3})
	if !hit {
		t.Fatal("Hit ship without hit being true")
	}
	if sunk {
		t.Fatal("Sunk ship without shot all cells")
	}
	hit, sunk = field.Shoot(Vector{0, 2})
	if !hit {
		t.Fatal("Hit ship without hit being true")
	}
	if sunk {
		t.Fatal("Sunk ship without shot all cells")
	}
	hit, sunk = field.Shoot(Vector{0, 1})
	if !hit {
		t.Fatal("Hit ship without hit being true")
	}
	if sunk {
		t.Fatal("Sunk ship without shot all cells")
	}
	hit, sunk = field.Shoot(Vector{0, 0})
	if !hit {
		t.Fatal("Hit ship without hit being true")
	}
	if !sunk {
		t.Fatal("Did not sunk ship with all cells shot")
	}
}

func TestDefeat(t *testing.T) {
	field := NewField(Vector{4, 4})
	ship := NewShip(Vector{0, 0}, Vector{0, 1}, 1)
	field.AddShip(&ship)
	field.Shoot(Vector{0, 0})
	if !field.IsDefeated() {
		t.Fatal("Is incorrectly undefeated")
	}
}
