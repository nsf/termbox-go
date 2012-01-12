#include <stdlib.h>
#include <string.h>

#include "term.h"

/* Eterm */
static const char *Eterm_keys[] = {
	"\033[11~","\033[12~","\033[13~","\033[14~","\033[15~","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033[7~","\033[8~","\033[5~","\033[6~","\033[A","\033[B","\033[D","\033[C", 0
};
static const char *Eterm_funcs[] = {
	"\0337\033[?47h","\033[2J\033[?47l\0338","\033[?25h","\033[?25l","\033[H\033[2J","\033[3%u;4%um","\033[m","\033[4m","\033[1m","\033[5m","\033[%u;%uH","",""
};

/* screen */
static const char *screen_keys[] = {
	"\033OP","\033OQ","\033OR","\033OS","\033[15~","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033[1~","\033[4~","\033[5~","\033[6~","\033OA","\033OB","\033OD","\033OC", 0
};
static const char *screen_funcs[] = {
	"\033[?1049h","\033[?1049l","\033[34h\033[?25h","\033[?25l","\033[H\033[J","\033[3%u;4%um","\033[m","\033[4m","\033[1m","\033[5m","\033[%u;%uH","\033[?1h\033=","\033[?1l\033>"
};

/* xterm */
static const char *xterm_keys[] = {
	"\033OP","\033OQ","\033OR","\033OS","\033[15~","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033OH","\033OF","\033[5~","\033[6~","\033OA","\033OB","\033OD","\033OC", 0
};
static const char *xterm_funcs[] = {
	"\033[?1049h","\033[?1049l","\033[?12l\033[?25h","\033[?25l","\033[H\033[2J","\033[3%u;4%um","\033(B\033[m","\033[4m","\033[1m","\033[5m","\033[%u;%uH","\033[?1h\033=","\033[?1l\033>"
};

/* rxvt-unicode */
static const char *rxvt_unicode_keys[] = {
	"\033[11~","\033[12~","\033[13~","\033[14~","\033[15~","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033[7~","\033[8~","\033[5~","\033[6~","\033[A","\033[B","\033[D","\033[C", 0
};
static const char *rxvt_unicode_funcs[] = {
	"\033[?1049h","\033[r\033[?1049l","\033[?25h","\033[?25l","\033[H\033[2J","\033[3%u;4%um","\033[m\033(B","\033[4m","\033[1m","\033[5m","\033[%u;%uH","\033=","\033>"
};

/* linux */
static const char *linux_keys[] = {
	"\033[[A","\033[[B","\033[[C","\033[[D","\033[[E","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033[1~","\033[4~","\033[5~","\033[6~","\033[A","\033[B","\033[D","\033[C", 0
};
static const char *linux_funcs[] = {
	"","","\033[?25h\033[?0c","\033[?25l\033[?1c","\033[H\033[J","\033[3%u;4%um","\033[0;10m","\033[4m","\033[1m","\033[5m","\033[%u;%uH","",""
};

/* rxvt-256color */
static const char *rxvt_256color_keys[] = {
	"\033[11~","\033[12~","\033[13~","\033[14~","\033[15~","\033[17~","\033[18~","\033[19~","\033[20~","\033[21~","\033[23~","\033[24~","\033[2~","\033[3~","\033[7~","\033[8~","\033[5~","\033[6~","\033[A","\033[B","\033[D","\033[C", 0
};
static const char *rxvt_256color_funcs[] = {
	"\0337\033[?47h","\033[2J\033[?47l\0338","\033[?25h","\033[?25l","\033[H\033[2J","\033[3%u;4%um","\033[m","\033[4m","\033[1m","\033[5m","\033[%u;%uH","\033=","\033>"
};

static struct term {
	const char *name;
	const char **keys;
	const char **funcs;
} terms[] = {
	{"Eterm", Eterm_keys, Eterm_funcs},
	{"screen", screen_keys, screen_funcs},
	{"xterm", xterm_keys, xterm_funcs},
	{"rxvt-unicode", rxvt_unicode_keys, rxvt_unicode_funcs},
	{"linux", linux_keys, linux_funcs},
	{"rxvt-256color", rxvt_256color_keys, rxvt_256color_funcs},
	{0,0,0}
};

const char **keys;
const char **funcs;

static int try_compatible(const char *term, const char *name,
			  const char **tkeys, const char **tfuncs)
{
	if (strstr(term, name)) {
		keys = tkeys;
		funcs = tfuncs;
		return 0;
	}

	return EUNSUPPORTED_TERM;
}

int init_term()
{
	int i;
	const char *term = getenv("TERM");
	for (i = 0; terms[i].name; i++) {
		if (!strcmp(terms[i].name, term)) {
			keys = terms[i].keys;
			funcs = terms[i].funcs;
			return 0;
		}
	}

	/* let's do some heuristic, maybe it's a compatible terminal */
	if (try_compatible(term, "xterm", xterm_keys, xterm_funcs) == 0)
		return 0;
	if (try_compatible(term, "rxvt", rxvt_unicode_keys, rxvt_unicode_funcs) == 0)
		return 0;
	if (try_compatible(term, "linux", linux_keys, linux_funcs) == 0)
		return 0;
	if (try_compatible(term, "Eterm", Eterm_keys, Eterm_funcs) == 0)
		return 0;
	if (try_compatible(term, "screen", screen_keys, screen_funcs) == 0)
		return 0;
	/* let's assume that 'cygwin' is xterm compatible */
	if (try_compatible(term, "cygwin", xterm_keys, xterm_funcs) == 0)
		return 0;

	return EUNSUPPORTED_TERM;
}
