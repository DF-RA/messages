package services

import (
	"messages/core/entity"
	"sync"
)

type MessageServiceImpl struct{}
type MessageService interface {
	ProcessMessage(
		messages []entity.Message,
		call func(entity.Message) entity.MessageOutput,
	) []entity.MessageOutput
}

func (impl MessageServiceImpl) ProcessMessage(
	messages []entity.Message,
	call func(entity.Message) entity.MessageOutput,
) []entity.MessageOutput {
	var wg sync.WaitGroup
	results := make([]entity.MessageOutput, len(messages))

	for i, message := range messages {
		wg.Add(1)
		go func(i int, message entity.Message) {
			defer wg.Done()
			results[i] = call(message)
		}(i, message)
	}
	wg.Wait()
	return results
}
