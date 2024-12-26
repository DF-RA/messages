package entity

type ActiveMQ struct {
	Queues []Message
	Topics []Message
}

type ActiveMQOutput struct {
	Queues []MessageOutput
	Topics []MessageOutput
}
