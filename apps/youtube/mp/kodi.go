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
		kodiLogger.Println("OnPause")
		eventChan <- STATE_PAUSED
	})
	kodi.client.Handle("Player.OnPlay", func(method string, data interface{}) {
		kodiLogger.Println("OnPlay")
		eventChan <- STATE_PLAYING
	})
	kodi.client.Handle("Player.OnStop", func(method string, data interface{}) {
		kodiLogger.Println("OnStop")
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
func (kodi *Kodi) sendCommand(command string, params map[string]interface{}) (interface{}, error) {
	kodiLogger.Println(command)
	kodiLogger.Println(params)
	resp, err := kodi.client.Call(command, params)
	if err != nil {
		kodiLogger.Println(err)
	}
	kodiLogger.Println(resp)
	return resp, err
}

func (kodi *Kodi) play(stream string, position time.Duration, volume int) {
	params := map[string]interface{}{
		"item": map[string]string{
			"file": "plugin://plugin.video.youtube/?action=play_video&videoid="+stream,
		},
	}
	resp, _ := kodi.sendCommand("Player.Open", params)
	result := resp.(string)
	kodiLogger.Println(result)
}

func (kodi *Kodi) getPlayerId() (int) {
	params := map[string]interface{}{
	}
	resp, err := kodi.sendCommand("Player.GetActivePlayers", params)
	if err != nil {
		return 0
	}

	result := resp.([]interface{})
	for _, i := range result {
		item := i.(map[string]interface{})
		playerType := item["type"].(string)
		if playerType == "video" {
			return int(item["playerid"].(float64))
		}
	}

	return -1
}

func (kodi *Kodi) pause() {
	params := map[string]interface{}{
		"playerid": kodi.getPlayerId(),
	}
	result, _ := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) resume() {
	params := map[string]interface{}{
		"playerid": kodi.getPlayerId(),
	}
	result, _ := kodi.sendCommand("Player.PlayPause", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) getPosition() (time.Duration) {
	params := map[string]interface{}{
		"playerid": kodi.getPlayerId(),
		"properties": [1]string{"time"},
	}
	resp, err := kodi.sendCommand("Player.GetProperties", params)
	if err != nil {
		return 0
	}

	result := resp.(map[string]interface{})
	timeData := result["time"].(map[string]interface{})
	hours := int64(timeData["hours"].(float64))
	minutes := int64(timeData["minutes"].(float64))
	seconds := int64(timeData["seconds"].(float64))

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	position := time.Duration(hours*hour + minutes*minute + seconds*second)

	return position
}

func (kodi *Kodi) setPosition(position time.Duration) {
	params := map[string]interface{}{
		"playerid": kodi.getPlayerId(),
		"value": map[string]int64{
			"hours": int64(position.Hours()) % 24,
			"minutes": int64(position.Minutes()) % 60,
			"seconds": int64(position.Seconds()) % 60,
		},
	}
	result, _ := kodi.sendCommand("Player.Seek", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) setVolume(volume int) {
	params := map[string]interface{}{
		"volume": volume,
	}
	result, _ := kodi.sendCommand("Application.SetVolume", params)
	kodiLogger.Println(result)
}

func (kodi *Kodi) stop() {
	params := map[string]interface{}{
		"playerid": kodi.getPlayerId(),
	}
	result, _ := kodi.sendCommand("Player.Stop", params)
	kodiLogger.Println(result)
}
