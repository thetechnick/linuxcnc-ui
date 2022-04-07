

#include "rs274.hh"
#include "canon.hh"
#include "inifile.hh"
#include "interp_return.hh"
#include "rs274ngc.hh"
#include "rs274ngc_interp.hh"

InterpBase *pinterp;
static int interp_error;
static int last_sequence_number;
static double _pos_x, _pos_y, _pos_z, _pos_a, _pos_b, _pos_c, _pos_u, _pos_v,
    _pos_w;
EmcPose tool_offset;
CANON_UNITS _length_unit_type = CANON_UNITS_MM;
int _task = 0; // control preview behaviour when remapping
struct _inittab builtin_modules[] = {};

//----------------------
// Callback Registration
//----------------------

static Callbacks callbacks;

// register callback functions for the parser to call.
void registerCallbacks(Callbacks cb) { callbacks = cb; }

#define RESULT_OK (result == INTERP_OK || result == INTERP_EXECUTE_FINISH)

int parseFile(char *file) {
  struct timeval t0, t1;
  int wait = 1;
  int error_line_offset = 0;

  if (pinterp) {
    delete pinterp;
    pinterp = 0;
  }
  if (!pinterp) {
    pinterp = new Interp;
  }

  // for (int i = 0; i < USER_DEFINED_FUNCTION_NUM; i++)
  //   USER_DEFINED_FUNCTION[i] = user_defined_function;

  gettimeofday(&t0, NULL);

  interp_error = 0;
  last_sequence_number = -1;

  _pos_x = _pos_y = _pos_z = _pos_a = _pos_b = _pos_c = 0;
  _pos_u = _pos_v = _pos_w = 0;

  pinterp->init();
  pinterp->open(file);

  int result = INTERP_OK;

  // if (initcodes) {
  //   for (int i = 0; i < PyList_Size(initcodes) && RESULT_OK; i++) {
  //     PyObject *item = PyList_GetItem(initcodes, i);
  //     if (!item)
  //       return NULL;
  //     const char *code = PyUnicode_AsUTF8(item);
  //     if (!code)
  //       return NULL;
  //     result = pinterp->read(code);
  //     if (!RESULT_OK)
  //       goto out_error;
  //     result = pinterp->execute();
  //   }
  // }
  while (!interp_error && RESULT_OK) {
    error_line_offset = 1;
    result = pinterp->read();
    gettimeofday(&t1, NULL);
    if (t1.tv_sec > t0.tv_sec + wait) {
      if (callbacks.abort())
        return 0;
      t0 = t1;
    }
    if (!RESULT_OK)
      break;
    error_line_offset = 0;
    result = pinterp->execute();
  }

out_error:
  if (pinterp) {
    auto interp = dynamic_cast<Interp *>(pinterp);
    if (interp)
      interp->_setup.use_lazy_close = false;
    pinterp->close();
  }
  if (interp_error) {
    callbacks.error(interp_error, last_sequence_number);
    return interp_error;
  }
  return 0;
}

//----------------------------
// RS274 Interpreter Functions
//----------------------------

void SET_XY_ROTATION(double t) {
  if (interp_error)
    return;
  callbacks.setXYRotation(t);
}

void SET_G5X_OFFSET(int index, double x, double y, double z, double a, double b,
                    double c, double u, double v, double w) {
  if (interp_error)
    return;
  callbacks.setG5XOffset(index, x, y, z, a, b, c, u, v, w);
}

void SET_G92_OFFSET(double x, double y, double z, double a, double b, double c,
                    double u, double v, double w) {
  if (interp_error)
    return;
  callbacks.setG92Offset(x, y, z, a, b, c, u, v, w);
}

void USE_LENGTH_UNITS(CANON_UNITS inUnit) {
  if (interp_error)
    return;
  callbacks.useLengthUnits((int)inUnit);
}

