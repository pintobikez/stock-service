package rabbitmq

import (
	"encoding/json"
	"fmt"

	cnfs "bitbucket.org/ricardomvpinto/stock-service/config/structures"
	"github.com/streadway/amqp"
)

type Rabbitmq struct {
	conn    *amqp.Connection
	config  *cnfs.PublisherConfig
	channel *amqp.Channel
}

// Creates a pointer to a new Rabbitmq struct
func New(cnfg *cnfs.PublisherConfig) (*Rabbitmq, error) {
	p := &Rabbitmq{config: cnfg}
	err := p.Connect()

	return p, err
}

// Connects to the RabbitServer and to the defined ExchangeQueue
func (p *Rabbitmq) Connect() error {

	var err error

	p.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", p.config.Pw, p.config.User, p.config.Host, p.config.Port))
	if err != nil {
		return err
	}

	if p.channel, err = p.conn.Channel(); err != nil {
		p.Close()
		return err
	}

	if err = p.channel.ExchangeDeclare(
		p.config.Exchange, // name
		"fanout",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	); err != nil {
		p.Close()
	}

	return err
}

// Closes both Rabbit and Exchange connection
func (p *Rabbitmq) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

// Publishes a message to the Defined ExchangeQueue
func (p *Rabbitmq) Publish(s interface{}) error {

	if p.channel == nil || p.conn == nil {
		if err := p.Connect(); err != nil {
			return err
		}
	}

	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err = p.channel.Publish(
		p.config.Exchange, // exchange
		"",                // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        body,
		}); err != nil {
		return err
	}

	return nil
}

// Health Endpoint of the Client
func (p *Rabbitmq) Health() error {

	if p.config == nil {
		return fmt.Errorf("Publisher configuration not loaded")
	}

	pt := &Rabbitmq{config: p.config}
	if err := pt.Connect(); err != nil {
		return err
	}

	pt.Close()
	return nil
}
