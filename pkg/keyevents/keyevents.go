package keyevents

import "gorm.io/gorm"

//https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
// normal keys
var NKeys = map[string]uint16{
	"`": 41,
	"1": 2,
	"2": 3,
	"3": 4,
	"4": 5,
	"5": 6,
	"6": 7,
	"7": 8,
	"8": 9,
	"9": 10,
	"0": 11,
	"-": 12,
	"+": 13,
	//
	"q":  16,
	"w":  17,
	"e":  18,
	"r":  19,
	"t":  20,
	"y":  21,
	"u":  22,
	"i":  23,
	"o":  24,
	"p":  25,
	"[":  26,
	"]":  27,
	"\\": 43,
	//
	"a": 30,
	"s": 31,
	"d": 32,
	"f": 33,
	"g": 34,
	"h": 35,
	"j": 36,
	"k": 37,
	"l": 38,
	";": 39,
	"'": 40,
	//
	"z": 44,
	"x": 45,
	"c": 46,
	"v": 47,
	"b": 48,
	"n": 49,
	"m": 50,
	",": 51,
	".": 52,
	"/": 53,
	//
	"f1":  59,
	"f2":  60,
	"f3":  61,
	"f4":  62,
	"f5":  63,
	"f6":  64,
	"f7":  65,
	"f8":  66,
	"f9":  67,
	"f10": 68,
	"f11": 69,
	"f12": 70,
	// more
	"esc":     1,
	"delete":  14,
	"tab":     15,
	"ctrl":    29,
	"control": 29,
	"alt":     56,
	"space":   57,
	"shift":   42,
	"rshift":  54,
	"enter":   28,
	"cmd":     3675,
	"command": 3675,
	"rcmd":    3676,
	"ralt":    3640,
	"up":      57416,
	"down":    57424,
	"left":    57419,
	"right":   57421,
}

// moonlander weird keys
var mKeys = map[string]uint16{
	"backspace": 14,
	"delele": 123,

}

var (
	KEYPRESS = "Keypress"
	KEYRELEASE = "Keyreleased"
	AUTOREPEAST = "Autorepeat"
	UNKNOWN = "Unknown"
)

// event value correlating human readable translation, Source: https://www.kernel.org/doc/Documentation/input/input.txt
var ValueMap = map[string]uint16{
	KEYPRESS: 0, 
	KEYRELEASE: 1,
	AUTOREPEAST: 2,
	// everything else is unknown for example there is no documentation what does the value 4 corralate to 
	UNKNOWN: 3,	
}

//https://www.kernel.org/doc/Documentation/input/input.txt
type KeyEvent struct {
	gorm.Model
	Time int64
	Code uint16
	Type uint16
	Value int32

	Char string
}