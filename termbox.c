#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <fcntl.h>
#include <signal.h>
#include <sys/ioctl.h>
#include <sys/time.h>
#include <termios.h>
#include <unistd.h>

#include "term.h"
#include "termbox.h"

struct cellbuf {
	unsigned int width;
	unsigned int height;
	struct tb_cell *cells;
};

#define CELL(buf, x, y) (buf)->cells[(y) * (buf)->width + (x)]
#define IS_CURSOR_HIDDEN(cx, cy) (cx == -1 || cy == -1)

#define LAST_COORD_INIT 0xFFFFFFFE

static struct termios orig_tios;

static struct cellbuf back_buffer;
static struct cellbuf front_buffer;
static unsigned int termw;
static unsigned int termh;

static int inputmode = TB_INPUT_ESC;

static struct ringbuffer inbuf;

static FILE *out;
static FILE *in;

static int out_fileno;
static int in_fileno;

static int winch_fds[2];

static unsigned int lastx = LAST_COORD_INIT;
static unsigned int lasty = LAST_COORD_INIT;
static int cursor_x = -1;
static int cursor_y = -1;

static uint16_t background = TB_BLACK;
static uint16_t foreground = TB_WHITE;

static void cellbuf_init(struct cellbuf *buf, unsigned int width, unsigned int height);
static void cellbuf_resize(struct cellbuf *buf, unsigned int width, unsigned int height);
static void cellbuf_clear(struct cellbuf *buf);
static void cellbuf_free(struct cellbuf *buf);

static void update_size();
static void update_term_size();
static void send_attr(uint16_t fg, uint16_t bg);
static void send_char(unsigned int x, unsigned int y, uint32_t c);
static void send_clear();
static void sigwinch_handler(int xxx);
static int wait_fill_event(struct tb_event *event, struct timeval *timeout);

/* may happen in a different thread */
static volatile int buffer_size_change_request;

/* -------------------------------------------------------- */

int tb_init()
{
	out = fopen("/dev/tty", "w");
	in = fopen("/dev/tty", "r");

	if (!out || !in)
		return TB_EFAILED_TO_OPEN_TTY;

	out_fileno = fileno(out);
	in_fileno = fileno(in);

	if (init_term() < 0)
		return TB_EUNSUPPORTED_TERMINAL;

	if (pipe(winch_fds) < 0)
		return TB_EPIPE_TRAP_ERROR;

	struct sigaction sa;
	sa.sa_handler = sigwinch_handler;
	sa.sa_flags = 0;
	sigaction(SIGWINCH, &sa, 0);

	tcgetattr(out_fileno, &orig_tios);

	struct termios tios;
	memcpy(&tios, &orig_tios, sizeof(tios));

	tios.c_iflag &= ~(IGNBRK | BRKINT | PARMRK | ISTRIP
                           | INLCR | IGNCR | ICRNL | IXON);
	tios.c_oflag &= ~OPOST;
	tios.c_lflag &= ~(ECHO | ECHONL | ICANON | ISIG | IEXTEN);
	tios.c_cflag &= ~(CSIZE | PARENB);
	tios.c_cflag |= CS8;
	tios.c_cc[VMIN] = 0;
	tios.c_cc[VTIME] = 0;
	tcsetattr(out_fileno, TCSAFLUSH, &tios);

	fputs(funcs[T_ENTER_CA], out);
	fputs(funcs[T_ENTER_KEYPAD], out);
	fputs(funcs[T_HIDE_CURSOR], out);
	send_clear();

	update_term_size();
	cellbuf_init(&back_buffer, termw, termh);
	cellbuf_init(&front_buffer, termw, termh);
	cellbuf_clear(&back_buffer);
	cellbuf_clear(&front_buffer);
	init_ringbuffer(&inbuf, 4096);

	return 0;
}

void tb_shutdown()
{
	fputs(funcs[T_SHOW_CURSOR], out);
	fputs(funcs[T_SGR0], out);
	fputs(funcs[T_CLEAR_SCREEN], out);
	fputs(funcs[T_EXIT_CA], out);
	fputs(funcs[T_EXIT_KEYPAD], out);
	fflush(out);
	tcsetattr(out_fileno, TCSAFLUSH, &orig_tios);

	fclose(out);
	fclose(in);
	close(winch_fds[0]);
	close(winch_fds[1]);

	cellbuf_free(&back_buffer);
	cellbuf_free(&front_buffer);
	free_ringbuffer(&inbuf);
}

void tb_present()
{
	unsigned int x,y;
	struct tb_cell *back, *front;

	/* invalidate cursor position */
	lastx = LAST_COORD_INIT;
	lasty = LAST_COORD_INIT;

	if (buffer_size_change_request) {
		update_size();
		buffer_size_change_request = 0;
	}

	for (y = 0; y < front_buffer.height; ++y) {
		for (x = 0; x < front_buffer.width; ++x) {
			back = &CELL(&back_buffer, x, y);
			front = &CELL(&front_buffer, x, y);
			if (memcmp(back, front, sizeof(struct tb_cell)) == 0)
				continue;
			send_attr(back->fg, back->bg);
			send_char(x, y, back->ch);
			memcpy(front, back, sizeof(struct tb_cell));
		}
	}
	if (!IS_CURSOR_HIDDEN(cursor_x, cursor_y))
		fprintf(out, funcs[T_MOVE_CURSOR], cursor_y+1, cursor_x+1);
	fflush(out);
}

