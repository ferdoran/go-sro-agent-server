package navmesh

import (
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Resource struct {
	OffsetMaterial           uint32
	OffsetMesh               uint32
	OffsetPrimBranch         uint32
	OffsetAnimation          uint32
	OffsetPrimMeshGroup      uint32
	OffsetPrimAnimationGroup uint32
	OffsetModPalette         uint32
	OffsetNavMeshObj         uint32
	Int0                     uint32
	Int1                     uint32
	Int2                     uint32
	Int3                     uint32
	Int4                     uint32
	NavMeshObjPath           string
}

func LoadResource(filename string, reader *pk2.Pk2Reader) *Resource {
	//file, err := os.Open(filename)
	filename = strings.ReplaceAll(filename, "\\", string(os.PathSeparator))
	filename = strings.ReplaceAll(filename, "/", string(os.PathSeparator))
	fileContent, err := reader.ReadFile(filename)
	if err != nil {
		logrus.Panic(err)
	}
	header := string(fileContent[:12])

	if header != "JMXVRES 0109" {
		logrus.Panicf("Invalid signature %s\n", header)
	}

	res := Resource{
		OffsetMaterial:           utils.ByteArrayToUint32(fileContent[12:16]),
		OffsetMesh:               utils.ByteArrayToUint32(fileContent[16:20]),
		OffsetPrimBranch:         utils.ByteArrayToUint32(fileContent[20:24]),
		OffsetAnimation:          utils.ByteArrayToUint32(fileContent[24:28]),
		OffsetPrimMeshGroup:      utils.ByteArrayToUint32(fileContent[28:32]),
		OffsetPrimAnimationGroup: utils.ByteArrayToUint32(fileContent[32:36]),
		OffsetModPalette:         utils.ByteArrayToUint32(fileContent[36:40]),
		OffsetNavMeshObj:         utils.ByteArrayToUint32(fileContent[40:44]),
		Int0:                     utils.ByteArrayToUint32(fileContent[44:48]),
		Int1:                     utils.ByteArrayToUint32(fileContent[48:52]),
		Int2:                     utils.ByteArrayToUint32(fileContent[52:56]),
		Int3:                     utils.ByteArrayToUint32(fileContent[56:60]),
		Int4:                     utils.ByteArrayToUint32(fileContent[60:64]),
	}

	strLen := utils.ByteArrayToUint32(fileContent[res.OffsetNavMeshObj : res.OffsetNavMeshObj+4])
	str := string(fileContent[res.OffsetNavMeshObj+4 : res.OffsetNavMeshObj+4+strLen])
	str = strings.ReplaceAll(str, "\\", string(os.PathSeparator))
	str = strings.ReplaceAll(str, "/", string(os.PathSeparator))

	res.NavMeshObjPath = str

	return &res
}
