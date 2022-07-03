package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
// normal keys
var nKeys = map[string]uint16{
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
	KEYPRESS = "KEYPRESS"
	KEYRELEASE = "KEYRELEASE"
	AUTOREPEAST = "AUTOREPEAT"
	UNKNOWN = "UNKNOWN"
)

// event value correlating human readable translation, Source: https://www.kernel.org/doc/Documentation/input/input.txt
var valueMap = map[string]uint16{
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

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

var (
	list []*KeyEvent
	m sync.RWMutex
	db *gorm.DB
) 

var quitCh = make(chan struct{})

func saveKeyEvents(){
	tic := time.NewTicker(10 * time.Second)
	loop:
		for {
			select {
			// Stop the execution of this goroutine when the script has loaded
			case <-quitCh:
				break loop
			case <-tic.C:
				n := time.Now()
				m.Lock()

				db.CreateInBatches(list, 100)

				fmt.Printf("Written %d events to database\n",len(list))

				list = *new([]*KeyEvent)


				m.Unlock()
				timeTrack(n,"time spent writing to database")
			}
		}
}

func main() {
	d, err := gorm.Open(sqlite.Open("log.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}else{
		db = d 
		fmt.Println("Connected to database")
	}
	db.AutoMigrate(&KeyEvent{})

	var r2key = make(map[uint16]string)
	
	for k,c := range nKeys {
		r2key[c] = k
	}
	// save changes to database
	go saveKeyEvents()


	// change according to your keyboards event file
	f, err := os.Open("/dev/input/event2")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := make([]byte, 24)
	for {
		f.Read(b)
		n := time.Now()
		sec := binary.LittleEndian.Uint64(b[0:8])
		usec := binary.LittleEndian.Uint64(b[8:16])
		t := time.Unix(int64(sec), int64(usec))
		var value int32
		typ := binary.LittleEndian.Uint16(b[16:18])
		// co := binary.LittleEndian.Uint64(b[18:20])
		code := binary.LittleEndian.Uint16(b[18:20])
		binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)
		
		c, ok := r2key[code]
		if ok {
			fmt.Println(t)
			evt := &KeyEvent{
				Time: t.UnixNano(),
				Char: c,
				Value: value,
				Code: code,
				Type: typ,
			}
			if value == 1 {
				fmt.Printf("KEYPRESSED\ntype: %x\ncode: %d\nvalue: %d\nchar: %s\n", typ, code, value,c)
				if m.TryLock() {
					list = append(list,evt)
					m.Unlock()				
					timeTrack(n, "Keypress")
				}else{
					continue
				} 
			}else {
				if value == 0 {
					fmt.Printf("KEYRELEASED\ntype: %x\ncode: %d\nvalue: %d\nchar: %s \n", typ, code, value,c)
					if m.TryLock() {
						list = append(list,evt)
						m.Unlock()				
						timeTrack(n, "Keyreleased")
					}else{
						continue
					} 	
				}else if value == 2{
					if m.TryLock() {
						list = append(list,evt)
						m.Unlock()				
						timeTrack(n, "Held")
					}else{
						continue
					} 
					fmt.Printf("HELD\ntype: %x\ncode: %d\nvalue: %d\nchar: %s \n", typ, code, value,c)
				}else{ 
					// FIXME: not sure if we want to track weird keyevents?
					fmt.Printf("Something weird\ntype: %x\ncode: %d\nvalue: %d\nchar: %s \n", typ, code, value,c)
				}
				timeTrack(n, "Others")
			}
		}
	}
}