void SET_TRAVERSE_RATE(double rate) {
  if (interp_error)
    return;
  callbacks.setTraverseRate(rate);
}

void STRAIGHT_TRAVERSE(int line_number, double x, double y, double z,
                       double a /*AA*/, double b /*BB*/, double c /*CC*/,
                       double u, double v, double w) {
  _pos_x = x;
  _pos_y = y;
  _pos_z = z;
  _pos_a = a;
  _pos_b = b;
  _pos_c = c;
  _pos_u = u;
  _pos_v = v;
  _pos_w = w;

  if (interp_error)
    return;

  callbacks.straightTraverse(line_number, x, y, z, a, b, c, u, v, w);
}

void SET_FEED_MODE(int spindle, int mode) {
  if (interp_error)
    return;
  callbacks.setFeedMode(spindle, mode);
}

void SET_FEED_RATE(double rate) {
  if (interp_error)
    return;
  callbacks.setFeedRate(rate);
}

void SET_FEED_REFERENCE(CANON_FEED_REFERENCE reference) {}

CANON_MOTION_MODE motion_mode;
void SET_MOTION_CONTROL_MODE(CANON_MOTION_MODE mode, double tolerance) {
  motion_mode = mode;
}

CANON_MOTION_MODE GET_EXTERNAL_MOTION_CONTROL_MODE() { return motion_mode; }

void SET_NAIVECAM_TOLERANCE(double tolerance) {}

void SELECT_PLANE(CANON_PLANE in_plane) {
  if (interp_error)
    return;
  callbacks.selectPlane((int)in_plane);
}

void SET_CUTTER_RADIUS_COMPENSATION(double radius) {}

void START_CUTTER_RADIUS_COMPENSATION(int side) {}

void STOP_CUTTER_RADIUS_COMPENSATION() {}

void START_SPEED_FEED_SYNCH() {}
void START_SPEED_FEED_SYNCH(int spindle, double sync, bool vel){};
void STOP_SPEED_FEED_SYNCH() {}

void MESSAGE(char *comment) {
  if (interp_error)
    return;
  callbacks.message(comment);
}

void COMMENT(const char *c) {
  if (interp_error)
    return;
  callbacks.comment(c);
}

void LOG(char *s) {}
void LOGOPEN(char *f) {}
void LOGAPPEND(char *f) {}
void LOGCLOSE() {}

//--------------------
// Machining Functions
//--------------------

void NURBS_FEED(int line_number,
                std::vector<CONTROL_POINT> nurbs_control_points,
                unsigned int k) {
  double u = 0.0;
  unsigned int n = nurbs_control_points.size() - 1;
  double umax = n - k + 2;
  unsigned int div = nurbs_control_points.size() * 15;
  std::vector<unsigned int> knot_vector = knot_vector_creator(n, k);
  PLANE_POINT P1;
  while (u + umax / div < umax) {
    PLANE_POINT P1 =
        nurbs_point(u + umax / div, k, nurbs_control_points, knot_vector);
    STRAIGHT_FEED(line_number, P1.X, P1.Y, _pos_z, _pos_a, _pos_b, _pos_c,
                  _pos_u, _pos_v, _pos_w);
    u = u + umax / div;
  }
  P1.X = nurbs_control_points[n].X;
  P1.Y = nurbs_control_points[n].Y;
  STRAIGHT_FEED(line_number, P1.X, P1.Y, _pos_z, _pos_a, _pos_b, _pos_c, _pos_u,
                _pos_v, _pos_w);
  knot_vector.clear();
}

void ARC_FEED(int line_number, double first_end, double second_end,
              double first_axis, double second_axis, int rotation,
              double axis_end_point, double a_position, double b_position,
              double c_position, double u_position, double v_position,
              double w_position) {
  if (interp_error)
    return;
  callbacks.arcFeed(first_end, second_end, first_axis, second_axis, rotation,
                    axis_end_point, a_position, b_position, c_position,
                    u_position, v_position, w_position);
}

