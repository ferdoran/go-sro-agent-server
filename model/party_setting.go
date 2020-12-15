package model

type PartySetting byte

const (
	ExpDistItemDistNoInvite   PartySetting = 0 // 0000 0000
	ExpShareItemDistNoInvite  PartySetting = 1 // 0000 0001
	ExpDistItemShareNoInvite  PartySetting = 2 // 0000 0010
	ExpShareItemShareNoInvite PartySetting = 3 // 0000	0011
	ExpDistItemDist           PartySetting = 4 // 0000 0100
	ExpShareItemDist          PartySetting = 5 // 0000 0101
	ExpDistItemShare          PartySetting = 6 // 0000 0110
	ExpShareItemShare         PartySetting = 7 // 0000 0111
)

func (p PartySetting) IsSharingExp() bool {
	return p&1 != 0
}

func (p PartySetting) IsSharingItem() bool {
	return p&2 != 0
}

func (p PartySetting) HasGuestInvite() bool {
	return p&4 != 0
}

func (p PartySetting) ToByte() byte {
	return byte(p)
}
