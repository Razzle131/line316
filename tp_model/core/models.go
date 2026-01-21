package core

import (
	"errors"
	"sync/atomic"
	"time"
)

type Puck struct {
	Color      string
	IsPackaged bool
}

func NewPuck(color string) Puck {
	return Puck{
		Color:      color,
		IsPackaged: false,
	}
}

type Start struct {
	PuckSlot *Puck
}

func NewStart() Start {
	return Start{
		PuckSlot: nil,
	}
}

func (s *Start) PlacePuck(puck Puck) error {
	if s.PuckSlot != nil {
		return ErrSlotOccupied
	}

	s.PuckSlot = &puck
	return nil
}

func (s *Start) TakePuck() (Puck, error) {
	if s.PuckSlot == nil {
		return Puck{}, ErrSlotEmpty
	}

	puck := *s.PuckSlot
	s.PuckSlot = nil

	return puck, nil
}

type Gripper struct {
	IsOpen                bool
	PuckSlot              *Puck
	IsMovingVerticly      bool
	IsMovingHorizontaly   bool
	IsWantToStopMoving    atomic.Bool
	CurHorizontalPosition float64
	CurVerticalPosition   float64
}

func NewGripper() Gripper {
	return Gripper{
		IsOpen:                false,
		PuckSlot:              nil,
		IsMovingVerticly:      false,
		IsMovingHorizontaly:   false,
		IsWantToStopMoving:    atomic.Bool{},
		CurHorizontalPosition: gripperStartPos,
		CurVerticalPosition:   gripperUpPos,
	}
}

func (g *Gripper) Stop() {
	g.IsWantToStopMoving.Store(true)
}

// to enable moving commands after stop
func (g *Gripper) EnableMoving() {
	g.IsWantToStopMoving.Store(false)
}

func (g *Gripper) Open() {
	g.IsOpen = true
}

func (g *Gripper) Close() {
	g.IsOpen = false
}

func (g *Gripper) TakePuck(puck Puck) error {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return ErrGripperAlreadyMoving
	}

	if g.PuckSlot != nil {
		return ErrSlotOccupied
	}

	if !g.IsOpen {
		return errors.New("need to open gripper to take puck")
	}

	g.PuckSlot = &puck

	return nil
}

func (g *Gripper) PlacePuck() (Puck, error) {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return Puck{}, ErrGripperAlreadyMoving
	}

	if g.PuckSlot == nil {
		return Puck{}, ErrSlotEmpty
	}

	if !g.IsOpen {
		return Puck{}, errors.New("need to open gripper to place puck")
	}

	puck := *g.PuckSlot
	g.PuckSlot = nil

	return puck, nil
}

func (g *Gripper) MoveLeft() error {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return ErrGripperAlreadyMoving
	}

	if g.IsWantToStopMoving.Load() {
		return ErrGripperWantToStop
	}

	g.IsMovingHorizontaly = true

	go func() {
		ticker := time.NewTicker(time.Second / tickrate)
		for range ticker.C {
			if g.IsWantToStopMoving.Load() {
				break
			}
			g.CurHorizontalPosition = max(gripperCarouselPos, g.CurHorizontalPosition-float64(gripperHorizontalSpeed)/tickrate)
		}
		g.IsMovingHorizontaly = false
	}()

	return nil
}

func (g *Gripper) MoveRight() error {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return ErrGripperAlreadyMoving
	}

	if g.IsWantToStopMoving.Load() {
		return ErrGripperWantToStop
	}

	g.IsMovingHorizontaly = true

	go func() {
		ticker := time.NewTicker(time.Second / tickrate)
		for range ticker.C {
			if g.IsWantToStopMoving.Load() {
				break
			}
			g.CurHorizontalPosition = min(gripperSortingPos, g.CurHorizontalPosition+float64(gripperHorizontalSpeed)/tickrate)
		}
		g.IsMovingHorizontaly = false
	}()

	return nil
}

func (g *Gripper) MoveUp() error {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return ErrGripperAlreadyMoving
	}

	if g.IsWantToStopMoving.Load() {
		return ErrGripperWantToStop
	}

	g.IsMovingVerticly = true

	go func() {
		ticker := time.NewTicker(time.Second / tickrate)
		for range ticker.C {
			if g.IsWantToStopMoving.Load() {
				break
			}
			g.CurVerticalPosition = min(gripperUpPos, g.CurVerticalPosition+float64(gripperVerticalSpeed)/tickrate)
		}
		g.IsMovingVerticly = false
	}()

	return nil
}

