#include "clipboardlistener.h"
#include <X11/Xlib.h>
#include <X11/extensions/Xfixes.h>
#include <X11/Xatom.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

Display* listenerdisplay;
Window root, listenerwindow, writerwindow2;

void set_writer(Window w) {
    writerwindow2 = w;
}

Item create_item() {
    return (Item){.len = 0, .values = NULL};
}

void free_item(Item* item) {
    if (item == NULL) {
        return;
    }
    for (int i = 0; i < item->len; i++) {
        if (item->values[i].format != NULL) {
            free(item->values[i].format);
        }
        if (item->values[i].data != NULL) {
            free(item->values[i].data);
        }
    }
    if (item->values != NULL) {
        free(item->values);
    }
}

void add_value(Item* item, const char* format, const uint8_t* data, int data_len) {
    item->values = (Value*)realloc(item->values, (item->len + 1) * sizeof(Value));
    Value* new_value = &item->values[item->len];
    new_value->format_len = strlen(format);
    new_value->format = (char*)malloc(new_value->format_len+1);
    strcpy(new_value->format, format);

    new_value->data_len = data_len;
    new_value->data = (uint8_t*)malloc(new_value->data_len);
    memcpy(new_value->data, data, data_len);

    item->len++;
}

int skip_target(char* name) {
    if (
        strcmp(name, "TARGETS") == 0 ||
        strcmp(name, "TIMESTAMP") == 0 ||
        strcmp(name, "SAVE_TARGETS") == 0 ||
        strcmp(name, "MULTIPLE") == 0
    ) {
        return 1;
    }
    return 0;
}

unsigned char *get_target_data(Display *dpy, Window w, Atom sel, Atom target, Atom property, unsigned long *size) {
    Atom type;
    int format;
    unsigned long nitems, bytes_after;
    unsigned char *prop_ret = NULL;

    // Request the content of the target
    XConvertSelection(dpy, sel, target, property, w, CurrentTime);

    // Wait for the SelectionNotify event
    for (;;) {
        XEvent ev;
        XNextEvent(dpy, &ev);
        if (ev.type == SelectionNotify) {
            XSelectionEvent *sev = (XSelectionEvent *)&ev.xselection;
            if (sev->property == None) {
                printf("Failed to retrieve data for target '%s'.\n", XGetAtomName(dpy, target));
                return NULL;
            }
            break;
        }
    }

    // Retrieve the property content
    XGetWindowProperty(dpy, w, property, 0, (~0L), False, AnyPropertyType,
                       &type, &format, &nitems, &bytes_after, &prop_ret);

    if (size) *size = nitems;

    XDeleteProperty(dpy, w, property);
    return prop_ret;
}

void process_targets(Display *dpy, Window w, Atom sel, Atom targets_property, Atom target_property) {
    //show_targets(dpy, w, sel, targets_property);

    Atom *targets;
    unsigned long nitems;
    int di;
    unsigned long dul;
    Atom type;
    unsigned char *prop_ret = NULL;

    // Request the TARGETS list
    XConvertSelection(dpy, sel, XInternAtom(dpy, "TARGETS", False), targets_property, w, CurrentTime);
    for (;;) {
        XEvent ev;
        XNextEvent(dpy, &ev);
        if (ev.type == SelectionNotify) break;
    }

    // Retrieve the TARGETS property
    XGetWindowProperty(dpy, w, targets_property, 0, 1024 * sizeof(Atom), False, XA_ATOM,
                       &type, &di, &nitems, &dul, &prop_ret);

    targets = (Atom *)prop_ret;
    Item item = create_item();
    for (unsigned long i = 0; i < nitems; i++) {
        char *name = XGetAtomName(dpy, targets[i]);

        if (!name) continue;

        // Skip unwanted TARGETS
        if (skip_target(name)) {
            XFree(name);
            continue;
        }

        unsigned long data_size;
        unsigned char *data = get_target_data(dpy, w, sel, targets[i], target_property, &data_size);
        if (data) {
            // Do something with the raw data (e.g., store or process)
            add_value(&item, XGetAtomName(dpy, targets[i]), data, data_size);
            XFree(data);
        }

        XFree(name);
    }
    handleChange(item);
    free_item(&item);

    XFree(prop_ret);
}

Window Init() {
    // Open the X display
    listenerdisplay = XOpenDisplay(NULL);
    if (!listenerdisplay) {
        fprintf(stderr, "Unable to open X display\n");
        return 0;
    }

    // Define atoms for CLIPBOARD and UTF8_STRING
    root = DefaultRootWindow(listenerdisplay);

    listenerwindow = XCreateSimpleWindow(listenerdisplay, root, -10, -10, 1, 1, 0, 0, 0);
    return listenerwindow;
}

void StartListening() {
    printf("starting listening\n");
    

    // Check if XFixes extension is available
    int event_base, error_base;
    if (!XFixesQueryExtension(listenerdisplay, &event_base, &error_base)) {
        fprintf(stderr, "XFixes extension not available\n");
        XCloseDisplay(listenerdisplay);
        return;
    }

    //Atom XA_UTF8_STRING = XInternAtom(display, "TARGETS", False);

    Atom clipboard = XInternAtom(listenerdisplay, "CLIPBOARD", False);
    Atom targets = XInternAtom(listenerdisplay, "TARGETS_PROPERTY", False);
    Atom target = XInternAtom(listenerdisplay, "TARGET_DATA", False);


    // Select the necessary XFixes input event to listen for selection ownership changes
    XFixesSelectSelectionInput(listenerdisplay, root, clipboard, XFixesSetSelectionOwnerNotifyMask);

    // Event loop to listen for changes
    while (1) {
        XEvent event;
        XNextEvent(listenerdisplay, &event);

        // Check for the selection ownership notification
        if (event.type == event_base + XFixesSelectionNotify) {
            XFixesSelectionNotifyEvent *notify_event = (XFixesSelectionNotifyEvent *)&event;
            if (notify_event->owner == writerwindow2) {
                printf("got writer request\n");
                continue;
            }
            // Check if the clipboard has been updated
            if (notify_event->selection == clipboard) {
                //printf("Clipboard content changed!\n");

                process_targets(listenerdisplay, listenerwindow, clipboard, targets, target);

                // Optionally, you can fetch the clipboard data here
                // For simplicity, we assume that the content was updated and we can print a message
            }
        }
    }

    // Close the display connection
    XCloseDisplay(listenerdisplay);
    return;
}
