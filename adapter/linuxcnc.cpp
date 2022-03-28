#define RTAPI

#include "linuxcnc.hh"

#include "config.h"
#include "emc.hh"
#include "emc_nml.hh"
#include "inifile.hh"
#include "kinematics.h"
#include "nml_oi.hh"
#include "rcs.hh"
#include "rcs_print.hh"
#include "rtapi_string.h"
#include "timer.hh"

#include "tooldata.hh"

// -------------
// Stat Handling
// -------------

struct statHandle {
  RCS_STAT_CHANNEL *c;
  EMC_STAT status;
  bool initialized;
};

StatHandle stat_newHandle() {
  statHandle *h = new statHandle;
  h->initialized = 0;
  return h;
}

int stat_initHandle(StatHandle handle) {
  statHandle *s = (statHandle *)handle;
  s->c = new RCS_STAT_CHANNEL(emcFormat, "emcStatus", "xemc",
                              EMC2_DEFAULT_NMLFILE);
  if (!s->c) {
    // error
    return -1;
  }
  return 0;
}

int stat_poll(StatHandle handle) {
  statHandle *s = (statHandle *)handle;
#ifdef TOOL_NML
  if (!s->initialized) {
    // fprintf(stderr,"%8d tool_nml_register\n",getpid());
    tool_nml_register((CANON_TOOL_TABLE *)&s->status.io.tool.toolTable);
    s->initialized = 1;
  }
#else
  static bool mmap_available = 1;
  if (!mmap_available)
    return 0;
  if (!s->initialized) {
    s->initialized = 1;
    if (tool_mmap_user()) {
      mmap_available = 0;
      fprintf(stderr, "mmap tool data not available, continuing %s\n",
              __FILE__);
    }
  }
#endif
  if (!s->c->valid()) {
    return -1;
  }
  if (s->c->peek() == EMC_STAT_TYPE) {
    EMC_STAT *emcStatus = static_cast<EMC_STAT *>(s->c->get_address());
    memcpy((char *)&s->status, emcStatus, sizeof(EMC_STAT));
  }
  return 0;
}

void pos(stat_Pos *sp, const EmcPose &p) {
  sp->x = p.tran.x;
  sp->y = p.tran.y;
  sp->z = p.tran.z;
  sp->a = p.a;
  sp->b = p.b;
  sp->c = p.c;
  sp->u = p.u;
  sp->v = p.v;
  sp->w = p.w;
}

void stats_global(StatHandle handle, stat_Global *stats) {
  statHandle *s = (statHandle *)handle;
  // State
  stats->estop = s->status.io.aux.estop;
  stats->cycleTime = s->status.motion.traj.cycleTime;
  stats->trajectoryPlannerEnabled = s->status.motion.traj.enabled;
  stats->file = s->status.io.source_file;
  stats->fileLineNumber = s->status.task.currentLine;

  stats->motionType = s->status.motion.traj.motion_type;

  // Feed
  stats->feedHoldEnabled = s->status.motion.traj.feed_hold_enabled;
  stats->feedOverrideEnabled = s->status.motion.traj.feed_override_enabled;
  stats->feedOverride = s->status.motion.traj.scale;

  // Movement
  stats->inpos = s->status.motion.traj.inpos;
  stats->motionPaused = s->status.motion.traj.paused;

  // Tool
  stats->toolInSpindle = s->status.io.tool.toolInSpindle;
  stats->pocketPrepped = s->status.io.tool.pocketPrepped;

  // Coolant
  stats->mist = s->status.io.coolant.mist;
  stats->flood = s->status.io.coolant.flood;

  // Machine
  stats->numberOfJoints = s->status.motion.traj.joints;
  stats->numberOfSpindles = s->status.motion.traj.spindles;

  // Positions
  pos(&stats->probedPosition, s->status.motion.traj.probedPosition);
  pos(&stats->inputPosition, s->status.motion.traj.position);
  pos(&stats->currentPosition, s->status.motion.traj.actualPosition);
  pos(&stats->distanceToGo, s->status.motion.traj.dtg);
}

void stats_spindles(StatHandle handle, stat_Spindle *stats) {
  statHandle *s = (statHandle *)handle;

  int spindles = s->status.motion.traj.spindles;
  for (int i = 0; i < spindles; i++) {
    stats[i].brake = s->status.motion.spindle[i].brake;
    stats[i].direction = s->status.motion.spindle[i].direction;
    stats[i].enabled = s->status.motion.spindle[i].enabled;
    stats[i].override = s->status.motion.spindle[i].spindle_scale;
    stats[i].overrideEnabled =
        s->status.motion.spindle[i].spindle_override_enabled;
    stats[i].speed = s->status.motion.spindle[i].speed;
  }
}

void stats_joints(StatHandle handle, stat_Joint *stats) {
  statHandle *s = (statHandle *)handle;

  int joints = s->status.motion.traj.joints;
  for (int i = 0; i < joints; i++) {
    stats[i].backlash = s->status.motion.joint[i].backlash;
    stats[i].enabled = s->status.motion.joint[i].enabled;
    stats[i].fault = s->status.motion.joint[i].fault;
    stats[i].ferrorCurrent = s->status.motion.joint[i].ferrorCurrent;
    stats[i].ferrorHighMark = s->status.motion.joint[i].ferrorHighMark;
    stats[i].homed = s->status.motion.joint[i].homed;
    stats[i].homing = s->status.motion.joint[i].homing;
    stats[i].inpos = s->status.motion.joint[i].inpos;
    stats[i].input = s->status.motion.joint[i].input;
    stats[i].jointType = s->status.motion.joint[i].jointType;
    stats[i].maxFerror = s->status.motion.joint[i].maxFerror;
    stats[i].maxHardLimit = s->status.motion.joint[i].maxHardLimit;
    stats[i].maxPositionLimit = s->status.motion.joint[i].maxPositionLimit;
    stats[i].maxSoftLimit = s->status.motion.joint[i].maxSoftLimit;
    stats[i].minFerror = s->status.motion.joint[i].minFerror;
    stats[i].minHardLimit = s->status.motion.joint[i].minHardLimit;
    stats[i].minPositionLimit = s->status.motion.joint[i].minPositionLimit;
    stats[i].minSoftLimit = s->status.motion.joint[i].minSoftLimit;
    stats[i].output = s->status.motion.joint[i].output;
    stats[i].overrideLimits = s->status.motion.joint[i].overrideLimits;
    stats[i].units = s->status.motion.joint[i].units;
    stats[i].velocity = s->status.motion.joint[i].velocity;
  }
}
