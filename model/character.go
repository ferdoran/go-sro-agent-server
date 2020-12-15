package model

type ICharacter interface {
	ISRObject
	GetLifeState() LifeState
	SetLifeState(state LifeState)
	GetMotionState() MotionState
	SetMotionState(state MotionState)
	GetBodyState() BodyState
	SetBodyState(state BodyState)
	GetWalkSpeed() float32
	SetWalkSpeed(speed float32)
	GetRunSpeed() float32
	SetRunSpeed(speed float32)
	GetHwanSpeed() float32
	SetHwanSpeed(speed float32)
	SetName(name string)
	MoveToPosition(position Position)
	UpdatePosition() bool
}

//func (c *Character) GetLifeState() LifeState {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.LifeState
//	return s
//}
//
//func (c *Character) SetLifeState(newState LifeState) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.LifeState = newState
//}
//
//func (c *Character) GetMotionState() MotionState {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.MotionState
//	return s
//}
//
//func (c *Character) SetMotionState(newState MotionState) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.MotionState = newState
//}
//
//func (c *Character) GetBodyState() BodyState {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.BodyState
//	return s
//}
//
//func (c *Character) SetBodyState(newState BodyState) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.BodyState = newState
//}
//
//func (c *Character) GetWalkSpeed() float32 {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.WalkSpeed
//	return s
//}
//
//func (c *Character) SetWalkSpeed(speed float32) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.WalkSpeed = speed
//}
//
//func (c *Character) GetRunSpeed() float32 {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.RunSpeed
//	return s
//}
//
//func (c *Character) SetRunSpeed(speed float32) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.RunSpeed = speed
//}
//
//func (c *Character) GetHwanSpeed() float32 {
//	c.RWMutex.RLock()
//	defer c.RWMutex.RUnlock()
//	s := c.HwanSpeed
//	return s
//}
//
//func (c *Character) SetHwanSpeed(speed float32) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.HwanSpeed = speed
//}
//
//func (c *Character) SetName(name string) {
//	c.RWMutex.Lock()
//	defer c.RWMutex.Unlock()
//	c.Name = name
//}
//
