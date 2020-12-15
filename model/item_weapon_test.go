package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeaponTypes(t *testing.T) {
	w := Weapon{}
	w.TypeID1 = 3
	w.TypeID2 = 1
	w.TypeID3 = 6
	w.TypeID4 = 2

	weaponType := w.GetWeaponType()
	assert.Equal(t, Sword, weaponType, "Weapon type not referenced correctly")
}
