package model

import (
	"github.com/ferdoran/go-sro-framework/db"
	log "github.com/sirupsen/logrus"
	"time"
	"database/sql"
)

type Char struct {
	ID          int
	RefObjID    int
	User        int
	Shard       int
	Name        string
	Scale       byte
	Level       int
	MaxLevel    int
	Exp         int64
	SkillExp    int64
	SkillPoints int64
	Str         int
	Int         int
	StatPoints  int
	HP          int
	MP          int
	IsDeleting  bool
	PosX        float32
	PosY        float32
	PosZ        float32
	Ctime       time.Time
	Utime       time.Time
	Region      int16
	IsGm        bool
}

const (
	SelectCharByName           string = "SELECT * FROM `SRO_SHARD`.`CHAR` WHERE CHAR_NAME=?;"
	SelectCharsByAccountId     string = "SELECT ID, REF_OBJ_ID, CHAR_NAME, CHAR_SCALE, CURRENT_LEVEL, EXP, SKILL_EXP, STRENGTH, INTELLECT, STAT_POINTS, HP, MP, DELETING, UTIME, FK_REGION FROM `SRO_SHARD`.`CHAR` WHERE FK_USER=? ORDER BY CTIME ASC"
	select_does_charname_exist string = "SELECT 1 FROM `SRO_SHARD`.`CHAR` WHERE CHAR_NAME=? LIMIT 1"
	update_is_deleting         string = "UPDATE `SRO_SHARD`.`CHAR` SET DELETING=? WHERE CHAR_NAME=?"
	insert_char                string = "INSERT INTO `SRO_SHARD`.`CHAR`(REF_OBJ_ID, FK_USER, FK_SHARD, CHAR_NAME, CHAR_SCALE, CURRENT_LEVEL, EXP, SKILL_EXP, STRENGTH, INTELLECT, STAT_POINTS, HP, MP, DELETING, POS_X, POS_Y, POS_Z, FK_REGION) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	SelectIsGm                 string = "SELECT IS_GM FROM `SRO_ACCOUNT`.`USER` WHERE ID=?"
)

func GetCharactersByUserId(userid int) []Char {
	if userid < 0 {
		return nil
	}

	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(SelectCharsByAccountId, userid)
	db.CheckError(err)

	var chars []Char
	for queryHandle.Next() {
		var id, refObjId, level, strength, intellect, statPoints, hp, mp int
		var scale byte
		var exp, skillExp int64
		var name string
		var isDeleting int
		var utime sql.NullTime
		var regionId int16

		err = queryHandle.Scan(&id, &refObjId, &name, &scale, &level, &exp, &skillExp, &strength, &intellect, &statPoints, &hp, &mp, &isDeleting, &utime, &regionId)
		db.CheckError(err)

		chars = append(chars, Char{
			ID:         id,
			RefObjID:   refObjId,
			Name:       name,
			Scale:      scale,
			Level:      level,
			Exp:        exp,
			SkillExp:   skillExp,
			Str:        strength,
			Int:        intellect,
			StatPoints: statPoints,
			HP:         hp,
			MP:         mp,
			Utime:      utime.Time,
			IsDeleting: isDeleting == 1,
			Region:     regionId,
		})
	}

	return chars
}

func CreateChar(char Char, weapon, chest, boots, pants uint32) (bool, int64) {
	conn := db.OpenConnShard()
	defer conn.Close()

	stmt, err := conn.Prepare(insert_char)
	db.CheckError(err)

	res, err := stmt.Exec(char.RefObjID, char.User, char.Shard, char.Name, char.Scale, char.Level, char.Exp, char.SkillExp, char.Str, char.Int, char.StatPoints, char.HP, char.MP, char.IsDeleting, char.PosX, char.PosY, char.PosZ, char.Region)
	db.CheckError(err)

	id, err := res.LastInsertId()
	db.CheckError(err)

	weaponId := CreateItem(weapon, 0)
	chestId := CreateItem(chest, 0)
	bootsId := CreateItem(boots, 0)
	pantsId := CreateItem(pants, 0)

	AddItemToInventory(id, weaponId, SlotPrimaryWeapon)
	AddItemToInventory(id, chestId, SlotChest)
	AddItemToInventory(id, bootsId, SlotBoots)
	AddItemToInventory(id, pantsId, SlotPants)

	return true, id
}

func DoesCharNameExist(name string) bool {
	if name == "" {
		return false
	}

	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(select_does_charname_exist, name)
	db.CheckError(err)

	var val int
	if !queryHandle.Next() {
		// No data available
		return false
	}
	err = queryHandle.Scan(&val)
	db.CheckError(err)

	return val == 1
}

func MarkCharIsDeletion(isDeleting int, name string) bool {
	if name == "" {
		return false
	}

	conn := db.OpenConnShard()
	defer conn.Close()

	stmt, err1 := conn.Prepare(update_is_deleting)
	db.CheckError(err1)

	_, err2 := stmt.Exec(isDeleting, name)
	db.CheckError(err2)

	return true
}

func GetCharacterByName(charName string) Char {
	if charName == "" {
		return Char{}
	}

	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(SelectCharByName, charName)
	db.CheckError(err)

	if queryHandle.Next() {
		char := Char{}
		queryHandle.Scan(&char.ID, &char.RefObjID, &char.User, &char.Shard, &char.Name, &char.Scale, &char.Level, &char.Exp, &char.SkillExp, &char.Str, &char.Int, &char.StatPoints, &char.HP, &char.MP, &char.IsDeleting, &char.PosX, &char.PosY, &char.PosZ, &char.Ctime, &char.Utime, &char.Region, &char.SkillPoints, &char.MaxLevel)

		connAcc := db.OpenConnAccount()
		defer connAcc.Close()
		queryHandleAcc, errAcc := conn.Query(SelectIsGm, char.User)
		db.CheckError(errAcc)
		if queryHandleAcc.Next() {
			var isGmVal int
			queryHandleAcc.Scan(&isGmVal)
			if isGmVal == 1 {
				char.IsGm = true
			} else {
				char.IsGm = false
			}
		}
		return char
	}
	log.Info("GET CHARACHTER ", charName)
	return Char{}
}
