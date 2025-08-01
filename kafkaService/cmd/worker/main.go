package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	// 1) Настройки из окружения
	brokers := []string{os.Getenv("KAFKA_ADDR")} // e.g. "kafka:9092"
	topic := os.Getenv("KAFKA_TOPIC")            // e.g. "events"
	group := os.Getenv("KAFKA_GROUP")            // e.g. "logger"

	// 2) Консюмер-группа
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokers, group, cfg)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}
	defer consumer.Close()

	// 3) Открываем файл для логов
	f, err := os.OpenFile("/app/logs/events.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	defer writer.Flush()

	// 4) Трэйс сигналы для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	// 5) Запускаем цикл Consume
	handler := consumerGroupHandler{writer: writer}
	for {
		if err := consumer.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("consume error: %v", err)
		}
		if ctx.Err() != nil {
			break
		}
	}
}

// consumerGroupHandler реализует интерфейс sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	writer *bufio.Writer
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Записываем в буфер
		if _, err := h.writer.Write(msg.Value); err != nil {
			log.Printf("failed to write message: %v", err)
		}
		if err := h.writer.WriteByte('\n'); err != nil {
			log.Printf("failed to write newline: %v", err)
		}
		// Сбрасываем буфер сразу, чтобы данные появились в файле
		if err := h.writer.Flush(); err != nil {
			log.Printf("failed to flush buffer: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
