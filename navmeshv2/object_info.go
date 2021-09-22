package navmeshv2

import (
	"bufio"
	"bytes"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	objectIfoHeader       = "JMXVOBJI1000"
	ObjectFilePathDataPk2 = "Data" + string(os.PathSeparator) + "navmesh" + string(os.PathSeparator) + "object.ifo"
)

type ObjectInfoEntry struct {
	Index    uint32
	Flag     uint32
	FilePath string
}

type ObjectInfo struct {
	Objects map[uint32]ObjectInfoEntry
}

func LoadObjectInfo(reader *pk2.Pk2Reader) ObjectInfo {
	objectInfoBytes, err := reader.ReadFile(ObjectFilePathDataPk2)
	if err != nil {
		logrus.Panicf("error loading file %s. Error = %v", MapFilePathDataPk2, err)
	}
	fileReader := bufio.NewReader(bytes.NewReader(objectInfoBytes))
	line, _, err := fileReader.ReadLine()

	if string(line) != objectIfoHeader {
		logrus.Panicf("File does not start with %s. Got = %s", objectIfoHeader, line)
	}

	entryCountBuffer, _, err := fileReader.ReadLine()
	if err != nil {
		logrus.Panic(err)
	}
	entryCount, _ := strconv.ParseInt(string(entryCountBuffer), 10, 64)
	logrus.Printf("Found %d entries", entryCount)
	objects := make(map[uint32]ObjectInfoEntry, entryCount)
	for i := 0; i < int(entryCount); i++ {
		line, _, err := fileReader.ReadLine()
		if err == io.EOF {
			logrus.Errorf("EOF on line %d", i)
			logrus.Panic(err)
		}

		index, err := strconv.ParseUint(string(line[:5]), 10, 64)
		if err != nil {
			logrus.Error(err)
			logrus.Panic("Failed to read index on entry %d", i)
		}
		flag, err := strconv.ParseInt(string(line[6:16]), 0, 32)
		if err != nil {
			logrus.Error(err)
			logrus.Panicf("Failed to read flag[%s] on entry %d", flag, i)
		}
		filepath := strings.ReplaceAll(string(line[18:]), "\"", "")

		objects[uint32(index)] = ObjectInfoEntry{
			Index:    uint32(index),
			Flag:     uint32(flag),
			FilePath: filepath,
		}
	}

	return ObjectInfo{Objects: objects}
}
