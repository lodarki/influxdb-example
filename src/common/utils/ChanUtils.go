package utils

import (
	"github.com/astaxie/beego/logs"
	"time"
)

func WatchChan(c *chan interface{}, f func(interface{}) error) {
	idleDuration := 1 * time.Minute
	idleDelay := time.NewTimer(idleDuration)
	defer idleDelay.Stop()
	for {
		select {
		case t := <-*c:
			err := f(t)
			if err != nil {
				logs.Error(err)
			}
		case <-idleDelay.C:
		}
	}
}
