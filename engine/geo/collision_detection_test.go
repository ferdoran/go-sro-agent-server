package geo

import (
	"github.com/ferdoran/go-sro-agent-server/navmesh"
	"github.com/ferdoran/go-sro-framework/math"
	"github.com/g3n/engine/math32"
	"reflect"
	"testing"
)

func TestHasInnerCollision2(t *testing.T) {
	type args struct {
		a      *math32.Vector3
		b      *math32.Vector3
		object *navmesh.Object
	}
	tests := []struct {
		name                 string
		args                 args
		wantCollision        bool
		wantCollisionVectors []Collision
		wantInObjectSpace    bool
	}{
		// TODO: Add test cases.
		{
			name: "given a blocked rectangle, when running through it, then a collision is detected",
			args: args{
				a:      math32.NewVector3(15, 0, 5),
				b:      math32.NewVector3(15, 0, 20),
				object: createObjectWithBlockedRectangle(),
			},
			wantCollision: true,
			wantCollisionVectors: []Collision{{
				EdgeFlag:     7,
				VectorGlobal: math32.NewVector3(15, 0, 15),
				VectorLocal:  math32.NewVector3(15, 0, 15),
			}},
			wantInObjectSpace: true,
		},
		{
			name: "given a passable rectangle, when running through it, then there is no collision",
			args: args{
				a:      math32.NewVector3(15, 0, 5),
				b:      math32.NewVector3(15, 0, 50),
				object: createObjectWithPassableRectangle(),
			},
			wantCollision:        false,
			wantCollisionVectors: nil,
			wantInObjectSpace:    true,
		},
		{
			name: "given a simple bridge, when running through it, then there is no collision",
			args: args{
				a:      math32.NewVector3(20, 0, 20),
				b:      math32.NewVector3(20, 0, 110),
				object: createBasicBridge(),
			},
			wantCollision:        false,
			wantCollisionVectors: nil,
			wantInObjectSpace:    true,
		},
		{
			name: "given a simple bridge, when running through it from the side, then there are no collisions",
			args: args{
				a:      math32.NewVector3(0, 0, 80),
				b:      math32.NewVector3(50, 0, 80),
				object: createBasicBridge(),
			},
			wantCollision:        false,
			wantCollisionVectors: nil,
			wantInObjectSpace:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasCollision, collisionVectors, inObjSpace := HasInnerCollision2(tt.args.a, tt.args.b, tt.args.object)
			if hasCollision != tt.wantCollision {
				t.Errorf("HasInnerCollision2() hasCollision = %v, wantCollision %v", hasCollision, tt.wantCollision)
			}
			if !reflect.DeepEqual(len(collisionVectors), len(tt.wantCollisionVectors)) {
				t.Errorf("HasInnerCollision2() collisionVectors = %v, wantCollision %v", len(collisionVectors), len(tt.wantCollisionVectors))
			}
			if len(collisionVectors) > 0 && len(tt.wantCollisionVectors) > 0 {
				for i := 0; i < len(collisionVectors); i++ {
					if !collisionVectors[i].Equals(tt.wantCollisionVectors[i]) {
						t.Errorf("HasInnerCollision2 collision vectors do not equal")
					}
				}
			}
			if inObjSpace != tt.wantInObjectSpace {
				t.Errorf("HasInnerCollision2() inObjSpace = %v, wantInObjSpace %v", inObjSpace, tt.wantInObjectSpace)
			}
		})
	}
}

