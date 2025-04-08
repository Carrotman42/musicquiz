package oslevelinput

import "fmt"

// Note: these were copy-paste-edited out of input-event-codes.h from my local
// workstation. Presumably they don't change much.

type EventType uint16

const (
	EV_SYN       EventType = 0x00
	EV_KEY       EventType = 0x01
	EV_REL       EventType = 0x02
	EV_ABS       EventType = 0x03
	EV_MSC       EventType = 0x04
	EV_SW        EventType = 0x05
	EV_LED       EventType = 0x11
	EV_SND       EventType = 0x12
	EV_REP       EventType = 0x14
	EV_FF        EventType = 0x15
	EV_PWR       EventType = 0x16
	EV_FF_STATUS EventType = 0x17
	EV_MAX       EventType = 0x1f

// EV_CNT EventType	=		(EV_MAX+1)
)

/*
 * Synchronization events.
 */

const (
	SYN_REPORT    EventCode = 0
	SYN_CONFIG    EventCode = 1
	SYN_MT_REPORT EventCode = 2
	SYN_DROPPED   EventCode = 3
	SYN_MAX       EventCode = 0xf

// SYN_CNT           EventCode      =  (SYN_MAX+1)
)

type EventCode uint16

const (
	KEY_RESERVED   EventCode = 0
	KEY_ESC        EventCode = 1
	KEY_1          EventCode = 2
	KEY_2          EventCode = 3
	KEY_3          EventCode = 4
	KEY_4          EventCode = 5
	KEY_5          EventCode = 6
	KEY_6          EventCode = 7
	KEY_7          EventCode = 8
	KEY_8          EventCode = 9
	KEY_9          EventCode = 10
	KEY_0          EventCode = 11
	KEY_MINUS      EventCode = 12
	KEY_EQUAL      EventCode = 13
	KEY_BACKSPACE  EventCode = 14
	KEY_TAB        EventCode = 15
	KEY_Q          EventCode = 16
	KEY_W          EventCode = 17
	KEY_E          EventCode = 18
	KEY_R          EventCode = 19
	KEY_T          EventCode = 20
	KEY_Y          EventCode = 21
	KEY_U          EventCode = 22
	KEY_I          EventCode = 23
	KEY_O          EventCode = 24
	KEY_P          EventCode = 25
	KEY_LEFTBRACE  EventCode = 26
	KEY_RIGHTBRACE EventCode = 27
	KEY_ENTER      EventCode = 28
	KEY_LEFTCTRL   EventCode = 29
	KEY_A          EventCode = 30
	KEY_S          EventCode = 31
	KEY_D          EventCode = 32
	KEY_F          EventCode = 33
	KEY_G          EventCode = 34
	KEY_H          EventCode = 35
	KEY_J          EventCode = 36
	KEY_K          EventCode = 37
	KEY_L          EventCode = 38
	KEY_SEMICOLON  EventCode = 39
	KEY_APOSTROPHE EventCode = 40
	KEY_GRAVE      EventCode = 41
	KEY_LEFTSHIFT  EventCode = 42
	KEY_BACKSLASH  EventCode = 43
	KEY_Z          EventCode = 44
	KEY_X          EventCode = 45
	KEY_C          EventCode = 46
	KEY_V          EventCode = 47
	KEY_B          EventCode = 48
	KEY_N          EventCode = 49
	KEY_M          EventCode = 50
	KEY_COMMA      EventCode = 51
	KEY_DOT        EventCode = 52
	KEY_SLASH      EventCode = 53
	KEY_RIGHTSHIFT EventCode = 54
	KEY_KPASTERISK EventCode = 55
	KEY_LEFTALT    EventCode = 56
	KEY_SPACE      EventCode = 57
	KEY_CAPSLOCK   EventCode = 58
	KEY_F1         EventCode = 59
	KEY_F2         EventCode = 60
	KEY_F3         EventCode = 61
	KEY_F4         EventCode = 62
	KEY_F5         EventCode = 63
	KEY_F6         EventCode = 64
	KEY_F7         EventCode = 65
	KEY_F8         EventCode = 66
	KEY_F9         EventCode = 67
	KEY_F10        EventCode = 68
	KEY_NUMLOCK    EventCode = 69
	KEY_SCROLLLOCK EventCode = 70
	KEY_KP7        EventCode = 71
	KEY_KP8        EventCode = 72
	KEY_KP9        EventCode = 73
	KEY_KPMINUS    EventCode = 74
	KEY_KP4        EventCode = 75
	KEY_KP5        EventCode = 76
	KEY_KP6        EventCode = 77
	KEY_KPPLUS     EventCode = 78
	KEY_KP1        EventCode = 79
	KEY_KP2        EventCode = 80
	KEY_KP3        EventCode = 81
	KEY_KP0        EventCode = 82
	KEY_KPDOT      EventCode = 83

	KEY_ZENKAKUHANKAKU   EventCode = 85
	KEY_102ND            EventCode = 86
	KEY_F11              EventCode = 87
	KEY_F12              EventCode = 88
	KEY_RO               EventCode = 89
	KEY_KATAKANA         EventCode = 90
	KEY_HIRAGANA         EventCode = 91
	KEY_HENKAN           EventCode = 92
	KEY_KATAKANAHIRAGANA EventCode = 93
	KEY_MUHENKAN         EventCode = 94
	KEY_KPJPCOMMA        EventCode = 95
	KEY_KPENTER          EventCode = 96
	KEY_RIGHTCTRL        EventCode = 97
	KEY_KPSLASH          EventCode = 98
	KEY_SYSRQ            EventCode = 99
	KEY_RIGHTALT         EventCode = 100
	KEY_LINEFEED         EventCode = 101
	KEY_HOME             EventCode = 102
	KEY_UP               EventCode = 103
	KEY_PAGEUP           EventCode = 104
	KEY_LEFT             EventCode = 105
	KEY_RIGHT            EventCode = 106
	KEY_END              EventCode = 107
	KEY_DOWN             EventCode = 108
	KEY_PAGEDOWN         EventCode = 109
	KEY_INSERT           EventCode = 110
	KEY_DELETE           EventCode = 111
	KEY_MACRO            EventCode = 112
	KEY_MUTE             EventCode = 113
	KEY_VOLUMEDOWN       EventCode = 114
	KEY_VOLUMEUP         EventCode = 115
	KEY_POWER            EventCode = 116 /* SC System Power Down */
	KEY_KPEQUAL          EventCode = 117
	KEY_KPPLUSMINUS      EventCode = 118
	KEY_PAUSE            EventCode = 119
	KEY_SCALE            EventCode = 120 /* AL Compiz Scale (Expose) */

	KEY_KPCOMMA   EventCode = 121
	KEY_HANGEUL   EventCode = 122
	KEY_HANGUEL   EventCode = KEY_HANGEUL
	KEY_HANJA     EventCode = 123
	KEY_YEN       EventCode = 124
	KEY_LEFTMETA  EventCode = 125
	KEY_RIGHTMETA EventCode = 126
	KEY_COMPOSE   EventCode = 127

	KEY_STOP           EventCode = 128 /* AC Stop */
	KEY_AGAIN          EventCode = 129
	KEY_PROPS          EventCode = 130 /* AC Properties */
	KEY_UNDO           EventCode = 131 /* AC Undo */
	KEY_FRONT          EventCode = 132
	KEY_COPY           EventCode = 133 /* AC Copy */
	KEY_OPEN           EventCode = 134 /* AC Open */
	KEY_PASTE          EventCode = 135 /* AC Paste */
	KEY_FIND           EventCode = 136 /* AC Search */
	KEY_CUT            EventCode = 137 /* AC Cut */
	KEY_HELP           EventCode = 138 /* AL Integrated Help Center */
	KEY_MENU           EventCode = 139 /* Menu (show menu) */
	KEY_CALC           EventCode = 140 /* AL Calculator */
	KEY_SETUP          EventCode = 141
	KEY_SLEEP          EventCode = 142 /* SC System Sleep */
	KEY_WAKEUP         EventCode = 143 /* System Wake Up */
	KEY_FILE           EventCode = 144 /* AL Local Machine Browser */
	KEY_SENDFILE       EventCode = 145
	KEY_DELETEFILE     EventCode = 146
	KEY_XFER           EventCode = 147
	KEY_PROG1          EventCode = 148
	KEY_PROG2          EventCode = 149
	KEY_WWW            EventCode = 150 /* AL Internet Browser */
	KEY_MSDOS          EventCode = 151
	KEY_COFFEE         EventCode = 152 /* AL Terminal Lock/Screensaver */
	KEY_SCREENLOCK     EventCode = KEY_COFFEE
	KEY_ROTATE_DISPLAY EventCode = 153 /* Display orientation for e.g. tablets */
	KEY_DIRECTION      EventCode = KEY_ROTATE_DISPLAY
	KEY_CYCLEWINDOWS   EventCode = 154
	KEY_MAIL           EventCode = 155
	KEY_BOOKMARKS      EventCode = 156 /* AC Bookmarks */
	KEY_COMPUTER       EventCode = 157
	KEY_BACK           EventCode = 158 /* AC Back */
	KEY_FORWARD        EventCode = 159 /* AC Forward */
	KEY_CLOSECD        EventCode = 160
	KEY_EJECTCD        EventCode = 161
	KEY_EJECTCLOSECD   EventCode = 162
	KEY_NEXTSONG       EventCode = 163
	KEY_PLAYPAUSE      EventCode = 164
	KEY_PREVIOUSSONG   EventCode = 165
	KEY_STOPCD         EventCode = 166
	KEY_RECORD         EventCode = 167
	KEY_REWIND         EventCode = 168
	KEY_PHONE          EventCode = 169 /* Media Select Telephone */
	KEY_ISO            EventCode = 170
	KEY_CONFIG         EventCode = 171 /* AL Consumer Control Configuration */
	KEY_HOMEPAGE       EventCode = 172 /* AC Home */
	KEY_REFRESH        EventCode = 173 /* AC Refresh */
	KEY_EXIT           EventCode = 174 /* AC Exit */
	KEY_MOVE           EventCode = 175
	KEY_EDIT           EventCode = 176
	KEY_SCROLLUP       EventCode = 177
	KEY_SCROLLDOWN     EventCode = 178
	KEY_KPLEFTPAREN    EventCode = 179
	KEY_KPRIGHTPAREN   EventCode = 180
	KEY_NEW            EventCode = 181 /* AC New */
	KEY_REDO           EventCode = 182 /* AC Redo/Repeat */

	KEY_F13 EventCode = 183
	KEY_F14 EventCode = 184
	KEY_F15 EventCode = 185
	KEY_F16 EventCode = 186
	KEY_F17 EventCode = 187
	KEY_F18 EventCode = 188
	KEY_F19 EventCode = 189
	KEY_F20 EventCode = 190
	KEY_F21 EventCode = 191
	KEY_F22 EventCode = 192
	KEY_F23 EventCode = 193
	KEY_F24 EventCode = 194

	KEY_PLAYCD           EventCode = 200
	KEY_PAUSECD          EventCode = 201
	KEY_PROG3            EventCode = 202
	KEY_PROG4            EventCode = 203
	KEY_ALL_APPLICATIONS EventCode = 204 /* AC Desktop Show All Applications */
	KEY_DASHBOARD        EventCode = KEY_ALL_APPLICATIONS
	KEY_SUSPEND          EventCode = 205
	KEY_CLOSE            EventCode = 206 /* AC Close */
	KEY_PLAY             EventCode = 207
	KEY_FASTFORWARD      EventCode = 208
	KEY_BASSBOOST        EventCode = 209
	KEY_PRINT            EventCode = 210 /* AC Print */
	KEY_HP               EventCode = 211
	KEY_CAMERA           EventCode = 212
	KEY_SOUND            EventCode = 213
	KEY_QUESTION         EventCode = 214
	KEY_EMAIL            EventCode = 215
	KEY_CHAT             EventCode = 216
	KEY_SEARCH           EventCode = 217
	KEY_CONNECT          EventCode = 218
	KEY_FINANCE          EventCode = 219 /* AL Checkbook/Finance */
	KEY_SPORT            EventCode = 220
	KEY_SHOP             EventCode = 221
	KEY_ALTERASE         EventCode = 222
	KEY_CANCEL           EventCode = 223 /* AC Cancel */
	KEY_BRIGHTNESSDOWN   EventCode = 224
	KEY_BRIGHTNESSUP     EventCode = 225
	KEY_MEDIA            EventCode = 226

	KEY_SWITCHVIDEOMODE EventCode = 227 /* Cycle between available video
	   outputs (Monitor/LCD/TV-out/etc) */
	KEY_KBDILLUMTOGGLE EventCode = 228
	KEY_KBDILLUMDOWN   EventCode = 229
	KEY_KBDILLUMUP     EventCode = 230

	KEY_SEND        EventCode = 231 /* AC Send */
	KEY_REPLY       EventCode = 232 /* AC Reply */
	KEY_FORWARDMAIL EventCode = 233 /* AC Forward Msg */
	KEY_SAVE        EventCode = 234 /* AC Save */
	KEY_DOCUMENTS   EventCode = 235

	KEY_BATTERY EventCode = 236

	KEY_BLUETOOTH EventCode = 237
	KEY_WLAN      EventCode = 238
	KEY_UWB       EventCode = 239

	KEY_UNKNOWN EventCode = 240

	KEY_VIDEO_NEXT       EventCode = 241 /* drive next video source */
	KEY_VIDEO_PREV       EventCode = 242 /* drive previous video source */
	KEY_BRIGHTNESS_CYCLE EventCode = 243 /* brightness up, after max is min */
	KEY_BRIGHTNESS_AUTO  EventCode = 244 /* Set Auto Brightness: manual
	brightness control is off,
	rely on ambient */
	KEY_BRIGHTNESS_ZERO EventCode = KEY_BRIGHTNESS_AUTO
	KEY_DISPLAY_OFF     EventCode = 245 /* display device to off state */

	KEY_WWAN   EventCode = 246 /* Wireless WAN (LTE, UMTS, GSM, etc.) */
	KEY_WIMAX  EventCode = KEY_WWAN
	KEY_RFKILL EventCode = 247 /* Key that controls all radios */

	KEY_MICMUTE EventCode = 248 /* Mute / unmute the microphone */

	// KEYBOARD
	BTN_0 EventCode = 0x100
	BTN_1 EventCode = 0x101
	BTN_2 EventCode = 0x102
	BTN_3 EventCode = 0x103
	BTN_4 EventCode = 0x104
	BTN_5 EventCode = 0x105
	BTN_6 EventCode = 0x106
	BTN_7 EventCode = 0x107
	BTN_8 EventCode = 0x108
	BTN_9 EventCode = 0x109

	BTN_MOUSE   EventCode = 0x110
	BTN_LEFT    EventCode = 0x110
	BTN_RIGHT   EventCode = 0x111
	BTN_MIDDLE  EventCode = 0x112
	BTN_SIDE    EventCode = 0x113
	BTN_EXTRA   EventCode = 0x114
	BTN_FORWARD EventCode = 0x115
	BTN_BACK    EventCode = 0x116
	BTN_TASK    EventCode = 0x117

	BTN_JOYSTICK EventCode = 0x120
	BTN_TRIGGER  EventCode = 0x120
	BTN_THUMB    EventCode = 0x121
	BTN_THUMB2   EventCode = 0x122
	BTN_TOP      EventCode = 0x123
	BTN_TOP2     EventCode = 0x124
	BTN_PINKIE   EventCode = 0x125
	BTN_BASE     EventCode = 0x126
	BTN_BASE2    EventCode = 0x127
	BTN_BASE3    EventCode = 0x128
	BTN_BASE4    EventCode = 0x129
	BTN_BASE5    EventCode = 0x12a
	BTN_BASE6    EventCode = 0x12b
	BTN_DEAD     EventCode = 0x12f

	BTN_GAMEPAD EventCode = 0x130
	BTN_SOUTH   EventCode = 0x130
	BTN_A       EventCode = BTN_SOUTH
	BTN_EAST    EventCode = 0x131
	BTN_B       EventCode = BTN_EAST
	BTN_C       EventCode = 0x132
	BTN_NORTH   EventCode = 0x133
	BTN_X       EventCode = BTN_NORTH
	BTN_WEST    EventCode = 0x134
	BTN_Y       EventCode = BTN_WEST
	BTN_Z       EventCode = 0x135
	BTN_TL      EventCode = 0x136
	BTN_TR      EventCode = 0x137
	BTN_TL2     EventCode = 0x138
	BTN_TR2     EventCode = 0x139
	BTN_SELECT  EventCode = 0x13a
	BTN_START   EventCode = 0x13b
	BTN_MODE    EventCode = 0x13c
	BTN_THUMBL  EventCode = 0x13d
	BTN_THUMBR  EventCode = 0x13e

	BTN_DIGI           EventCode = 0x140
	BTN_TOOL_PEN       EventCode = 0x140
	BTN_TOOL_RUBBER    EventCode = 0x141
	BTN_TOOL_BRUSH     EventCode = 0x142
	BTN_TOOL_PENCIL    EventCode = 0x143
	BTN_TOOL_AIRBRUSH  EventCode = 0x144
	BTN_TOOL_FINGER    EventCode = 0x145
	BTN_TOOL_MOUSE     EventCode = 0x146
	BTN_TOOL_LENS      EventCode = 0x147
	BTN_TOOL_QUINTTAP  EventCode = 0x148 /* Five fingers on trackpad */
	BTN_STYLUS3        EventCode = 0x149
	BTN_TOUCH          EventCode = 0x14a
	BTN_STYLUS         EventCode = 0x14b
	BTN_STYLUS2        EventCode = 0x14c
	BTN_TOOL_DOUBLETAP EventCode = 0x14d
	BTN_TOOL_TRIPLETAP EventCode = 0x14e
	BTN_TOOL_QUADTAP   EventCode = 0x14f /* Four fingers on trackpad */

	BTN_WHEEL     EventCode = 0x150
	BTN_GEAR_DOWN EventCode = 0x150
	BTN_GEAR_UP   EventCode = 0x151
)

