package consumer

import (
	"context"
)

type Consumer interface {
	// Consume закрывает канал, как только вызывается метод Close. С гарантией что
	// все сообщения которые мы успели получить из кафки, будут переданы.
	Consume() <-chan []byte
	Close(context.Context) error
}

// type Consumer interface {
// 	Consume() <-chan []Message
// }

// type Message struct {
// 	CurrencyID string
//  OldPrice float64
//  CurrentPrice float64
// 	Timestamp time.Time `json:"timestamp"`
// }

func New() (Consumer, error) {
	// Подключение к кафке
}

type kafkaConsumer struct{}

func (c *kafkaConsumer) Consume() <-chan []byte {
	// Подключается и начинает забирать сообщения и перекладыыать в выходной канал.
}
