// Stub implementation for testing C interop.
// this file is removed while building in favour
// of linking against rs274.so from /adapter.
#include "rs274.hh"

//----------------------
// Callback Registration
//----------------------

static Callbacks callbacks;

// register callback functions for the parser to call.
void registerCallbacks(Callbacks cb) { callbacks = cb; }

int parseFile(char *file) {
  callbacks.error(1, 2);
  return 0;
}
