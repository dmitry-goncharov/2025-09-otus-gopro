//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	internalhttp "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/server/http"
	"github.com/stretchr/testify/suite"

	// Posgresql driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type CalendarSuite struct {
	suite.Suite
	ctx    context.Context
	client *http.Client
	db     *sqlx.DB
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) SetupSuite() {
	s.ctx = context.Background()
	s.client = &http.Client{
		Timeout: 10 * time.Second,
	}
	db, err := sqlx.Open("pgx", "host=db port=5432 user=postgres password=postgres dbname=otus sslmode=disable")
	if err != nil {
		s.Require().NoError(err)
	}
	err = db.PingContext(s.ctx)
	if err != nil {
		s.Require().NoError(err)
	}
	s.db = db
}

func (s *CalendarSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			s.Require().NoError(err)
		}
	}
}

func (s *CalendarSuite) TestCreateEvent() {
	reqDto := internalhttp.EventTitleDto{
		Title: "some title",
		Date:  time.Now(),
	}
	body, err := json.Marshal(reqDto)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", "http://app:8080/events", bytes.NewReader(body))
	if err != nil {
		s.Require().NoError(err)
	}
	req.Header.Add("User-Agent", "Go-http-client/integrationtest")
	req.Header.Add("X-User-Id", "00112233-4455-6677-8899-ccddeeffaabb")

	resp, err := s.client.Do(req)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respDto := &internalhttp.EventDto{}
	err = json.NewDecoder(resp.Body).Decode(respDto)
	s.Require().NoError(err)

	s.Require().NotEmpty(respDto.ID)
	s.Require().Equal(reqDto.Title, respDto.Title)
	s.Require().NotEmpty(respDto.Date)
	s.Require().NotEmpty(respDto.UserID)

	query := `select title from events where id = :id;`
	args := map[string]any{
		"id": respDto.ID,
	}
	rows, err := s.db.NamedQueryContext(s.ctx, query, args)
	if err != nil {
		s.Require().NoError(err)
	}
	var title string
	rows.Next()
	rows.Scan(&title)
	s.Require().Equal(reqDto.Title, title)
}

func (s *CalendarSuite) TestCreateEventError() {
	reqDto := internalhttp.EventTitleDto{
		Title: "some title",
		Date:  time.Now(),
	}
	body, err := json.Marshal(reqDto)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", "http://app:8080/events", bytes.NewReader(body))
	if err != nil {
		s.Require().NoError(err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp.Body.Close()

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *CalendarSuite) TestGetEvents() {
	reqDto := internalhttp.EventTitleDto{
		Title: "some title",
		Date:  time.Now(),
	}
	body, err := json.Marshal(reqDto)
	s.Require().NoError(err)

	req0, err := http.NewRequest("POST", "http://app:8080/events", bytes.NewReader(body))
	if err != nil {
		s.Require().NoError(err)
	}
	req0.Header.Add("User-Agent", "Go-http-client/integrationtest")
	req0.Header.Add("X-User-Id", "00112233-4455-6677-8899-ccddeeffaabb")

	resp0, err := s.client.Do(req0)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp0.Body.Close()

	s.Require().Equal(http.StatusOK, resp0.StatusCode)

	now := time.Now()
	req, err := http.NewRequest("GET", "http://app:8080/events?p=d&y="+strconv.Itoa(now.Year())+"&m="+strconv.Itoa(int(now.Month()))+"&d="+strconv.Itoa(now.Day()), nil)
	if err != nil {
		s.Require().NoError(err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var respDto []internalhttp.EventDto
	err = json.NewDecoder(resp.Body).Decode(&respDto)
	s.Require().NoError(err)

	s.Require().True(len(respDto) > 0)
}

func (s *CalendarSuite) TestGetEventsEmpty() {
	now := time.Now()
	req, err := http.NewRequest("GET", "http://app:8080/events?p=d&y="+strconv.Itoa(now.Year()-1)+"&m="+strconv.Itoa(int(now.Month()))+"&d="+strconv.Itoa(now.Day()), nil)
	if err != nil {
		s.Require().NoError(err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var respDto []internalhttp.EventDto
	err = json.NewDecoder(resp.Body).Decode(&respDto)
	s.Require().NoError(err)

	s.Require().Len(respDto, 0)
}

func (s *CalendarSuite) TestGetEventNotifications() {
	reqDto := internalhttp.EventTitleDto{
		Title: "some title",
		Date:  time.Now().Add(1 * time.Minute),
	}
	body, err := json.Marshal(reqDto)
	s.Require().NoError(err)

	req0, err := http.NewRequest("POST", "http://app:8080/events", bytes.NewReader(body))
	if err != nil {
		s.Require().NoError(err)
	}
	req0.Header.Add("User-Agent", "Go-http-client/integrationtest")
	req0.Header.Add("X-User-Id", "00112233-4455-6677-8899-ccddeeffaabb")

	resp0, err := s.client.Do(req0)
	if err != nil {
		s.Require().NoError(err)
	}
	defer resp0.Body.Close()

	s.Require().Equal(http.StatusOK, resp0.StatusCode)

	respDto := &internalhttp.EventDto{}
	err = json.NewDecoder(resp0.Body).Decode(respDto)
	s.Require().NoError(err)

	time.Sleep(15 * time.Second)

	query := `select count(*) from event_notifications where event_id = :event_id;`
	args := map[string]any{
		"event_id": respDto.ID,
	}
	rows, err := s.db.NamedQueryContext(s.ctx, query, args)
	if err != nil {
		s.Require().NoError(err)
	}
	var count int
	rows.Next()
	rows.Scan(&count)
	s.Require().True(count > 0)
}