void tb_set_cursor(int cx, int cy)
{
	if (IS_CURSOR_HIDDEN(cursor_x, cursor_y) && !IS_CURSOR_HIDDEN(cx, cy))
		fputs(funcs[T_SHOW_CURSOR], out);

	if (!IS_CURSOR_HIDDEN(cursor_x, cursor_y) && IS_CURSOR_HIDDEN(cx, cy))
		fputs(funcs[T_HIDE_CURSOR], out);

	cursor_x = cx;
	cursor_y = cy;
	if (!IS_CURSOR_HIDDEN(cursor_x, cursor_y))
		fprintf(out, funcs[T_MOVE_CURSOR], cursor_y+1, cursor_x+1);
}

void tb_put_cell(unsigned int x, unsigned int y, const struct tb_cell *cell)
{
	if (x >= back_buffer.width || y >= back_buffer.height)
		return;
	CELL(&back_buffer, x, y) = *cell;
}

void tb_change_cell(unsigned int x, unsigned int y, uint32_t ch, uint16_t fg, uint16_t bg)
{
	struct tb_cell c = {ch, fg, bg};
	tb_put_cell(x, y, &c);
}

void tb_blit(unsigned int x, unsigned int y, unsigned int w, unsigned int h, const struct tb_cell *cells)
{
	if (x+w > back_buffer.width || y+h > back_buffer.height)
		return;

	unsigned int sy;
	struct tb_cell *dst = &CELL(&back_buffer, x, y);
	size_t size = sizeof(struct tb_cell) * w;

	for (sy = 0; sy < h; ++sy) {
		memcpy(dst, cells, size);
		dst += back_buffer.width;
		cells += w;
	}
}

int tb_poll_event(struct tb_event *event)
{
	return wait_fill_event(event, 0);
}

int tb_peek_event(struct tb_event *event, unsigned int timeout)
{
	struct timeval tv;
	tv.tv_sec = timeout / 1000;
	tv.tv_usec = (timeout - (tv.tv_sec * 1000)) * 1000;
	return wait_fill_event(event, &tv);
}

unsigned int tb_width()
{
	return termw;
}

unsigned int tb_height()
{
	return termh;
}

void tb_clear()
{
	if (buffer_size_change_request) {
		update_size();
		buffer_size_change_request = 0;
	}
	cellbuf_clear(&back_buffer);
}

int tb_select_input_mode(int mode)
{
	if (mode)
		inputmode = mode;
	return inputmode;
}

void tb_set_clear_attributes(uint16_t fg, uint16_t bg)
{
	foreground = fg;
	background = bg;
}

/* -------------------------------------------------------- */

static void cellbuf_init(struct cellbuf *buf, unsigned int width, unsigned int height)
{
	buf->cells = (struct tb_cell*)malloc(sizeof(struct tb_cell) * width * height);
	assert(buf->cells);
	buf->width = width;
	buf->height = height;
}

static void cellbuf_resize(struct cellbuf *buf, unsigned int width, unsigned int height)
{
	if (buf->width == width && buf->height == height)
		return;

	unsigned int oldw = buf->width;
	unsigned int oldh = buf->height;
	struct tb_cell *oldcells = buf->cells;

	cellbuf_init(buf, width, height);
	cellbuf_clear(buf);

	unsigned int minw = (width < oldw) ? width : oldw;
	unsigned int minh = (height < oldh) ? height : oldh;
	unsigned int i;

	for (i = 0; i < minh; ++i) {
		struct tb_cell *csrc = oldcells + (i * oldw);
		struct tb_cell *cdst = buf->cells + (i * width);
		memcpy(cdst, csrc, sizeof(struct tb_cell) * minw);
	}

	free(oldcells);
}

static void cellbuf_clear(struct cellbuf *buf)
{
	unsigned int i;
	unsigned int ncells = buf->width * buf->height;

	for (i = 0; i < ncells; ++i) {
		buf->cells[i].ch = ' ';
		buf->cells[i].fg = foreground;
		buf->cells[i].bg = background;
	}
}

static void cellbuf_free(struct cellbuf *buf)
{
	free(buf->cells);
}

static void get_term_size(int *w, int *h)
{
	struct winsize sz;
	memset(&sz, 0, sizeof(sz));

	ioctl(out_fileno, TIOCGWINSZ, &sz);

	if (w) *w = sz.ws_col;
	if (h) *h = sz.ws_row;
}

