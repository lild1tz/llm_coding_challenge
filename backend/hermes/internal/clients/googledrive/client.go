package googledrive

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gomutex/godocx"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	JSONKey  string `json:"DRIVE_JSON_KEY"`
	FolderID string `json:"DRIVE_FOLDER_ID"`
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

	if fileID == "" {
		return "", fmt.Errorf("file not found")
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

func (c *Client) SaveImage(ctx context.Context, fileName string, image []byte) error {
	mime := mimetype.Detect(image)

	file := &drive.File{
		Name:     fileName,
		Parents:  []string{c.folderID},
		MimeType: mime.String(),
	}

	_, err := c.Drive.Files.Create(file).Media(bytes.NewReader(image)).Do()
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

	vr := &sheets.ValueRange{Values: values}

	response, err := c.Sheets.Spreadsheets.Values.Append(spreadsheetID, "Sheet1!A1", vr).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		return err
	}

	sheetName := "Sheet1"
	sheetId, err := c.getSheetID(spreadsheetID, sheetName)
	if err != nil {
		return err
	}

	updatedRange := response.Updates.UpdatedRange
	parts := strings.Split(updatedRange, "!")
	if len(parts) != 2 {
		return fmt.Errorf("invalid updated range: %s", updatedRange)
	}
	rangePart := parts[1]
	rangeParts := strings.Split(rangePart, ":")
	if len(rangeParts) != 2 {
		return fmt.Errorf("invalid range part: %s", rangePart)
	}
	startCell, endCell := rangeParts[0], rangeParts[1]

	startRow, err := parseRowNumber(startCell)
	if err != nil {
		return err
	}
	endRow, err := parseRowNumber(endCell)
	if err != nil {
		return err
	}

	if endRow-startRow+1 != len(table) {
		return fmt.Errorf("mismatch between added rows and table length")
	}

	requests := []*sheets.Request{}
	for i := range table {
		row := table[i]
		rowNumber := startRow + i

		if row.DivisionYellow {
			gridRange := createGridRange(sheetId, rowNumber, 1) // Column B
			requests = append(requests, createUpdateCellRequest(gridRange))
		}
		if row.OperationYellow {
			gridRange := createGridRange(sheetId, rowNumber, 2) // Column C
			requests = append(requests, createUpdateCellRequest(gridRange))
		}
		if row.CultureYellow {
			gridRange := createGridRange(sheetId, rowNumber, 3) // Column D
			requests = append(requests, createUpdateCellRequest(gridRange))
		}
	}

	if len(requests) > 0 {
		batchReq := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}
		_, err = c.Sheets.Spreadsheets.BatchUpdate(spreadsheetID, batchReq).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) getSheetID(spreadsheetID, sheetName string) (int64, error) {
	resp, err := c.Sheets.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return 0, err
	}

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("sheet %s not found", sheetName)
}

func parseRowNumber(cell string) (int, error) {
	re := regexp.MustCompile(`[A-Za-z]+(\d+)`)
	matches := re.FindStringSubmatch(cell)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid cell address: %s", cell)
	}
	return strconv.Atoi(matches[1])
}

func createGridRange(sheetId int64, rowNumber, columnIndex int) *sheets.GridRange {
	return &sheets.GridRange{
		SheetId:          sheetId,
		StartRowIndex:    int64(rowNumber - 1),
		EndRowIndex:      int64(rowNumber),
		StartColumnIndex: int64(columnIndex),
		EndColumnIndex:   int64(columnIndex + 1),
	}
}

func createUpdateCellRequest(gridRange *sheets.GridRange) *sheets.Request {
	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Range:  gridRange,
			Fields: "userEnteredFormat.backgroundColor",
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{Red: 1.0, Green: 1.0, Blue: 0.0},
							},
						},
					},
				},
			},
		},
	}
}
