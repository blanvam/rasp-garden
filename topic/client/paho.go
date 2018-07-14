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
	timeout     time.Duration
	clientID    string
	username    string
	password    string
	servers     []string
	storePath   string
	client      paho.Client
	certificate CredentialsProvider
}

func NewPahoClient(t time.Duration, cid string, u string, p string, s []string) topic.Client {

	return &pahoClient{
		options:  paho.NewClientOptions(),
		timeout:  t,
		clientID: cid,
		username: u,
		password: p,
		servers:  s,
	}
}

// IsConnected return true if the client is connected
func (p *pahoClient) IsConnected(c chan bool) {
	c <- p.client.IsConnected()
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
	p.options.SetUsername(p.username)
	p.options.SetPassword(p.password)

	p.options.SetCleanSession(false)
	p.options.SetAutoReconnect(true)
	p.options.SetProtocolVersion(4)
	p.options.SetStore(store)
	// p.options.SetCredentialsProvider(func() (string, string) { return p.credentialsProvider() })
	p.options.SetOnConnectHandler(func(i paho.Client) { log.Println("Connected") })
	p.options.SetConnectionLostHandler(func(client paho.Client, e error) { log.Printf("Connection Lost. Error: %v", e) })
	p.options.SetOnConnectHandler(func(client paho.Client) { log.Println("Handler connected") })

	for _, server := range p.servers {
		p.options.AddBroker(server)
	}

	p.client = paho.NewClient(p.options)

	token := p.client.Connect()
	c <- p.waitForToken(token)
}

func (p *pahoClient) Disconnect(c chan error, quiesce uint) {
	p.client.Disconnect(quiesce)
	p.client = nil
	c <- nil
}

func (p *pahoClient) Publish(c chan error, topic string, qos uint8, payload interface{}) {
	token := p.client.Publish(topic, qos, true, payload)
	c <- p.waitForToken(token)
}

func (p *pahoClient) Subscribe(c chan error, topic string, qos uint8, callback topic.CallbackHandler) {
	handler := func(i paho.Client, message paho.Message) {
		callback(topic, p.clientID, message.Payload())
	}
	token := p.client.Subscribe(topic, qos, handler)
	c <- p.waitForToken(token)
}

func (p *pahoClient) Unsubscribe(c chan error, topic string) {
	token := p.client.Unsubscribe(topic)
	c <- p.waitForToken(token)
}

func (p *pahoClient) waitForToken(token paho.Token) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
