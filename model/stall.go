package model

import (
	"sync"
)

type StallEntry struct {
	StallSlot byte
	InventorySlot byte
	StackCount uint16
	Price uint64
	FleaMarketnetworkTidGroup uint32
	UnkUshort0 uint16
	Mutex *sync.Mutex
}

type Stall struct {
	PlayerId uint32 // Player.UniqueID
	Entries []StallEntry
}

var Stalls map[uint32]*Stall = make(map[uint32]*Stall)

func (s *StallEntry) AddItem(playerId uint32) *Stall {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if stall, ok := Stalls[playerId]; !ok {
		var newStallEntry []StallEntry
		newStallEntry = append(newStallEntry, *s)
		newStall := Stall {
			PlayerId: playerId,
			Entries: newStallEntry,
		}
		Stalls[playerId] = &newStall
	} else {
		stall.Entries = append(stall.Entries, *s)
	}
	return Stalls[playerId]
}

func (s *StallEntry) RemoveItem(playerId uint32) *Stall {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	var newEntries []StallEntry
	if stalls, ok := Stalls[playerId]; ok {
		for _, v := range stalls.Entries {
	        if v.StallSlot != s.StallSlot {
	            newEntries = append(newEntries, v)
	        }
    	}
    	if len(newEntries) > 0 {
    		Stalls[playerId].Entries = newEntries
    	} else {
    		delete(Stalls, playerId)
    	}
    	return Stalls[playerId]
	} else {
		return nil
	}
}

func (s *StallEntry) UpdateItem(playerId uint32) StallEntry {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	var entry StallEntry
	if stalls, ok := Stalls[playerId]; ok {
		for _, v := range stalls.Entries {
	        if v.StallSlot == s.StallSlot {
	            v.StackCount = s.StackCount
	            v.Price = s.Price
	            v.UnkUshort0 = s.UnkUshort0
	            entry = v
	            break
	        }
    	}
    }
    return entry
}

func GetStall(playerId uint32) (*Stall, bool) {
	stall, ok := Stalls[playerId]
	return stall, ok
}