package main

import (
	"log"

	"github.com/solafide-dev/gobible"
	"github.com/solafide-dev/gobible/bible"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) initBible() {
	a.Bible = gobible.NewGoBible()

	// Leverage august's ability to monitor filesystem changes
	a.Data.Register("bibles", bible.Bible{})
	rt.EventsOn(a.ctx, "data-mutate", func(d ...interface{}) {
		data := d[0].(DataMutationEvent)
		if data.DataType == "bibles" {
			switch data.Type {
			case "create":
				// Load from store
				store, err := a.Data.GetStore("bibles")
				if err != nil {
					rt.LogError(a.ctx, err.Error())
				}
				b, err := store.Get(data.Id)
				if err != nil {
					rt.LogError(a.ctx, err.Error())
				}
				bible := b.(bible.Bible)
				err = a.Bible.LoadObject(bible)
				if err != nil {
					log.Println(err.Error())
				}
			case "update":
				// Unload, then load from store
				a.Bible.Unload(data.Id)
				store, err := a.Data.GetStore("bibles")
				if err != nil {
					rt.LogError(a.ctx, err.Error())
				}
				b, err := store.Get(data.Id)
				if err != nil {
					rt.LogError(a.ctx, err.Error())
				}
				bible := b.(bible.Bible)
				err = a.Bible.LoadObject(bible)
				if err != nil {
					log.Println(err.Error())
				}
			case "delete":
				// Just unload
				a.Bible.Unload(data.Id)
			}
		}
	})

	// Leverage August for data monitoring, and register bibles
	bibleStore, err := a.Data.GetStore("bibles")
	if err != nil {
		rt.LogError(a.ctx, err.Error())
	}

	// Load our bibles
	existingBibles, err := bibleStore.GetAll()
	if err != nil {
		rt.LogError(a.ctx, err.Error())
	}

	for _, b := range existingBibles {
		bible := b.(bible.Bible)
		err := a.Bible.LoadObject(bible)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
