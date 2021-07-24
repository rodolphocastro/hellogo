package main

import (
	"testing"
	"time"
)

// A struct is a collection of related fields
type Person struct {
	ID     int
	Name   string
	BornAt time.Time
}

// "Constructors" must be created as functions
func birthPerson(name string) *Person {
	p := Person{
		ID:     -1,
		Name:   name,
		BornAt: time.Now(),
	}

	return &p
}

// "Methods" must be created as functions with receiver params
func (p Person) IsMinor() bool {
	return p.BornAt.Year()-time.Now().Year() > 18
}

// "Methods" that modify the state have to use pointers, yikes
func (p *Person) Rename(newName string) {
	p.Name = newName
}

func TestStructs(t *testing.T) {
	rodolpho := Person{
		ID:     1,
		Name:   "Rodolpho's Clone",
		BornAt: time.Date(2021, time.April, 30, 0, 0, 0, 0, time.UTC),
	}

	if rodolpho.ID != 1 {
		t.Errorf("expected ID be %v, but found %v", 1, rodolpho.ID)
	}

	baby := birthPerson("A name")
	if baby.Name != "A name" {
		t.Errorf("expected Name to be %v, but found %v", "A name", baby.Name)
	}

	otherBaby := birthPerson("Another name")
	if baby == otherBaby {
		t.Error("expected otherBaby to be different than baby")
	}

	if baby.IsMinor() {
		t.Error("expected baby to be a minor but it is not")
	}

	baby.Rename("Fus roh dah")
	if baby.Name == "A name" {
		t.Error("expected baby to have been renamed")
	}
}
