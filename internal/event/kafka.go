package event

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

type KafkaProducer struct {
	brokers string
	topic   string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		brokers: strings.Join(brokers, ","),
		topic:   topic,
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, key uuid.UUID, eventType string, payload []byte) error {
	message := fmt.Sprintf("%s|%s|%s", key.String(), eventType, string(payload))
	cmd := exec.CommandContext(ctx, "kafka-console-producer", "--broker-list", p.brokers, "--topic", p.topic, "--property", "parse.key=true", "--property", "key.separator=|")
	cmd.Stdin = strings.NewReader(message)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("publish kafka message failed: %w: %s", err, string(out))
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return nil
}
