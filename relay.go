package main

import (
	"strconv"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// Relay models the circuit used to control the turnout switch
type Relay struct {
	relayDriver *gpio.RelayDriver
}

// NewRelay creates a new turnout switch
func NewRelay(rpi *raspi.Adaptor, onOffPin int) *Relay {
	relay := &Relay{
		relayDriver: gpio.NewRelayDriver(rpi, strconv.Itoa(onOffPin)),
	}
	return relay
}

// Status returns true if the relay is on
func (relay *Relay) Status() bool {
	// TODO
	return false
}

// On turns the relay on
func (relay *Relay) On() error {
	if err := relay.relayDriver.On(); err != nil {
		return err
	}
	return nil
}

// Off turns the relay off
func (relay *Relay) Off() error {
	if err := relay.relayDriver.Off(); err != nil {
		return err
	}
	return nil
}