func (g *Gripper) MoveDown() error {
	if g.IsMovingHorizontaly || g.IsMovingVerticly {
		return ErrGripperAlreadyMoving
	}

	if g.IsWantToStopMoving.Load() {
		return ErrGripperWantToStop
	}

	g.IsMovingVerticly = true

	go func() {
		ticker := time.NewTicker(time.Second / tickrate)
		for range ticker.C {
			if g.IsWantToStopMoving.Load() {
				break
			}
			g.CurVerticalPosition = max(gripperDownPos, g.CurVerticalPosition-float64(gripperVerticalSpeed)/tickrate)
		}
		g.IsMovingVerticly = false
	}()

	return nil
}

type Carousel struct {
	Slots []*Puck
}

func NewCarousel() Carousel {
	return Carousel{
		Slots: make([]*Puck, carouselTotalSlots),
	}
}

func (c *Carousel) PlacePuck(puck Puck) error {
	if c.Slots[0] != nil {
		return ErrSlotOccupied
	}

	if puck.IsPackaged {
		return ErrPuckPackaged
	}

	c.Slots[0] = &puck
	return nil
}

func (c *Carousel) TakePuck() (Puck, error) {
	if c.Slots[0] == nil {
		return Puck{}, ErrSlotEmpty
	}

	p := *c.Slots[0]
	c.Slots[0] = nil

	return p, nil
}

func (c *Carousel) InspectPuck() (Puck, error) {
	if carouselInspectSlot >= len(c.Slots) {
		return Puck{}, errors.New("bad inspect slot param")
	}

	if c.Slots[carouselInspectSlot] == nil {
		return Puck{}, ErrSlotEmpty
	}

	return *c.Slots[carouselInspectSlot], nil
}

func (c *Carousel) DrillPuck() error {
	if carouselDrillSlot >= len(c.Slots) {
		return errors.New("bad inspect slot param")
	}

	if c.Slots[carouselDrillSlot] == nil {
		return ErrSlotEmpty
	}

	return nil
}

func (c *Carousel) RotateOnce() {
	start := time.Now()
	end := start.Add(carouselNextSlotTime)

	res := make([]*Puck, len(c.Slots))
	for i := range c.Slots {
		res[(i+1)%len(c.Slots)] = c.Slots[i]
	}

	for time.Since(end) < 0 {
		time.Sleep(time.Millisecond)
	}

	c.Slots = res
}

type PackagingLine struct {
	PuckSlot *Puck
}

func NewPackagingLine() PackagingLine {
	return PackagingLine{
		PuckSlot: nil,
	}
}

func (p *PackagingLine) PlacePuck(puck Puck) error {
	if p.PuckSlot != nil {
		return ErrSlotOccupied
	}

	if puck.IsPackaged {
		return ErrPuckPackaged
	}

	p.PuckSlot = &puck

	return nil
}

func (p *PackagingLine) TakePuck() (Puck, error) {
	if p.PuckSlot == nil {
		return Puck{}, ErrSlotEmpty
	}

	if !p.PuckSlot.IsPackaged {
		return Puck{}, errors.New("need to package puck before taking")
	}

	puck := *p.PuckSlot
	p.PuckSlot = nil

	return puck, nil
}

func (p *PackagingLine) PackagePuck() error {
	if p.PuckSlot == nil {
		return ErrSlotEmpty
	}

	if p.PuckSlot.IsPackaged {
		return ErrPuckPackaged
	}

	time.Sleep(packagingTime)

	p.PuckSlot.IsPackaged = true

	return nil
}

type SortingLine struct {
	PuckSlot *Puck
	Produced map[string][]Puck // color -> pucks
}

func NewSortingLine() SortingLine {
	return SortingLine{
		PuckSlot: nil,
		Produced: make(map[string][]Puck),
	}
}

func (s *SortingLine) PlacePuck(puck Puck) error {
	if s.PuckSlot != nil {
		return ErrSlotOccupied
	}

	if !puck.IsPackaged {
		return errors.New("need to package puck before placing")
	}

	s.PuckSlot = &puck

	return nil
}

func (s *SortingLine) TakePuck() (Puck, error) {
	if s.PuckSlot == nil {
		return Puck{}, ErrSlotEmpty
	}

	puck := *s.PuckSlot
	s.PuckSlot = nil

	return puck, nil
}

func (s *SortingLine) SortPuck() error {
	if s.PuckSlot == nil {
		return ErrSlotEmpty
	}

	time.Sleep(sortingTime)

	s.Produced[s.PuckSlot.Color] = append(s.Produced[s.PuckSlot.Color], *s.PuckSlot)
	s.PuckSlot = nil

	return nil
}
