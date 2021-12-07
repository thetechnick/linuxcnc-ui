package linuxcnc

// #include <stdlib.h>
// #include "linuxcnc.hh"
// #cgo CPPFLAGS: -I${SRCDIR}/../../adapter
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

var (
	ErrStatsHandleInit = errors.New(
		"stats handle could not be initialized")
	ErrStatsPoll = errors.New(
		"stats could not be polled")
)

type Client struct {
	statHandle C.StatHandle
	closeOnce  sync.Once
}

type ClientOptions struct{}

type ClientOption func(o *ClientOptions)

func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		statHandle: C.stat_newHandle(),
	}
	if r := C.stat_initHandle(c.statHandle); r != 0 {
		c.Close()
		return nil, ErrStatsHandleInit
	}
	return c, nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		C.free(unsafe.Pointer(c.statHandle))
	})
}

// Polls for new status and loads it into the given pointer.
func (c *Client) PollStatus(s *Status) error {
	if r := C.stat_poll(c.statHandle); r != 0 {
		return ErrStatsPoll
	}

	// Global
	global := C.struct_stat_Global{}
	C.stats_global(c.statHandle, &global)

	s.EstopEnabled = global.estop != 0
	s.InPosition = bool(global.inpos)
	s.MotionPaused = bool(global.motionPaused)
	s.Coolant.Flood = global.flood != 0
	s.Coolant.Mist = global.mist != 0

	// Joints
	if len(s.Joints) != int(global.numberOfJoints) {
		s.Joints = make([]JointStatus, global.numberOfJoints)
	}
	joints := make([]C.struct_stat_Joint, global.numberOfJoints)
	C.stats_joints(c.statHandle, &joints[0])
	for i, joint := range joints {
		var jointType JointType
		switch joint.jointType {
		case 1:
			jointType = JointTypeLinear
		case 2:
			jointType = JointTypeAngular
		}

		var jointHomingPhase JointHomingPhase
		switch {
		case joint.homed != 0:
			jointHomingPhase = JointHomingPhaseHomed
		case joint.homing != 0:
			jointHomingPhase = JointHomingPhaseHoming
		default:
			jointHomingPhase = JointHomingPhaseUnknown
		}

		s.Joints[i] = JointStatus{
			Type:           jointType,
			HomingPhase:    jointHomingPhase,
			Backlash:       float64(joint.backlash),
			Enabled:        joint.enabled != 0,
			InPosition:     joint.inpos != 0,
			InputPosition:  float64(joint.input),
			OutputPosition: float64(joint.output),
			Velocity:       float64(joint.velocity),

			Limits: JointLimitsStatus{
				OverrideEnabled: joint.overrideLimits != 0,

				HardMaxExceeded: joint.maxHardLimit != 0,
				SoftMaxExceeded: joint.maxSoftLimit != 0,
				MaxPosition:     float64(joint.maxPositionLimit),

				HardMinExceeded: joint.minHardLimit != 0,
				SoftMinExceeded: joint.minSoftLimit != 0,
				MinPosition:     float64(joint.minPositionLimit),
			},
		}
	}

	// Spindles
	if len(s.Spindles) != int(global.numberOfSpindles) {
		s.Spindles = make([]SpindleStatus, global.numberOfSpindles)
	}
	spindles := make([]C.struct_stat_Spindle, global.numberOfSpindles)
	C.stats_spindles(c.statHandle, &spindles[0])
	for i, spindle := range spindles {
		s.Spindles[i] = SpindleStatus{
			Break:   spindle.brake != 0,
			Enabled: spindle.enabled != 0,
			Speed:   float64(spindle.speed),

			OverrideScale:   float64(spindle.override),
			OverrideEnabled: bool(spindle.overrideEnabled),
		}
	}
	return nil
}
