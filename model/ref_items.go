package model

import (
	"github.com/ferdoran/go-sro-framework/db"
	"github.com/sirupsen/logrus"
)

const (
	SelectAllRefItems     = "SELECT * FROM `SRO_SHARD`.`ITEMDATA`;"
	SelectCountOfRefItems = "SELECT COUNT(*) FROM `SRO_SHARD`.`ITEMDATA`;"
)

type RefItem struct {
	RefObject
	StackSize                  int
	RequiredGender             int
	RequiredStr                int
	RequiredInt                int
	ItemClass                  int
	SetID                      int
	DurabilityLower            float64
	DurabilityUpper            float64
	PhyDefLower                float64
	PhyDefUpper                float64
	PhyDefIncrease             float64
	EvasionRateLower           float64
	EvasionRateUpper           float64
	EvasionRateIncrease        float64
	PhyAbsRateLower            float64
	PhyAbsRateUpper            float64
	PhyAbsRateIncrease         float64
	BlockingRateLower          float64
	BlockingRateUpper          float64
	MagDefLower                float64
	MagDefUpper                float64
	MagDefIncrease             float64
	MagAbsRateLower            float64
	MagAbsRateUpper            float64
	MagAbsRateIncrease         float64
	PhyDefReinforceLower       float64
	PhyDefReinforceUpper       float64
	MagDefReinforceLower       float64
	MagDefReinforceUpper       float64
	Quivered                   int
	Ammo1TypeID4               int
	Ammo2TypeID4               int
	Ammo3TypeID4               int
	Ammo4TypeID4               int
	Ammo5TypeID4               int
	SpeedClass                 int
	IsTwoHanded                bool
	Range                      int
	PhyAttackMinLower          float64
	PhyAttackMinUpper          float64
	PhyAttackMaxLower          float64
	PhyAttackMaxUpper          float64
	PhyAttackIncrease          float64
	MagAttackMinLower          float64
	MagAttackMinUpper          float64
	MagAttackMaxLower          float64
	MagAttackMaxUpper          float64
	MagAttackIncrease          float64
	PhyAttackReinforceMinLower float64
	PhyAttackReinforceMinUpper float64
	PhyAttackReinforceMaxLower float64
	PhyAttackReinforceMaxUpper float64
	MagAttackReinforceMinLower float64
	MagAttackReinforceMinUpper float64
	MagAttackReinforceMaxLower float64
	MagAttackReinforceMaxUpper float64
	HitRateLower               float64
	HitRateUpper               float64
	HitRateIncrease            float64
	CritRateLower              float64
	CritRateUpper              float64
	Param1                     int
	Param2                     int
	Param3                     int
	Param4                     int
	Param5                     int
	Param6                     int
	Param7                     int
	Param8                     int
	Param9                     int
	Param10                    int
	Param11                    int
	Param12                    int
	Param13                    int
	Param14                    int
	Param15                    int
	Param16                    int
	Param17                    int
	Param18                    int
	Param19                    int
	Param20                    int
	Param1Desc                 string
	Param2Desc                 string
	Param3Desc                 string
	Param4Desc                 string
	Param5Desc                 string
	Param6Desc                 string
	Param7Desc                 string
	Param8Desc                 string
	Param9Desc                 string
	Param10Desc                string
	Param11Desc                string
	Param12Desc                string
	Param13Desc                string
	Param14Desc                string
	Param15Desc                string
	Param16Desc                string
	Param17Desc                string
	Param18Desc                string
	Param19Desc                string
	Param20Desc                string
	MaxMagicOptCount           int
	ChildItemCount             int
}

var RefItems map[uint32]RefItem

