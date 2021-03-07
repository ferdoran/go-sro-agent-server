package model

import (
	"github.com/ferdoran/go-sro-framework/db"
	"github.com/sirupsen/logrus"
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

	var refCharCount int
	queryHandle, err := dbConn.Query(SelectAllRefCharsCount)
	db.CheckError(err)
	if queryHandle.Next() {
		queryHandle.Scan(&refCharCount)
	}
	queryHandle.Close()

	logrus.Infof("loading %d ref chars from database", refCharCount)

	queryHandle, err = dbConn.Query(SelectAllRefChars)
	db.CheckError(err)

	counter := 1
	var refChars = make(map[uint32]RefChar)
	for queryHandle.Next() {
		var refChar RefChar
		queryHandle.Scan(
			&refChar.ID,
			&refChar.CodeName,
			&refChar.ObjName,
			&refChar.OrgObjCodeName,
			&refChar.NameStrID,
			&refChar.DescStrID,
			&refChar.IsCashItem,
			&refChar.IsBionic,
			&refChar.TypeID1,
			&refChar.TypeID2,
			&refChar.TypeID3,
			&refChar.TypeID4,
			&refChar.DecayTime,
			&refChar.Country,
			&refChar.Rarity,
			&refChar.IsTradable,
			&refChar.IsSellable,
			&refChar.IsBuyable,
			&refChar.IsBorrowable,
			&refChar.IsDropable,
			&refChar.IsPickable,
			&refChar.IsRepairable,
			&refChar.IsRevivable,
			&refChar.IsUseable,
			&refChar.IsThrowable,
			&refChar.Price,
			&refChar.RepairCost,
			&refChar.ReviveCost,
			&refChar.BorrowCost,
			&refChar.KeepingFee,
			&refChar.SellPrice,
			&refChar.ReqLevelType1,
			&refChar.ReqLevel1,
			&refChar.ReqLevelType2,
			&refChar.ReqLevel2,
			&refChar.ReqLevelType3,
			&refChar.ReqLevel3,
			&refChar.ReqLevelType4,
			&refChar.ReqLevel4,
			&refChar.MaxContain,
			&refChar.RegionID,
			&refChar.Dir,
			&refChar.OffsetX,
			&refChar.OffsetY,
			&refChar.OffsetZ,
			&refChar.Speed1,
			&refChar.Speed2,
			&refChar.Scale,
			&refChar.BCHeight,
			&refChar.BCRadius,
			&refChar.EventID,
			&refChar.AssocFileObj,
			&refChar.AssocFileDrop,
			&refChar.AssocFileIcon,
			&refChar.AssocFile1,
			&refChar.AssocFile2,
			&refChar.Level,
			&refChar.CharGender,
			&refChar.MaxHP,
			&refChar.MaxMP,
			&refChar.InventorySize,
			&refChar.CanStoreTID1,
			&refChar.CanStoreTID2,
			&refChar.CanStoreTID3,
			&refChar.CanStoreTID4,
			&refChar.CanBeVehicle,
			&refChar.CanControl,
			&refChar.DamagePortion,
			&refChar.MaxPassenger,
			&refChar.AssocTactics,
			&refChar.PhyDef,
			&refChar.MagDef,
			&refChar.PhyAbsorbRate,
			&refChar.MagAbsorbRate,
			&refChar.EvasionRate,
			&refChar.BlockRate,
			&refChar.HitRate,
			&refChar.CriticalHitRate,
			&refChar.ExpToGive,
			&refChar.CreepType,
			&refChar.Knockdown,
			&refChar.KORecoveryTime,
			&refChar.DefaultSkill1,
			&refChar.DefaultSkill2,
			&refChar.DefaultSkill3,
			&refChar.DefaultSkill4,
			&refChar.DefaultSkill5,
			&refChar.DefaultSkill6,
			&refChar.DefaultSkill7,
			&refChar.DefaultSkill8,
			&refChar.DefaultSkill9,
			&refChar.DefaultSkill10,
			&refChar.TextureType,
			&refChar.Except1,
			&refChar.Except2,
			&refChar.Except3,
			&refChar.Except4,
			&refChar.Except5,
			&refChar.Except6,
			&refChar.Except7,
			&refChar.Except8,
			&refChar.Except9,
			&refChar.Except10,
		)
		counter++
		refChars[refChar.ID] = refChar
	}
	logrus.Infoln("finished loading ref chars")
	return refChars
}
