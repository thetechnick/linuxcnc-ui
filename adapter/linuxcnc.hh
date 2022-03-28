#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// ----------------
// Stats
// ----------------

typedef void *StatHandle;

// Returns a new stat handle to LinuxCNC.
StatHandle stat_newHandle();

// Initialize the connection to LinuxCNC.
int stat_initHandle(StatHandle handle);

// Load status over established connection.
int stat_poll(StatHandle handle);

typedef struct stat_Pos {
  double x, y, z;
  double a, b, c;
  double u, v, w;
} stat_Pos;

typedef struct stat_Global {
  // estop off=0 or on=1
  int estop;
  // thread period.
  double cycleTime;
  int trajectoryPlannerEnabled;
  // Currently loaded gcode filename with path.
  char *file;
  // source line number motion is currently executing. Relation to id unclear.
  int fileLineNumber;

  // Machine in position flag.
  bool inpos;
  // Motion paused.
  bool motionPaused;
  int motionType;

  // mist coolant off=0 or on=1
  int mist;
  // flood coolant off=0 or on=1
  int flood;

  // current tool number.
  int toolInSpindle;
  // The index into the stat.tool_table list of the tool currently prepped
  // for
  // tool change, or -1 no tool is prepped.  On a Random toolchanger this is the
  // same as the tool's pocket number.  On a Non-random toolchanger it's a
  // random small integer.
  int pocketPrepped;

  // Number of defined spindles. Reflects [TRAJ]SPINDLES ini value.
  int numberOfSpindles;
  // number of defined joints. Reflects [KINS]JOINTS ini value.
  int numberOfJoints;

  int feedOverrideEnabled;
  // current feedrate override, 1.0 = 100%
  double feedOverride;
  int feedHoldEnabled;

  // Positions

  // Last position where the probe was tripped.
  stat_Pos probedPosition;
  // current commanded position
  stat_Pos inputPosition;
  // current actual position, from forward kins
  stat_Pos currentPosition;
  // Distance to go in current move.
  stat_Pos distanceToGo;
} stat_Global;

void stats_global(StatHandle handle, stat_Global *stats);

typedef struct stat_Spindle {
  // spindle brake flag.
  int brake;
  // rotational direction of the spindle. forward=1, reverse=-1.
  int direction;
  // value of the spindle enabled flag.
  int enabled;
  // spindle speed override scale.
  double override;
  // value of the spindle override enabled flag.
  bool overrideEnabled;
  // spindle speed value, rpm, > 0: clockwise, < 0: counterclockwise.
  double speed;
} stat_Spindle;

void stats_spindles(StatHandle handle, stat_Spindle *stats);

typedef struct stat_Joint {
  //  Backlash in machine units. configuration parameter, reflects
  //  [JOINT_n]BACKLASH.
  double backlash;
  // non-zero means enabled.
  int enabled;
  // non-zero means axis amp fault.
  int fault;
  // current following error.
  double ferrorCurrent;
  // magnitude of max following error.
  double ferrorHighMark;
  // non-zero means has been homed.
  int homed;
  // non-zero means homing in progress.
  int homing;
  // non-zero means in position.
  int inpos;
  // current input position.
  double input;
  // type of axis configuration parameter, reflects
  // [JOINT_n]TYPE. LINEAR=1, ANGULAR=2. See Joint ini configuration for
  // details.
  int jointType;
  // maximum following error. configuration parameter,
  // reflects [JOINT_n]FERROR.
  double maxFerror;
  // non-zero means max hard limit exceeded.
  int maxHardLimit;
  // maximum limit (soft limit) for joint motion, in machine
  // units. configuration parameter, reflects [JOINT_n]MAX_LIMIT.
  double maxPositionLimit;
  // non-zero means maxPositionLimit was exceeded.
  int maxSoftLimit;
  // configuration parameter, reflects [JOINT_n]MIN_FERROR.
  double minFerror;
  // non-zero means min hard limit exceeded.
  int minHardLimit;
  // minimum limit (soft limit) for joint motion, in machine
  // units. configuration parameter, reflects [JOINT_n]MIN_LIMIT.
  double minPositionLimit;
  // non-zero means minPositionLimit was exceeded.
  int minSoftLimit;
  // commanded output position.
  double output;
  // non-zero means limits are overridden.
  int overrideLimits;
  // joint units per mm, or per degree for angular joints.
  // (joint units are the same as machine units, unless set otherwise by the
  // configuration parameter [JOINT_n]UNITS)
  double units;
  // current velocity.
  double velocity;
} stat_Joint;

void stats_joints(StatHandle handle, stat_Joint *stats);

#ifdef __cplusplus
}
#endif
