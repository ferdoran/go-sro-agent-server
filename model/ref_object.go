package model

type RefObject struct {
	ID             uint32
	CodeName       string
	ObjName        string
	OrgObjCodeName string
	NameStrID      string
	DescStrID      string
	IsCashItem     bool
	IsBionic       bool
	TypeInfo
	DecayTime     int
	Country       int
	Rarity        int
	IsTradable    bool
	IsSellable    bool
	IsBuyable     bool
	IsBorrowable  bool
	IsDropable    bool
	IsPickable    bool
	IsRepairable  bool
	IsRevivable   bool
	IsUseable     bool
	IsThrowable   bool
	Price         uint64
	RepairCost    uint64
	ReviveCost    uint64
	BorrowCost    uint64
	KeepingFee    uint64
	SellPrice     uint64
	ReqLevelType1 int
	ReqLevelType2 int
	ReqLevelType3 int
	ReqLevelType4 int
	ReqLevel1     int
	ReqLevel2     int
	ReqLevel3     int
	ReqLevel4     int
	MaxContain    int
	RegionID      int16
	Dir           int
	OffsetX       int
	OffsetY       int
	OffsetZ       int
	Speed1        int
	Speed2        int
	Scale         int
	BCHeight      int
	BCRadius      int
	EventID       int
	AssocFileObj  string
	AssocFileDrop string
	AssocFileIcon string
	AssocFile1    string
	AssocFile2    string
}

var RefObjects map[uint32]RefObject = make(map[uint32]RefObject)
