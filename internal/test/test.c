#include "test.h"
#include "stdlib.h"

void traverse(char *filename, Callbacks cbs) {
  // Simulate some traversal that calls the start callback and then the end
  // callback, if they are defined.
  if (cbs.start != NULL) {
    cbs.start(100);
  }
  if (cbs.end != NULL) {
    cbs.end(2, 3);
  }
}
