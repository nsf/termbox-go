#ifndef TERMBOX_TERM_H
#define TERMBOX_TERM_H

#include "termbox.h"
#include "ringbuffer.h"

enum {
	T_ENTER_CA,
	T_EXIT_CA,
	T_SHOW_CURSOR,
	T_HIDE_CURSOR,
	T_CLEAR_SCREEN,
	T_SGR,
	T_SGR0,
	T_UNDERLINE,
	T_BOLD,
	T_BLINK,
	T_MOVE_CURSOR,
	T_ENTER_KEYPAD,
	T_EXIT_KEYPAD
};

#define EUNSUPPORTED_TERM -1

int init_term();

/* 0 on success, -1 on failure */
int extract_event(struct tb_event *event, struct ringbuffer *inbuf, int inputmode);

extern const char **keys;
extern const char **funcs;

#endif
