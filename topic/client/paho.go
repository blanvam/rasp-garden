package client

import (
	"context"
	"log"
	"time"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/topic"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type pahoClient struct {
	options     *paho.ClientOptions
	client      paho.Client
	username 	string
	password	string
	servers		...string
	clientID    string
	storePath   string
	timeout     time.Duration
	certificate MQTTCredentialsProvider
}

func NewPahoClient(cid string) topic.Client {

	return &pahoClient{
		options:  paho.NewClientOptions(),
		clientID: cid,
	}
}

// Connect try to connect to the given MQTT server
func (p *pahoClient) Connect(c chan error) {

	var store paho.Store
	if p.storePath == "" {
		store = paho.NewMemoryStore()
	} else {
		store = paho.NewFileStore(p.storePath)
	}

	// p.options.SetTLSConfig(&tls.Config{
	//	Certificates:       []tls.Certificate{p.options.Credentials.Certificate},
	//	InsecureSkipVerify: true,
	//})

	p.options.SetClientID(p.clientID)
	p.options.SetUsername(username)
	p.options.SetPassword(password)

	p.options.SetCleanSession(false)
	p.options.SetAutoReconnect(true)
	p.options.SetProtocolVersion(4)
	p.options.SetStore(store)
	// p.options.SetCredentialsProvider(func() (string, string) { return p.credentialsProvider() })
	p.options.SetOnConnectHandler(func(i paho.Client) { log.Println("Connected") })
	p.options.SetConnectionLostHandler(func(client paho.Client, e error) { log.Printf("Connection Lost. Error: %v", e) })
	p.options.SetOnConnectHandler(func(client paho.Client) { log.Println("Handler connected") })

	for _, server := range servers {
		p.options.AddBroker(server)
	}

	p.client = paho.NewClient(p.options)

	token := p.client.Connect()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c <- p.waitForToken(ctx, token)
}

func (p *pahoClient) Publish(ctx context.Context, topic string, qos uint8, payload interface{}) error {
	token := p.client.Publish(topic, qos, true, payload)
	return p.waitForToken(ctx, token)
}

func (p *pahoClient) Subscribe(c chan []byte, ctx context.Context, topic string, qos uint8) error {
	token := p.client.Subscribe(topic, qos, handler)
	return p.waitForToken(ctx, token)
}

func (p *pahoClient) Unsubscribe(ctx context.Context, topic string) error {
	token := p.client.Unsubscribe(topic)
	return p.waitForToken(ctx, token)
}

func (p *pahoClient) waitForToken(ctx context.Context, token paho.Token) error {
	result := make(chan error)
	cancelled := false
	go func() {
		defer func() { result <- token.Error() }()
		for {
			if (token.WaitTimeout(p.timeout)) || cancelled {
				return
			}
		}
	}()
	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		cancelled = true
	}
	return entity.ErrCancelled
}
