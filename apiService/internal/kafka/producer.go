package kafka

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	_ "go.uber.org/zap"
)

type Producer struct {
	Async sarama.AsyncProducer
	Topic string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true //нахуя?
	//говорит кафке что если произошла ошибка продюсера то есть апи
	//надо уведомить
	cfg.Producer.RequiredAcks = sarama.WaitForLocal //нахуя?
	//апиха не будем выполнять код дальше, пока не
	//Дождеться ответа, что один из броуков
	//получил наше событие
	async, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{async, topic}, nil
}

func (p *Producer) Send(evt interface{}) {
	buf, _ := json.Marshal(evt)
	msg := &sarama.ProducerMessage{
		Topic:     p.Topic,
		Value:     sarama.ByteEncoder(buf),
		Timestamp: time.Now(),
	}
	p.Async.Input() <- msg
}

func (p *Producer) Close() error { return p.Async.Close() }