void STRAIGHT_FEED(int line_number, double x, double y, double z, double a,
                   double b, double c, double u, double v, double w) {
  _pos_x = x;
  _pos_y = y;
  _pos_z = z;
  _pos_a = a;
  _pos_b = b;
  _pos_c = c;
  _pos_u = u;
  _pos_v = v;
  _pos_w = w;

  if (interp_error)
    return;

  callbacks.straightFeed(line_number, x, y, z, a, b, c, u, v, w);
}

void STRAIGHT_PROBE(int line_number, double x, double y, double z, double a,
                    double b, double c, double u, double v, double w,
                    unsigned char probe_type) {
  _pos_x = x;
  _pos_y = y;
  _pos_z = z;
  _pos_a = a;
  _pos_b = b;
  _pos_c = c;
  _pos_u = u;
  _pos_v = v;
  _pos_w = w;

  if (interp_error)
    return;
  callbacks.straightProbe(line_number, x, y, z, a, b, c, u, v, w);
}

void RIGID_TAP(int line_number, double x, double y, double z, double scale) {
  if (interp_error)
    return;

  callbacks.rigidTap(line_number, x, y, z, scale);
}

void DWELL(double seconds) {
  if (interp_error)
    return;

  callbacks.dwell(seconds);
}

//------------------
// Spindle Functions
//------------------

void SPINDLE_RETRACT_TRAVERSE() {}

void SET_SPINDLE_MODE(int spindle, double arg) {}

void START_SPINDLE_CLOCKWISE(int spindle, int wait_for_atspeed) {}

void START_SPINDLE_COUNTERCLOCKWISE(int spindle, int wait_for_atspeed) {}

void SET_SPINDLE_SPEED(int spindle, double rpm) {}

void STOP_SPINDLE_TURNING(int spindle) {}

void SPINDLE_RETRACT() {}

void ORIENT_SPINDLE(int spindle, double orientation, int mode) {}

void WAIT_SPINDLE_ORIENT_COMPLETE(int spindle, double timeout) {}

void USE_NO_SPINDLE_FORCE() {}

//---------------
// Tool Functions
//---------------

void SET_TOOL_TABLE_ENTRY(int idx, int toolno, EmcPose offset, double diameter,
                          double frontangle, double backangle,
                          int orientation) {}

void USE_TOOL_LENGTH_OFFSET(EmcPose offset) {
  tool_offset = offset;

  if (interp_error)
    return;

  callbacks.useToolLengthOffset(offset.tran.x, offset.tran.y, offset.tran.z,
                                offset.a, offset.b, offset.c, offset.u,
                                offset.v, offset.w);
}

void CHANGE_TOOL(int slot) {
  if (interp_error)
    return;

  callbacks.changeTool(slot);
}

// ------
// Unused
// ------

void SELECT_TOOL(int tool) {}

void CHANGE_TOOL_NUMBER(int tool) {}

void RELOAD_TOOLDATA(void) {}
void TURN_PROBE_ON() {}
void TURN_PROBE_OFF() {}
int GET_EXTERNAL_TC_FAULT() { return 0; }
int GET_EXTERNAL_TC_REASON() { return 0; }

