package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"io/ioutil"

	ke "github.com/ferealqq/koki/pkg/keyevents"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

var (
	list []*ke.KeyEvent
	m sync.RWMutex
	_db *gorm.DB
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
				// set char type values to every item in a list
				for i, v := range list {
					if _, ok := ke.CKeys[v.Char]; ok {
						list[i].CharType = string(ke.CHAR)
					}
					if _, ok := ke.SKeys[v.Char]; ok {
						list[i].CharType = string(ke.SPECIAL)
					}
					if _, ok := ke.FKeys[v.Char]; ok {
						list[i].CharType = string(ke.FUNCTION)
					}
				}
				_db.CreateInBatches(list, 100)
				fmt.Printf("Written %d events to database\n",len(list))

				list = *new([]*ke.KeyEvent)


				m.Unlock()
				timeTrack(n,"time spent writing to database")
			}
		}
}

const INPUT_BY_ID_DIR = "/dev/input/by-id"
const KEYBOARD_INPUT_EXTENSION = "-event-kbd"

func getKeyboardInputPath(id string) (string, error) {
	files, e := ioutil.ReadDir(INPUT_BY_ID_DIR)
	if e != nil {
		return "", e
	}

	var inputPath string
	for _, file := range files { 
		if strings.Contains(file.Name(), id+KEYBOARD_INPUT_EXTENSION) {
			inputPath = file.Name()
		}
		fmt.Printf("file name %s\n",file.Name())
	}

	return INPUT_BY_ID_DIR+"/"+inputPath, nil 
}

func main() {
	// Remember to use this procces own db variable and not the shared one.
	d, err := gorm.Open(sqlite.Open("log.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}else{
		_db = d 
		fmt.Println("Connected to database")
	}
	_db.AutoMigrate(&ke.KeyEvent{})

	var r2key = make(map[uint16]string)
	
	for k,c := range ke.NKeys() {
		r2key[c] = k
	}
	// save changes to database
	go saveKeyEvents()

	inputPath, err := getKeyboardInputPath("usb-ZSA_Technology_Labs_Moonlander_Mark_I")
	if err != nil { 
		panic(err)
	}
	// change according to your keyboards event file
	f, err := os.Open(inputPath)
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
		code := binary.LittleEndian.Uint16(b[18:20])
		binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)
		
		c, ok := r2key[code]
		if ok {
			fmt.Println(t)
			evt := &ke.KeyEvent{
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