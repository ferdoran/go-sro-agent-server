package navmeshv2

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/engine/geo/math"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

type Loader struct {
	Pk2Reader      *pk2.Pk2Reader
	DataPk2Path    string
	NavMeshPath    string
	MapProjectInfo MapProjectInfo
	ObjectInfo     ObjectInfo
	DungeonInfo    DungeonInfo
	ObjectData     map[uint32]RtNavmeshObj
	RegionData     map[int16]RtNavmeshTerrain
}

func NewLoader(dataPk2Path string) *Loader {

	reader := pk2.NewPk2Reader(dataPk2Path)
	reader.IndexArchive()
	return &Loader{
		Pk2Reader:      &reader,
		DataPk2Path:    "Data",
		NavMeshPath:    "Data" + string(os.PathSeparator) + "navmesh",
		MapProjectInfo: MapProjectInfo{},
		ObjectInfo:     ObjectInfo{},
		DungeonInfo:    DungeonInfo{},
		ObjectData:     make(map[uint32]RtNavmeshObj),
		RegionData:     make(map[int16]RtNavmeshTerrain),
	}
}

func (l *Loader) LoadNavMeshInfos() {
	l.MapProjectInfo = LoadMapProjectInfo(l.Pk2Reader)
	l.ObjectInfo = LoadObjectInfo(l.Pk2Reader)
	//l.DungeonInfo = LoadDungeonInfo(l.Pk2Reader)
}

func (l *Loader) LoadObjectMeshes() {
	for _, objectEntry := range l.ObjectInfo.Objects {
		switch path.Ext(objectEntry.FilePath) {
		case ".bms":
			l.LoadObjectFromPrimMesh(objectEntry)
		case ".bsr":
			l.LoadObjectFromResource(objectEntry)
		case ".cpd":
			l.LoadObjectFromCompound(objectEntry)
		default:
			logrus.Panicf("file is not an object mesh: %s", objectEntry.FilePath)
		}
	}
}

func (l *Loader) LoadObjectMesh(objectIndex uint32) RtNavmeshObj {
	if object, exists := l.ObjectData[objectIndex]; exists {
		return object
	}
	objectEntry := l.ObjectInfo.Objects[objectIndex]
	switch path.Ext(objectEntry.FilePath) {
	case ".bms":
		l.LoadObjectFromPrimMesh(objectEntry)
	case ".bsr":
		l.LoadObjectFromResource(objectEntry)
	case ".cpd":
		l.LoadObjectFromCompound(objectEntry)
	default:
		logrus.Panicf("file is not an object mesh: %s", objectEntry.FilePath)
	}

	return l.ObjectData[objectIndex]
}

func (l *Loader) LoadTerrainMeshes() {
	counter := 0
	for regionId, enabled := range l.MapProjectInfo.EnabledRegions {
		if !enabled {
			continue
		}
		regionShortHex := fmt.Sprintf("%x", regionId)
		navMeshHex := fmt.Sprintf("nv_%s.nvm", regionShortHex)
		counter++
		fmt.Printf("\rReading %s. Finished [%d / %d] files", navMeshHex, counter, l.MapProjectInfo.ActiveRegionsCount)
		err := l.LoadTerrainMesh(l.NavMeshPath+string(os.PathSeparator)+navMeshHex, regionId)
		if err != nil {
			logrus.Panic(errors.Wrap(err, "failed to load terrain mesh"))
		}
	}
}

