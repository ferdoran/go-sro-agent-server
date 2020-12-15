package model

type PartyPurpose byte

const (
	MonsterHuntingParty PartyPurpose = 0
	QuestingParty       PartyPurpose = 1
	HunterParty         PartyPurpose = 2
	ThiefParty          PartyPurpose = 3
)

func (p PartyPurpose) IsMonsterHuntingParty() bool {
	return p == MonsterHuntingParty
}

func (p PartyPurpose) IsQuestingParty() bool {
	return p == QuestingParty
}

func (p PartyPurpose) IsHunterParty() bool {
	return p == HunterParty
}

func (p PartyPurpose) IsThiefParty() bool {
	return p == ThiefParty
}

func (p PartyPurpose) ToByte() byte {
	return byte(p)
}
