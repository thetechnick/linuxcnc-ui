#include "linuxcnc.hh"
#include <stdlib.h>

StatHandle stat_newHandle() { return NULL; }

int stat_initHandle(StatHandle handle) { return 0; }
int stat_poll(StatHandle handle) { return 0; }
void stats_global(StatHandle handle, stat_Global *stats) {}
void stats_spindles(StatHandle handle, stat_Spindle *stats) {}
void stats_joints(StatHandle handle, stat_Joint *stats) {}
