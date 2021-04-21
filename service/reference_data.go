package service

import (
	"errors"
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/db"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	SelectAllRefItems     = "SELECT * FROM `SRO_SHARD`.`ITEMDATA`;"
	SelectCountOfRefItems = "SELECT COUNT(*) FROM `SRO_SHARD`.`ITEMDATA`;"

	SelectAllRefChars      = "SELECT * FROM `SRO_SHARD`.`CHAR_REF_DATA`"
	SelectAllRefCharsCount = "SELECT COUNT(*) FROM `SRO_SHARD`.`CHAR_REF_DATA`"
)

type ReferenceDataService struct {
	referenceItems      map[uint32]model.RefItem
	referenceCharacters map[uint32]model.RefChar
	mutex               sync.RWMutex
}

var referenceDataServiceInstance *ReferenceDataService
var referenceDataServiceOnce sync.Once

func GetReferenceDataServiceInstance() *ReferenceDataService {
	referenceDataServiceOnce.Do(func() {
		referenceDataServiceInstance = &ReferenceDataService{
			referenceItems:      loadReferenceItems(),
			referenceCharacters: loadReferenceCharacters(),
			mutex:               sync.RWMutex{},
		}
	})
	return referenceDataServiceInstance
}

func (r *ReferenceDataService) GetReferenceItem(id uint32) (model.RefItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if referenceItem, exists := r.referenceItems[id]; exists {
		return referenceItem, nil
	} else {
		return model.RefItem{}, errors.New(fmt.Sprintf("reference item with id %d does not exist", id))
	}
}

func (r *ReferenceDataService) GetReferenceCharacter(id uint32) (model.RefChar, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if referenceCharacter, exists := r.referenceCharacters[id]; exists {
		return referenceCharacter, nil
	} else {
		return model.RefChar{}, errors.New(fmt.Sprintf("reference character with id %d does not exist", id))
	}
}

func loadReferenceItems() map[uint32]model.RefItem {
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
	var refItems = make(map[uint32]model.RefItem)
	for queryHandle.Next() {
		var refItem model.RefItem
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

func loadReferenceCharacters() map[uint32]model.RefChar {
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
	var refChars = make(map[uint32]model.RefChar)
	for queryHandle.Next() {
		var refChar model.RefChar
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