func (l *Loader) LoadTerrainMesh(filepath string, regionId int16) error {
	fileContent, err := l.Pk2Reader.ReadFile(filepath)
	if err != nil {
		logrus.Panic(err)
	}

	readIndex := 0
	signature := string(fileContent[:12])

	if signature != "JMXVNVM 1000" {
		log.Panic("Invalid signature")
	}
	readIndex += 12

	terrain := NewRtNavmeshTerrain(filepath, NewRegionFromInt16(regionId))

	// 1. Read Object Instances
	objectInstanceCount := binary.LittleEndian.Uint16(fileContent[readIndex : readIndex+2])
	readIndex += 2
	objectInstances := make([]RtNavmeshInstObj, objectInstanceCount)
	for i := 0; i < int(objectInstanceCount); i++ {
		objectInstances[i] = l.loadTerrainObjectInstance(fileContent, &readIndex)
	}
	terrain.Objects = objectInstances

	// 2. Read Cells
	cellCount := binary.LittleEndian.Uint32(fileContent[readIndex : readIndex+4])
	// openCellCount := binary.LittleEndian.Uint32(fileContent[readIndex+4 : readIndex+8])
	readIndex += 8

	cells := make([]RtNavmeshCellQuad, cellCount)
	for i := 0; i < int(cellCount); i++ {
		cells[i] = l.loadTerrainNavigationCell(fileContent, &readIndex, i, &terrain)
	}
	terrain.Cells = cells

	// 3. Read Global Edges
	globalEdgeCount := binary.LittleEndian.Uint32(fileContent[readIndex : readIndex+4])
	globalEdges := make([]RtNavmeshEdgeGlobal, globalEdgeCount)
	readIndex += 4
	for i := 0; i < int(globalEdgeCount); i++ {
		globalEdges[i] = l.loadTerrainGlobalEdge(fileContent, &readIndex, i, &terrain)
	}
	terrain.GlobalEdges = globalEdges

	// 4. Read Internal Edges
	internalEdgeCount := binary.LittleEndian.Uint32(fileContent[readIndex : readIndex+4])
	internalEdges := make([]RtNavmeshEdgeInternal, internalEdgeCount)
	readIndex += 4
	for i := 0; i < int(internalEdgeCount); i++ {
		internalEdges[i] = l.loadTerrainInternalEdge(fileContent, &readIndex, i, &terrain)
	}
	terrain.InternalEdges = internalEdges

	// 5. Read TileMap
	terrain.tileMap = l.loadTerrainTileMap(fileContent, &readIndex)

	// 6. Read Height Map
	terrain.heightMap = l.loadTerrainHeightMap(fileContent, &readIndex)

	// 7. Read Plane Map
	terrain.planeMap = l.loadTerrainPlaneMap(fileContent, &readIndex)

	l.RegionData[regionId] = terrain

	return nil
}

func (l *Loader) loadTerrainPlaneMap(content []byte, readIndex *int) [BlocksTotal]RtNavmeshPlane {
	planeMap := [BlocksTotal]RtNavmeshPlane{}

	for i := 0; i < BlocksTotal; i++ {
		planeMap[i] = RtNavmeshPlane{
			SurfaceType: RtNavmeshSurfaceType(content[*readIndex]),
		}
		*readIndex++
	}
	for i := 0; i < BlocksTotal; i++ {
		planeMap[i].Height = utils.Float32FromByteArray(content[*readIndex : *readIndex+4])
		*readIndex += 4
	}

	return planeMap
}

func (l *Loader) loadTerrainHeightMap(content []byte, readIndex *int) [VerticesTotal]float32 {
	heightMap := [VerticesTotal]float32{}
	for i := 0; i < VerticesTotal; i++ {
		heightMap[i] = utils.Float32FromByteArray(content[*readIndex : *readIndex+4])
		*readIndex += 4
	}
	return heightMap
}

func (l *Loader) loadTerrainTileMap(content []byte, readIndex *int) [TilesTotal]RtNavmeshTile {
	tileMap := [TilesTotal]RtNavmeshTile{}
	for i := 0; i < TilesTotal; i++ {
		cellIdOffset := *readIndex
		flagOffset := cellIdOffset + 4
		textureIndexOffset := flagOffset + 2

		*readIndex += 8

		tileMap[i] = RtNavmeshTile{
			CellIndex: int(binary.LittleEndian.Uint32(content[cellIdOffset:flagOffset])),
			Flag:      RtNavmeshTileFlag(binary.LittleEndian.Uint16(content[flagOffset:textureIndexOffset])),
			TextureID: int16(binary.LittleEndian.Uint16(content[textureIndexOffset : textureIndexOffset+2])),
		}
	}
	return tileMap
}

