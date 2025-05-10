#include "clipboardwriter.h"
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <time.h>
#include <stdlib.h>
#include <stdint.h>

Item clipboard_item;
Display *writerdisplay;
Window writerwindow, listenerwindow2;
Atom clipboard, targets_atom;

void set_listener(Window l) {
    listenerwindow2 = l;
}

int take_clipboard_ownership() {
    int x = XSetSelectionOwner(writerdisplay, clipboard, writerwindow, CurrentTime);
    XFlush(writerdisplay);
    return x;
}

void free_clipboard() {
    for (int i = 0; i < clipboard_item.len; i++) {
        free(clipboard_item.values[i].format);
        free(clipboard_item.values[i].data);
    }
    free(clipboard_item.values);
}

void set_clipboard_item(Item new_item) {
    free_clipboard();
    
    clipboard_item = new_item;
    take_clipboard_ownership();
}

Window init_clipboard() {
    writerdisplay = XOpenDisplay(NULL);
    if (!writerdisplay) {
        printf("Unable to open X display\n");
        return 0;
    }

    writerwindow = XCreateSimpleWindow(writerdisplay, RootWindow(writerdisplay, 0), 0, 0, 1, 1, 0, 0, 0);
    clipboard = XInternAtom(writerdisplay, "CLIPBOARD", False);
    targets_atom = XInternAtom(writerdisplay, "TARGETS", False);

    char *data = "TEST\n";
    Item c;
    c.len = 1;
    c.values = malloc(sizeof(Value));
    c.values[0].format = strdup("STRING");
    c.values[0].format_len = strlen("STRING");
    c.values[0].data = (uint8_t*)strdup(data);
    c.values[0].data_len = strlen(data);
    set_clipboard_item(c);

    return writerwindow;
}

void send_no(Display *dpy, XSelectionRequestEvent *sev) {
    XSelectionEvent ssev;
    char *an = XGetAtomName(dpy, sev->target);
    
    if (an)
        XFree(an);

    ssev.type = SelectionNotify;
    ssev.requestor = sev->requestor;
    ssev.selection = sev->selection;
    ssev.target = sev->target;
    ssev.property = None;
    ssev.time = sev->time;

    XSendEvent(dpy, sev->requestor, True, NoEventMask, (XEvent *)&ssev);
}

void send_targets(Display *dpy, XSelectionRequestEvent *sev) {
    Atom *supported_targets = malloc((clipboard_item.len + 1) * sizeof(Atom));
    supported_targets[0] = XInternAtom(dpy, "TARGETS", False);
    
    for (int i = 0; i < clipboard_item.len; i++) {
        supported_targets[i + 1] = XInternAtom(dpy, clipboard_item.values[i].format, False);
    }
    
    XChangeProperty(dpy, sev->requestor, sev->property, XA_ATOM, 32, PropModeReplace,
                    (unsigned char *)supported_targets, clipboard_item.len + 1);
    
    free(supported_targets);
    
    XSelectionEvent ssev = {
        .type = SelectionNotify,
        .requestor = sev->requestor,
        .selection = sev->selection,
        .target = sev->target,
        .property = sev->property,
        .time = sev->time
    };
    
    XSendEvent(dpy, sev->requestor, True, NoEventMask, (XEvent *)&ssev);
}

void send_format(Display *dpy, XSelectionRequestEvent *sev) {
    char* format_name = XGetAtomName(dpy, sev->target);
    for (int i = 0; i < clipboard_item.len; i++) {
        if (strcmp(clipboard_item.values[i].format, format_name) == 0) {
            XChangeProperty(dpy, sev->requestor, sev->property, sev->target, 8, PropModeReplace,
                            clipboard_item.values[i].data, clipboard_item.values[i].data_len);

            XSelectionEvent ssev = {
                .type = SelectionNotify,
                .requestor = sev->requestor,
                .selection = sev->selection,
                .target = sev->target,
                .property = sev->property,
                .time = sev->time
            };

            XSendEvent(dpy, sev->requestor, True, NoEventMask, (XEvent *)&ssev);
            return;
        }
    }
    send_no(writerdisplay, sev);
}

void start_clipboard() {
    XEvent ev;
    XSelectionRequestEvent* sev;
    for (;;) {
        XNextEvent(writerdisplay, &ev);
        if (ev.type == SelectionRequest) {
            sev = (XSelectionRequestEvent*)&ev.xselectionrequest;
            if (sev->requestor == listenerwindow2) {
                printf("got listener request\n");
                continue;
            }
            
            if (sev->target == targets_atom)
                send_targets(writerdisplay, sev);
            else
                send_format(writerdisplay, sev);
        }
    }
    free_clipboard();
}