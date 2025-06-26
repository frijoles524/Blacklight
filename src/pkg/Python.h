// Custom Blacklight file. Defines essential functions for python
void Py_Initialize(void);
void Py_Finalize(void);
int PyRun_SimpleString(const char *);
int PyRun_SimpleFile(FILE *, const char *);
void PyErr_Print(void);
