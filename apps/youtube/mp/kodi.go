package mp

import (
	"os"
	"errors"
	"sync"
	"time"

	"github.com/sargo/kodicast/log"
	"github.com/pdf/kodirpc"
	"github.com/Sirupsen/logrus"
)

var KODI_PROPERTY_UNAVAILABLE = errors.New("kodi: property unavailable")

// Kodi is an implementation of Backend.
type Kodi struct {
	client       *kodirpc.Client
	running      bool
	runningMutex sync.Mutex
}

var kodiLogger = log.New("kodi", "log Kodi wrapper output")

func (kodi *Kodi) initialize() chan State {
	if kodi.running {
		panic("already initialized")
	}
	kodiLogger.Println("connecting")
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	kodirpc.SetLogger(logger)

	config := kodirpc.NewConfig()
	config.ReadTimeout = 20 * time.Second
	client, err := kodirpc.NewClient("127.0.0.1:9090", config)
	if err != nil {
		panic(err)
	}
	kodiLogger.Println("connected")
	kodi.client = client

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
	kodiLogger.Println("initialized")

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
func (kodi *Kodi) sendCommand(command string, params map[string]interface{}) (map[string]interface{}) {
	kodiLogger.Println(command)
	kodiLogger.Println(params)
	resp, err := kodi.client.Call(command, params)
	if err != nil {
		kodiLogger.Fatal(err)
	}
	result := resp.(map[string]interface{})
	kodiLogger.Println(result)
	return result
}

func (kodi *Kodi) play(stream string, position time.Duration, volume int) {
	params := map[string]interface{}{
		"item": map[string]string{
			"file": "plugin://plugin.video.youtube/?action=play_video&videoid="+stream,
		},
	}
	result := kodi.sendCommand("Player.Open", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) pause() {
	params := map[string]interface{}{
		"playerid": 1,
	}
	result := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) resume() {
	params := map[string]interface{}{
		"playerid": 1,
	}
	result := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) getPosition() (time.Duration, error) {
	params := map[string]interface{}{
		"playerid": 1,
		"properties": [1]string{"time"},
	}
	result := kodi.sendCommand("Player.GetProperties", params)

	timeData := result["time"].(map[string]int64)
	hours := timeData["hours"]
	minutes := timeData["minutes"]
	seconds := timeData["seconds"]

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	position := time.Duration(hours*hour + minutes*minute + seconds*second)

	return position, nil
}

func (kodi *Kodi) setPosition(position time.Duration) {
	params := map[string]interface{}{
		"playerid": 1,
		"value": map[string]int64{
			"hours": int64(position.Hours()),
			"milliseconds": 0,
			"minutes": int64(position.Minutes()),
			"seconds": int64(position.Seconds()),
		},
	}
	result := kodi.sendCommand("Player.Seek", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) getVolume() int {
	params := map[string]interface{}{
		"properties": [1]string{"volume"},
	}
	result := kodi.sendCommand("Application.GetProperties", params)
	volume := result["volume"].(int)
	return int(volume)
}

func (kodi *Kodi) setVolume(volume int) {
	params := map[string]interface{}{
		"volume": volume,
	}
	result := kodi.sendCommand("Application.SetVolume", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) stop() {
	params := map[string]interface{}{
		"playerid": 1,
	}
	result := kodi.sendCommand("Player.Stop", params)
	kodiLogger.Println(result)
}
