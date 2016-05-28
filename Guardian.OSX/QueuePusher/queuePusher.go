package QueuePusher

import (
	"fmt"
	"time"
)

type QueuePusher struct{}

func (queuePusher QueuePusher) Start() {
	fmt.Printf("QueuePusher initialized at %s", time.Now())
}
