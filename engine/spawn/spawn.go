package spawn

//
//import (
//	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
//	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/packetutils"
//	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
//	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
//	"sync"
//)
//
//const (
//	PVPCapeNone byte = iota
//	PVPCapeRed
//	PVPCapeGray
//	PVPCapeBlue
//	PVPCapeWhite
//	PVPCapeOrange
//)
//
//const VisibleObjectRange = 1000
//
//type SpawnEngine struct {
//	VisibleObjects map[uint32]map[uint32]model.ISRObject
//	mutex          *sync.Mutex
//}
//
//var spawnEngine *SpawnEngine
//var spawnEngineOnce sync.Once
//
//func GetSpawnEngineInstance() *SpawnEngine {
//	spawnEngineOnce.Do(func() {
//		spawnEngine = &SpawnEngine{
//			VisibleObjects: make(map[uint32]map[uint32]model.ISRObject),
//			mutex:          &sync.Mutex{},
//		}
//	})
//	return spawnEngine
//}
//
//func (s *SpawnEngine) PositionChanged(p *model.Player) {
//	world := model.GetSroWorldInstance()
//
//	objsToDespawn := make([]uint32, 0)
//	objsToSpawn := make([]uint32, 0)
//
//	if s.VisibleObjects[p.UniqueID] == nil {
//		s.VisibleObjects[p.UniqueID] = make(map[uint32]model.ISRObject, 0)
//	}
//
//	// 1. despawn objects that are out of sight
//
//	for uid := range s.VisibleObjects[p.UniqueID] {
//		obj := world.VisibleObjects[uid]
//
//		if obj.GetType() != model.PlayerType {
//			delete(s.VisibleObjects[uid], p.UniqueID)
//			continue
//		}
//
//		player := world.PlayersByUniqueId[uid]
//
//		if obj.GetPosition().DistanceTo(p.Position) > VisibleObjectRange {
//			objsToDespawn = append(objsToDespawn, uid)
//
//			// TODO despawn packet to other player
//			p1 := network.EmptyPacket()
//			p1.MessageID = opcode.EntityGroupSpawnBegin
//			p1.WriteByte(2)
//			p1.WriteUInt16(1)
//
//			player.Session.Conn.Write(p1.ToBytes())
//
//			p2 := network.EmptyPacket()
//			p2.MessageID = opcode.EntityGroupSpawnData
//			p2.WriteUInt32(p.UniqueID)
//			player.Session.Conn.Write(p2.ToBytes())
//
//			p3 := network.EmptyPacket()
//			p3.MessageID = opcode.EntityGroupSpawnEnd
//			player.Session.Conn.Write(p3.ToBytes())
//
//			delete(s.VisibleObjects[uid], p.UniqueID)
//		}
//	}
//
//	// 2. spawn player to other players if visible
//	for uid := range s.VisibleObjects {
//
//		if uid == p.UniqueID {
//			continue
//		}
//
//		//object := world.VisibleObjects[uid]
//		player := world.PlayersByUniqueId[uid]
//
//		if player.GetPosition().DistanceTo(p.Position) <= VisibleObjectRange {
//
//			if s.VisibleObjects[p.UniqueID][uid] == nil {
//				// in visible objects already
//				objsToSpawn = append(objsToSpawn, uid)
//			}
//
//			if s.VisibleObjects[uid][p.UniqueID] == nil {
//				p1 := network.EmptyPacket()
//				p1.MessageID = opcode.EntityGroupSpawnBegin
//				p1.WriteByte(1)
//				p1.WriteUInt16(1)
//
//				player.Session.Conn.Write(p1.ToBytes())
//
//				p2 := network.EmptyPacket()
//				p2.MessageID = opcode.EntityGroupSpawnData
//				packetutils.WriteEntitySpawn(&p2, p)
//				player.Session.Conn.Write(p2.ToBytes())
//
//				p3 := network.EmptyPacket()
//				p3.MessageID = opcode.EntityGroupSpawnEnd
//				player.Session.Conn.Write(p3.ToBytes())
//
//				s.VisibleObjects[uid][p.UniqueID] = p
//			}
//		}
//
//		//if object.GetPosition().DistanceTo(p.Position) <= VisibleObjectRange {
//		//	if s.VisibleObjects[p.UniqueID][uid] == nil {
//		//		// in visible objects already
//		//		objsToSpawn = append(objsToSpawn, uid)
//		//	}
//		//}
//	}
//
//	// 3. despawn other players to player
//	if len(objsToDespawn) > 0 {
//		p1 := network.EmptyPacket()
//		p1.MessageID = opcode.EntityGroupSpawnBegin
//		p1.WriteByte(2)
//		p1.WriteUInt16(uint16(len(objsToDespawn)))
//
//		p.Session.Conn.Write(p1.ToBytes())
//
//		for _, uid := range objsToDespawn {
//			p2 := network.EmptyPacket()
//			p2.MessageID = opcode.EntityGroupSpawnData
//			p2.WriteUInt32(uid)
//			p.Session.Conn.Write(p2.ToBytes())
//			delete(s.VisibleObjects[p.UniqueID], uid)
//		}
//
//		p3 := network.EmptyPacket()
//		p3.MessageID = opcode.EntityGroupSpawnEnd
//		p.Session.Conn.Write(p3.ToBytes())
//	}
//
//	// 4. spawn other players to player
//	if len(objsToSpawn) > 0 {
//		p1 := network.EmptyPacket()
//		p1.MessageID = opcode.EntityGroupSpawnBegin
//		p1.WriteByte(1)
//		p1.WriteUInt16(uint16(len(objsToSpawn)))
//
//		p.Session.Conn.Write(p1.ToBytes())
//
//		for _, uid := range objsToSpawn {
//			player := world.PlayersByUniqueId[uid]
//			p2 := network.EmptyPacket()
//			p2.MessageID = opcode.EntityGroupSpawnData
//			packetutils.WriteEntitySpawn(&p2, player)
//			p.Session.Conn.Write(p2.ToBytes())
//			s.VisibleObjects[p.UniqueID][uid] = player
//		}
//
//		p3 := network.EmptyPacket()
//		p3.MessageID = opcode.EntityGroupSpawnEnd
//		p.Session.Conn.Write(p3.ToBytes())
//	}
//	//TODO: if len > 1  use opcode.GroupSpawnBegin, opcode.GroupSpawnData, opcode.GroupSpawnEnd
//}
//
//func (s *SpawnEngine) StartedMoving(p *model.Player, targetPos model.Position) {
//	world := model.GetSroWorldInstance()
//
//	for uid, _ := range s.VisibleObjects[p.UniqueID] {
//		packet := network.EmptyPacket()
//		packet.MessageID = opcode.EntityMovementResponse
//		packet.WriteUInt32(p.UniqueID)
//		packet.WriteByte(1)
//		packet.WriteUInt16(uint16(targetPos.Region.ID))
//		packet.WriteUInt16(uint16(targetPos.X) + 0xFFFF)
//		packet.WriteUInt16(uint16(targetPos.Y))
//		packet.WriteUInt16(uint16(targetPos.Z) + 0xFFFF)
//		packet.WriteByte(1)
//		packet.WriteUInt16(uint16(p.Position.Region.ID))
//		packet.WriteUInt16(uint16(p.Position.X) * 10)
//		packet.WriteFloat32(p.Position.Y)
//		packet.WriteUInt16(uint16(p.Position.Z) * 10)
//		player := world.PlayersByUniqueId[uid]
//		player.Session.Conn.Write(packet.ToBytes())
//	}
//}
//
//func (s *SpawnEngine) PlayerObjectCollision(p *model.Player) {
//	world := model.GetSroWorldInstance()
//	packet := network.EmptyPacket()
//	packet.MessageID = opcode.MovementPositionUpdate
//	packet.WriteUInt32(p.UniqueID)
//	packet.WriteUInt16(uint16(p.Position.Region.ID))
//	packet.WriteFloat32(p.Position.X)
//	packet.WriteFloat32(p.Position.Y)
//	packet.WriteFloat32(p.Position.Z)
//	packet.WriteUInt16(uint16(p.Position.Heading))
//
//	for uid, _ := range s.VisibleObjects[p.UniqueID] {
//		player := world.PlayersByUniqueId[uid]
//		player.Session.Conn.Write(packet.ToBytes())
//	}
//}
//
//func (s *SpawnEngine) SpawnNPC() {
//	//logrus.Debug("Spawning an NPC")
//	//
//	//p := network.EmptyPacket()
//	//p.MessageID = opcode.SingleSpawn
//
//	// 1. Add to visible objects
//	// 2. Fetch players in visible range
//	// 3. Notify those players about spawned NPC
//
//	//TODO: if len > 1  use opcode.GroupSpawnBegin, opcode.GroupSpawnData, opcode.GroupSpawnEnd
//}
//
//func (s *SpawnEngine) UpdatedMovementSpeed(p *model.Player) {
//	world := model.GetSroWorldInstance()
//	packet := network.EmptyPacket()
//	packet.MessageID = opcode.EntityUpdateMovementSpeed
//	packet.WriteUInt32(p.UniqueID)
//	packet.WriteFloat32(p.GetWalkSpeed()) // WalkSpeed
//	packet.WriteFloat32(p.GetRunSpeed())
//
//	for uid, _ := range s.VisibleObjects[p.UniqueID] {
//		player := world.PlayersByUniqueId[uid]
//		player.Session.Conn.Write(packet.ToBytes())
//	}
//}
//
//func (s *SpawnEngine) PlayerDisconnected(uid uint32, charName string) {
//	world := model.GetSroWorldInstance()
//
//	p1 := network.EmptyPacket()
//	p1.MessageID = opcode.EntityGroupSpawnBegin
//	p1.WriteByte(2)
//	p1.WriteUInt16(1)
//
//	p2 := network.EmptyPacket()
//	p2.MessageID = opcode.EntityGroupSpawnData
//	p2.WriteUInt32(uid)
//
//	p3 := network.EmptyPacket()
//	p3.MessageID = opcode.EntityGroupSpawnEnd
//
//	for uid2, _ := range s.VisibleObjects[uid] {
//		player := world.PlayersByUniqueId[uid2]
//		player.Session.Conn.Write(p1.ToBytes())
//		player.Session.Conn.Write(p2.ToBytes())
//		player.Session.Conn.Write(p3.ToBytes())
//		delete(s.VisibleObjects[player.UniqueID], uid)
//	}
//
//	delete(s.VisibleObjects, uid)
//	world.PlayerDisconnected(uid, charName)
//}
//
///*
//4   uint    RefObjID
//if(obj.TypeID1 == 1)
//{
//    //BIONIC:
//    //  CHARACTER
//    //  NPC
//    //      NPC_FORTRESS_STRUCT
//    //      NPC_MOB
//    //      NPC_COS
//    //      NPC_FORTRESS_COS
//
//    if(obj.TypeID2 == 1)
//    {
//        //CHARACTER
//        1   byte    Scale
//        1   byte    HwanLevel
//        1   byte    PVPCape         //0 = None, 1 = Red, 2 = Gray, 3 = Blue, 4 = White, 5 = Orange
//        1   byte    AutoInverstExp  //1 = Beginner Icon, 2 = Helpful, 3 = Beginner & Helpful
//
//        //Inventory
//        1   byte    Inventory.Size
//        1   byte    Inventory.ItemCount
//        for (int i = 0; i < Inventory.ItemCount; i++)
//        {
//            4   uint    item.RefItemID
//            if(item.TypeID1 == 3 && item.TypeID2 == 1)
//            {
//                1   byte    item.OptLevel
//            }
//        }
//        //AvatarInventory
//        1   byte    Inventory.Size
//        1   byte    Inventory.ItemCount
//        for (int i = 0; i < Inventory.ItemCount; i++)
//        {
//            4   uint    item.RefItemID
//            if(item.TypeID1 == 3 && item.TypeID2 == 1)
//            {
//                1   byte    item.OptLevel
//            }
//        }
//
//        //Mask
//        1   byte    HasMask
//        if(HasMask)
//        {
//            4   uint    mask.RefObjID
//            if (maskObject.TypeID1 == obj.TypeID1 &&
//                maskObject.TypeID2 == obj.TypeID2)
//            {
//                //Duplicate
//                1   byte    Mask.Scale
//                1   byte    Mask.ItemCount
//                for(int i = 0; i < Mask.ItemCount; i++)
//                {
//                    4   uint    item.RefItemID
//                }
//            }
//        }
//    }
//    else if(obj.TypeID2 == 2 && obj.TypeID3 == 5)
//    {
//        //NPC_FORTRESS_STRUCT
//        4   uint    HP
//        4   uint    RefEventStructID
//        2   ushort  State
//    }
//    4   uint    UniqueID
//
//    //Position
//    2   ushort  Position.RegionID
//    4   float   Position.X
//    4   float   Position.Y
//    4   float   Position.Z
//    2   ushort  Position.Angle
//
//    //Movement
//    1   byte    Movement.HasDestination
//    1   byte    Movement.Type
//    if(Movement.HasDestination)
//    {
//        //MD (Mouse destination)
//        2   ushort  Movement.DestinationRegion
//        if(Position.RegionID < short.MaxValue)
//        {
//            //World
//            2   ushort  Movement.DestinationOffsetX
//            2   ushort  Movement.DestinationOffsetY
//            2   ushort  Movement.DestinationOffsetZ
//        }
//        else
//        {
//            //Dungeon
//            4   uint  Movement.DestinationOffsetX
//            4   uint  Movement.DestinationOffsetY
//            4   uint  Movement.DestinationOffsetZ
//        }
//    }
//    else
//    {
//        1   byte    Movement.Source  //0 = Spinning, 1 = Sky-/Key-walking
//        2   ushort  Movement.Angle   //Represents the new angle, character is looking at
//        //FC_GO_FORWARD
//        //
//    }
//
//    //State
//    1   byte    State.LifeState     //1 = Alive, 2 = Dead
//    1   byte    State.unkByte0
//    1   byte    State.MotionState   //0 = None, 2 = Walking, 3 = Running, 4 = Sitting
//    1   byte    State.Status        //0 = None, 1 = Hwan, 2 = Untouchable, 3 = GameMasterInvincible, 5 = GameMasterInvisible, 5 = ?, 6 = Stealth, 7 = Invisible
//    4   float   State.WalkSpeed
//    4   float   State.RunSpeed
//    4   float   State.HwanSpeed
//    1   byte    State.BuffCount
//    for (int i = 0; i < State.BuffCount; i++)
//    {
//        4   uint    Buff.RefSkillID
//        4   uint    Buff.Duration
//        if(skill.Params.Contains(1701213281))
//        {
//            //1701213281 -> atfe -> "auto transfer effect" like Recovery Division
//            1   bool    IsCreator
//        }
//    }
//
//    if(obj.TypeID2 == 1)
//    {
//        //CHARACTER
//        2   ushort  Name.Length
//        2   string  Name
//
//        1   byte    JobType            //0 = None, 1 = Trader, 2 = Tief, 3 = Hunter
//        1   byte    JobLevel
//        1   byte    PVPState         //0 = White (Neutral), 1 = Purple (Assaulter), 2 = Red
//        1   byte    TransportFlag
//        1   byte    InCombat
//        if(TransportFlag)
//        {
//            4   uint    Transport.UniqueID
//        }
//        1   byte    ScrollMode         //0 = None, 1 = Return Scroll, 2 = Bandit Return Scroll
//        1   byte    InteractMode       //0 = None 2 = P2P, 4 = P2N_TALK, 6 = OPNMKT_DEAL
//        1   byte    unkByte4
//
//        //Guild
//        2   ushort  Guild.Name.Length
//        *   string  Guild.Name
//        if(Inventory.ContainsJobEquipment == false)
//        {
//            4   uint    Guild.ID
//            2   ushort  GuildMember.Nickname.Length
//            *   string  GuildMember.Nickname
//            4   uint    Guild.LastCrestRev
//            4   uint    Union.ID
//            4   uint    Union.LastCrestRev
//            1   byte    Guild.IsFriendly            //0 = Hostile, 1 = Friendly
//            1   byte    GuildMember.SiegeAuthority  //See SiegeAuthority.cs
//        }
//
//        if(InteractMode == 4) //P2N_TALK
//        {
//            2   ushort  Stall.Name.Length
//            *   string  Stall.Name
//            4   uint    Stall.RefItemID         //Decoration
//        }
//        1   byte    EquipmentCooldown
//        1   byte    PKFlag
//    }
//    else if(obj.TypeID2 == 2)
//    {
//        //NPC
//        1   byte    TalkFlag
//        if(TalkFlag == 2)
//        {
//            1   byte    TalkOptionCount
//            *   byte[]  TalkOptions
//        }
//
//        if(obj.TypeID3 == 2)
//        {
//            //NPC_MOB
//            1   byte    Rarity
//            if(obj.TypeID4 == 2 || obj.TypeID4 == 3)
//            {
//                1   byte    Appearance  //Randomized by server.
//            }
//        }
//        else if(obj.TypeID3 == 3)
//        {
//            //NPC_COS
//            if(obj.TypeID4 == 3 || obj.TypeID4 == 4)
//            {
//                //NPC_COS_P (Growth)
//                //NPC_COS_P (Ability)
//                2   ushort  Name.Length
//                2   string  Name
//            }
//
//            if(obj.TypeID4 == 5)
//            {
//                //NPC_COS_GUILD
//                2   ushort  GuildName.Length
//                2   string  GuildName
//            }
//            else
//            {
//                2   ushort  Owner.Name.Length
//                2   string  Owner.Name
//            }
//
//            if(obj.TypeID4 == 2 ||      //NPC_COS_T
//               obj.TypeID4 == 3 ||      //NPC_COS_P (Growth)
//               obj.TypeID4 == 4 ||      //NPC_COS_P (Ability)
//               obj.TypeID4 == 5)        //NPC_COS_GUILD
//            {
//               1    byte    JobType
//               if(obj.TypeID4 != 4) //NO NPC_COS_P (Ability)
//               {
//                   1    byte    MurderFlag    //0 = White, 1 = Purple, 2 = Red
//               }
//
//               if(obj.TypeID4 == 5)
//               {
//                    //NPC_COS_GUILD
//                   4    uint    Owner.RefObjID
//               }
//            }
//
//            4    uint    Owner.UniqueID
//        }
//        else if (obj.TypeID3 == 4)
//        {
//            //NPC_FORTRESS_COS
//            4   uint    Guild.ID
//            2   ushort  Guild.Name.Length
//            *   string  Guild.Name
//        }
//    }
//}
//else if(obj.TypeID1 == 3)
//{
//    //ITEM:
//    //  ITEM_EQUIP
//    //  ITEM_ETC
//    //      ITEM_ETC_MONEY_GOLD
//    //      ITEM_ETC_TRADE
//    //      ITEM_ETC_QUEST
//
//    if(obj.TypeID2 == 1)
//    {
//        //ITEM_EQUIP
//        1   byte    OptLevel
//    }
//    else if(obj.TypeID2 == 3)
//    {
//        //ITEM_ETC
//        if(obj.TypeID3 == 5 && obj.TypeID4 == 0)
//        {
//            //ITEM_ETC_MONEY_GOLD
//            4   uint    Amount
//        }
//        else if(obj.TypeID3 == 8 || obj.TypeID3 == 9)
//        {
//            //ITEM_ETC_TRADE
//            //ITEM_ETC_QUEST
//            4   ushort  Owner.Name.Length
//            *   string  Owner.Name
//        }
//    }
//    4   uint    UniqueID
//    2   ushort  Position.RegionID
//    4   float   Position.X
//    4   float   Position.Y
//    4   float   Position.Z
//    2   ushort  Position.Angle
//    1   bool    hasOwner
//    if(hasOwner)
//    {
//        4   uint    Owner.JID
//    }
//    1   byte    Rarity
//}
//else if(obj.TypeID1 == 4)
//{
//    //PORTALS:
//    //  STORE
//    //  INS_TELEPORTER
//
//    4   uint    UniqueID
//    2   ushort  Position.RegionID
//    4   float   Position.X
//    4   float   Position.Y
//    4   float   Position.Z
//    2   ushort  Position.Angle
//
//    1   byte    unkByte0
//    1   byte    unkByte1
//    1   byte    unkByte2
//    1   byte    unkByte3
//    if(unkByte3 == 1)
//    {
//        //Regular
//        4   uint    unkUInt0
//        4   uint    unkUInt1
//    }
//    else if(unkByte3 == 6)
//    {
//        //Dimension Hole
//        2   ushort  Owner.Name.Length
//        *   string  Owner.Name
//        4   uint    Owner.UniqueID
//    }
//
//    if(unkByte1 == 1)
//    {
//        //STORE_EVENTZONE_DEFAULT
//        4   uint    unkUInt3
//        1   byte    unkByte5
//    }
//}
//else if(RefObjID == uint.MaxValue)
//{
//    //EVENT_ZONE (Traps, Buffzones, ...)
//    2   ushort  unkUShort0
//    4   uint    RefSkillID
//    4   uint    UniqueID
//    2   ushort  Position.RegionID
//    4   float   Position.X
//    4   float   Position.Y
//    4   float   Position.Z
//    2   ushort  Position.Angle
//}
//
//if(packet.Opcode == 0x3015)
//{
//    if(obj.TypeID1 == 0x01 || obj.TypeID1 == 0x4)
//    {
//        //BIONIC and STORE
//        1    byte    unkByte6
//    }
//    else if(obj.TypeID1 == 0x03)
//    {
//        1    byte    DropSource
//        4    uint    Dropper.UniqueID
//    }
//}
////EOP
//*/
