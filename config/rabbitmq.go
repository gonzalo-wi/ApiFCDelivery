package config

import (
	"fmt"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Queue    string
}

func LoadRabbitMQConfig() *RabbitMQConfig {
	port, _ := strconv.Atoi(getEnvOrDefault("RABBITMQ_PORT", "5672"))

	return &RabbitMQConfig{
		Host:     getEnvOrDefault("RABBITMQ_HOST", "192.168.0.250"),
		Port:     port,
		User:     getEnvOrDefault("RABBITMQ_USER", "guest"),
		Password: getEnvOrDefault("RABBITMQ_PASSWORD", "guest"),
		Queue:    getEnvOrDefault("RABBITMQ_QUEUE", "q.workorder.generate"),
	}
}

func (c *RabbitMQConfig) GetConnectionURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}

// ConnectRabbitMQ establece la conexión con RabbitMQ
func ConnectRabbitMQ(config *RabbitMQConfig) (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.GetConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

// DeclareQueue asegura que la cola existe
func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue: %w", err)
	}
	return q, nil
}
