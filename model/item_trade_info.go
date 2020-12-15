package model

type TradeInfo struct {
	IsTradable   bool
	IsSellable   bool
	IsBuyable    bool
	IsBorrowable bool
	IsDropable   bool
	IsPickable   bool
	IsRepairable bool // equipment
	IsRevivable  bool // equipment, chars
	IsUsable     bool
	IsThrowable  bool
}
