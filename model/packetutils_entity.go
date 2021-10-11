package model

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/sirupsen/logrus"
)

func WriteEntitySpawnData(p *network.Packet, object ISRObject) {
	p.WriteUInt32(object.GetRefObjectID())

	typeInfo := object.GetTypeInfo()

	// Load RefData
	if typeInfo.IsCharacter() {
		// Bionic
		if typeInfo.IsPlayerCharacter() {
			// Player Character
			player, isPlayer := object.(IPlayer)

			if !isPlayer {
				logrus.Warnf("object %s is not a player. Type: %T", object.GetName(), object)
				return
			}

			p.WriteByte(byte(player.GetScale()))
			p.WriteByte(0) // HwanLevel
			p.WriteByte(byte(player.GetPVPCape()))
			p.WriteByte(1) // AutoInverstXP
			p.WriteByte(BaseInventorySize)
			p.WriteByte(byte(len(player.GetInventory().Items)))

			for _, item := range player.GetInventory().Items {
				p.WriteUInt32(item.GetRefObjectID())

				if item.IsEquipment() {
					// TODO: write real opt level or just always 0?
					p.WriteByte(0) // OptLevel
				}
			}

			p.WriteByte(5) // Avatar Inventory Slots
			p.WriteByte(0) // Avatar Items

			p.WriteByte(0) // HasMask

		} else if typeInfo.IsSiegeStruct() {
			// NPC_FORTRESS_STRUCT
		}
		p.WriteUInt32(object.GetUniqueID())
		WritePosition(p, object.GetNavmeshPosition())

		// TODO WriteMovementData
		if movementData := object.GetMovementData(); movementData != nil {
			p.WriteBool(movementData.HasDestination)
			p.WriteByte(1)

			if movementData.HasDestination {
				p.WriteUInt16(uint16(movementData.TargetPosition.Region.ID))
				if movementData.TargetPosition.Region.ID > 0 {
					p.WriteUInt16(uint16(movementData.TargetPosition.Offset.X))
					p.WriteUInt16(uint16(movementData.TargetPosition.Offset.Y))
					p.WriteUInt16(uint16(movementData.TargetPosition.Offset.Z))
				} else {
					// TODO Dungeon
				}
			} else {
				p.WriteByte(1) // 0 = Spinning, 1 = Sky-/Key-Walking
				p.WriteUInt16(uint16(movementData.TargetPosition.Heading))
			}
		} else {
			p.WriteByte(0)
			p.WriteByte(1)
			p.WriteByte(1)
			p.WriteUInt16(uint16(object.GetNavmeshPosition().Heading))
		}
		character := object.(ICharacter)

		p.WriteByte(byte(character.GetLifeState()))
		p.WriteByte(0) // Unknown
		p.WriteByte(byte(character.GetMotionState()))
		p.WriteByte(byte(character.GetBodyState()))
		p.WriteFloat32(character.GetWalkSpeed())
		p.WriteFloat32(character.GetRunSpeed())
		p.WriteFloat32(character.GetHwanSpeed())

		p.WriteByte(0) // TODO BuffCount
		//	for (int i = 0; i < State.BuffCount; i++)
		//	{
		//	4   uint    Buff.RefSkillID
		//	4   uint    Buff.Duration
		//	if(skill.Params.Contains(1701213281))
		//	{
		//	//1701213281 -> atfe -> "auto transfer effect" like Recovery Division
		//	1   bool    IsCreator
		//	}
		//	}

		if typeInfo.IsPlayerCharacter() {
			player := character.(IPlayer)
			p.WriteString(player.GetName())
			p.WriteByte(0) // TODO JobType
			p.WriteByte(1) // TODO JobLevel
			p.WriteByte(byte(player.GetPKState()))
			p.WriteByte(0) // TODO Transport Flag
			p.WriteByte(0) // TODO in combat?
			// if transportFlag {
			//     p.WriteUInt32(transport.UniqueID)
			// }
			p.WriteByte(0) // TODO ScrollMode / 0 = None / 1 = Return Scroll / 2 = Bandit Return Scroll
			p.WriteByte(0) // TODO Interact Mode / 0 = None / 2 = P2P / 4 = P2N_Talk / 6 = OPNMKT_DEAL
			p.WriteByte(0) // unknown

			WriteGuild(p, player)

			p.WriteByte(0)    // Equipment Cooldown, probably when equipping pvp capes or job eq
			p.WriteByte(0xFF) // PK Flag
		} else if typeInfo.IsNPC() {
			p.WriteByte(0) // TODO TalkFlag

			if typeInfo.IsNPCMob() {
				p.WriteByte(1) // TODO Rarity
			}
			//	//NPC
			//	1   byte    TalkFlag
			//	if(TalkFlag == 2)
			//	{
			//		1   byte    TalkOptionCount
			//		*   byte[]  TalkOptions
			//	}
			//
			//	if(obj.TypeID3 == 1)
			//	{
			//		//NPC_MOB
			//		1   byte    Rarity
			//		if(obj.TypeID4 == 2 || obj.TypeID4 == 3)
			//		{
			//			// NPC_MOB_THIEF, NPC_MOB_HUNTER
			//			1   byte    Appearance  //Randomized by server.
			//		}
			//	}
			//	else if(obj.TypeID3 == 3)
			//	{
			//		//NPC_COS
			//		if(obj.TypeID4 == 3 || obj.TypeID4 == 4)
			//		{
			//			//NPC_COS_P (Growth)
			//			//NPC_COS_P (Ability)
			//			2   ushort  Name.Length
			//			2   string  Name
			//		}
			//
			//		if(obj.TypeID4 == 5)
			//		{
			//			//NPC_COS_GUILD
			//			2   ushort  GuildName.Length
			//			2   string  GuildName
			//		}
			//		else
			//		{
			//			2   ushort  Owner.Name.Length
			//			2   string  Owner.Name
			//		}
			//
			//		if(obj.TypeID4 == 2 ||      //NPC_COS_T
			//			obj.TypeID4 == 3 ||      //NPC_COS_P (Growth)
			//			obj.TypeID4 == 4 ||      //NPC_COS_P (Ability)
			//			obj.TypeID4 == 5)        //NPC_COS_GUILD
			//		{
			//			1    byte    JobType
			//			if(obj.TypeID4 != 4) //NO NPC_COS_P (Ability)
			//			{
			//				1    byte    MurderFlag    //0 = White, 1 = Purple, 2 = Red
			//			}
			//
			//			if(obj.TypeID4 == 5)
			//			{
			//				//NPC_COS_GUILD
			//				4    uint    Owner.RefObjID
			//			}
			//		}
			//
			//		4    uint    Owner.UniqueID
			//	}
			//	else if (obj.TypeID3 == 4)
			//	{
			//		//NPC_FORTRESS_COS
			//		4   uint    Guild.ID
			//		2   ushort  Guild.Name.Length
			//		*   string  Guild.Name
			//	}
			//}
		}

	} else if typeInfo.IsItem() {
		// Item
		item := object.(IItem)
		if item.GetTypeInfo().IsEquipment() {
			eq := item.(IEquipment)
			p.WriteByte(eq.GetOptLevel())
		} else if typeInfo.IsExpendable() {
			if typeInfo.IsGold() {
				p.WriteUInt32(1_000_000) // FIXME correct amount
			} else if typeInfo.IsTradeItem() || typeInfo.IsQuestItem() {
				p.WriteString("Owner") // FIXME correct owner name
			}
		}

		p.WriteUInt32(object.GetUniqueID())
		WritePosition(p, object.GetNavmeshPosition())

		if ownerJID := item.GetOwner(); ownerJID > 0 {
			p.WriteUInt32(item.GetOwner())
		}
		p.WriteByte(item.GetRarity())
	} else if typeInfo.IsStructure() {
		p.WriteUInt32(object.GetUniqueID())
		WritePosition(p, object.GetNavmeshPosition())

		p.WriteByte(0) // TODO unkByte0
		p.WriteByte(0) // TODO unkByte1
		p.WriteByte(0) // TODO unkByte2
		p.WriteByte(0) // TODO unkByte3

		// if unkByte3 == 1 {
		//     p.WriteUint32(0) // TODO unkUint0
		//     p.WriteUint32(0) // TODO unkUint1
		// } else if unkByte3 == 6 {
		//     // Dimension Hole
		//     p.WriteString("owner")
		//     p.WriteUint32(owner.UniqueID)
		// }
		//
		// if unkByte1 == 1 {
		//     p.WriteUint32(0) // TODO unkUint2
		// 	   p.WriteByte(0)	// TODO unkByte4
		// }
	} else if object.GetRefObjectID() == 0xFFFFFFFF {
		p.WriteUInt16(0) // TODO unkUint160
		p.WriteUInt32(0) // TODO RefSkillID
		p.WriteUInt32(object.GetUniqueID())
		WritePosition(p, object.GetNavmeshPosition())
	}

	if p.MessageID == opcode.EntitySingleSpawn {
		if typeInfo.IsCharacter() || typeInfo.IsStructure() {
			p.WriteByte(0) // TODO unkByte5
		} else if typeInfo.IsItem() {
			p.WriteByte(0)   // TODO DropSource
			p.WriteUInt32(0) // TODO Dropper.UniqueID
		}
	}
}

