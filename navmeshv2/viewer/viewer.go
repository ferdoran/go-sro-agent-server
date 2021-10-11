package viewer

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var screenWidth = float32(1920)
var screenHeight = float32(1080)
var loader *navmeshv2.Loader
var terrainList map[int16]navmeshv2.RtNavmeshTerrain
var loaderProgress = float32(0)
var loaderProgressAbsolute = 0
var loaderProgressChannel = make(chan int)
var loadingMeshes = false

func Main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "GoSro Navmesh Viewer")

	defer rl.CloseWindow()

	monitorWidth := float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))
	monitorHeight := float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))
	rl.SetWindowSize(int(monitorWidth*0.8), int(monitorHeight*0.8))
	rl.SetWindowPosition(int(monitorWidth*0.1), int(monitorHeight*0.1))

	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	raygui.LoadGuiStyle("cmd/navmesh-viewer/zahnrad.style")
	rl.SetCameraMode(cam, rl.CameraFree)
	rl.SetCameraAltControl(rl.KeyLeftAlt)
	rl.SetCameraPanControl(rl.MouseMiddleButton)
	rl.SetCameraSmoothZoomControl(rl.KeyLeftShift)
	rl.SetExitKey(0)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&cam)
		doFrame()
	}
}

func doFrame() {
	screenWidth = float32(rl.GetScreenWidth())
	screenHeight = float32(rl.GetScreenHeight())

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(raygui.BackgroundColor())

	drawToolbar()
	drawTerrainListPane()
	drawTerrain()
}
