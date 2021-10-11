package navmeshv2

import "github.com/pkg/errors"

var ErrCellNotInObject = errors.New("cell not in object")
var ErrNoObjectCellForPosition = errors.New("no cell found for position")
var ErrEmptyCollisionList = errors.New("collision list is empty")
var ErrCellCast = errors.New("cell is not of that type")
