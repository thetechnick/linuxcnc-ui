package main

/*
#include <stdlib.h>

extern void goStart(int);
extern void goEnd(int, int);

void startCgo(int i) {
  goStart(i);
}

void endCgo(int a, int b) {
  goEnd(a, b);
}
*/
import "C"
