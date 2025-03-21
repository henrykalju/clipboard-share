#ifndef CLIPBOARDWRITER_H
#define CLIPBOARDWRITER_H

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

void set_listener(Window);
void set_clipboard_item(Item);
Window init_clipboard();
void start_clipboard();

#endif