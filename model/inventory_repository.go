package model

import (
	"github.com/ferdoran/go-sro-framework/db"
)

/*
	Slots:

	0 - Helm
	1 - Chest
	2 - Shoulder
	3 - Gauntlet
	4 - Pants
	5 - Boots
	6 - Primary Weapon
	7 - Secondary Weapon (Ammo, Shield)
	8 - Earring
	9 - Necklace
	10 - Left Ring
	11 - Right Ring
	12 - Extra Slot (PVP Cape, Job Suit)

Avatar Inventory
	0 - Flag
	1 - Chest / Dress
	2 - Attachment
	3 - Helm
	4 - Devil Spirit ?
*/

const (
	SelectCharacterInventory string = "SELECT inv.SLOT, inv.FK_ITEM, it.VARIANCE, it.FK_REF_ITEM FROM `SRO_SHARD`.`INVENTORY` inv INNER JOIN `SRO_SHARD`.`ITEM` AS it ON inv.FK_ITEM = it.ID WHERE inv.FK_CHAR=?"
	InsertInventoryItem      string = "INSERT INTO `SRO_SHARD`.`INVENTORY` (FK_CHAR, SLOT, FK_ITEM) VALUES(?, ?, ?);"
)

// The DB should save CharID, Slot, ItemID
//type Inventory struct {
//	CharID int64
//	Items  map[byte]model.Item
//}

func AddItemToInventory(characterId, itemid int64, slot byte) int64 {
	conn := db.OpenConnShard()
	res, err := conn.Exec(InsertInventoryItem, characterId, slot, itemid)
	db.CheckError(err)
	id, err := res.LastInsertId()
	db.CheckError(err)
	return id
}

func GetCharacterInventory(characterId int64) Inventory {
	conn := db.OpenConnShard()
	queryHandle, err := conn.Query(SelectCharacterInventory, characterId)
	db.CheckError(err)

	items := make(map[byte]Item)
	for queryHandle.Next() {
		var slot, itemId, variance int64
		var refItemId uint32
		queryHandle.Scan(
			&slot,
			&itemId,
			&variance,
			&refItemId)
		refItem := RefItems[refItemId]
		item := Item{
			ID:       int(itemId),
			Name:     refItem.CodeName,
			Variance: uint64(variance),
		}
		item.RefObjectID = refItem.ID
		item.LevelInfo = LevelInfo{
			RequiredLevelType1: refItem.ReqLevelType1,
			RequiredLevelType2: refItem.ReqLevelType2,
			RequiredLevelType3: refItem.ReqLevelType3,
			RequiredLevelType4: refItem.ReqLevelType4,
			RequiredLevel1:     refItem.ReqLevel1,
			RequiredLevel2:     refItem.ReqLevel2,
			RequiredLevel3:     refItem.ReqLevel3,
			RequiredLevel4:     refItem.ReqLevel4,
		}
		item.SetTypeInfo(TypeInfo{
			TypeID1: refItem.TypeID1, TypeID2: refItem.TypeID2, TypeID3: refItem.TypeID3, TypeID4: refItem.TypeID4,
		})
		items[byte(slot)] = item
	}

	return Inventory{Items: items}
}