double GET_EXTERNAL_MOTION_CONTROL_TOLERANCE() { return 0.1; }
double GET_EXTERNAL_MOTION_CONTROL_NAIVECAM_TOLERANCE() { return 0.1; }
double GET_EXTERNAL_PROBE_POSITION_X() { return _pos_x; }
double GET_EXTERNAL_PROBE_POSITION_Y() { return _pos_y; }
double GET_EXTERNAL_PROBE_POSITION_Z() { return _pos_z; }
double GET_EXTERNAL_PROBE_POSITION_A() { return _pos_a; }
double GET_EXTERNAL_PROBE_POSITION_B() { return _pos_b; }
double GET_EXTERNAL_PROBE_POSITION_C() { return _pos_c; }
double GET_EXTERNAL_PROBE_POSITION_U() { return _pos_u; }
double GET_EXTERNAL_PROBE_POSITION_V() { return _pos_v; }
double GET_EXTERNAL_PROBE_POSITION_W() { return _pos_w; }
double GET_EXTERNAL_PROBE_VALUE() { return 0.0; }
int GET_EXTERNAL_PROBE_TRIPPED_VALUE() { return 0; }
double GET_EXTERNAL_POSITION_X() { return _pos_x; }
double GET_EXTERNAL_POSITION_Y() { return _pos_y; }
double GET_EXTERNAL_POSITION_Z() { return _pos_z; }
double GET_EXTERNAL_POSITION_A() { return _pos_a; }
double GET_EXTERNAL_POSITION_B() { return _pos_b; }
double GET_EXTERNAL_POSITION_C() { return _pos_c; }
double GET_EXTERNAL_POSITION_U() { return _pos_u; }
double GET_EXTERNAL_POSITION_V() { return _pos_v; }
double GET_EXTERNAL_POSITION_W() { return _pos_w; }
void INIT_CANON() {}

int GET_EXTERNAL_QUEUE_EMPTY() { return true; }
CANON_DIRECTION GET_EXTERNAL_SPINDLE(int) { return CANON_STOPPED; }
int GET_EXTERNAL_TOOL_SLOT() { return 0; }
int GET_EXTERNAL_SELECTED_TOOL_SLOT() { return 0; }
double GET_EXTERNAL_FEED_RATE() { return 1; }
double GET_EXTERNAL_TRAVERSE_RATE() { return 0; }
int GET_EXTERNAL_FLOOD() { return 0; }
int GET_EXTERNAL_MIST() { return 0; }
CANON_PLANE GET_EXTERNAL_PLANE() { return CANON_PLANE_XY; }
double GET_EXTERNAL_SPEED(int spindle) { return 0; }
void DISABLE_ADAPTIVE_FEED() {}
void ENABLE_ADAPTIVE_FEED() {}

int GET_EXTERNAL_FEED_OVERRIDE_ENABLE() { return 1; }
int GET_EXTERNAL_SPINDLE_OVERRIDE_ENABLE(int spindle) { return 1; }
int GET_EXTERNAL_ADAPTIVE_FEED_ENABLE() { return 0; }
int GET_EXTERNAL_FEED_HOLD_ENABLE() { return 1; }

double GET_EXTERNAL_TOOL_LENGTH_XOFFSET() { return tool_offset.tran.x; }
double GET_EXTERNAL_TOOL_LENGTH_YOFFSET() { return tool_offset.tran.y; }
double GET_EXTERNAL_TOOL_LENGTH_ZOFFSET() { return tool_offset.tran.z; }
double GET_EXTERNAL_TOOL_LENGTH_AOFFSET() { return tool_offset.a; }
double GET_EXTERNAL_TOOL_LENGTH_BOFFSET() { return tool_offset.b; }
double GET_EXTERNAL_TOOL_LENGTH_COFFSET() { return tool_offset.c; }
double GET_EXTERNAL_TOOL_LENGTH_UOFFSET() { return tool_offset.u; }
double GET_EXTERNAL_TOOL_LENGTH_VOFFSET() { return tool_offset.v; }
double GET_EXTERNAL_TOOL_LENGTH_WOFFSET() { return tool_offset.w; }

void START_CHANGE() {}
void MIST_OFF() {}
void FLOOD_OFF() {}
void MIST_ON() {}
void FLOOD_ON() {}
void PROGRAM_STOP() {}
void PROGRAM_END() {}
void FINISH() {}

bool GET_OPTIONAL_PROGRAM_STOP() { return false; }
void SET_OPTIONAL_PROGRAM_STOP(bool state) {}
void OPTIONAL_PROGRAM_STOP() {}
void PALLET_SHUTTLE() {}
void ENABLE_FEED_HOLD() {}

