package mp

import (
	"errors"
	"strconv"
	"sync"
	"time"

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
	resp, err := kodi.sendCommand("play", stream)
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}

func (kodi *Kodi) pause() {
	resp, err := kodi.sendCommand("pause", "yes")
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}

func (kodi *Kodi) resume() {
	resp, err := kodi.sendCommand("pause", "no")
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}

func (kodi *Kodi) getPosition() (time.Duration, error) {
	resp, err := kodi.sendCommand("get-position", "")
	if err != nil {
		kodiLogger.Fatal(err)
	}
	position, err := strconv.ParseFloat(resp, 64)

	if position < 0 {
		// Sometimes, the position appears to be slightly off.
		position = 0
	}

	return time.Duration(position * float64(time.Second)), nil
}

func (kodi *Kodi) setPosition(position time.Duration) {
	resp, err := kodi.sendCommand("set-position", position.String())
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}

func (kodi *Kodi) getVolume() int {
	resp, err := kodi.sendCommand("get-volume", "")
	if err != nil {
		kodiLogger.Fatal(err)
	}
	volume, err := strconv.ParseFloat(resp, 64)
	if err != nil {
		kodiLogger.Fatal(err)
	}

	return int(volume + 0.5)
}

func (kodi *Kodi) setVolume(volume int) {
	resp, err := kodi.sendCommand("set-volume", strconv.Itoa(volume))
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}

func (kodi *Kodi) stop() {
	resp, err := kodi.sendCommand("stop", "")
	if err != nil {
		kodiLogger.Fatal(err)
	}
	kodiLogger.Println(resp)
}
