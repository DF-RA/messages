package entity

import "time"

type ConsumerEvent struct {
	ID         string
	Source     string // "kafka" | "activemq"
	Name       string // topic or queue name
	Payload    string
	ReceivedAt time.Time
}
