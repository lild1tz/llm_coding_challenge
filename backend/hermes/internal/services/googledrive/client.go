package googledrive

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gomutex/godocx"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Config struct {
	JSONKey  string `json:"JSON_KEY"`
	FolderID string `json:"FOLDER_ID"`
}

func NewClient(cfg Config) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(cfg.JSONKey), drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt config: %w", err)
	}

	client := jwtConfig.Client(context.Background())
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %w", err)
	}

	return &Client{Service: driveService, folderID: cfg.FolderID}, nil
}

type Client struct {
	*drive.Service
	folderID string
}

func (c *Client) Release() error {
	return nil
}

func (c *Client) SaveMessage(ctx context.Context, name string, number int, timestamp time.Time, text string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		doc.AddParagraph(line)
	}

	// doc.SaveToFile(fmt.Sprintf("%s_%d.docx", name, number)) - optional think how to save in host fs

	// f, err := os.Create("generated.docx")
	// if err != nil {
	// 	return fmt.Errorf("failed to create file: %w", err)
	// }

	// _, err = doc.WriteTo(f)
	// if err != nil {
	// 	return fmt.Errorf("failed to write docx: %w", err)
	// }

	// f.Close()

	var buf bytes.Buffer
	_, err = doc.WriteTo(&buf)
	if err != nil {
		return fmt.Errorf("failed to write docx: %w", err)
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("failed to load location: %w", err)
	}

	timestamp = timestamp.In(loc)

	fileName := fmt.Sprintf("%s_%d_%s.docx", strings.ReplaceAll(name, " ", "_"), number, timestamp.Format("04-15-02-01-2006"))

	file := &drive.File{
		Name:     fileName,
		Parents:  []string{c.folderID},
		MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	_, err = c.Files.Create(file).Media(&buf).Do()
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (c *Client) SaveTable(ctx context.Context, createdAt time.Time, table models.Table) error {
	return nil
}
