package mp

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/sargo/kodicast/log"
	"github.com/pdf/kodirpc"
)

var KODI_PROPERTY_UNAVAILABLE = errors.New("kodi: property unavailable")

// Kodi is an implementation of Backend.
type Kodi struct {
	client       kodirpc.Client
	running      bool
	runningMutex sync.Mutex
}

var kodiLogger = log.New("kodi", "log Kodi wrapper output")

func (kodi *Kodi) initialize() chan State {
	if kodi.running {
		panic("already initialized")
	}

	kodi.client, err := kodirpc.NewClient("127.0.0.1:8080", kodirpc.NewConfig())
	if err != nil {
		panic(err)
	}

	eventChan := make(chan State)
	kodi.client.Handle("Player.OnPause", func(method string, data interface{}) {
		eventChan <- STATE_PAUSED
	})
	kodi.client.Handle("Player.OnPlay", func(method string, data interface{}) {
		eventChan <- STATE_PLAYING
	})
	kodi.client.Handle("Player.OnStop", func(method string, data interface{}) {
		eventChan <- STATE_STOPPED
	})

	kodi.running = true

	return eventChan
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
func (kodi *Kodi) sendCommand(command string, params map[string]interface{}) (string) {
	resp, err := client.Call(command, params)
	if err != nil {
		kodiLogger.Fatal(err)
	}
	return resp
}

func (kodi *Kodi) play(stream string, position time.Duration, volume int) {
	params := map[string]interface{}{
		"item": map[string]interface{}{
			"file": "plugin://plugin.video.youtube/?action=play_video&videoid="+stream
		}
	}
	resp := kodi.sendCommand("Player.Open", params)
	kodiLogger.Println(resp)
}

func (kodi *Kodi) pause() {
	params := map[string]interface{}{
		"playerid": 1
	}
	resp := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(resp)
}

func (kodi *Kodi) resume() {
	params := map[string]interface{}{
		"playerid": 1
	}
	resp := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(resp)
}

func (kodi *Kodi) getPosition() (time.Duration, error) {
	params := map[string]interface{}{
		"playerid": 1
		"properties": [1]string{"time"}
	}
	resp := kodi.sendCommand("Player.GetProperties", "")

	hours := ParseInt64(resp["result"]["time"]["hours"])
	minutes := ParseInt64(resp["result"]["time"]["minutes"])
	seconds := ParseInt64(resp["result"]["time"]["seconds"])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	position := time.Duration(hours*hour + minutes*minute + seconds*second)

	return position, nil
}

func (kodi *Kodi) setPosition(position time.Duration) {
	params := map[string]interface{}{
		"playerid": 1
		"value": map[string]interface{}{
			"hours": position.Hours()
			"milliseconds": position.Milliseconds()
			"minutes": position.Minutes()
			"seconds": position.Seconds()
		}
	}
	resp := kodi.sendCommand("Player.Seek", position.String())
	kodiLogger.Println(resp)
}

func (kodi *Kodi) getVolume() int {
	params := map[string]interface{}{
		"properties": [1]string{"volume"}
	}
	resp := kodi.sendCommand("Application.GetProperties", params)
	return int(resp["result"]["volume"])
}

func (kodi *Kodi) setVolume(volume int) {
	params := map[string]interface{}{
		"volume": volume
	}
	resp := kodi.sendCommand("Application.SetVolume", params)
	kodiLogger.Println(resp)
}

func (kodi *Kodi) stop() {
	params := map[string]interface{}{
		"playerid": 1
	}
	resp := kodi.sendCommand("Player.Stop", params)
	kodiLogger.Println(resp)
}