void SET_BLOCK_DELETE(bool enabled) {}
bool GET_BLOCK_DELETE(void) { return false; }

void IO_PLUGIN_CALL(int len, const char *call) {}
void CANON_ERROR(const char *fmt, ...){};
void DISABLE_FEED_OVERRIDE() {}
void ENABLE_FEED_OVERRIDE() {}
void ENABLE_SPEED_OVERRIDE(int spindle) {}

void UPDATE_TAG(StateTag tag) {}
void CLEAR_MOTION_OUTPUT_BIT(int bit) {}

void CLAMP_AXIS(CANON_AXIS axis) {}
int LOCK_ROTARY(int line_no, int joint_num) { return 0; }
int WAIT(int index, int input_type, int wait_type, double timeout) { return 0; }
void DISABLE_SPEED_OVERRIDE(int spindle) {}
double GET_EXTERNAL_ANALOG_INPUT(int index, double def) { return def; }
int GET_EXTERNAL_OFFSET_APPLIED() { return 0; }

char _parameter_file_name[PARAMETER_FILE_NAME_LENGTH];
extern void SET_PARAMETER_FILE_NAME(const char *name) {
  strncpy(_parameter_file_name, name, PARAMETER_FILE_NAME_LENGTH);
}

void GET_EXTERNAL_PARAMETER_FILE_NAME(
    char *file_name, /* string: to copy file name into       */
    int max_size)    /* maximum number of characters to copy */
{
  // Paranoid checks
  if (0 == file_name)
    return;

  if (max_size < 0)
    return;

  if (strlen(_parameter_file_name) < (size_t)max_size)
    strcpy(file_name, _parameter_file_name);
  else
    file_name[0] = 0;
}

EmcPose GET_EXTERNAL_OFFSETS() {
  EmcPose e;
  e.tran.x = 0;
  e.tran.y = 0;
  e.tran.z = 0;
  e.a = 0;
  e.b = 0;
  e.c = 0;
  e.u = 0;
  e.v = 0;
  e.w = 0;
  return e;
};

void PLUGIN_CALL(int len, const char *call) {}
double GET_EXTERNAL_LENGTH_UNITS() { return 0.03937007874016; }
int GET_EXTERNAL_AXIS_MASK() { return 0x3f; } // XYZABC machine
double GET_EXTERNAL_ANGLE_UNITS() { return 1.0; }
void ON_RESET(void) {}
void CLEAR_AUX_OUTPUT_BIT(int index) {}
int UNLOCK_ROTARY(int line_no, int joint_num) { return 0; }
void SET_MOTION_OUTPUT_VALUE(int index, double value) {}
int GET_EXTERNAL_DIGITAL_INPUT(int index, int def) { return def; }

CANON_UNITS GET_EXTERNAL_LENGTH_UNIT_TYPE() { return _length_unit_type; }
void SET_AUX_OUTPUT_BIT(int index) {}
void DISABLE_FEED_HOLD() {}

USER_DEFINED_FUNCTION_TYPE USER_DEFINED_FUNCTION[USER_DEFINED_FUNCTION_NUM] = {
    0};

int USER_DEFINED_FUNCTION_ADD(USER_DEFINED_FUNCTION_TYPE func, int num) {
  if (num < 0 || num >= USER_DEFINED_FUNCTION_NUM) {
    return -1;
  }

  USER_DEFINED_FUNCTION[num] = func;

  return 0;
}

void SET_AUX_OUTPUT_VALUE(int index, double value) {}

extern CANON_TOOL_TABLE GET_EXTERNAL_TOOL_TABLE(int idx) {
#ifdef TOOL_NML //{
  return _sai._tools[idx];
#else  //}{
  CANON_TOOL_TABLE tdata;
  return tdata;
#endif //}
}
void SET_MOTION_OUTPUT_BIT(int index) {}