func TestHasOuterCollision2(t *testing.T) {
	type args struct {
		a      *math32.Vector3
		b      *math32.Vector3
		object *navmesh.Object
	}
	tests := []struct {
		name                 string
		args                 args
		wantCollision        bool
		wantCollisionVectors []Collision
		wantInObjectSpace    bool
	}{
		// TODO: Add test cases.
		{
			name: "given a blocked rectangle, when entering it, then a single collision is detected",
			args: args{
				a:      math32.NewVector3(15, 0, 5),
				b:      math32.NewVector3(15, 0, 15),
				object: createObjectWithBlockedRectangle(),
			},
			wantCollision: true,
			wantCollisionVectors: []Collision{
				{
					EdgeFlag:     3,
					VectorLocal:  math32.NewVector3(15, 0, 10),
					VectorGlobal: math32.NewVector3(15, 0, 10),
				}},
			wantInObjectSpace: false,
		},
		{
			name: "given a passable rectangle, when running through it, then there is no collision",
			args: args{
				a:      math32.NewVector3(15, 0, 5),
				b:      math32.NewVector3(15, 0, 50),
				object: createObjectWithPassableRectangle(),
			},
			wantCollision:        false,
			wantCollisionVectors: nil,
			wantInObjectSpace:    true,
		},
		{
			name: "given a simple bridge, when running through it, then there is no collision",
			args: args{
				a:      math32.NewVector3(20, 0, 20),
				b:      math32.NewVector3(20, 0, 110),
				object: createBasicBridge(),
			},
			wantCollision:        false,
			wantCollisionVectors: nil,
			wantInObjectSpace:    true,
		},
		{
			name: "given a simple bridge, when running through it from the side, then there are 2 collisions",
			args: args{
				a:      math32.NewVector3(0, 0, 80),
				b:      math32.NewVector3(50, 0, 80),
				object: createBasicBridge(),
			},
			wantCollision: true,
			wantCollisionVectors: []Collision{
				{
					EdgeFlag:     3,
					VectorLocal:  math32.NewVector3(10, 0, 80),
					VectorGlobal: math32.NewVector3(10, 0, 80),
				},
				{
					EdgeFlag:     3,
					VectorLocal:  math32.NewVector3(30, 0, 80),
					VectorGlobal: math32.NewVector3(30, 0, 80),
				},
			},
			wantInObjectSpace: false,
		},
		{
			name: "given a simple bridge with bridge outline, when running through it from the side, then there are 2 collisions",
			args: args{
				a:      math32.NewVector3(0, 0, 80),
				b:      math32.NewVector3(50, 0, 80),
				object: createBasicBridgeWithBridgeOutline(),
			},
			wantCollision: true,
			wantCollisionVectors: []Collision{
				{
					EdgeFlag:     16,
					VectorGlobal: math32.NewVector3(10, 0, 80),
					VectorLocal:  math32.NewVector3(10, 0, 80),
				},
				{
					EdgeFlag:     16,
					VectorLocal:  math32.NewVector3(30, 0, 80),
					VectorGlobal: math32.NewVector3(30, 0, 80),
				},
			},
			wantInObjectSpace: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasCollision, collisionVectors, inObjSpace := HasOuterCollision2(tt.args.a, tt.args.b, tt.args.object)
			if hasCollision != tt.wantCollision {
				t.Errorf("HasOuterCollision2() hasCollision = %v, wantCollision %v", hasCollision, tt.wantCollision)
			}
			if !reflect.DeepEqual(len(collisionVectors), len(tt.wantCollisionVectors)) {
				t.Errorf("HasOuterCollision2() collisionVectors = %v, wantCollision %v", len(collisionVectors), len(tt.wantCollisionVectors))
			}
			if len(collisionVectors) > 0 && len(tt.wantCollisionVectors) > 0 {
				for i := 0; i < len(collisionVectors); i++ {
					if !collisionVectors[i].Equals(tt.wantCollisionVectors[i]) {
						t.Errorf("HasOuterCollision2() collision vectors do not equal: got %v, want %v", collisionVectors[i], tt.wantCollisionVectors[i])
					}
				}
			}
			if inObjSpace != tt.wantInObjectSpace {
				t.Errorf("HasOuterCollision2() inObjSpace = %v, wantInObjSpace %v", inObjSpace, tt.wantInObjectSpace)
			}
		})
	}
}

