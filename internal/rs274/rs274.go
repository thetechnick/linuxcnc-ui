package rs274

import (
	"fmt"
	"time"
)

type Plane string

const (
	PlaneXY Plane = "XY"
	PlaneYZ Plane = "YZ"
	PlaneXZ Plane = "XZ"
	PlaneUV Plane = "UV"
	PlaneVW Plane = "VW"
	PlaneUW Plane = "UW"
)

type LengthUnit string

const (
	LengthInches     LengthUnit = "Inches"
	LengthMillimeter LengthUnit = "Millimeter"
	LengthCentimeter LengthUnit = "Centimeter"
)

type Position struct {
	X, Y, Z float64
	A, B, C float64
	U, V, W float64
}

type ArcMove struct {
	FirstEnd, SecondEnd, FirstAxis, SecondAxis float64
	Rotation                                   int
	AxisEndPoint                               float64
	A, B, C, U, V, W                           float64
}

type InterpreterError struct {
	InterpreterError   int
	LastSequenceNumber int
}

func (e *InterpreterError) Error() string {
	return fmt.Sprintf("interpreter error %d at sequence number %d",
		e.InterpreterError, e.LastSequenceNumber)
}

type InterpreterSink interface {
	Message(message string)
	Comment(comment string)
	ChangeTool(toolNumber int)

	UseLengthUnit(unit LengthUnit)
	UseToolLengthOffset(pos Position)
	SelectPlane(planeNumber Plane)
	SetXYRotation(position float64)
	SetG5XOffset(index int, pos Position)
	SetG92Offset(pos Position)
	SetTraverseRate(rate float64)
	SetFeedMode(spindleNumber int, mode int)
	SetFeedRate(rate float64)

	StraightTraverse(lineNumber int, pos Position)
	ArcFeed(m ArcMove)
	StraightFeed(lineNumber int, pos Position)
	StraightProbe(lineNumber int, pos Position)
	RigidTap(lineNumber int, X, Y, Z, scale float64)
	Dwell(time.Duration)
}

type InterpreterSinkFuncs struct {
	MessageFn    func(message string)
	CommentFn    func(comment string)
	ChangeToolFn func(toolNumber int)

	UseLengthUnitFn       func(unit LengthUnit)
	UseToolLengthOffsetFn func(pos Position)
	SelectPlaneFn         func(planeNumber Plane)
	SetXYRotationFn       func(position float64)
	SetG5XOffsetFn        func(index int, pos Position)
	SetG92OffsetFn        func(pos Position)
	SetTraverseRateFn     func(rate float64)
	SetFeedModeFn         func(spindleNumber int, mode int)
	SetFeedRateFn         func(rate float64)

	StraightTraverseFn func(lineNumber int, pos Position)
	ArcFeedFn          func(m ArcMove)
	StraightFeedFn     func(lineNumber int, pos Position)
	StraightProbeFn    func(lineNumber int, pos Position)
	RigidTapFn         func(lineNumber int, X, Y, Z, scale float64)
	DwellFn            func(time.Duration)
}

func (f *InterpreterSinkFuncs) Message(message string) {
	if f.MessageFn != nil {
		f.MessageFn(message)
	}
}

func (f *InterpreterSinkFuncs) Comment(comment string) {
	if f.CommentFn != nil {
		f.CommentFn(comment)
	}
}

func (f *InterpreterSinkFuncs) ChangeTool(toolNumber int) {
	if f.ChangeToolFn != nil {
		f.ChangeToolFn(toolNumber)
	}
}

func (f *InterpreterSinkFuncs) UseLengthUnit(unit LengthUnit) {
	if f.UseLengthUnitFn != nil {
		f.UseLengthUnitFn(unit)
	}
}

func (f *InterpreterSinkFuncs) UseToolLengthOffset(pos Position) {
	if f.UseToolLengthOffsetFn != nil {
		f.UseToolLengthOffsetFn(pos)
	}
}

func (f *InterpreterSinkFuncs) SelectPlane(planeNumber Plane) {
	if f.SelectPlaneFn != nil {
		f.SelectPlaneFn(planeNumber)
	}
}

