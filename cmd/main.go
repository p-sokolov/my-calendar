package main

import (
	"net/http"
	"sync"
	"log"

	"github.com/gorilla/mux"
	
	"my-calendar/internal/config"
	"my-calendar/internal/calendar"
	"my-calendar/internal/handler"
)

func main() {

	// load config
	cfg := config.LoadCfg()

	var calendar = calendar.NewCalendar()

	h := handler.NewHandler(*calendar)

	r := mux.NewRouter()
	r.HandleFunc("/events_for_day", h.GetDaily).Methods("get")

	// srv := &http.Server{
	// 	Addr:    cfg.HttpPort,
	// 	Handler: r,
	// }

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("server started")
		if err := http.ListenAndServe(cfg.HttpPort, r); err != nil && err != http.ErrServerClosed {
			log.Fatal("http.ListenAndServe failed")
		}
	}()

	wg.Wait()	
}