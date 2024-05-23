package firebase

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
)

const (
	UiPort        = "4000/tcp"
	HubPort       = "4400/tcp"
	LoggingPort   = "4600/tcp"
	FunctionsPort = "5001/tcp"
	FirestorePort = "8080/tcp"
	PubsubPort    = "8085/tcp"
	DatabasePort  = "9000/tcp"
	AuthPort      = "9099/tcp"
	StoragePort   = "9199/tcp"
	HostingPort   = "6000/tcp"
)

func (c *FirebaseContainer) connectionString(ctx context.Context, portName nat.Port) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := c.MappedPort(ctx, portName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

func (c *FirebaseContainer) UIConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, UiPort)
}

func (c *FirebaseContainer) HubConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, HubPort)
}

func (c *FirebaseContainer) LoggingConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, LoggingPort)
}

func (c *FirebaseContainer) FunctionsConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, FunctionsPort)
}

func (c *FirebaseContainer) FirestoreConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, FirestorePort)
}

func (c *FirebaseContainer) PubSubConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, PubsubPort)
}

func (c *FirebaseContainer) DatabaseConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, DatabasePort)
}

func (c *FirebaseContainer) AuthConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, AuthPort)
}

func (c *FirebaseContainer) StorageConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, StoragePort)
}

func (c *FirebaseContainer) HostingConnectionString(ctx context.Context) (string, error) {
	return c.connectionString(ctx, HostingPort)
}
