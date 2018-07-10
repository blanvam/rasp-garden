package client

// MQTTCredentialsProvider should return the current username and password for the MQTT client to use.
type MQTTCredentialsProvider func() (username string, password string)

// MQTTOnConnectHandler will be called after the client connects.
// It should be used to resubscribe to topics and perform other connection related tasks.
type MQTTOnConnectHandler func(client MQTTClient)