func createObjectWithBlockedRectangle() *navmesh.Object {
	obj := &navmesh.Object{
		ID:                  0,
		Position:            nil,
		Type:                0,
		Yaw:                 0,
		LocalUID:            0,
		Short0:              0,
		IsLarge:             false,
		IsStructure:         false,
		RegionID:            0,
		GlobalEdgeLinkCount: 0,
		GlobalEdgeLinks:     nil,
		Vertices:            nil,
		Cells:               nil,
		GlobalEdges:         nil,
		InternalEdges:       nil,
		Events:              nil,
		Grid:                nil,
		Rotation:            nil,
		LocalToWorld:        nil,
		WorldToLocal:        nil,
	}

	obj.LocalToWorld = math32.NewMatrix4()

	pointA := math32.NewVector3(10, 0, 10)
	pointB := math32.NewVector3(10, 0, 20)
	pointC := math32.NewVector3(20, 0, 10)
	pointD := math32.NewVector3(20, 0, 20)

	gEdge1 := createGlobalBlockedEdge(pointA, pointB)
	gEdge2 := createGlobalBlockedEdge(pointA, pointC)
	gEdge3 := createGlobalBlockedEdge(pointB, pointD)
	gEdge4 := createGlobalBlockedEdge(pointC, pointD)

	iEdge := createInternalBlockedEdge(pointB, pointC)

	obj.GlobalEdges = []*navmesh.ObjectGlobalEdge{gEdge1, gEdge2, gEdge3, gEdge4}
	obj.InternalEdges = []*navmesh.ObjectInternalEdge{iEdge}

	cell1 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointB,
			C: pointC,
		},
		Index: 0,
		Flag:  0,
	}

	cell2 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointB,
			B: pointC,
			C: pointD,
		},
		Index: 1,
		Flag:  0,
	}

	obj.Cells = []*navmesh.ObjectCell{cell1, cell2}
	return obj

}

func createObjectWithPassableRectangle() *navmesh.Object {
	obj := &navmesh.Object{
		ID:                  0,
		Position:            nil,
		Type:                0,
		Yaw:                 0,
		LocalUID:            0,
		Short0:              0,
		IsLarge:             false,
		IsStructure:         false,
		RegionID:            0,
		GlobalEdgeLinkCount: 0,
		GlobalEdgeLinks:     nil,
		Vertices:            nil,
		Cells:               nil,
		GlobalEdges:         nil,
		InternalEdges:       nil,
		Events:              nil,
		Grid:                nil,
		Rotation:            nil,
		LocalToWorld:        nil,
		WorldToLocal:        nil,
	}

	obj.LocalToWorld = math32.NewMatrix4()

	pointA := math32.NewVector3(10, 0, 10)
	pointB := math32.NewVector3(10, 0, 20)
	pointC := math32.NewVector3(20, 0, 10)
	pointD := math32.NewVector3(20, 0, 20)

	gEdge1 := createGlobalPassableEdge(pointA, pointB)
	gEdge2 := createGlobalPassableEdge(pointA, pointC)
	gEdge3 := createGlobalPassableEdge(pointB, pointD)
	gEdge4 := createGlobalPassableEdge(pointC, pointD)

	iEdge := createInternalPassableEdge(pointB, pointC)

	obj.GlobalEdges = []*navmesh.ObjectGlobalEdge{gEdge1, gEdge2, gEdge3, gEdge4}
	obj.InternalEdges = []*navmesh.ObjectInternalEdge{iEdge}

	cell1 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointB,
			C: pointC,
		},
		Index: 0,
		Flag:  0,
	}

	cell2 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointB,
			B: pointC,
			C: pointD,
		},
		Index: 1,
		Flag:  0,
	}

	obj.Cells = []*navmesh.ObjectCell{cell1, cell2}
	return obj

}