func (ec EventCode) String() string {
	switch ec {
	case KEY_RESERVED:
		return "RESERVED"
	case KEY_ESC:
		return "ESC"
	case KEY_1:
		return "1"
	case KEY_2:
		return "2"
	case KEY_3:
		return "3"
	case KEY_4:
		return "4"
	case KEY_5:
		return "5"
	case KEY_6:
		return "6"
	case KEY_7:
		return "7"
	case KEY_8:
		return "8"
	case KEY_9:
		return "9"
	case KEY_0:
		return "0"
	case KEY_MINUS:
		return "MINUS"
	case KEY_EQUAL:
		return "EQUAL"
	case KEY_BACKSPACE:
		return "BACKSPACE"
	case KEY_TAB:
		return "TAB"
	case KEY_Q:
		return "Q"
	case KEY_W:
		return "W"
	case KEY_E:
		return "E"
	case KEY_R:
		return "R"
	case KEY_T:
		return "T"
	case KEY_Y:
		return "Y"
	case KEY_U:
		return "U"
	case KEY_I:
		return "I"
	case KEY_O:
		return "O"
	case KEY_P:
		return "P"
	case KEY_LEFTBRACE:
		return "LEFTBRACE"
	case KEY_RIGHTBRACE:
		return "RIGHTBRACE"
	case KEY_ENTER:
		return "ENTER"
	case KEY_LEFTCTRL:
		return "LEFTCTRL"
	case KEY_A:
		return "A"
	case KEY_S:
		return "S"
	case KEY_D:
		return "D"
	case KEY_F:
		return "F"
	case KEY_G:
		return "G"
	case KEY_H:
		return "H"
	case KEY_J:
		return "J"
	case KEY_K:
		return "K"
	case KEY_L:
		return "L"
	case KEY_SEMICOLON:
		return "SEMICOLON"
	case KEY_APOSTROPHE:
		return "APOSTROPHE"
	case KEY_GRAVE:
		return "GRAVE"
	case KEY_LEFTSHIFT:
		return "LEFTSHIFT"
	case KEY_BACKSLASH:
		return "BACKSLASH"
	case KEY_Z:
		return "Z"
	case KEY_X:
		return "X"
	case KEY_C:
		return "C"
	case KEY_V:
		return "V"
	case KEY_B:
		return "B"
	case KEY_N:
		return "N"
	case KEY_M:
		return "M"
	case KEY_COMMA:
		return "COMMA"
	case KEY_DOT:
		return "DOT"
	case KEY_SLASH:
		return "SLASH"
	case KEY_RIGHTSHIFT:
		return "RIGHTSHIFT"
	case KEY_KPASTERISK:
		return "KPASTERISK"
	case KEY_LEFTALT:
		return "LEFTALT"
	case KEY_SPACE:
		return "SPACE"
	case KEY_CAPSLOCK:
		return "CAPSLOCK"
	case KEY_F1:
		return "F1"
	case KEY_F2:
		return "F2"
	case KEY_F3:
		return "F3"
	case KEY_F4:
		return "F4"
	case KEY_F5:
		return "F5"
	case KEY_F6:
		return "F6"
	case KEY_F7:
		return "F7"
	case KEY_F8:
		return "F8"
	case KEY_F9:
		return "F9"
	case KEY_F10:
		return "F10"
	case KEY_NUMLOCK:
		return "NUMLOCK"
	case KEY_SCROLLLOCK:
		return "SCROLLLOCK"
	case KEY_KP7:
		return "KP7"
	case KEY_KP8:
		return "KP8"
	case KEY_KP9:
		return "KP9"
	case KEY_KPMINUS:
		return "KPMINUS"
	case KEY_KP4:
		return "KP4"
	case KEY_KP5:
		return "KP5"
	case KEY_KP6:
		return "KP6"
	case KEY_KPPLUS:
		return "KPPLUS"
	case KEY_KP1:
		return "KP1"
	case KEY_KP2:
		return "KP2"
	case KEY_KP3:
		return "KP3"
	case KEY_KP0:
		return "KP0"
	case KEY_KPDOT:
		return "KPDOT"

	case KEY_ZENKAKUHANKAKU:
		return "ZENKAKUHANKAKU"
	case KEY_102ND:
		return "102ND"
	case KEY_F11:
		return "F11"
	case KEY_F12:
		return "F12"
	case KEY_RO:
		return "RO"
	case KEY_KATAKANA:
		return "KATAKANA"
	case KEY_HIRAGANA:
		return "HIRAGANA"
	case KEY_HENKAN:
		return "HENKAN"
	case KEY_KATAKANAHIRAGANA:
		return "KATAKANAHIRAGANA"
	case KEY_MUHENKAN:
		return "MUHENKAN"
	case KEY_KPJPCOMMA:
		return "KPJPCOMMA"
	case KEY_KPENTER:
		return "KPENTER"
	case KEY_RIGHTCTRL:
		return "RIGHTCTRL"
	case KEY_KPSLASH:
		return "KPSLASH"
	case KEY_SYSRQ:
		return "SYSRQ"
	case KEY_RIGHTALT:
		return "RIGHTALT"
	case KEY_LINEFEED:
		return "LINEFEED"
	case KEY_HOME:
		return "HOME"
	case KEY_UP:
		return "UP"
	case KEY_PAGEUP:
		return "PAGEUP"
	case KEY_LEFT:
		return "LEFT"
	case KEY_RIGHT:
		return "RIGHT"
	case KEY_END:
		return "END"
	case KEY_DOWN:
		return "DOWN"
	case KEY_PAGEDOWN:
		return "PAGEDOWN"
	case KEY_INSERT:
		return "INSERT"
	case KEY_DELETE:
		return "DELETE"
	case KEY_MACRO:
		return "MACRO"
	case KEY_MUTE:
		return "MUTE"
	case KEY_VOLUMEDOWN:
		return "VOLUMEDOWN"
	case KEY_VOLUMEUP:
		return "VOLUMEUP"
	case KEY_POWER:
		return "POWER"
	case KEY_KPEQUAL:
		return "KPEQUAL"
	case KEY_KPPLUSMINUS:
		return "KPPLUSMINUS"
	case KEY_PAUSE:
		return "PAUSE"
	case KEY_SCALE:
		return "SCALE"

	case KEY_KPCOMMA:
		return "KPCOMMA"
	case KEY_HANGEUL:
		return "HANGEUL"
		//case KEY_HANGUEL: return "HANGUEL"
	case KEY_HANJA:
		return "HANJA"
	case KEY_YEN:
		return "YEN"
	case KEY_LEFTMETA:
		return "LEFTMETA"
	case KEY_RIGHTMETA:
		return "RIGHTMETA"
	case KEY_COMPOSE:
		return "COMPOSE"

	case KEY_STOP:
		return "STOP"
	case KEY_AGAIN:
		return "AGAIN"
	case KEY_PROPS:
		return "PROPS"
	case KEY_UNDO:
		return "UNDO"
	case KEY_FRONT:
		return "FRONT"
	case KEY_COPY:
		return "COPY"
	case KEY_OPEN:
		return "OPEN"
	case KEY_PASTE:
		return "PASTE"
	case KEY_FIND:
		return "FIND"
	case KEY_CUT:
		return "CUT"
	case KEY_HELP:
		return "HELP"
	case KEY_MENU:
		return "MENU"
	case KEY_CALC:
		return "CALC"
	case KEY_SETUP:
		return "SETUP"
	case KEY_SLEEP:
		return "SLEEP"
	case KEY_WAKEUP:
		return "WAKEUP"
	case KEY_FILE:
		return "FILE"
	case KEY_SENDFILE:
		return "SENDFILE"
	case KEY_DELETEFILE:
		return "DELETEFILE"
	case KEY_XFER:
		return "XFER"
	case KEY_PROG1:
		return "PROG1"
	case KEY_PROG2:
		return "PROG2"
	case KEY_WWW:
		return "WWW"
	case KEY_MSDOS:
		return "MSDOS"
	case KEY_COFFEE:
		return "COFFEE"
		//case KEY_SCREENLOCK: return "SCREENLOCK"
	case KEY_ROTATE_DISPLAY:
		return "ROTATE_DISPLAY"
		//case KEY_DIRECTION: return "DIRECTION"
	case KEY_CYCLEWINDOWS:
		return "CYCLEWINDOWS"
	case KEY_MAIL:
		return "MAIL"
	case KEY_BOOKMARKS:
		return "BOOKMARKS"
	case KEY_COMPUTER:
		return "COMPUTER"
	case KEY_BACK:
		return "BACK"
	case KEY_FORWARD:
		return "FORWARD"
	case KEY_CLOSECD:
		return "CLOSECD"
	case KEY_EJECTCD:
		return "EJECTCD"
	case KEY_EJECTCLOSECD:
		return "EJECTCLOSECD"
	case KEY_NEXTSONG:
		return "NEXTSONG"
	case KEY_PLAYPAUSE:
		return "PLAYPAUSE"
	case KEY_PREVIOUSSONG:
		return "PREVIOUSSONG"
	case KEY_STOPCD:
		return "STOPCD"
	case KEY_RECORD:
		return "RECORD"
	case KEY_REWIND:
		return "REWIND"
	case KEY_PHONE:
		return "PHONE"
	case KEY_ISO:
		return "ISO"
	case KEY_CONFIG:
		return "CONFIG"
	case KEY_HOMEPAGE:
		return "HOMEPAGE"
	case KEY_REFRESH:
		return "REFRESH"
	case KEY_EXIT:
		return "EXIT"
	case KEY_MOVE:
		return "MOVE"
	case KEY_EDIT:
		return "EDIT"
	case KEY_SCROLLUP:
		return "SCROLLUP"
	case KEY_SCROLLDOWN:
		return "SCROLLDOWN"
	case KEY_KPLEFTPAREN:
		return "KPLEFTPAREN"
	case KEY_KPRIGHTPAREN:
		return "KPRIGHTPAREN"
	case KEY_NEW:
		return "NEW"
	case KEY_REDO:
		return "REDO"

	case KEY_F13:
		return "F13"
	case KEY_F14:
		return "F14"
	case KEY_F15:
		return "F15"
	case KEY_F16:
		return "F16"
	case KEY_F17:
		return "F17"
	case KEY_F18:
		return "F18"
	case KEY_F19:
		return "F19"
	case KEY_F20:
		return "F20"
	case KEY_F21:
		return "F21"
	case KEY_F22:
		return "F22"
	case KEY_F23:
		return "F23"
	case KEY_F24:
		return "F24"

	case KEY_PLAYCD:
		return "PLAYCD"
	case KEY_PAUSECD:
		return "PAUSECD"
	case KEY_PROG3:
		return "PROG3"
	case KEY_PROG4:
		return "PROG4"
	case KEY_ALL_APPLICATIONS:
		return "ALL_APPLICATIONS"
		//case KEY_DASHBOARD: return "DASHBOARD"
	case KEY_SUSPEND:
		return "SUSPEND"
	case KEY_CLOSE:
		return "CLOSE"
	case KEY_PLAY:
		return "PLAY"
	case KEY_FASTFORWARD:
		return "FASTFORWARD"
	case KEY_BASSBOOST:
		return "BASSBOOST"
	case KEY_PRINT:
		return "PRINT"
	case KEY_HP:
		return "HP"
	case KEY_CAMERA:
		return "CAMERA"
	case KEY_SOUND:
		return "SOUND"
	case KEY_QUESTION:
		return "QUESTION"
	case KEY_EMAIL:
		return "EMAIL"
	case KEY_CHAT:
		return "CHAT"
	case KEY_SEARCH:
		return "SEARCH"
	case KEY_CONNECT:
		return "CONNECT"
	case KEY_FINANCE:
		return "FINANCE"
	case KEY_SPORT:
		return "SPORT"
	case KEY_SHOP:
		return "SHOP"
	case KEY_ALTERASE:
		return "ALTERASE"
	case KEY_CANCEL:
		return "CANCEL"
	case KEY_BRIGHTNESSDOWN:
		return "BRIGHTNESSDOWN"
	case KEY_BRIGHTNESSUP:
		return "BRIGHTNESSUP"
	case KEY_MEDIA:
		return "MEDIA"

	case KEY_SWITCHVIDEOMODE:
		return "SWITCHVIDEOMODE"
	case KEY_KBDILLUMTOGGLE:
		return "KBDILLUMTOGGLE"
	case KEY_KBDILLUMDOWN:
		return "KBDILLUMDOWN"
	case KEY_KBDILLUMUP:
		return "KBDILLUMUP"

	case KEY_SEND:
		return "SEND"
	case KEY_REPLY:
		return "REPLY"
	case KEY_FORWARDMAIL:
		return "FORWARDMAIL"
	case KEY_SAVE:
		return "SAVE"
	case KEY_DOCUMENTS:
		return "DOCUMENTS"
	case KEY_BATTERY:
		return "BATTERY"
	case KEY_BLUETOOTH:
		return "BLUETOOTH"
	case KEY_WLAN:
		return "WLAN"
	case KEY_UWB:
		return "UWB"
	case KEY_UNKNOWN:
		return "UNKNOWN"
	case KEY_VIDEO_NEXT:
		return "VIDEO_NEXT"
	case KEY_VIDEO_PREV:
		return "VIDEO_PREV"
	case KEY_BRIGHTNESS_CYCLE:
		return "BRIGHTNESS_CYCLE"
	case KEY_BRIGHTNESS_AUTO:
		return "BRIGHTNESS_AUTO"
		//case KEY_BRIGHTNESS_ZERO: return "BRIGHTNESS_ZERO"
	case KEY_DISPLAY_OFF:
		return "DISPLAY_OFF"
	case KEY_WWAN:
		return "WWAN"
		//case KEY_WIMAX: return "WIMAX"
	case KEY_RFKILL:
		return "RFKILL"
	case KEY_MICMUTE:
		return "MICMUTE"

	case BTN_0:
		return "BTN_0"
	case BTN_1:
		return "BTN_1"
	case BTN_2:
		return "BTN_2"
	case BTN_3:
		return "BTN_3"
	case BTN_4:
		return "BTN_4"
	case BTN_5:
		return "BTN_5"
	case BTN_6:
		return "BTN_6"
	case BTN_7:
		return "BTN_7"
	case BTN_8:
		return "BTN_8"
	case BTN_9:
		return "BTN_9"
	//case BTN_MOUSE:
	//	return "BTN_MOUSE"
	case BTN_LEFT:
		return "BTN_LEFT"
	case BTN_RIGHT:
		return "BTN_RIGHT"
	case BTN_MIDDLE:
		return "BTN_MIDDLE"
	case BTN_SIDE:
		return "BTN_SIDE"
	case BTN_EXTRA:
		return "BTN_EXTRA"
	case BTN_FORWARD:
		return "BTN_FORWARD"
	case BTN_BACK:
		return "BTN_BACK"
	case BTN_TASK:
		return "BTN_TASK"
	case BTN_JOYSTICK:
		return "BTN_JOYSTICK"
	//case BTN_TRIGGER:
	//	return "BTN_TRIGGER"
	case BTN_THUMB:
		return "BTN_THUMB"
	case BTN_THUMB2:
		return "BTN_THUMB2"
	case BTN_TOP:
		return "BTN_TOP"
	case BTN_TOP2:
		return "BTN_TOP2"
	case BTN_PINKIE:
		return "BTN_PINKIE"
	case BTN_BASE:
		return "BTN_BASE"
	case BTN_BASE2:
		return "BTN_BASE2"
	case BTN_BASE3:
		return "BTN_BASE3"
	case BTN_BASE4:
		return "BTN_BASE4"
	case BTN_BASE5:
		return "BTN_BASE5"
	case BTN_BASE6:
		return "BTN_BASE6"
	case BTN_DEAD:
		return "BTN_DEAD"
	//case BTN_GAMEPAD:
	//	return "BTN_GAMEPAD"
	//case BTN_SOUTH:
	//	return "BTN_SOUTH"
	case BTN_A:
		return "BTN_A"
	//case BTN_EAST:
	//	return "BTN_EAST"
	case BTN_B:
		return "BTN_B"
	case BTN_C:
		return "BTN_C"
	//case BTN_NORTH:
	//	return "BTN_NORTH"
	case BTN_X:
		return "BTN_X"
	//case BTN_WEST:
	//	return "BTN_WEST"
	case BTN_Y:
		return "BTN_Y"
	case BTN_Z:
		return "BTN_Z"
	case BTN_TL:
		return "BTN_TL"
	case BTN_TR:
		return "BTN_TR"
	case BTN_TL2:
		return "BTN_TL2"
	case BTN_TR2:
		return "BTN_TR2"
	case BTN_SELECT:
		return "BTN_SELECT"
	case BTN_START:
		return "BTN_START"
	case BTN_MODE:
		return "BTN_MODE"
	case BTN_THUMBL:
		return "BTN_THUMBL"
	case BTN_THUMBR:
		return "BTN_THUMBR"
	//case BTN_DIGI:
	//	return "BTN_DIGI"
	case BTN_TOOL_PEN:
		return "BTN_TOOL_PEN"
	case BTN_TOOL_RUBBER:
		return "BTN_TOOL_RUBBER"
	case BTN_TOOL_BRUSH:
		return "BTN_TOOL_BRUSH"
	case BTN_TOOL_PENCIL:
		return "BTN_TOOL_PENCIL"
	case BTN_TOOL_AIRBRUSH:
		return "BTN_TOOL_AIRBRUSH"
	case BTN_TOOL_FINGER:
		return "BTN_TOOL_FINGER"
	case BTN_TOOL_MOUSE:
		return "BTN_TOOL_MOUSE"
	case BTN_TOOL_LENS:
		return "BTN_TOOL_LENS"
	case BTN_TOOL_QUINTTAP:
		return "BTN_TOOL_QUINTTAP"
	case BTN_STYLUS3:
		return "BTN_STYLUS3"
	case BTN_TOUCH:
		return "BTN_TOUCH"
	case BTN_STYLUS:
		return "BTN_STYLUS"
	case BTN_STYLUS2:
		return "BTN_STYLUS2"
	case BTN_TOOL_DOUBLETAP:
		return "BTN_TOOL_DOUBLETAP"
	case BTN_TOOL_TRIPLETAP:
		return "BTN_TOOL_TRIPLETAP"
	case BTN_TOOL_QUADTAP:
		return "BTN_TOOL_QUADTAP"
	//case BTN_WHEEL:
	//	return "BTN_WHEEL"
	case BTN_GEAR_DOWN:
		return "BTN_GEAR_DOWN"
	case BTN_GEAR_UP:
		return "BTN_GEAR_UP"
	default:
		return fmt.Sprintf("EventCode(unknown: 0x%X)", uint16(ec))
	}
}