func (l *Loader) loadTerrainInternalEdge(content []byte, readIndex *int, index int, terrain *RtNavmeshTerrain) RtNavmeshEdgeInternal {
	minVectorOffset := *readIndex
	maxVectorOffset := minVectorOffset + 8
	edgeFlagOffset := maxVectorOffset + 8
	assocDir0Offset := edgeFlagOffset + 1
	assocDir1Offset := assocDir0Offset + 1
	assocCell0Offset := assocDir1Offset + 1
	assocCell1Offset := assocCell0Offset + 2

	*readIndex += 23

	if val := content[edgeFlagOffset] & 32; val != 0 {
		log.Printf("Internal EdgeType & Bit5 %v", val)
	}
	if val := content[edgeFlagOffset] & 64; val != 0 {
		log.Printf("Internal EdgeType & Bit6 %v", val)
	}

	internalEdge := RtNavmeshEdgeInternal{
		RtNavmeshEdgeBase: RtNavmeshEdgeBase{
			RtNavmeshEdgeMeshType: RtNavmeshEdgeMeshTypeTerrain,
			Mesh:                  terrain,
			Index:                 index,
			Line: LineSegment{
				A: math32.NewVector3(utils.Float32FromByteArray(content[minVectorOffset:minVectorOffset+4]), 0, utils.Float32FromByteArray(content[minVectorOffset+4:maxVectorOffset])),
				B: math32.NewVector3(utils.Float32FromByteArray(content[maxVectorOffset:maxVectorOffset+4]), 0, utils.Float32FromByteArray(content[maxVectorOffset+4:edgeFlagOffset])),
			},
			Flag:         RtNavmeshEdgeFlag(content[edgeFlagOffset]),
			SrcDirection: RtNavmeshEdgeDirection(content[assocDir0Offset]),
			DstDirection: RtNavmeshEdgeDirection(content[assocDir1Offset]),
			SrcCellIndex: int(binary.LittleEndian.Uint16(content[assocCell0Offset:assocCell1Offset])),
			DstCellIndex: int(binary.LittleEndian.Uint16(content[assocCell1Offset : assocCell1Offset+2])),
			SrcCell:      nil,
			DstCell:      nil,
		},
	}

	return internalEdge
}

func (l *Loader) loadTerrainGlobalEdge(content []byte, readIndex *int, index int, terrain *RtNavmeshTerrain) RtNavmeshEdgeGlobal {
	minVectorOffset := *readIndex
	maxVectorOffset := minVectorOffset + 8
	edgeFlagOffset := maxVectorOffset + 8
	assocDir0Offset := edgeFlagOffset + 1
	assocDir1Offset := assocDir0Offset + 1
	assocCell0Offset := assocDir1Offset + 1
	assocCell1Offset := assocCell0Offset + 2
	assocRgn0Offset := assocCell1Offset + 2
	assocRgn1Offset := assocRgn0Offset + 2

	*readIndex += 27

	if val := content[edgeFlagOffset] & 32; val != 0 {
		log.Printf("Internal EdgeType & Bit5 %v", val)
	}
	if val := content[edgeFlagOffset] & 64; val != 0 {
		log.Printf("Internal EdgeType & Bit6 %v", val)
	}

	globalEdge := RtNavmeshEdgeGlobal{
		RtNavmeshEdgeBase: RtNavmeshEdgeBase{
			RtNavmeshEdgeMeshType: RtNavmeshEdgeMeshTypeTerrain,
			Mesh:                  terrain,
			Index:                 index,
			Line: LineSegment{
				A: math32.NewVector3(utils.Float32FromByteArray(content[minVectorOffset:minVectorOffset+4]), 0, utils.Float32FromByteArray(content[minVectorOffset+4:maxVectorOffset])),
				B: math32.NewVector3(utils.Float32FromByteArray(content[maxVectorOffset:maxVectorOffset+4]), 0, utils.Float32FromByteArray(content[maxVectorOffset+4:edgeFlagOffset])),
			},
			Flag:         RtNavmeshEdgeFlag(content[edgeFlagOffset]),
			SrcDirection: RtNavmeshEdgeDirection(content[assocDir0Offset]),
			DstDirection: RtNavmeshEdgeDirection(content[assocDir1Offset]),
			SrcCellIndex: int(binary.LittleEndian.Uint16(content[assocCell0Offset:assocCell1Offset])),
			DstCellIndex: int(binary.LittleEndian.Uint16(content[assocCell1Offset:assocRgn0Offset])),
			SrcCell:      nil,
			DstCell:      nil,
		},
		SrcMeshIndex: int(binary.LittleEndian.Uint16(content[assocRgn0Offset:assocRgn1Offset])),
		DstMeshIndex: int(binary.LittleEndian.Uint16(content[assocRgn1Offset : assocRgn1Offset+2])),
	}

	return globalEdge
}

