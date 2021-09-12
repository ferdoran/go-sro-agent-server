package main

import (
	"github.com/ferdoran/go-sro-agent-server/navmesh"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/fogleman/gg"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/sirupsen/logrus"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"strings"
)

const DataPk2File = "E:\\Silkroad_TestIn3\\Data.pk2"

type Pk2File struct {
	pk2.PackFileEntry
}
type Pk2Dir struct {
	pk2.Directory
}

func (dir Pk2Dir) String() string {
	pathSegments := strings.Split(dir.Name, string(os.PathSeparator))
	return pathSegments[len(pathSegments)-1]
}

func (dir Pk2Dir) Expand() []*widgets.TreeNode {
	nodes := make([]*widgets.TreeNode, 0)

	for _, d := range dir.DirectoriesByName {
		di := Pk2Dir{d}
		nodes = append(nodes, &widgets.TreeNode{
			Value:    di,
			Expanded: false,
			Nodes:    nil,
		})
	}

	for _, file := range dir.Files {
		f := Pk2File{file}
		nodes = append(nodes, &widgets.TreeNode{
			Value:    f,
			Expanded: false,
			Nodes:    nil,
		})
	}

	return nodes
}

func (entry Pk2File) String() string {
	return entry.Name
}

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	navmeshLoader := navmesh.NewLoader(DataPk2File)
	reader := navmeshLoader.Pk2Reader
	fileTree := widgets.NewTree()
	fileTree.WrapText = false
	fileTree.TitleStyle = ui.NewStyle(ui.ColorYellow)
	fileTree.TextStyle = ui.NewStyle(ui.ColorWhite)
	fileTree.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	fileTree.Title = DataPk2File

	rootDir := Pk2Dir{reader.Directory}
	fileTree.SetNodes(rootDir.Expand())
	termWidth, termHeight := ui.TerminalDimensions()

	//image.SetRect(termWidth / 2, 0, termWidth/2, termHeight)
	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewCol(1.0, fileTree),
	)
	ui.Render(grid)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "j", "<Down>":
				fileTree.ScrollDown()
			case "k", "<Up>":
				fileTree.ScrollUp()
			case "<Right>":
				selectedNode := fileTree.SelectedNode()
				dir, ok := selectedNode.Value.(Pk2Dir)
				if ok {
					selectedNode.Nodes = dir.Expand()
					fileTree.Expand()
				}
			case "<Left>":
				fileTree.Collapse()
			case "<Enter>":
				selectedNode := fileTree.SelectedNode()
				switch t := selectedNode.Value.(type) {
				case Pk2File:
					extension := path.Ext(t.Name)
					switch extension {
					case ".nvm":
						// Load navmesh
						navmeshLoader.LoadNavMeshInfos()
						navmeshData := LoadNavmeshData(t.Name, navmeshLoader)
						CreateImageFromNavmeshData(t.Name, navmeshData)

					default:
						logrus.Infof("extension: %s", extension)
					}
				}
			case "q":
				return
			}
		}
		ui.Render(grid)
	}
}

func CreateImageFromNavmeshData(filename string, navmeshData navmesh.NavMeshData) image.Image {
	context := gg.NewContext(128+1920, 128+1920)
	context.SetStrokeStyle(gg.NewSolidPattern(color.White))
	context.SetLineWidth(8)
	context.DrawRectangle(64, 64, 1920, 1920)
	context.Stroke()
	context.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 169, G: 169, B: 169, A: 0x7F}))
	context.SetLineWidth(2)
	for _, cell := range navmeshData.Cells {
		context.DrawRectangle(64+float64(cell.Min.X), 64+float64(cell.Min.Y), float64(cell.Max.X-cell.Min.X), float64(cell.Max.Y-cell.Min.Y))
		context.Stroke()
	}
	context.SavePNG(filename + ".png")
	return context.Image()
}

func LoadNavmeshData(filename string, navmeshLoader *navmesh.Loader) navmesh.NavMeshData {
	fileContent, err := navmeshLoader.Pk2Reader.ReadFile(strings.Join([]string{"Data", "navmesh", filename}, string(os.PathSeparator)))
	if err != nil {
		logrus.Panic(err)
	}
	navmeshData := navmesh.ParseNavMeshFile(filename, fileContent)
	counter := 0
	for _, o := range navmeshData.Objects {
		counter++
		obj := navmeshLoader.ObjectInfo.Objects[o.ID]
		var res *navmesh.Resource
		var mesh *navmesh.MeshFile
		if strings.HasSuffix(obj.FilePath, "cpd") {
			cpd := navmesh.LoadCompoundFile(navmeshLoader.DataPk2Path+string(os.PathSeparator)+obj.FilePath, navmeshLoader.Pk2Reader)
			res = navmesh.LoadResource(navmeshLoader.DataPk2Path+string(os.PathSeparator)+cpd.NavMeshObjPath, navmeshLoader.Pk2Reader)
			mesh = navmesh.LoadMeshFile(navmeshLoader.DataPk2Path+string(os.PathSeparator)+res.NavMeshObjPath, navmeshLoader.Pk2Reader)
		} else if strings.HasSuffix(obj.FilePath, "bsr") {
			res = navmesh.LoadResource(navmeshLoader.DataPk2Path+string(os.PathSeparator)+obj.FilePath, navmeshLoader.Pk2Reader)
			mesh = navmesh.LoadMeshFile(navmeshLoader.DataPk2Path+string(os.PathSeparator)+res.NavMeshObjPath, navmeshLoader.Pk2Reader)
		} else if strings.HasSuffix(obj.FilePath, "bms") {
			mesh = navmesh.LoadMeshFile(navmeshLoader.DataPk2Path+string(os.PathSeparator)+obj.FilePath, navmeshLoader.Pk2Reader)
		} else {
			logrus.Panicf("unsupported file: %s\n", obj.FilePath)
		}
		mesh.LoadMeshObject(o)
	}
	navmeshLoader.NavMeshData[filename] = navmeshData
	return navmeshData
}
