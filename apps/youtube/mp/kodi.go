package mp

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sargo/kodicast/config"
	"github.com/sargo/kodicast/log"
)

var KODI_PROPERTY_UNAVAILABLE = errors.New("kodi: property unavailable")

// Kodi is an implementation of Backend.
type Kodi struct {
	running      bool
	runningMutex sync.Mutex
	mainloopExit chan struct{}
}

var kodiLogger = log.New("kodi", "log Kodi wrapper output")

func (kodi *Kodi) initialize() {
	kodi.running = true
}

// Function quit quits the player.
// WARNING: This MUST be the last call on this media player.
func (kodi *Kodi) quit() {
	kodi.runningMutex.Lock()
	if !kodi.running {
		panic("quit called twice")
	}
	kodi.running = false
	kodi.runningMutex.Unlock()
}

// sendCommand sends a command to the Kodi player
func (kodi *Kodi) sendCommand(command string, value string) (string, error) {
	return "", nil
}

func (kodi *Kodi) play(stream string, position time.Duration, volume int) {
	reply, err := kodi.sendCommand("play", stream)
}

func (kodi *Kodi) pause() {
	reply, err := kodi.sendCommand("pause", "yes")
}

func (kodi *Kodi) resume() {
	reply, err := kodi.sendCommand("pause", "no")
}

func (kodi *Kodi) getPosition() (time.Duration, error) {
	position, err := kodi.sendCommand("get-position", "")

	if position < 0 {
		// Sometimes, the position appears to be slightly off.
		position = 0
	}

	return time.Duration(position * float64(time.Second)), nil
}

func (kodi *Kodi) setPosition(position time.Duration) {
	reply, err := kodi.sendCommand("set-position", position)
}

func (kodi *Kodi) getVolume() int {
	volume, err := kodi.sendCommand("get-volume")
	if err != nil {
		// should not happen
		panic(err)
	}

	return int(volume + 0.5)
}

func (kodi *Kodi) setVolume(volume int) {
	reply, err := kodi.sendCommand("set-volume", volume)
}

func (kodi *Kodi) stop() {
	kodi.sendCommand("stop", "")
}
