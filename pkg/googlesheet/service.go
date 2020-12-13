package googlesheet

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Service ...
type Service struct {
	svc              *sheets.Service
	spreadSheetID    string
	spreadSheetRange string
}

// New ...
func New(credentials, sheetID, sheetRange string) (*Service, error) {
	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return &Service{
		svc:              sheetsService,
		spreadSheetRange: sheetRange,
		spreadSheetID:    sheetID,
	}, nil
}

// AppendMembersCount ...
func (s *Service) AppendMembersCount(data int) error {
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{time.Now().UTC(), time.Now().Month(), data}},
	}
	resp, err := s.svc.Spreadsheets.Values.
		Append(s.spreadSheetID, s.spreadSheetRange, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.ServerResponse.HTTPStatusCode != http.StatusOK {
		log.Println(resp.ServerResponse.HTTPStatusCode)
		msg := fmt.Sprintf("error with status code: %d", resp.ServerResponse.HTTPStatusCode)
		return errors.New(msg)
	}
	return nil
}
