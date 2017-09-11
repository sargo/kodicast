package mp

import (
	"time"
)

type Backend interface {
	initialize()
	quit()
	play(string, time.Duration, int)
	pause()
	resume()
	getPosition() (time.Duration, error)
	setPosition(time.Duration)
	setVolume(int)
	stop()
}