func (l *Loader) loadTerrainNavigationCell(content []byte, readIndex *int, index int, terrain *RtNavmeshTerrain) RtNavmeshCellQuad {

	minVectorOffset := *readIndex
	maxVectorOffset := minVectorOffset + 8
	objCountOffset := maxVectorOffset + 8
	*readIndex += 17

	// cell data
	cell := RtNavmeshCellQuad{
		RtNavmeshCellBase: RtNavmeshCellBase{
			Index: index,
			Mesh:  terrain,
		},
		Rect: Rectangle{
			Min: &math32.Vector2{
				X: utils.Float32FromByteArray(content[minVectorOffset : minVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[minVectorOffset+4 : maxVectorOffset]),
			},
			Max: &math32.Vector2{
				X: utils.Float32FromByteArray(content[maxVectorOffset : maxVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[maxVectorOffset+4 : objCountOffset]),
			},
		},
		Objects: make([]RtNavmeshInstObj, content[objCountOffset]),
		edges:   make([]RtNavmeshEdge, 0),
	}

	// Object instances
	for i := 0; i < int(content[objCountOffset]); i++ {
		cell.Objects[i] = terrain.Objects[binary.LittleEndian.Uint16(content[*readIndex:*readIndex+2])]
		*readIndex += 2
	}

	return cell
}

func (l *Loader) loadTerrainObjectInstance(fileContent []byte, readIndex *int) RtNavmeshInstObj {
	objIdOffset := *readIndex
	positionOffset := objIdOffset + 4
	objTypeOffset := positionOffset + 12
	yawOffset := objTypeOffset + 2
	localUidOffset := yawOffset + 4
	short0Offset := localUidOffset + 2
	isLargeOffset := short0Offset + 2
	isStructureOffset := isLargeOffset + 1
	regionIdOffset := isStructureOffset + 1
	globalEdgeLinkCountOffset := regionIdOffset + 2
	*readIndex += 32

	objectIndex := binary.LittleEndian.Uint32(fileContent[objIdOffset:positionOffset])
	object := l.LoadObjectMesh(objectIndex)
	objectInstance := RtNavmeshInstObj{
		RtNavmeshInstBase: RtNavmeshInstBase{
			Mesh:   object,
			Object: object,
			ID:     int16(binary.LittleEndian.Uint16(fileContent[localUidOffset:short0Offset])),
			Position: math32.NewVector3(
				utils.Float32FromByteArray(fileContent[positionOffset:positionOffset+4]),
				utils.Float32FromByteArray(fileContent[positionOffset+4:positionOffset+8]),
				utils.Float32FromByteArray(fileContent[positionOffset+8:objTypeOffset]),
			),
			Rotation:     math.NewQuaternion(-utils.Float32FromByteArray(fileContent[yawOffset:localUidOffset]), 0, 0),
			Scale:        math32.NewVector3(1, 1, 1),
			LocalToWorld: nil,
			WorldToLocal: nil,
		},
		Region: NewRegionFromUint16(binary.LittleEndian.Uint16(fileContent[regionIdOffset:globalEdgeLinkCountOffset])),
	}
	objectInstance.WorldID = int(objectInstance.Region.ID<<16 | objectInstance.ID)
	objectInstance.LocalToWorld = math32.NewMatrix4().Compose(objectInstance.GetPosition(), objectInstance.GetRotation(), objectInstance.GetScale())
	objectInstance.WorldToLocal = objectInstance.GetLocalToWorld().Clone()

	err := objectInstance.GetWorldToLocal().GetInverse(objectInstance.GetLocalToWorld())
	if err != nil {
		logrus.Panic(err)
	}

	globalEdgeLinkCount := binary.LittleEndian.Uint16(fileContent[globalEdgeLinkCountOffset : globalEdgeLinkCountOffset+2])

	for j := 0; j < int(globalEdgeLinkCount); j++ {
		// TODO find out what this is for
		//linkedObjIdOffset := *readIndex
		//linkedObjEdgeIdOffset := linkedObjIdOffset + 2
		//edgeIdOffset := linkedObjEdgeIdOffset + 2
		*readIndex += 6
	}
	return objectInstance
}

func (l *Loader) SaveNavmeshDataAsGOB(filepath string) {
	logrus.Infoln("saving navmesh data as gob")
	gob.Register(RtNavmeshObj{})
	gob.Register(RtNavmeshTerrain{})
	//jsonData, err := json.Marshal(l.NavMeshData)
	//if err != nil {
	//	logrus.Errorf("failed to marshal json data. Err = %v\n", err)
	//}
	logrus.Infoln("creating precomuted file")
	f, err := os.Create(filepath + "_regions")
	if err != nil {
		logrus.Errorf("failed to create file. Err = %v\n", err)
	}

	logrus.Debugf("writing compressed data")
	//w, err := zlib.NewWriterLevel(f, zlib.BestCompression)
	w := gob.NewEncoder(f)
	err = w.Encode(l.RegionData)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("saved object navmesh data as gob")
	f.Close()

	f2, err := os.Create(filepath + "_objects")
	if err != nil {
		logrus.Errorf("f2ailed to create f2ile. Err = %v\n", err)
	}

	logrus.Debugf("writing compressed data")
	//w, err := zlib.NewWriterLevel(f2, zlib.BestCompression)
	w2 := gob.NewEncoder(f2)
	err = w2.Encode(l.RegionData)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("saved region navmesh data as gob")
	f2.Close()
}

func (l *Loader) LoadPrecomputedNavmeshDataFromGOB(filepath string) {
	logrus.Debugf("loading precomputed navmeshdata\n")
	f, err := os.Open(filepath + "_regions")

	if err != nil {
		logrus.Error(err)
	}

	decodeStartTime := time.Now()
	r := gob.NewDecoder(f)
	err = r.Decode(&l.RegionData)
	decodeStopTime := time.Now()

	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("finished loading region navmesh data after %d ms\n", decodeStopTime.Sub(decodeStartTime).Milliseconds())

	f2, err := os.Open(filepath + "_objects")

	if err != nil {
		logrus.Error(err)
	}

	decodeStartTime2 := time.Now()
	r2 := gob.NewDecoder(f2)
	err = r2.Decode(&l.ObjectData)
	decodeStopTime2 := time.Now()

	if err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("finished loading object navmesh data after %d ms\n", decodeStopTime2.Sub(decodeStartTime2).Milliseconds())

}

func (l *Loader) LoadObjectFromPrimMesh(entry ObjectInfoEntry) {
	logrus.Tracef("loading object entry %d: %s", entry.Index, entry.FilePath)
	meshFile := LoadMeshFile(l.DataPk2Path+string(os.PathSeparator)+entry.FilePath, l.Pk2Reader)
	object := meshFile.LoadMeshObject()
	l.ObjectData[entry.Index] = object
}

func (l *Loader) LoadObjectFromResource(entry ObjectInfoEntry) {
	logrus.Tracef("loading object entry %d: %s", entry.Index, entry.FilePath)
	res := LoadResource(l.DataPk2Path+string(os.PathSeparator)+entry.FilePath, l.Pk2Reader)
	mesh := LoadMeshFile(l.DataPk2Path+string(os.PathSeparator)+res.NavMeshObjPath, l.Pk2Reader)
	object := mesh.LoadMeshObject()
	l.ObjectData[entry.Index] = object
}

func (l *Loader) LoadObjectFromCompound(entry ObjectInfoEntry) {
	logrus.Tracef("loading object entry %d: %s", entry.Index, entry.FilePath)
	cpd := LoadCompoundFile(l.DataPk2Path+string(os.PathSeparator)+entry.FilePath, l.Pk2Reader)
	res := LoadResource(l.DataPk2Path+string(os.PathSeparator)+cpd.NavMeshObjPath, l.Pk2Reader)
	mesh := LoadMeshFile(l.DataPk2Path+string(os.PathSeparator)+res.NavMeshObjPath, l.Pk2Reader)
	object := mesh.LoadMeshObject()
	l.ObjectData[entry.Index] = object
}
