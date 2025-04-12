package whatsapp

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Config struct {
	DatabaseURL string `json:"DATABASE_URL"`

	LoggingLevel string `json:"LOGGING_LEVEL" cfgDefault:"ERROR"`
}

func NewClient(cfg Config) (*Client, error) {
	dbLog := waLog.Stdout("Database", cfg.LoggingLevel, true)

	container, err := sqlstore.New("pgx", cfg.DatabaseURL, dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore: %w", err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get first device: %w", err)
	}
	clientLog := waLog.Stdout("Client", cfg.LoggingLevel, true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	return &Client{client}, nil
}

type Client struct {
	*whatsmeow.Client
}

func (c *Client) Release() error {
	c.Disconnect()
	return nil
}

func (c *Client) Start() error {
	if c.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := c.GetQRChannel(context.Background())
		err := c.Connect()
		if err != nil {
			log.Fatal(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := c.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	return nil
}
