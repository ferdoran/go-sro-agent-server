package navmeshv2

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"strconv"
)

const (
	DungeonInfoEntryPattern    = `(?P<service>0|1)\t(?P<id>\d*)\t\"(?P<path>.*)\"`
	DungeonInfoFilePathDataPk2 = "Data" + string(os.PathSeparator) + "dungeon" + string(os.PathSeparator) + "dungeoninfo.txt"
)

type DungeonInfo struct {
	Dungeons []DungeonInfoEntry
}

type DungeonInfoEntry struct {
	Service  int
	ID       int
	FilePath string
}

var entryPattern = regexp.MustCompile(DungeonInfoEntryPattern)
var entrySubexpNames = entryPattern.SubexpNames()

func LoadDungeonInfo(reader *pk2.Pk2Reader) DungeonInfo {
	objectInfoBytes, err := reader.ReadFile(DungeonInfoFilePathDataPk2)
	if err != nil {
		logrus.Panicf("error loading file %s. Error = %v", DungeonInfoFilePathDataPk2, err)
	}
	fileReader := bufio.NewReader(bytes.NewReader(objectInfoBytes))
	line, _, err := fileReader.ReadLine()
	if err != nil {
		logrus.Panic(err)
	}

	if len(line) < 1 {
		logrus.Panicf("line is empty: %v", line)
	}
	entries := make([]DungeonInfoEntry, 0)
	for {
		lineBytes, _, err := fileReader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Panic(err)
		}

		line := string(lineBytes)
		matches := entryPattern.FindStringSubmatch(line)
		matchMap := make(map[string]string)
		for i, v := range matches {
			matchMap[entrySubexpNames[i]] = v
		}

		service := matchMap["service"]
		id := matchMap["id"]
		path := matchMap["path"]

		serviceInt, err := strconv.Atoi(service)
		if err != nil {
			logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to parse dungeon service as int. got %x", lineBytes)))
			logrus.Panicf("failed to read dungeon info")
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to parse dungeon id as int. got %s", service)))
			logrus.Panicf("failed to read dungeon info")
		}

		entry := DungeonInfoEntry{
			Service:  serviceInt,
			ID:       idInt,
			FilePath: path,
		}

		entries = append(entries, entry)
	}

	return DungeonInfo{Dungeons: entries}
}
