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


func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

var (
	list []*KeyEvent
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

				_db.CreateInBatches(list, 100)

				fmt.Printf("Written %d events to database\n",len(list))

				list = *new([]*KeyEvent)


				m.Unlock()
				timeTrack(n,"time spent writing to database")
			}
		}
}

func logger() {
	// Remember to use this procces own db variable and not the shared one.
	d, err := gorm.Open(sqlite.Open("log.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}else{
		_db = d 
		fmt.Println("Connected to database")
	}
	_db.AutoMigrate(&KeyEvent{})

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