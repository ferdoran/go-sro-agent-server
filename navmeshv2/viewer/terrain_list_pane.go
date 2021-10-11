package viewer

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
)

const TerrainListPaneWidthFactor = 0.2
const ScrollSpeed = 10.0

var searchRegionText string
var searchedRegionId int
var selectedRegionId int
var terrainListCam = rl.Camera2D{
	Offset:   rl.Vector2{X: 0, Y: 20},
	Target:   rl.Vector2{},
	Rotation: 0,
	Zoom:     1,
}

type TerrainListPane struct {
	MeshFiles  map[int16]navmeshv2.RtNavmeshTerrain
	SearchText string
}

func NewPane(regions map[int16]navmeshv2.RtNavmeshTerrain) *TerrainListPane {
	return &TerrainListPane{
		MeshFiles:  regions,
		SearchText: "",
	}
}

func drawTerrainListPane() {
	paneWidth := float32(rl.GetScreenWidth()) * TerrainListPaneWidthFactor
	paneHeight := rl.GetScreenHeight()
	if mousePos := rl.GetMousePosition(); mousePos.X > 0 && mousePos.X < paneWidth && mousePos.Y > ToolbarHeight && mousePos.Y < float32(rl.GetScreenHeight()) {
		if scroll := rl.GetMouseWheelMove(); scroll != 0 {
			terrainListCam.Target.Y += float32(scroll) * ScrollSpeed
		}
	}
	rl.DrawLine(int32(paneWidth), ToolbarHeight, int32(paneWidth), int32(paneHeight), rl.DarkGray)
	// Toolbar Height is 64
	rl.BeginScissorMode(0, ToolbarHeight, int32(paneWidth), int32(paneHeight))
	defer rl.EndScissorMode()

	if loadingMeshes {
		select {
		case c := <-loaderProgressChannel:
			loaderProgressAbsolute = c
			loaderProgress = float32(c) / float32(loader.MapProjectInfo.ActiveRegionsCount)
			if loaderProgress >= 1.0 {
				loadingMeshes = false
				loaderProgress = 0
				loaderProgressAbsolute = 0
				close(loaderProgressChannel)
				loaderProgressChannel = make(chan int)
				terrainList = loader.RegionData
			}
		default:
			break
		}
		raygui.ProgressBar(rl.NewRectangle(0, ToolbarHeight+1, paneWidth-1, 20), loaderProgress)
		raygui.Label(rl.NewRectangle(paneWidth*0.25, ToolbarHeight, paneWidth*0.5, 20), fmt.Sprintf("%d/%d", loaderProgressAbsolute, loader.MapProjectInfo.ActiveRegionsCount))
	} else if len(terrainList) == 0 {
		raygui.Label(rl.NewRectangle(0, ToolbarHeight, paneWidth, 20), "There are no terrains!")
	} else {
		searchRegionText = raygui.TextBox(rl.NewRectangle(0, ToolbarHeight, paneWidth*0.75, 20), searchRegionText)
		if raygui.Button(rl.NewRectangle(paneWidth*0.75, ToolbarHeight, paneWidth*0.25, 20), "Search") {
			regex := regexp.MustCompile("^\\d{5}$")
			if regex.MatchString(searchRegionText) {
				searchedRegionId, _ = strconv.Atoi(searchRegionText)
				logrus.Infof("searching for region with id %d", searchedRegionId)
			} else if searchRegionText == "" {
				searchedRegionId = 0
			}
		}
		raygui.Label(rl.NewRectangle(0, ToolbarHeight+20, paneWidth*0.5, 20), "Region ID")
		raygui.Label(rl.NewRectangle(paneWidth*0.5, ToolbarHeight+20, paneWidth*0.25, 20), "X")
		raygui.Label(rl.NewRectangle(paneWidth*0.75, ToolbarHeight+20, paneWidth*0.25, 20), "Y")
		rl.DrawLine(0, ToolbarHeight+40, int32(paneWidth), ToolbarHeight+40, rl.DarkGray)
		rl.DrawLine(int32(paneWidth*0.5), ToolbarHeight+20, int32(paneWidth*0.5), int32(paneHeight), rl.DarkGray)
		rl.DrawLine(int32(paneWidth*0.75), ToolbarHeight+20, int32(paneWidth*0.75), int32(paneHeight), rl.DarkGray)
		rl.BeginMode2D(terrainListCam)
		defer rl.EndMode2D()

		var keys []int
		if _, exists := terrainList[int16(searchedRegionId)]; searchedRegionId != 0 && exists {
			keys = []int{searchedRegionId}
		} else if searchedRegionId != 0 && !exists {
			raygui.Label(rl.NewRectangle(0, ToolbarHeight+40, paneWidth, 20), fmt.Sprintf("There are no regions with id %d", searchedRegionId))
		} else {
			for regionId := range terrainList {
				keys = append(keys, int(regionId))
			}
			sort.Ints(keys)
		}
		nextY := float32(ToolbarHeight + 40)

		// TODO draw scrollbar?
		for _, regionId := range keys {
			if mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), terrainListCam); mousePos.X > 0 && mousePos.X < paneWidth && mousePos.Y < nextY+20 && mousePos.Y > nextY {
				// hover effect
				rl.DrawRectangleRec(rl.NewRectangle(0, nextY, paneWidth, 20), rl.Gray)
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					// TODO Select this element
					if selectedRegionId != regionId {
						selectedRegionId = regionId
					} else {
						selectedRegionId = 0
					}
				}
			}

			if selectedRegionId != 0 && selectedRegionId == regionId {
				rl.DrawRectangleRec(rl.NewRectangle(0, nextY, paneWidth, 20), rl.Gray)
			}

			raygui.Label(rl.NewRectangle(0, nextY, paneWidth*0.5, 20), fmt.Sprintf("%d", regionId))

			x, y := utils.Int16ToXAndZ(int16(regionId))
			raygui.Label(rl.NewRectangle(paneWidth*0.5, nextY, paneWidth*0.25, 20), fmt.Sprintf("%d", x))
			raygui.Label(rl.NewRectangle(paneWidth*0.75, nextY, paneWidth*0.25, 20), fmt.Sprintf("%d", y))
			nextY += 20
		}
	}
}
