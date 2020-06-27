package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	"gobot.io/x/gobot/platforms/raspi"
)

// Options is program options
type Options struct {
	Port        int
	RelayConfig map[string]RelayConfiguration `yaml:"relays"`
}

// RelayConfiguration defines the gpio pin configuration
type RelayConfiguration struct {
	OnOffPin int `yaml:"on_off_pin"`
}

// RelayControl is the container for the HTTP server
type RelayControl struct {
	relays map[string]*Relay
}

// TurnOn turns the relay with :id on
func (rc *RelayControl) TurnOn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	relay, ok := rc.relays[vars["id"]]
	if ok {
		if err := relay.On(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error turning relay on: %s", err.Error())
			return
		}
		fmt.Fprint(w, "1")
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// TurnOff turns the relay with :id off
func (rc *RelayControl) TurnOff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	relay, ok := rc.relays[vars["id"]]
	if ok {
		if err := relay.Off(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error turning relay off: %s", err.Error())
			return
		}
		fmt.Fprint(w, "0")
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// RelayStatus returns the given id's status
func (rc *RelayControl) RelayStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	relay, ok := rc.relays[vars["id"]]
	if ok {
		if relay.Status() {
			fmt.Fprint(w, "1")
		} else {
			fmt.Fprint(w, "0")
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	opts := Options{
		Port: 8080,
	}

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s /path/to/config", os.Args[0])
	}

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	if err := yaml.Unmarshal(b, &opts); err != nil {
		log.Fatalf("error parsing config file: %s", err.Error())
	}

	rpi := raspi.NewAdaptor()
	if err := rpi.Connect(); err != nil {
		log.Fatalf("error connecting to GPIO: %s", err.Error())
	}

	rtr := mux.NewRouter()

	relays := make(map[string]*Relay)
	for key, item := range opts.RelayConfig {
		relays[key] = NewRelay(rpi, item.OnOffPin)
	}

	rc := &RelayControl{
		relays: relays,
	}

	rtr.HandleFunc("/relay/{id}/status", rc.RelayStatus).Methods("GET")
	rtr.HandleFunc("/relay/{id}/on", rc.TurnOn).Methods("GET")
	rtr.HandleFunc("/relay/{id}/off", rc.TurnOff).Methods("GET")

	srv := &http.Server{
		Handler: rtr,
		Addr:    fmt.Sprintf(":%d", opts.Port),
	}

	log.Printf("Listening on :%d...\n", opts.Port)
	log.Fatal(srv.ListenAndServe())
}