func createBasicBridge() *navmesh.Object {
	obj := &navmesh.Object{
		ID:                  0,
		Position:            nil,
		Type:                0,
		Yaw:                 0,
		LocalUID:            0,
		Short0:              0,
		IsLarge:             false,
		IsStructure:         false,
		RegionID:            0,
		GlobalEdgeLinkCount: 0,
		GlobalEdgeLinks:     nil,
		Vertices:            nil,
		Cells:               nil,
		GlobalEdges:         nil,
		InternalEdges:       nil,
		Events:              nil,
		Grid:                nil,
		Rotation:            nil,
		LocalToWorld:        nil,
		WorldToLocal:        nil,
	}

	obj.LocalToWorld = math32.NewMatrix4()
	/*

	   A____B
	   |____|
	   C    D
	   |____|
	   E    F
	   |____|
	   G    H
	   |____|
	   I    J
	   |____|
	   K    L
	*/

	pointA := math32.NewVector3(10, 0, 100)
	pointB := math32.NewVector3(30, 0, 100)
	pointC := math32.NewVector3(10, 0, 90)
	pointD := math32.NewVector3(30, 0, 90)
	pointE := math32.NewVector3(10, 0, 80)
	pointF := math32.NewVector3(30, 0, 80)
	pointG := math32.NewVector3(10, 0, 70)
	pointH := math32.NewVector3(30, 0, 70)
	pointI := math32.NewVector3(10, 0, 60)
	pointJ := math32.NewVector3(30, 0, 60)
	pointK := math32.NewVector3(10, 0, 50)
	pointL := math32.NewVector3(30, 0, 50)

	gEdge1 := createGlobalPassableEdge(pointA, pointB)
	gEdge2 := createGlobalPassableEdge(pointK, pointL)
	gEdge3 := createGlobalBlockedEdge(pointA, pointK)
	gEdge4 := createGlobalBlockedEdge(pointB, pointL)

	iEdge1 := createInternalPassableEdge(pointA, pointD)
	iEdge2 := createInternalPassableEdge(pointC, pointF)
	iEdge3 := createInternalPassableEdge(pointE, pointH)
	iEdge4 := createInternalPassableEdge(pointG, pointJ)
	iEdge5 := createInternalPassableEdge(pointI, pointL)

	obj.GlobalEdges = []*navmesh.ObjectGlobalEdge{gEdge1, gEdge2, gEdge3, gEdge4}
	obj.InternalEdges = []*navmesh.ObjectInternalEdge{iEdge1, iEdge2, iEdge3, iEdge4, iEdge5}

	cell1 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointB,
			C: pointD,
		},
		Index: 0,
		Flag:  0,
	}

	cell2 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointC,
			C: pointD,
		},
		Index: 1,
		Flag:  0,
	}

	cell3 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointC,
			B: pointE,
			C: pointF,
		},
		Index: 2,
		Flag:  0,
	}

	cell4 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointC,
			B: pointD,
			C: pointF,
		},
		Index: 3,
		Flag:  0,
	}

	cell5 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointE,
			B: pointG,
			C: pointH,
		},
		Index: 4,
		Flag:  0,
	}

	cell6 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointE,
			B: pointF,
			C: pointH,
		},
		Index: 5,
		Flag:  0,
	}

	cell7 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointG,
			B: pointI,
			C: pointJ,
		},
		Index: 6,
		Flag:  0,
	}

	cell8 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointG,
			B: pointH,
			C: pointJ,
		},
		Index: 7,
		Flag:  0,
	}

	cell9 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointI,
			B: pointK,
			C: pointL,
		},
		Index: 8,
		Flag:  0,
	}

	cell10 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointI,
			B: pointJ,
			C: pointL,
		},
		Index: 9,
		Flag:  0,
	}

	obj.Cells = []*navmesh.ObjectCell{cell1, cell2, cell3, cell4, cell5, cell6, cell7, cell8, cell9, cell10}
	return obj
}

