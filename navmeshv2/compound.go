package navmeshv2

import (
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Compound struct {
	OffsetNavMeshObj uint32
	OffsetResObjList uint32
	Int0             uint32
	Int1             uint32
	Int2             uint32
	Int3             uint32
	Int4             uint32
	NavMeshObjPath   string
}

func LoadCompoundFile(filename string, reader *pk2.Pk2Reader) *Compound {
	filename = strings.ReplaceAll(filename, "\\", string(os.PathSeparator))
	filename = strings.ReplaceAll(filename, "/", string(os.PathSeparator))
	fileContent, err := reader.ReadFile(filename)
	if err != nil {
		logrus.Panic(err)
	}

	header := string(fileContent[:12])

	if header != "JMXVCPD 0101" {
		logrus.Panicf("Header did not start with JMXVCPD 0101. Got %s\n", header)
	}

	cpd := Compound{
		OffsetNavMeshObj: utils.ByteArrayToUint32(fileContent[12:16]),
		OffsetResObjList: utils.ByteArrayToUint32(fileContent[16:20]),
		Int0:             utils.ByteArrayToUint32(fileContent[20:24]),
		Int1:             utils.ByteArrayToUint32(fileContent[24:28]),
		Int2:             utils.ByteArrayToUint32(fileContent[28:32]),
		Int3:             utils.ByteArrayToUint32(fileContent[32:36]),
		Int4:             utils.ByteArrayToUint32(fileContent[36:40]),
	}

	strLen := utils.ByteArrayToUint32(fileContent[cpd.OffsetNavMeshObj : cpd.OffsetNavMeshObj+4])
	str := string(fileContent[cpd.OffsetNavMeshObj+4 : cpd.OffsetNavMeshObj+4+strLen])
	str = strings.ReplaceAll(str, "\\", string(os.PathSeparator))
	str = strings.ReplaceAll(str, "/", string(os.PathSeparator))

	cpd.NavMeshObjPath = str

	return &cpd
}
