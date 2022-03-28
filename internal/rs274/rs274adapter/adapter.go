package rs274adapter

/*
#include <stdlib.h>
#include "rs274.hh"

extern void ErrorAdpt(int, int);
extern bool AbortAdpt();
extern void MessageAdpt(char *comment);
extern void CommentAdpt(char *comment);
extern void ChangeToolAdpt(int pocket);

extern void UseLengthUnitsAdpt(int units); // CANON_UNITS enum
extern void UseToolLengthOffsetAdpt(double x, double y, double z, double a,
                                    double b, double c, double u, double v,
                                    double w);
extern void SelectPlaneAdpt(int plane);
extern void SetXYRotationAdpt(double);
extern void SetG5XOffsetAdpt(int index, double x, double y, double z, double a,
                             double b, double c, double u, double v, double w);
extern void SetG92OffsetAdpt(double x, double y, double z, double a, double b,
                             double c, double u, double v, double w);
extern void SetTraverseRateAdpt(double rate);
extern void SetFeedModeAdpt(int spindle, int mode);
extern void SetFeedRateAdpt(double rate);

extern void StraightTraverseAdpt(int lineNo, double x, double y, double z,
                                 double a, double b, double c, double u,
                                 double v, double w);
extern void ArcFeedAdpt(double firstEnd, double secondEnd, double firstAxis,
                        double secondAxis, int rotation, double axisEndPoint,
                        double aPosition, double bPosition, double cPosition,
                        double uPosition, double vPosition, double wPosition);
extern void StraightFeedAdpt(int lineNo, double x, double y, double z, double a,
                             double b, double c, double u, double v, double w);
extern void StraightProbeAdpt(int lineNo, double x, double y, double z,
                              double a, double b, double c, double u, double v,
                              double w);
extern void RigidTapAdpt(int lineNo, double x, double y, double z,
                         double scale);
extern void DwellAdpt(double seconds);
*/
// #cgo CPPFLAGS: -I${SRCDIR}/../../../adapter
import "C"
import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/thetechnick/linuxcnc-ui/internal/rs274"
)

var d = &Interpreter{
	sink: &rs274.DiscardInterpreterSink{},
}

func init() {
	cb := C.Callbacks{
		error:      C.ErrorFn(C.ErrorAdpt),
		abort:      C.AbortFn(C.AbortAdpt),
		message:    C.MessageFn(C.MessageAdpt),
		comment:    C.CommentFn(C.CommentAdpt),
		changeTool: C.ChangeToolFn(C.ChangeToolAdpt),

		useLengthUnits:      C.UseLengthUnitsFn(C.UseLengthUnitsAdpt),
		useToolLengthOffset: C.UseToolLengthOffsetFn(C.UseToolLengthOffsetAdpt),
		selectPlane:         C.SelectPlaneFn(C.SelectPlaneAdpt),
		setXYRotation:       C.SetXYRotationFn(C.SetXYRotationAdpt),
		setG5XOffset:        C.SetG5XOffsetFn(C.SetG5XOffsetAdpt),
		setG92Offset:        C.SetG92OffsetFn(C.SetG92OffsetAdpt),
		setTraverseRate:     C.SetTraverseRateFn(C.SetTraverseRateAdpt),
		setFeedMode:         C.SetFeedModeFn(C.SetFeedModeAdpt),
		setFeedRate:         C.SetFeedRateFn(C.SetFeedRateAdpt),

		straightTraverse: C.StraightTraverseFn(C.StraightTraverseAdpt),
		arcFeed:          C.ArcFeedFn(C.ArcFeedAdpt),
		straightFeed:     C.StraightFeedFn(C.StraightFeedAdpt),
		straightProbe:    C.StraightProbeFn(C.StraightProbeAdpt),
		rigidTap:         C.RigidTapFn(C.RigidTapAdpt),
		dwell:            C.DwellFn(C.DwellAdpt),
	}
	C.registerCallbacks(cb)
}

type Interpreter struct {
	lock sync.Mutex

	sink      rs274.InterpreterSink
	abort     bool
	abortLock sync.RWMutex

	lastErr error
}

func (p *Interpreter) Parse(
	ctx context.Context,
	sink rs274.InterpreterSink, filename string,
) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.lastErr = nil
	p.abort = false

	doneCh := make(chan struct{})
	defer close(doneCh)
	go func() {
		select {
		case <-doneCh:
			return
		case <-ctx.Done():
			p.abortLock.Lock()
			defer p.abortLock.Lock()
			p.abort = true
		}
	}()

	// C.parseFile(cfilename)
	return p.lastErr
}

func (p *Interpreter) reportError(interpreterError, lastSequenceNumber int) {
	p.lastErr = &rs274.InterpreterError{
		InterpreterError:   interpreterError,
		LastSequenceNumber: lastSequenceNumber,
	}
}

func (p *Interpreter) checkAbort() bool {
	p.abortLock.RLock()
	defer p.abortLock.RUnlock()
	return p.abort
}

