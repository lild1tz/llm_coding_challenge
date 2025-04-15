package googledrive

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gomutex/godocx"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	JSONKey  string `json:"JSON_KEY"`
	FolderID string `json:"FOLDER_ID"`
}

func NewClient(cfg Config) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(cfg.JSONKey), drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt config: %w", err)
	}

	client := jwtConfig.Client(context.Background())
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %w", err)
	}

	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &Client{Drive: driveService, Sheets: sheetsService, folderID: cfg.FolderID}, nil
}

type Client struct {
	Drive  *drive.Service
	Sheets *sheets.Service

	folderID string

	tableMutex sync.Mutex
}

func (c *Client) Release() error {
	return nil
}

func (c *Client) GetTableURL(ctx context.Context, name string) (string, error) {
	fileID, err := c.findFileInFolder(name)
	if err != nil {
		return "", fmt.Errorf("find file: %w", err)
	}

	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit?usp=sharing", fileID), nil
}

func (c *Client) SaveMessage(ctx context.Context, fileName string, text string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		doc.AddParagraph(line)
	}

	var buf bytes.Buffer
	_, err = doc.WriteTo(&buf)
	if err != nil {
		return fmt.Errorf("failed to write docx: %w", err)
	}

	file := &drive.File{
		Name:     fileName,
		Parents:  []string{c.folderID},
		MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	_, err = c.Drive.Files.Create(file).Media(&buf).Do()
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (c *Client) SaveTable(ctx context.Context, name string, table models.Table) error {
	c.tableMutex.Lock()
	defer c.tableMutex.Unlock()

	fileID, err := c.findFileInFolder(name)
	if err != nil {
		return fmt.Errorf("find file: %w", err)
	}

	if fileID == "" {
		fileID, err = c.createSpreadsheet(name)
		if err != nil {
			return fmt.Errorf("create spreadsheet: %w", err)
		}

		if err := c.addHeaders(fileID); err != nil {
			return fmt.Errorf("add headers: %w", err)
		}
	}

	return c.appendData(fileID, table)
}

func (c *Client) findFileInFolder(name string) (string, error) {
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.spreadsheet' and parents in '%s'",
		name, c.folderID)

	resp, err := c.Drive.Files.List().Q(query).Do()
	if err != nil {
		return "", err
	}

	if len(resp.Files) > 0 {
		return resp.Files[0].Id, nil
	}
	return "", nil
}

func (c *Client) createSpreadsheet(name string) (string, error) {
	file := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.spreadsheet",
		Parents:  []string{c.folderID},
	}

	resp, err := c.Drive.Files.Create(file).Do()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *Client) addHeaders(spreadsheetID string) error {
	headers := []interface{}{"Дата", "Подразделение", "Операция", "Культура",
		"За день, га", "С начала операции, га", "Вал за день, ц", "Вал с начала, ц"}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{headers},
	}

	_, err := c.Sheets.Spreadsheets.Values.Update(spreadsheetID, "Sheet1!A1:H1", vr).
		ValueInputOption("RAW").Do()
	return err
}

func (c *Client) appendData(spreadsheetID string, table models.Table) error {
	values := make([][]interface{}, len(table))

	for i, row := range table {
		values[i] = []interface{}{
			row.Date,
			row.Division,
			row.Operation,
			row.Culture,
			row.PerDay,
			row.PerOperation,
			row.ValDay,
			row.ValBeginning,
		}
	}

	vr := &sheets.ValueRange{
		Values: values,
	}

	_, err := c.Sheets.Spreadsheets.Values.Append(spreadsheetID, "Sheet1!A1", vr).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	return err
}
