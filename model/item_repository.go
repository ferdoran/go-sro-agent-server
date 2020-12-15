package model

import "gitlab.ferdoran.de/game-dev/go-sro/framework/db"

const (
	InsertItem = "INSERT INTO `SRO_SHARD`.`ITEM` (FK_REF_ITEM, VARIANCE) VALUES (?, ?);"
)

func CreateItem(referenceItemId uint32, variance uint64) int64 {
	conn := db.OpenConnShard()
	res, err := conn.Exec(InsertItem, referenceItemId, variance)
	db.CheckError(err)

	id, err := res.LastInsertId()
	db.CheckError(err)
	return id
}
