package model

import (
	"github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/db"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/utils"
)

const (
	SelectAllRefChars      = "SELECT * FROM `SRO_SHARD`.`CHAR_REF_DATA`"
	SelectAllRefCharsCount = "SELECT COUNT(*) FROM `SRO_SHARD`.`CHAR_REF_DATA`"
)

type RefChar struct {
	RefObject
	Level           int
	CharGender      int
	MaxHP           int
	MaxMP           int
	InventorySize   int
	CanStoreTID1    bool
	CanStoreTID2    bool
	CanStoreTID3    bool
	CanStoreTID4    bool
	CanBeVehicle    bool
	CanControl      bool
	DamagePortion   int
	MaxPassenger    int
	AssocTactics    int
	PhyDef          int
	MagDef          int
	PhyAbsorbRate   int
	MagAbsorbRate   int
	EvasionRate     int
	BlockRate       int
	HitRate         int
	CriticalHitRate int
	ExpToGive       int
	CreepType       int
	Knockdown       int
	KORecoveryTime  int
	DefaultSkill1   int
	DefaultSkill2   int
	DefaultSkill3   int
	DefaultSkill4   int
	DefaultSkill5   int
	DefaultSkill6   int
	DefaultSkill7   int
	DefaultSkill8   int
	DefaultSkill9   int
	DefaultSkill10  int
	TextureType     int
	Except1         int
	Except2         int
	Except3         int
	Except4         int
	Except5         int
	Except6         int
	Except7         int
	Except8         int
	Except9         int
	Except10        int
}

var RefChars map[uint32]RefChar

func GetAllRefChars() map[uint32]RefChar {
	dbConn := db.OpenConnShard()
	defer dbConn.Close()

	var refItemCount int
	queryHandle, err := dbConn.Query(SelectAllRefCharsCount)
	db.CheckError(err)
	if queryHandle.Next() {
		queryHandle.Scan(&refItemCount)
	}
	queryHandle.Close()

	logrus.Infoln("loading ref chars from database")

	queryHandle, err = dbConn.Query(SelectAllRefChars)
	db.CheckError(err)

	counter := 1
	var refChars = make(map[uint32]RefChar)
	for queryHandle.Next() {
		var refItem RefChar
		queryHandle.Scan(
			&refItem.ID,
			&refItem.CodeName,
			&refItem.ObjName,
			&refItem.OrgObjCodeName,
			&refItem.NameStrID,
			&refItem.DescStrID,
			&refItem.IsCashItem,
			&refItem.IsBionic,
			&refItem.TypeID1,
			&refItem.TypeID2,
			&refItem.TypeID3,
			&refItem.TypeID4,
			&refItem.DecayTime,
			&refItem.Country,
			&refItem.Rarity,
			&refItem.IsTradable,
			&refItem.IsSellable,
			&refItem.IsBuyable,
			&refItem.IsBorrowable,
			&refItem.IsDropable,
			&refItem.IsPickable,
			&refItem.IsRepairable,
			&refItem.IsRevivable,
			&refItem.IsUseable,
			&refItem.IsThrowable,
			&refItem.Price,
			&refItem.RepairCost,
			&refItem.ReviveCost,
			&refItem.BorrowCost,
			&refItem.KeepingFee,
			&refItem.SellPrice,
			&refItem.ReqLevelType1,
			&refItem.ReqLevel1,
			&refItem.ReqLevelType2,
			&refItem.ReqLevel2,
			&refItem.ReqLevelType3,
			&refItem.ReqLevel3,
			&refItem.ReqLevelType4,
			&refItem.ReqLevel4,
			&refItem.MaxContain,
			&refItem.RegionID,
			&refItem.Dir,
			&refItem.OffsetX,
			&refItem.OffsetY,
			&refItem.OffsetZ,
			&refItem.Speed1,
			&refItem.Speed2,
			&refItem.Scale,
			&refItem.BCHeight,
			&refItem.BCRadius,
			&refItem.EventID,
			&refItem.AssocFileObj,
			&refItem.AssocFileDrop,
			&refItem.AssocFileIcon,
			&refItem.AssocFile1,
			&refItem.AssocFile2,
			&refItem.Level,
			&refItem.CharGender,
			&refItem.MaxHP,
			&refItem.MaxMP,
			&refItem.InventorySize,
			&refItem.CanStoreTID1,
			&refItem.CanStoreTID2,
			&refItem.CanStoreTID3,
			&refItem.CanStoreTID4,
			&refItem.CanBeVehicle,
			&refItem.CanControl,
			&refItem.DamagePortion,
			&refItem.MaxPassenger,
			&refItem.AssocTactics,
			&refItem.PhyDef,
			&refItem.MagDef,
			&refItem.PhyAbsorbRate,
			&refItem.MagAbsorbRate,
			&refItem.EvasionRate,
			&refItem.BlockRate,
			&refItem.HitRate,
			&refItem.CriticalHitRate,
			&refItem.ExpToGive,
			&refItem.CreepType,
			&refItem.Knockdown,
			&refItem.KORecoveryTime,
			&refItem.DefaultSkill1,
			&refItem.DefaultSkill2,
			&refItem.DefaultSkill3,
			&refItem.DefaultSkill4,
			&refItem.DefaultSkill5,
			&refItem.DefaultSkill6,
			&refItem.DefaultSkill7,
			&refItem.DefaultSkill8,
			&refItem.DefaultSkill9,
			&refItem.DefaultSkill10,
			&refItem.TextureType,
			&refItem.Except1,
			&refItem.Except2,
			&refItem.Except3,
			&refItem.Except4,
			&refItem.Except5,
			&refItem.Except6,
			&refItem.Except7,
			&refItem.Except8,
			&refItem.Except9,
			&refItem.Except10,
		)
		utils.PrintProgress(counter, refItemCount)
		counter++
		refChars[refItem.ID] = refItem
	}
	logrus.Infoln("finished loading ref chars")
	return refChars
}
