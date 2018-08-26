package client

// CredentialsProvider should return the current username and password for the MQTT client to use.
type CredentialsProvider func() (username string, password string)