func WriteEntitySpawn(p *network.Packet, player *Player) {
	// TODO not just players supported
	p.WriteUInt32(uint32(player.GetRefObjectID()))

	p.WriteByte(byte(player.Scale))
	p.WriteByte(0) // HwanLevel
	p.WriteByte(byte(player.PVPCape))
	p.WriteByte(1) // AutoInverstXP
	p.WriteByte(BaseInventorySize)
	p.WriteByte(byte(len(player.Inventory.Items)))

	for _, item := range player.Inventory.Items {
		p.WriteUInt32(item.GetRefObjectID())

		if item.IsEquipment() {
			// TODO: write real opt level or just always 0?
			p.WriteByte(0) // OptLevel
		}
	}

	p.WriteByte(5) // Avatar Inventory Slots
	p.WriteByte(0) // Avatar Items

	p.WriteByte(0) // HasMask
	p.WriteUInt32(player.UniqueID)
	WritePosition(p, player.GetNavmeshPosition())

	// TODO movement still does not save a target location
	p.WriteByte(0)                                             // Movement.HasDestination
	p.WriteByte(1)                                             // Movement.Type
	p.WriteByte(1)                                             // Movement Source | 0 - Spinning | 1 - Skywalking
	p.WriteUInt16(uint16(player.GetNavmeshPosition().Heading)) // Angle

	p.WriteByte(byte(player.LifeState))
	p.WriteByte(0) // unknown
	p.WriteByte(byte(player.MotionState))
	p.WriteByte(byte(player.BodyState))
	p.WriteFloat32(player.GetWalkSpeed()) // WalkSpeed
	p.WriteFloat32(player.GetRunSpeed())
	p.WriteFloat32(player.GetHwanSpeed()) // HwanSpeed TODO: is this correct?
	p.WriteByte(0)                        // TODO BuffCount

	p.WriteString(player.GetName())
	p.WriteByte(0) // TODO JobType
	p.WriteByte(1) // TODO JobLevel
	p.WriteByte(byte(player.GetPKState()))
	p.WriteByte(0) // TODO Transport Flag
	p.WriteByte(0) // TODO in combat?
	p.WriteByte(0) // TODO ScrollMode / 0 = None / 1 = Return Scroll / 2 = Bandit Return Scroll
	p.WriteByte(0) // TODO Interact Mode / 0 = None / 2 = P2P / 4 = P2N_Talk / 6 = OPNMKT_DEAL
	p.WriteByte(0) // unknown

	WriteGuild(p, player)

	p.WriteByte(0)    // Equipment Cooldown, probably when equipping pvp capes or job eq
	p.WriteByte(0xFF) // PK Flag

	// TODO:
	//  - Check if inventory contains job equipment
	//  - Check interact mode os P2N_Talk
	//  - If player has transport(camel, horse, etc.) send transport's unique id

	if p.MessageID == opcode.EntitySingleSpawn {
		p.WriteByte(0) // TODO find out what this is
	}
}

func WritePosition(p *network.Packet, position navmeshv2.RtNavmeshPosition) {
	p.WriteUInt16(uint16(position.Region.ID))
	p.WriteFloat32(position.Offset.X)
	p.WriteFloat32(position.Offset.Y)
	p.WriteFloat32(position.Offset.Z)
	p.WriteUInt16(uint16(position.Heading))
}

func WriteGuild(p *network.Packet, player IPlayer) {
	// todo
	p.WriteUInt16(0) // GuildNameLength

	p.WriteUInt32(0) // Guild ID
	p.WriteUInt16(0) // Guild Member Nickname length
	p.WriteUInt32(0) // Guild last crest rev
	p.WriteUInt32(0) // Union ID
	p.WriteUInt32(0) // Union last crest rev
	p.WriteByte(0)   // Guild IsFriendly
	p.WriteByte(0)   // GuildMember.SiegeAuthority
}
