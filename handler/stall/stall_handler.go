package stall

import (
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

const (
	StallActionLeave byte = iota + 1
	StallActionEnter
	StallActionBuy
)

const (
	StallUpdateItem byte = iota + 1
	StallAddItem
	StallRemoveItem
	StallFleaMarketMode
	StallState
	StallMessage
	StallName
)

const (
	StallErrorInvalidOperation         uint16 = 0x5
	StallErrorInvalidPrice             uint16 = 0x3C08
	StallErrorNothingToSell            uint16 = 0x3C0C
	StallErrorMarketClosed             uint16 = 0x3C0E
	StallErrorNotEnoughGold            uint16 = 0x3C11
	StallErrorInventoryFull            uint16 = 0x3C12
	StallErrorHostLeft                 uint16 = 0x3C15
	StallErrorImBusy                   uint16 = 0x3C16
	StallErrorClosedByHost             uint16 = 0x3C17
	StallErrorMarketFull               uint16 = 0x3C18
	StallErrorInvalidHostState         uint16 = 0x3C2B
	StallErrorCustomerBanned           uint16 = 0x3C2C
	StallErrorCannotOpenFromHorse      uint16 = 0x3C34
	StallErrorMarketnameNotAllowed     uint16 = 0x3C38
	StallErrorCannotOpenMarketMurderer uint16 = 0x3C39
	StallErrorNotUseJob                uint16 = 0x3C3B
	StallErrorWareNetworkFail          uint16 = 0x3C41
)

type StallHandler struct {
}

func (s *StallHandler) Handle(data server.PacketChannelData) {

}
