package main

import (
	"context"
	"fmt"
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
	
	for k,c := range valueMap {
		rw2name[c] = k
	}

	var res []KeyData


	DB().Table("key_events").
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

	DB().Table("key_events").
		Select("char", "count(*) as count").
		Group("char").
		Order("count desc").
		Limit(1).Scan(&m);

	if len(m) > 0 {
		return m[0]
	}

	return CharData{}
}
