package linuxcnc

type Status struct {
	EstopEnabled bool
	InPosition   bool
	MotionPaused bool
	Coolant      CoolantStatus
	Joints       []JointStatus
	Spindles     []SpindleStatus
}

type CoolantStatus struct {
	// Mist coolant enabled.
	Mist bool
	// Flood coolant enabled.
	Flood bool
}

type SpindleStatus struct {
	// Spindle break enabled.
	Break bool
	// Spindle enabled.
	Enabled bool
	// Spindle speed in rpm.
	// >0 clockwise
	// <0 counterclockwise
	Speed float64
	// Spindle override scale.
	OverrideScale float64
	// Spindle override enabled.
	OverrideEnabled bool
}

type JointType string

const (
	JointTypeLinear  JointType = "Linear"
	JointTypeAngular JointType = "Angular"
)

type JointHomingPhase string

const (
	// Unkown homing state of the joint.
	JointHomingPhaseUnknown JointHomingPhase = "Unknown"
	// Homing of the joint is in progress.
	JointHomingPhaseHoming JointHomingPhase = "Homing"
	// Joint is homed.
	JointHomingPhaseHomed JointHomingPhase = "Homed"
)

type JointStatus struct {
	// Type of joint configuration parameter, reflects [JOINT_n]TYPE.
	// Values: Linear, Angular
	// See Joint ini configuration for details.
	Type JointType

	// Joint homing status.
	// Values: Unkown, Homing, Homed.
	HomingPhase JointHomingPhase

	//  Backlash in machine units. configuration parameter, reflects
	//  [JOINT_n]BACKLASH.
	Backlash float64

	// Joint enabled.
	Enabled bool

	// Joint in position.
	InPosition bool
	// Current input position.
	InputPosition float64
	// Commanded output position.
	OutputPosition float64
	// Current joint velocity.
	Velocity float64

	// Joint limits configuration and status.
	Limits JointLimitsStatus
}

type JointLimitsStatus struct {
	// Limit override enabled.
	OverrideEnabled bool

	// Maximum hard limit exceeded.
	HardMaxExceeded bool
	// Maximum soft limit exceeded.
	SoftMaxExceeded bool
	// Maximum soft limit for joint motion.
	// Reflects [JOINT_n]MAX_LIMIT.
	MaxPosition float64

	// Minimum hard limit exceeded.
	HardMinExceeded bool
	// Minimum soft limit exceeded.
	SoftMinExceeded bool
	// Minimum soft limit for joint motion.
	// Reflects [JOINT_n]MIN_LIMIT.
	MinPosition float64
}