func GetAllRefItems() map[uint32]RefItem {
	dbConn := db.OpenConnShard()
	defer dbConn.Close()

	var refItemCount int
	queryHandle, err := dbConn.Query(SelectCountOfRefItems)
	db.CheckError(err)
	if queryHandle.Next() {
		queryHandle.Scan(&refItemCount)
	}
	queryHandle.Close()

	logrus.Infof("loading %d ref items from database", refItemCount)

	queryHandle, err = dbConn.Query(SelectAllRefItems)
	db.CheckError(err)

	counter := 1
	var refItems = make(map[uint32]RefItem)
	for queryHandle.Next() {
		var refItem RefItem
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
			&refItem.StackSize,
			&refItem.RequiredGender,
			&refItem.RequiredStr,
			&refItem.RequiredInt,
			&refItem.ItemClass,
			&refItem.SetID,
			&refItem.DurabilityLower,
			&refItem.DurabilityUpper,
			&refItem.PhyDefLower,
			&refItem.PhyDefUpper,
			&refItem.PhyDefIncrease,
			&refItem.EvasionRateLower,
			&refItem.EvasionRateUpper,
			&refItem.EvasionRateIncrease,
			&refItem.PhyAbsRateLower,
			&refItem.PhyAbsRateUpper,
			&refItem.PhyAbsRateIncrease,
			&refItem.BlockingRateLower,
			&refItem.BlockingRateUpper,
			&refItem.MagDefLower,
			&refItem.MagDefUpper,
			&refItem.MagDefIncrease,
			&refItem.MagAbsRateLower,
			&refItem.MagAbsRateUpper,
			&refItem.MagAbsRateIncrease,
			&refItem.PhyDefReinforceLower,
			&refItem.PhyDefReinforceUpper,
			&refItem.MagDefReinforceLower,
			&refItem.MagDefReinforceUpper,
			&refItem.Quivered,
			&refItem.Ammo1TypeID4,
			&refItem.Ammo2TypeID4,
			&refItem.Ammo3TypeID4,
			&refItem.Ammo4TypeID4,
			&refItem.Ammo5TypeID4,
			&refItem.SpeedClass,
			&refItem.IsTwoHanded,
			&refItem.Range,
			&refItem.PhyAttackMinLower,
			&refItem.PhyAttackMinUpper,
			&refItem.PhyAttackMaxLower,
			&refItem.PhyAttackMaxUpper,
			&refItem.PhyAttackIncrease,
			&refItem.MagAttackMinLower,
			&refItem.MagAttackMinUpper,
			&refItem.MagAttackMaxLower,
			&refItem.MagAttackMaxUpper,
			&refItem.MagAttackIncrease,
			&refItem.PhyAttackReinforceMinLower,
			&refItem.PhyAttackReinforceMinUpper,
			&refItem.PhyAttackReinforceMaxLower,
			&refItem.PhyAttackReinforceMaxUpper,
			&refItem.MagAttackReinforceMinLower,
			&refItem.MagAttackReinforceMinUpper,
			&refItem.MagAttackReinforceMaxLower,
			&refItem.MagAttackReinforceMaxUpper,
			&refItem.HitRateLower,
			&refItem.HitRateUpper,
			&refItem.HitRateIncrease,
			&refItem.CritRateLower,
			&refItem.CritRateUpper,
			&refItem.Param1,
			&refItem.Param1Desc,
			&refItem.Param2,
			&refItem.Param2Desc,
			&refItem.Param3,
			&refItem.Param3Desc,
			&refItem.Param4,
			&refItem.Param4Desc,
			&refItem.Param5,
			&refItem.Param5Desc,
			&refItem.Param6,
			&refItem.Param6Desc,
			&refItem.Param7,
			&refItem.Param7Desc,
			&refItem.Param8,
			&refItem.Param8Desc,
			&refItem.Param9,
			&refItem.Param9Desc,
			&refItem.Param10,
			&refItem.Param10Desc,
			&refItem.Param11,
			&refItem.Param11Desc,
			&refItem.Param12,
			&refItem.Param12Desc,
			&refItem.Param13,
			&refItem.Param13Desc,
			&refItem.Param14,
			&refItem.Param14Desc,
			&refItem.Param15,
			&refItem.Param15Desc,
			&refItem.Param16,
			&refItem.Param16Desc,
			&refItem.Param17,
			&refItem.Param17Desc,
			&refItem.Param18,
			&refItem.Param18Desc,
			&refItem.Param19,
			&refItem.Param19Desc,
			&refItem.Param20,
			&refItem.Param20Desc,
			&refItem.MaxMagicOptCount,
			&refItem.ChildItemCount,
		)
		counter++
		refItems[refItem.ID] = refItem
	}
	logrus.Infoln("finished loading ref items")
	return refItems
}
