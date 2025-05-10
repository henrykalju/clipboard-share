#ifndef CLIPBOARDLISTENER_H
#define CLIPBOARDLISTENER_H

#include <stdint.h>

#include <X11/Xlib.h>

typedef struct {
    char* format;
    int format_len;
    uint8_t* data;
    int data_len;
} Value;

typedef struct {
    Value* values;
    int len;
} Item;

extern void handleChange(Item);

Window Init();
void set_writer(Window);
void StartListening();

#endif