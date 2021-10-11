package navmeshv2

import (
	"github.com/g3n/engine/math32"
)

type RtNavmeshInst interface {
	GetMesh() RtNavmesh
	GetObject() RtNavmeshObj
	GetID() int16
	GetPosition() *math32.Vector3
	GetRotation() *math32.Quaternion
	GetScale() *math32.Vector3
	GetLocalToWorld() *math32.Matrix4
	GetWorldToLocal() *math32.Matrix4
}

type RtNavmeshInstBase struct {
	Mesh         RtNavmesh
	Object       RtNavmeshObj
	ID           int16
	Position     *math32.Vector3
	Rotation     *math32.Quaternion
	Scale        *math32.Vector3
	LocalToWorld *math32.Matrix4
	WorldToLocal *math32.Matrix4
}

func (base *RtNavmeshInstBase) GetMesh() RtNavmesh {
	return base.Mesh
}

func (base *RtNavmeshInstBase) GetObject() RtNavmeshObj {
	return base.Object
}

func (base *RtNavmeshInstBase) GetID() int16 {
	return base.ID
}

func (base *RtNavmeshInstBase) GetPosition() *math32.Vector3 {
	return base.Position
}

func (base *RtNavmeshInstBase) GetRotation() *math32.Quaternion {
	return base.Rotation
}

func (base *RtNavmeshInstBase) GetScale() *math32.Vector3 {
	return base.Scale
}

func (base *RtNavmeshInstBase) GetLocalToWorld() *math32.Matrix4 {
	return base.LocalToWorld
}

func (base *RtNavmeshInstBase) GetWorldToLocal() *math32.Matrix4 {
	return base.WorldToLocal
}
