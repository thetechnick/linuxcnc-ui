#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Util functions
typedef void (*ErrorFn)(int interpError, int lastSequenceNumber);
typedef bool (*AbortFn)();
typedef void (*MessageFn)(char *comment);
typedef void (*CommentFn)(const char *comment);
typedef void (*ChangeToolFn)(int pocket);

// Settings functions
typedef void (*UseLengthUnitsFn)(int units); // CANON_UNITS enum
typedef void (*UseToolLengthOffsetFn)(double x, double y, double z, double a,
                                      double b, double c, double u, double v,
                                      double w);
typedef void (*SelectPlaneFn)(int plane);
typedef void (*SetXYRotationFn)(double);
typedef void (*SetG5XOffsetFn)(int index, double x, double y, double z,
                               double a, double b, double c, double u, double v,
                               double w);
typedef void (*SetG92OffsetFn)(double x, double y, double z, double a, double b,
                               double c, double u, double v, double w);
typedef void (*SetTraverseRateFn)(double rate);
typedef void (*SetFeedModeFn)(int spindle, int mode);
typedef void (*SetFeedRateFn)(double rate);

// Movement functions
typedef void (*StraightTraverseFn)(int lineNo, double x, double y, double z,
                                   double a, double b, double c, double u,
                                   double v, double w);
typedef void (*ArcFeedFn)(double firstEnd, double secondEnd, double firstAxis,
                          double secondAxis, int rotation, double axisEndPoint,
                          double aPosition, double bPosition, double cPosition,
                          double uPosition, double vPosition, double wPosition);
typedef void (*StraightFeedFn)(int lineNo, double x, double y, double z,
                               double a, double b, double c, double u, double v,
                               double w);
typedef void (*StraightProbeFn)(int lineNo, double x, double y, double z,
                                double a, double b, double c, double u,
                                double v, double w);
typedef void (*RigidTapFn)(int lineNo, double x, double y, double z,
                           double scale);
typedef void (*DwellFn)(double seconds);

typedef struct {
  // Util
  ErrorFn error;
  AbortFn abort;
  MessageFn message;
  CommentFn comment;
  ChangeToolFn changeTool;

  // Settings
  UseLengthUnitsFn useLengthUnits;
  UseToolLengthOffsetFn useToolLengthOffset;
  SelectPlaneFn selectPlane;
  SetXYRotationFn setXYRotation;
  SetG5XOffsetFn setG5XOffset;
  SetG92OffsetFn setG92Offset;
  SetTraverseRateFn setTraverseRate;
  SetFeedModeFn setFeedMode;
  SetFeedRateFn setFeedRate;

  // Movement
  StraightTraverseFn straightTraverse;
  ArcFeedFn arcFeed;
  StraightFeedFn straightFeed;
  StraightProbeFn straightProbe;
  RigidTapFn rigidTap;
  DwellFn dwell;
} Callbacks;

void registerCallbacks(Callbacks cb);

int parseFile(char *file);

#ifdef __cplusplus
}
#endif