func (f *InterpreterSinkFuncs) SetXYRotation(position float64) {
	if f.SetXYRotationFn != nil {
		f.SetXYRotationFn(position)
	}
}

func (f *InterpreterSinkFuncs) SetG5XOffset(index int, pos Position) {
	if f.SetG5XOffsetFn != nil {
		f.SetG5XOffsetFn(index, pos)
	}
}

func (f *InterpreterSinkFuncs) SetG92Offset(pos Position) {
	if f.SetG92OffsetFn != nil {
		f.SetG92OffsetFn(pos)
	}
}

func (f *InterpreterSinkFuncs) SetTraverseRate(rate float64) {
	if f.SetTraverseRateFn != nil {
		f.SetTraverseRateFn(rate)
	}
}

func (f *InterpreterSinkFuncs) SetFeedMode(spindleNumber int, mode int) {
	if f.SetFeedModeFn != nil {
		f.SetFeedModeFn(spindleNumber, mode)
	}
}

func (f *InterpreterSinkFuncs) SetFeedRate(rate float64) {
	if f.SetFeedRateFn != nil {
		f.SetFeedRateFn(rate)
	}
}

func (f *InterpreterSinkFuncs) StraightTraverse(lineNumber int, pos Position) {
	if f.StraightTraverseFn != nil {
		f.StraightTraverseFn(lineNumber, pos)
	}
}

func (f *InterpreterSinkFuncs) ArcFeed(m ArcMove) {
	if f.ArcFeedFn != nil {
		f.ArcFeedFn(m)
	}
}

func (f *InterpreterSinkFuncs) StraightFeed(lineNumber int, pos Position) {
	if f.StraightFeedFn != nil {
		f.StraightFeedFn(lineNumber, pos)
	}
}

func (f *InterpreterSinkFuncs) StraightProbe(lineNumber int, pos Position) {
	if f.StraightProbeFn != nil {
		f.StraightProbeFn(lineNumber, pos)
	}
}

func (f *InterpreterSinkFuncs) RigidTap(lineNumber int, x, y, z, scale float64) {
	if f.RigidTapFn != nil {
		f.RigidTapFn(lineNumber, x, y, z, scale)
	}
}

func (f *InterpreterSinkFuncs) Dwell(d time.Duration) {
	if f.DwellFn != nil {
		f.DwellFn(d)
	}
}

type DiscardInterpreterSink struct{}

func (f *DiscardInterpreterSink) Message(message string) {}

func (f *DiscardInterpreterSink) Comment(comment string) {}

func (f *DiscardInterpreterSink) ChangeTool(toolNumber int) {}

func (f *DiscardInterpreterSink) UseLengthUnit(unit LengthUnit) {}

func (f *DiscardInterpreterSink) UseToolLengthOffset(pos Position) {}

func (f *DiscardInterpreterSink) SelectPlane(planeNumber Plane) {}

func (f *DiscardInterpreterSink) SetXYRotation(position float64) {}

func (f *DiscardInterpreterSink) SetG5XOffset(index int, pos Position) {}

func (f *DiscardInterpreterSink) SetG92Offset(pos Position) {}

func (f *DiscardInterpreterSink) SetTraverseRate(rate float64) {}

func (f *DiscardInterpreterSink) SetFeedMode(spindleNumber int, mode int) {}

func (f *DiscardInterpreterSink) SetFeedRate(rate float64) {}

func (f *DiscardInterpreterSink) StraightTraverse(lineNumber int, pos Position) {}

func (f *DiscardInterpreterSink) ArcFeed(m ArcMove) {}

func (f *DiscardInterpreterSink) StraightFeed(lineNumber int, pos Position) {}

func (f *DiscardInterpreterSink) StraightProbe(lineNumber int, pos Position) {}

func (f *DiscardInterpreterSink) RigidTap(lineNumber int, x, y, z, scale float64) {}

func (f *DiscardInterpreterSink) Dwell(d time.Duration) {}
