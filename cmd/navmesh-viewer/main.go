package main

import "github.com/ferdoran/go-sro-agent-server/navmeshv2/viewer"

func main() {
	viewer.Main()
	//var width int32 = 800
	//var height int32 = 450
	//rl.InitWindow(width, height, "Go SRO Navmesh Viewer")
	//rl.SetTargetFPS(60)
	//raygui.LoadGuiStyle("zahnrad.style")
	//
	//buttonClicked := false
	//var navmeshLoader *navmeshv2.Loader
	//progressChan := make(chan int)
	//progress := float32(0)
	//progressVal := 0
	//loadingNavmeshes := false
	//
	//for !rl.WindowShouldClose() {
	//	rl.BeginDrawing()
	//	rl.ClearBackground(raygui.BackgroundColor())
	//	if buttonClicked {
	//		logrus.Info("start loading navmesh files")
	//		navmeshLoader = navmeshv2.NewLoader("C:\\Data.pk2")
	//		navmeshLoader.LoadNavMeshInfos()
	//		loadingNavmeshes = true
	//		go navmeshLoader.LoadTerrainMeshes(progressChan)
	//	}
	//	rl.DrawText("Go SRO Navmesh Viewer", width/4, height/5, 20, rl.LightGray)
	//	buttonClicked = raygui.Button(rl.NewRectangle(float32(width/4), float32(height/3), float32(width/2), 50), "Load Navmesh Files")
	//
	//	if loadingNavmeshes {
	//		raygui.ProgressBar(rl.NewRectangle(float32(width/4), float32(height/2), float32(width/2), 50), progress)
	//		raygui.Label(rl.NewRectangle(float32(width/2)-50, float32(height/2), 100, 50), fmt.Sprintf("%d/%d", progressVal, len(navmeshLoader.MapProjectInfo.EnabledRegions)))
	//	}
	//
	//	select {
	//	case c := <-progressChan:
	//		progressVal = c
	//		progress = float32(c) / float32(navmeshLoader.MapProjectInfo.ActiveRegionsCount)
	//		if progress >= 1.0 {
	//			close(progressChan)
	//			loadingNavmeshes = false
	//			break
	//		}
	//	default:
	//		break
	//	}
	//	rl.EndDrawing()
	//}
	//
	//rl.CloseWindow()
}