static void update_term_size()
{
	struct winsize sz;
	memset(&sz, 0, sizeof(sz));

	ioctl(out_fileno, TIOCGWINSZ, &sz);

	termw = sz.ws_col;
	termh = sz.ws_row;
}

static void send_attr(uint16_t fg, uint16_t bg)
{
#define LAST_ATTR_INIT 0xFFFF
	static uint16_t lastfg = LAST_ATTR_INIT, lastbg = LAST_ATTR_INIT;
	if (fg != lastfg || bg != lastbg) {
		fputs(funcs[T_SGR0], out);
		/* TODO: get rid of fprintf */
		fprintf(out, funcs[T_SGR], fg & 0x0F, bg & 0x0F);
		if (fg & TB_BOLD)
			fputs(funcs[T_BOLD], out);
		if (bg & TB_BOLD)
			fputs(funcs[T_BLINK], out);
		if (fg & TB_UNDERLINE)
			fputs(funcs[T_UNDERLINE], out);

		lastfg = fg;
		lastbg = bg;
	}
}

static void send_char(unsigned int x, unsigned int y, uint32_t c)
{
	char buf[7];
	int bw = utf8_unicode_to_char(buf, c);
	buf[bw] = '\0';
	if (x-1 != lastx || y != lasty)
		fprintf(out, funcs[T_MOVE_CURSOR], y+1, x+1); /* TODO: get rid of fprintf */
	lastx = x; lasty = y;
	fputs(buf, out);
}

static void send_clear()
{
	send_attr(foreground, background);
	fputs(funcs[T_CLEAR_SCREEN], out);
	if (!IS_CURSOR_HIDDEN(cursor_x, cursor_y))
		fprintf(out, funcs[T_MOVE_CURSOR], cursor_y+1, cursor_x+1);
	fflush(out);

	/* we need to invalidate cursor position too and these two vars are
	 * used only for simple cursor positioning optimization, cursor
	 * actually may be in the correct place, but we simply discard
	 * optimization once and it gives us simple solution for the case when
	 * cursor moved */
	lastx = LAST_COORD_INIT;
	lasty = LAST_COORD_INIT;
}

static void sigwinch_handler(int xxx)
{
	const int zzz = 1;
	write(winch_fds[1], &zzz, sizeof(int));
}

static void update_size()
{
	update_term_size();
	cellbuf_resize(&back_buffer, termw, termh);
	cellbuf_resize(&front_buffer, termw, termh);
	cellbuf_clear(&front_buffer);
	send_clear();
}

static int wait_fill_event(struct tb_event *event, struct timeval *timeout)
{
	/* ;-) */
#define ENOUGH_DATA_FOR_INPUT_PARSING 128
	int result;
	char buf[ENOUGH_DATA_FOR_INPUT_PARSING];
	fd_set events;
	memset(event, 0, sizeof(struct tb_event));

	/* try to extract event from input buffer, return on success */
	event->type = TB_EVENT_KEY;
	if (extract_event(event, &inbuf, inputmode) == 0)
		return TB_EVENT_KEY;

	/* it looks like input buffer is incomplete, let's try the short path */
	size_t r = fread(buf, 1, ENOUGH_DATA_FOR_INPUT_PARSING, in);
	if (r < ENOUGH_DATA_FOR_INPUT_PARSING && feof(in))
		clearerr(in);
	if (r > 0) {
		if (ringbuffer_free_space(&inbuf) < r)
			return -1;
		ringbuffer_push(&inbuf, buf, r);
		if (extract_event(event, &inbuf, inputmode) == 0)
			return TB_EVENT_KEY;
	}

	/* no stuff in FILE's internal buffer, block in select */
	while (1) {
		FD_ZERO(&events);
		FD_SET(in_fileno, &events);
		FD_SET(winch_fds[0], &events);
		int maxfd = (winch_fds[0] > in_fileno) ? winch_fds[0] : in_fileno;
		result = select(maxfd+1, &events, 0, 0, timeout);
		if (!result)
			return 0;

		if (FD_ISSET(in_fileno, &events)) {
			event->type = TB_EVENT_KEY;
			size_t r = fread(buf, 1, ENOUGH_DATA_FOR_INPUT_PARSING, in);
			if (r < ENOUGH_DATA_FOR_INPUT_PARSING && feof(in))
				clearerr(in);
			if (r == 0)
				continue;
			/* if there is no free space in input buffer, return error */
			if (ringbuffer_free_space(&inbuf) < r)
				return -1;
			/* fill buffer */
			ringbuffer_push(&inbuf, buf, r);
			if (extract_event(event, &inbuf, inputmode) == 0)
				return TB_EVENT_KEY;
		}
		if (FD_ISSET(winch_fds[0], &events)) {
			event->type = TB_EVENT_RESIZE;
			int zzz = 0;
			read(winch_fds[0], &zzz, sizeof(int));
			buffer_size_change_request = 1;
			get_term_size(&event->w, &event->h);
			return TB_EVENT_RESIZE;
		}
	}
}