// ----
// Util
// ----

//export errorGo
func errorGo(interpError, lastSequenceNumber C.int) {
	d.reportError(int(interpError), int(lastSequenceNumber))
}

//export abortGo
func abortGo() C.bool {
	return C.bool(d.checkAbort())
}

//export messageGo
func messageGo(message *C.char) {
	d.sink.Message(C.GoString(message))
}

//export commentGo
func commentGo(message *C.char) {
	d.sink.Comment(C.GoString(message))
}

//export changeToolGo
func changeToolGo(pocket C.int) {
	d.sink.ChangeTool(int(pocket))
}

// --------
// Settings
// --------

//export useLengthUnitsGo
func useLengthUnitsGo(units C.int) {
	var u rs274.LengthUnit
	switch int(units) {
	case 1:
		u = rs274.LengthInches
	case 2:
		u = rs274.LengthMillimeter
	case 3:
		u = rs274.LengthCentimeter
	default:
		panic("unknown length unit")
	}
	d.sink.UseLengthUnit(u)
}

//export useToolLengthOffsetGo
func useToolLengthOffsetGo(x, y, z, a, b, c, u, v, w C.double) {
	d.sink.UseToolLengthOffset(rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export selectPlaneGo
func selectPlaneGo(plane C.int) {
	var p rs274.Plane
	switch int(plane) {
	case 1:
		p = rs274.PlaneXY
	case 2:
		p = rs274.PlaneYZ
	case 3:
		p = rs274.PlaneXZ
	case 4:
		p = rs274.PlaneUV
	case 5:
		p = rs274.PlaneVW
	case 6:
		p = rs274.PlaneUW
	default:
		panic("unknown plane")
	}
	d.sink.SelectPlane(p)
}

//export setXYRotationGo
func setXYRotationGo(rotation C.double) {
	d.sink.SetXYRotation(float64(rotation))
}

//export setG5XOffsetGo
func setG5XOffsetGo(index C.int, x, y, z, a, b, c, u, v, w C.double) {
	d.sink.SetG5XOffset(int(index), rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export setG92OffsetGo
func setG92OffsetGo(x, y, z, a, b, c, u, v, w C.double) {
	d.sink.SetG92Offset(rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export setTraverseRateGo
func setTraverseRateGo(rate C.double) {
	d.sink.SetTraverseRate(float64(rate))
}

//export setFeedModeGo
func setFeedModeGo(spindle, mode C.int) {
	d.sink.SetFeedMode(int(spindle), int(mode))
}

//export setFeedRateGo
func setFeedRateGo(rate C.double) {
	d.sink.SetFeedRate(float64(rate))
}

// --------
// Movement
// --------

//export straightTraverseGo
func straightTraverseGo(lineNo C.int, x, y, z, a, b, c, u, v, w C.double) {
	d.sink.StraightTraverse(int(lineNo), rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export arcFeedGo
func arcFeedGo(firstEnd, secondEnd, firstAxis, secondAxis C.double, rotation C.int, axisEndPoint, aPosition, bPosition, cPosition, uPosition, vPosition, wPosition C.double) {
	d.sink.ArcFeed(rs274.ArcMove{
		FirstEnd:     float64(firstEnd),
		SecondEnd:    float64(secondEnd),
		FirstAxis:    float64(firstAxis),
		SecondAxis:   float64(secondAxis),
		Rotation:     int(rotation),
		AxisEndPoint: float64(axisEndPoint),

		A: float64(aPosition), B: float64(bPosition), C: float64(cPosition),
		U: float64(uPosition), V: float64(vPosition), W: float64(wPosition),
	})
}

//export straightFeedGo
func straightFeedGo(lineNo C.int, x, y, z, a, b, c, u, v, w C.double) {
	d.sink.StraightFeed(int(lineNo), rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export straightProbeGo
func straightProbeGo(lineNo C.int, x, y, z, a, b, c, u, v, w C.double) {
	d.sink.StraightProbe(int(lineNo), rs274.Position{
		X: float64(x), Y: float64(y), Z: float64(z),
		A: float64(a), B: float64(b), C: float64(c),
		U: float64(u), V: float64(v), W: float64(w),
	})
}

//export rigidTapGo
func rigidTapGo(lineNo C.int, x, y, z, scale C.double) {
	d.sink.RigidTap(
		int(lineNo), float64(x), float64(y), float64(z), float64(scale))
}

//export dwellGo
func dwellGo(seconds C.double) {
	d.sink.Dwell(time.Duration(float64(seconds) * float64(time.Second)))
}

func DoStuff() {
	cb := C.Callbacks{}
	cb.error = C.ErrorFn(C.ErrorAdpt)
	C.registerCallbacks(cb)

	f := "xxx"
	var cfilename *C.char = C.CString(f)
	defer C.free(unsafe.Pointer(cfilename))

	C.parseFile(cfilename)
}