const (
	/*
	 * Relative axes
	 */

	REL_X      EventCode = 0x00
	REL_Y      EventCode = 0x01
	REL_Z      EventCode = 0x02
	REL_RX     EventCode = 0x03
	REL_RY     EventCode = 0x04
	REL_RZ     EventCode = 0x05
	REL_HWHEEL EventCode = 0x06
	REL_DIAL   EventCode = 0x07
	REL_WHEEL  EventCode = 0x08
	REL_MISC   EventCode = 0x09
	/*
	 * 0x0a is reserved and should not be used in input drivers.
	 * It was used by HID as REL_MISC+1 and userspace needs to detect if
	 * the next REL_* event is correct or is just REL_MISC + n.
	 * We define here REL_RESERVED so userspace can rely on it and detect
	 * the situation described above.
	 */
	REL_RESERVED      EventCode = 0x0a
	REL_WHEEL_HI_RES  EventCode = 0x0b
	REL_HWHEEL_HI_RES EventCode = 0x0c
	REL_MAX           EventCode = 0x0f
	REL_CNT           EventCode = (REL_MAX + 1)

	/*
	 * Absolute axes
	 */

	ABS_X          EventCode = 0x00
	ABS_Y          EventCode = 0x01
	ABS_Z          EventCode = 0x02
	ABS_RX         EventCode = 0x03
	ABS_RY         EventCode = 0x04
	ABS_RZ         EventCode = 0x05
	ABS_THROTTLE   EventCode = 0x06
	ABS_RUDDER     EventCode = 0x07
	ABS_WHEEL      EventCode = 0x08
	ABS_GAS        EventCode = 0x09
	ABS_BRAKE      EventCode = 0x0a
	ABS_HAT0X      EventCode = 0x10
	ABS_HAT0Y      EventCode = 0x11
	ABS_HAT1X      EventCode = 0x12
	ABS_HAT1Y      EventCode = 0x13
	ABS_HAT2X      EventCode = 0x14
	ABS_HAT2Y      EventCode = 0x15
	ABS_HAT3X      EventCode = 0x16
	ABS_HAT3Y      EventCode = 0x17
	ABS_PRESSURE   EventCode = 0x18
	ABS_DISTANCE   EventCode = 0x19
	ABS_TILT_X     EventCode = 0x1a
	ABS_TILT_Y     EventCode = 0x1b
	ABS_TOOL_WIDTH EventCode = 0x1c

	ABS_VOLUME  EventCode = 0x20
	ABS_PROFILE EventCode = 0x21

	ABS_MISC EventCode = 0x28

	/*
	 * 0x2e is reserved and should not be used in input drivers.
	 * It was used by HID as ABS_MISC+6 and userspace needs to detect if
	 * the next ABS_* event is correct or is just ABS_MISC + n.
	 * We define here ABS_RESERVED so userspace can rely on it and detect
	 * the situation described above.
	 */
	ABS_RESERVED EventCode = 0x2e

	ABS_MT_SLOT        EventCode = 0x2f /* MT slot being modified */
	ABS_MT_TOUCH_MAJOR EventCode = 0x30 /* Major axis of touching ellipse */
	ABS_MT_TOUCH_MINOR EventCode = 0x31 /* Minor axis (omit if circular) */
	ABS_MT_WIDTH_MAJOR EventCode = 0x32 /* Major axis of approaching ellipse */
	ABS_MT_WIDTH_MINOR EventCode = 0x33 /* Minor axis (omit if circular) */
	ABS_MT_ORIENTATION EventCode = 0x34 /* Ellipse orientation */
	ABS_MT_POSITION_X  EventCode = 0x35 /* Center X touch position */
	ABS_MT_POSITION_Y  EventCode = 0x36 /* Center Y touch position */
	ABS_MT_TOOL_TYPE   EventCode = 0x37 /* Type of touching device */
	ABS_MT_BLOB_ID     EventCode = 0x38 /* Group a set of packets as a blob */
	ABS_MT_TRACKING_ID EventCode = 0x39 /* Unique ID of initiated contact */
	ABS_MT_PRESSURE    EventCode = 0x3a /* Pressure on contact area */
	ABS_MT_DISTANCE    EventCode = 0x3b /* Contact hover distance */
	ABS_MT_TOOL_X      EventCode = 0x3c /* Center X tool position */
	ABS_MT_TOOL_Y      EventCode = 0x3d /* Center Y tool position */

	ABS_MAX EventCode = 0x3f
	ABS_CNT EventCode = (ABS_MAX + 1)
)
