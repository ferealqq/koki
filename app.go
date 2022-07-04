package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"time"

	db "github.com/ferealqq/koki/pkg/database"
	ke "github.com/ferealqq/koki/pkg/keyevents"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type KeyData struct{
	Char string
	Count int
	Value uint16 
	ValueName string
}


func (a *App) GetKeyEventData(c string) []KeyData{
	var rw2name = make(map[uint16]string)
	
	for k,c := range ke.ValueMap {
		rw2name[c] = k
	}

	var res []KeyData


	db.Conn().Table(db.KeyTable).
		Select("char", "value", "count(value) as count").
		Where("char = ?",c).
		Group("value, char").
		Order("value").Scan(&res)
		

	for i,v := range res {
		res[i].ValueName = rw2name[v.Value]
	}

	return res
}

type CharData struct {
	Char string
	Count int 
} 

func (a *App) GetMostPressedKey() CharData{
	var m []CharData

	db.Conn().Table(db.KeyTable).
		Select("char", "count(*) as count").
		Group("char").
		Order("count desc").
		Limit(1).Scan(&m);

	if len(m) > 0 {
		return m[0]
	}

	return CharData{}
}

func (a *App) IsLoggerActive() bool {
	cmd := exec.Command("sudo", "systemctl", "status", "koki-logger")

	if out, err := cmd.Output(); err != nil {
		log.Println(err)
		return false
	}else{
		if matched, err := regexp.MatchString("Active: active", string(out)); err == nil && matched {
			return true
		}
	}
	return false
} 

func (a *App) ToggleLoggerDaemon() (bool,error){
	toggle := "start"
	if a.IsLoggerActive() {
		toggle = "stop"
	}
	cmd := exec.Command("sudo", "systemctl", toggle, "koki-logger")

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return false, nil
	}

	return true, nil
}

type HourEvents struct {
	Hour int
	Count int
}

func (a *App) GetKeysPressedIn(hours int) []HourEvents {
	n := time.Now().UnixNano()
	n -= time.Hour.Nanoseconds() * int64(hours)

	var events []HourEvents

	fmt.Println(n)

	db.Conn().Table(db.KeyTable).
		Select("count(*) as count, strftime ('%H', created_at) hour").
		Where("time > ?", n).
		Where("type = 0").
		Group("strftime ('%H',created_at)").
		Order("hour").
		Scan(&events);

	return events
}