#include <stdbool.h>
#include <stdlib.h>

// ----
// Util
// ----

extern void errorGo(int, int);
void ErrorAdpt(int interpError, int lastSequenceNumber) {
  errorGo(interpError, lastSequenceNumber);
}

extern bool abortGo();
bool AbortAdpt() { return abortGo(); }

extern void messageGo(char *comment);
void MessageAdpt(char *comment) { messageGo(comment); }

extern void commentGo(char *comment);
void CommentAdpt(char *comment) { commentGo(comment); }

extern void changeToolGo(int pocket);
void ChangeToolAdpt(int pocket) { changeToolGo(pocket); }

// --------
// Settings
// --------

extern void useLengthUnitsGo(int units);
void UseLengthUnitsAdpt(int units) { useLengthUnitsGo(units); }

extern void useToolLengthOffsetGo(double x, double y, double z, double a,
                                  double b, double c, double u, double v,
                                  double w);
void UseToolLengthOffsetAdpt(double x, double y, double z, double a, double b,
                             double c, double u, double v, double w) {
  useToolLengthOffsetGo(x, y, z, a, b, c, u, v, w);
}

extern void selectPlaneGo(int plane);
void SelectPlaneAdpt(int plane) { selectPlaneGo(plane); }

extern void setXYRotationGo(double);
void SetXYRotationAdpt(double rotation) { setXYRotationGo(rotation); }

extern void setG5XOffsetGo(int index, double x, double y, double z, double a,
                           double b, double c, double u, double v, double w);
void SetG5XOffsetAdpt(int index, double x, double y, double z, double a,
                      double b, double c, double u, double v, double w) {
  setG5XOffsetGo(index, x, y, z, a, b, c, u, v, w);
}

extern void setG92OffsetGo(double x, double y, double z, double a, double b,
                           double c, double u, double v, double w);
void SetG92OffsetAdpt(double x, double y, double z, double a, double b,
                      double c, double u, double v, double w) {
  setG92OffsetGo(x, y, z, a, b, c, u, v, w);
}

extern void setTraverseRateGo(double rate);
void SetTraverseRateAdpt(double rate) { setTraverseRateGo(rate); }

extern void setFeedModeGo(int spindle, int mode);
void SetFeedModeAdpt(int spindle, int mode) { setFeedModeGo(spindle, mode); }

extern void setFeedRateGo(double rate);
void SetFeedRateAdpt(double rate) { setFeedRateGo(rate); }

// --------
// Movement
// --------
extern void straightTraverseGo(int lineNo, double x, double y, double z,
                               double a, double b, double c, double u, double v,
                               double w);
void StraightTraverseAdpt(int lineNo, double x, double y, double z, double a,
                          double b, double c, double u, double v, double w) {
  straightTraverseGo(lineNo, x, y, z, a, b, c, u, v, w);
}

extern void arcFeedGo(double firstEnd, double secondEnd, double firstAxis,
                      double secondAxis, int rotation, double axisEndPoint,
                      double aPosition, double bPosition, double cPosition,
                      double uPosition, double vPosition, double wPosition);
void ArcFeedAdpt(double firstEnd, double secondEnd, double firstAxis,
                 double secondAxis, int rotation, double axisEndPoint,
                 double aPosition, double bPosition, double cPosition,
                 double uPosition, double vPosition, double wPosition) {
  arcFeedGo(firstEnd, secondEnd, firstAxis, secondAxis, rotation, axisEndPoint,
            aPosition, bPosition, cPosition, uPosition, vPosition, wPosition);
}

extern void straightFeedGo(int lineNo, double x, double y, double z, double a,
                           double b, double c, double u, double v, double w);
void StraightFeedAdpt(int lineNo, double x, double y, double z, double a,
                      double b, double c, double u, double v, double w) {
  straightFeedGo(lineNo, x, y, z, a, b, c, u, v, w);
}

extern void straightProbeGo(int lineNo, double x, double y, double z, double a,
                            double b, double c, double u, double v, double w);
void StraightProbeAdpt(int lineNo, double x, double y, double z, double a,
                       double b, double c, double u, double v, double w) {
  straightProbeGo(lineNo, x, y, z, a, b, c, u, v, w);
}

extern void rigidTapGo(int lineNo, double x, double y, double z, double scale);
void RigidTapAdpt(int lineNo, double x, double y, double z, double scale) {
  rigidTapGo(lineNo, x, y, z, scale);
}

extern void dwellGo(double seconds);
void DwellAdpt(double seconds) { dwellGo(seconds); }
