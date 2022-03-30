typedef void (*StartCallbackFn)(int i);
typedef void (*EndCallbackFn)(int a, int b);

typedef struct {
  StartCallbackFn start;
  EndCallbackFn end;
} Callbacks;

// Processes the file and invokes callbacks from cbs on events found in the
// file, each with its own relevant data. user_data is passed through to the
// callbacks.
void traverse(char *filename, Callbacks cbs);