func createBasicBridgeWithBridgeOutline() *navmesh.Object {
	obj := &navmesh.Object{
		ID:                  0,
		Position:            nil,
		Type:                0,
		Yaw:                 0,
		LocalUID:            0,
		Short0:              0,
		IsLarge:             false,
		IsStructure:         false,
		RegionID:            0,
		GlobalEdgeLinkCount: 0,
		GlobalEdgeLinks:     nil,
		Vertices:            nil,
		Cells:               nil,
		GlobalEdges:         nil,
		InternalEdges:       nil,
		Events:              nil,
		Grid:                nil,
		Rotation:            nil,
		LocalToWorld:        nil,
		WorldToLocal:        nil,
	}

	obj.LocalToWorld = math32.NewMatrix4()
	/*

	   A____B
	   |____|
	   C    D
	   |____|
	   E    F
	   |____|
	   G    H
	   |____|
	   I    J
	   |____|
	   K    L
	*/

	pointA := math32.NewVector3(10, 0, 100)
	pointB := math32.NewVector3(30, 0, 100)
	pointC := math32.NewVector3(10, 0, 90)
	pointD := math32.NewVector3(30, 0, 90)
	pointE := math32.NewVector3(10, 0, 80)
	pointF := math32.NewVector3(30, 0, 80)
	pointG := math32.NewVector3(10, 0, 70)
	pointH := math32.NewVector3(30, 0, 70)
	pointI := math32.NewVector3(10, 0, 60)
	pointJ := math32.NewVector3(30, 0, 60)
	pointK := math32.NewVector3(10, 0, 50)
	pointL := math32.NewVector3(30, 0, 50)

	gEdge1 := createGlobalPassableEdge(pointA, pointB)
	gEdge2 := createGlobalPassableEdge(pointK, pointL)
	gEdge3 := createGlobalBridgeEdge(pointA, pointK)
	gEdge4 := createGlobalBridgeEdge(pointB, pointL)

	iEdge1 := createInternalPassableEdge(pointA, pointD)
	iEdge2 := createInternalPassableEdge(pointC, pointF)
	iEdge3 := createInternalPassableEdge(pointE, pointH)
	iEdge4 := createInternalPassableEdge(pointG, pointJ)
	iEdge5 := createInternalPassableEdge(pointI, pointL)

	obj.GlobalEdges = []*navmesh.ObjectGlobalEdge{gEdge1, gEdge2, gEdge3, gEdge4}
	obj.InternalEdges = []*navmesh.ObjectInternalEdge{iEdge1, iEdge2, iEdge3, iEdge4, iEdge5}

	cell1 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointB,
			C: pointD,
		},
		Index: 0,
		Flag:  0,
	}

	cell2 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointA,
			B: pointC,
			C: pointD,
		},
		Index: 1,
		Flag:  0,
	}

	cell3 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointC,
			B: pointE,
			C: pointF,
		},
		Index: 2,
		Flag:  0,
	}

	cell4 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointC,
			B: pointD,
			C: pointF,
		},
		Index: 3,
		Flag:  0,
	}

	cell5 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointE,
			B: pointG,
			C: pointH,
		},
		Index: 4,
		Flag:  0,
	}

	cell6 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointE,
			B: pointF,
			C: pointH,
		},
		Index: 5,
		Flag:  0,
	}

	cell7 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointG,
			B: pointI,
			C: pointJ,
		},
		Index: 6,
		Flag:  0,
	}

	cell8 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointG,
			B: pointH,
			C: pointJ,
		},
		Index: 7,
		Flag:  0,
	}

	cell9 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointI,
			B: pointK,
			C: pointL,
		},
		Index: 8,
		Flag:  0,
	}

	cell10 := &navmesh.ObjectCell{
		Triangle: &math.Triangle{
			A: pointI,
			B: pointJ,
			C: pointL,
		},
		Index: 9,
		Flag:  0,
	}

	obj.Cells = []*navmesh.ObjectCell{cell1, cell2, cell3, cell4, cell5, cell6, cell7, cell8, cell9, cell10}
	return obj
}

func createInternalBlockedEdge(a, b *math32.Vector3) *navmesh.ObjectInternalEdge {
	return &navmesh.ObjectInternalEdge{
		A:    a,
		B:    b,
		Flag: 7,
	}
}

func createInternalBridgeEdge(a, b *math32.Vector3) *navmesh.ObjectInternalEdge {
	return &navmesh.ObjectInternalEdge{
		A:    a,
		B:    b,
		Flag: 20,
	}
}

func createGlobalBridgeEdge(a, b *math32.Vector3) *navmesh.ObjectGlobalEdge {
	return &navmesh.ObjectGlobalEdge{
		A:    a,
		B:    b,
		Flag: 16,
	}
}

func createInternalPassableEdge(a, b *math32.Vector3) *navmesh.ObjectInternalEdge {
	return &navmesh.ObjectInternalEdge{
		A:    a,
		B:    b,
		Flag: 4,
	}
}

func createGlobalBlockedEdge(a, b *math32.Vector3) *navmesh.ObjectGlobalEdge {
	return &navmesh.ObjectGlobalEdge{
		A:    a,
		B:    b,
		Flag: 3,
	}
}

func createGlobalPassableEdge(a, b *math32.Vector3) *navmesh.ObjectGlobalEdge {
	return &navmesh.ObjectGlobalEdge{
		A:    a,
		B:    b,
		Flag: 0,
	}
}
