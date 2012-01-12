#ifndef TERMBOX_H
#define TERMBOX_H

#include <stdint.h>

/* for shared objects */
#if __GNUC__ >= 4
 #define SO_IMPORT __attribute__((visibility("default")))
#else
 #define SO_IMPORT
#endif

#ifdef __cplusplus
extern "C" {
#endif

/* --------------- keys ---------------- */

/* These are safe subset of terminfo keys, which exists on all popular terminals.
   I think it's important to use only these and not others.
*/
#define TB_KEY_F1		(0xFFFF-0)
#define TB_KEY_F2		(0xFFFF-1)
#define TB_KEY_F3		(0xFFFF-2)
#define TB_KEY_F4		(0xFFFF-3)
#define TB_KEY_F5		(0xFFFF-4)
#define TB_KEY_F6		(0xFFFF-5)
#define TB_KEY_F7		(0xFFFF-6)
#define TB_KEY_F8		(0xFFFF-7)
#define TB_KEY_F9		(0xFFFF-8)
#define TB_KEY_F10		(0xFFFF-9)
#define TB_KEY_F11		(0xFFFF-10)
#define TB_KEY_F12		(0xFFFF-11)
#define TB_KEY_INSERT		(0xFFFF-12)
#define TB_KEY_DELETE		(0xFFFF-13)
#define TB_KEY_HOME		(0xFFFF-14)
#define TB_KEY_END		(0xFFFF-15)
#define TB_KEY_PGUP		(0xFFFF-16)
#define TB_KEY_PGDN		(0xFFFF-17)
#define TB_KEY_ARROW_UP		(0xFFFF-18)
#define TB_KEY_ARROW_DOWN	(0xFFFF-19)
#define TB_KEY_ARROW_LEFT	(0xFFFF-20)
#define TB_KEY_ARROW_RIGHT	(0xFFFF-21)

/* These are all keys below SPACE character and BACKSPACE key */
#define TB_KEY_CTRL_TILDE	0x00
#define TB_KEY_CTRL_2		0x00 /* clash with 'CTRL_TILDE' */
#define TB_KEY_CTRL_A		0x01
#define TB_KEY_CTRL_B		0x02
#define TB_KEY_CTRL_C		0x03
#define TB_KEY_CTRL_D		0x04
#define TB_KEY_CTRL_E		0x05
#define TB_KEY_CTRL_F		0x06
#define TB_KEY_CTRL_G		0x07
#define TB_KEY_BACKSPACE	0x08
#define TB_KEY_CTRL_H		0x08 /* clash with 'CTRL_BACKSPACE' */
#define TB_KEY_TAB		0x09
#define TB_KEY_CTRL_I		0x09 /* clash with 'TAB' */
#define TB_KEY_CTRL_J		0x0A
#define TB_KEY_CTRL_K		0x0B
#define TB_KEY_CTRL_L		0x0C
#define TB_KEY_ENTER		0x0D
#define TB_KEY_CTRL_M		0x0D /* clash with 'ENTER' */
#define TB_KEY_CTRL_N		0x0E
#define TB_KEY_CTRL_O		0x0F
#define TB_KEY_CTRL_P		0x10
#define TB_KEY_CTRL_Q		0x11
#define TB_KEY_CTRL_R		0x12
#define TB_KEY_CTRL_S		0x13
#define TB_KEY_CTRL_T		0x14
#define TB_KEY_CTRL_U		0x15
#define TB_KEY_CTRL_V		0x16
#define TB_KEY_CTRL_W		0x17
#define TB_KEY_CTRL_X		0x18
#define TB_KEY_CTRL_Y		0x19
#define TB_KEY_CTRL_Z		0x1A
#define TB_KEY_ESC		0x1B
#define TB_KEY_CTRL_LSQ_BRACKET	0x1B /* clash with 'ESC' */
#define TB_KEY_CTRL_3		0x1B /* clash with 'ESC' */
#define TB_KEY_CTRL_4		0x1C
#define TB_KEY_CTRL_BACKSLASH	0x1C /* clash with 'CTRL_4' */
#define TB_KEY_CTRL_5		0x1D
#define TB_KEY_CTRL_RSQ_BRACKET 0x1D /* clash with 'CTRL_5' */
#define TB_KEY_CTRL_6		0x1E
#define TB_KEY_CTRL_7		0x1F
#define TB_KEY_CTRL_SLASH	0x1F /* clash with 'CTRL_7' */
#define TB_KEY_CTRL_UNDERSCORE	0x1F /* clash with 'CTRL_7' */
#define TB_KEY_SPACE		0x20
#define TB_KEY_BACKSPACE2	0x7F
#define TB_KEY_CTRL_8		0x7F /* clash with 'DELETE' */

/* These are fail ones (do not exist) */
/* #define TB_KEY_CTRL_1 clash with '1' */
/* #define TB_KEY_CTRL_9 clash with '9' */
/* #define TB_KEY_CTRL_0 clash with '0' */

/* Currently there is only one modificator */
#define TB_MOD_ALT		0x01

/* colors */
#define TB_BLACK	0x00
#define TB_RED		0x01
#define TB_GREEN	0x02
#define TB_YELLOW	0x03
#define TB_BLUE		0x04
#define TB_MAGENTA	0x05
#define TB_CYAN		0x06
#define TB_WHITE	0x07

/* attributes */
#define TB_BOLD		0x10
#define TB_UNDERLINE	0x20

struct tb_cell {
	uint32_t ch;
	uint16_t fg;
	uint16_t bg;
};

struct tb_event {
	uint8_t type;
	uint8_t mod;
	uint16_t key;
	uint32_t ch;
	int32_t w;
	int32_t h;
};

#define TB_EUNSUPPORTED_TERMINAL	-1
#define TB_EFAILED_TO_OPEN_TTY		-2
#define TB_EPIPE_TRAP_ERROR		-3

SO_IMPORT int tb_init();
SO_IMPORT void tb_shutdown();

SO_IMPORT unsigned int tb_width();
SO_IMPORT unsigned int tb_height();

SO_IMPORT void tb_clear();
SO_IMPORT void tb_present();

#define TB_HIDE_CURSOR -1
SO_IMPORT void tb_set_cursor(int cx, int cy);

SO_IMPORT void tb_put_cell(unsigned int x, unsigned int y, const struct tb_cell *cell);
SO_IMPORT void tb_change_cell(unsigned int x, unsigned int y, uint32_t ch, uint16_t fg, uint16_t bg);
SO_IMPORT void tb_blit(unsigned int x, unsigned int y, unsigned int w, unsigned int h, const struct tb_cell *cells);

#define TB_INPUT_ESC 1
#define TB_INPUT_ALT 2
/* with 0 returns current input mode */
SO_IMPORT int tb_select_input_mode(int mode);

#define TB_EVENT_KEY 1
#define TB_EVENT_RESIZE 2

SO_IMPORT int tb_peek_event(struct tb_event *event, unsigned int timeout);
SO_IMPORT int tb_poll_event(struct tb_event *event);
/* these return:
	0 - no events, no errors,
	1 - key event
	2 - resize event
	-1 - error (input buffer overflow, discarded input)

   timeout in milliseconds
*/


/* glib based interface */
/* TODO */

/* utility utf8 functions */
#define TB_EOF -1
SO_IMPORT int utf8_char_length(char c);
SO_IMPORT int utf8_char_to_unicode(uint32_t *out, const char *c);
SO_IMPORT int utf8_unicode_to_char(char *out, uint32_t c);

#ifdef __cplusplus
}
#endif

#endif /* ifndef TERMBOX_H */
