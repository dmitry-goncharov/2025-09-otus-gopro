package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var handler = NewHandler(logger.NewMock(), app.NewApplication(logger.NewMock(), memorystorage.New()))

func TestAddEventHandler(t *testing.T) {
	req := EventTitleDto{
		Title: "some title",
	}
	body, err := json.Marshal(req)
	require.NoError(t, err)

	request := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	request.Header.Add(UserID, uuid.NewString())
	responseWriter := httptest.NewRecorder()

	handler.AddEvent(responseWriter, request)

	require.Equal(t, http.StatusOK, responseWriter.Code)

	resp := &EventDto{}
	err = json.NewDecoder(responseWriter.Body).Decode(resp)
	require.NoError(t, err)

	require.NotEmpty(t, resp.ID)
	require.Equal(t, req.Title, resp.Title)
	require.NotEmpty(t, resp.Date)
	require.NotEmpty(t, resp.UserID)

	deleteEvent(resp.ID)
}

func TestUpdateEventHandler(t *testing.T) {
	event := addEvent()

	req := &EventDto{
		ID:     event.ID,
		Title:  "some title 111",
		Date:   event.Date,
		UserID: event.UserID,
	}
	body, err := json.Marshal(req)
	require.NoError(t, err)

	request := httptest.NewRequest("POST", "/events/{id}", bytes.NewReader(body))
	request.SetPathValue("id", event.ID)

	responseWriter := httptest.NewRecorder()

	handler.UpdateEvent(responseWriter, request)

	require.Equal(t, http.StatusOK, responseWriter.Code)

	resp := &EventDto{}
	err = json.NewDecoder(responseWriter.Body).Decode(resp)
	require.NoError(t, err)

	require.Equal(t, req.ID, resp.ID)
	require.Equal(t, req.Title, resp.Title)
	require.Equal(t, req.Date, resp.Date)
	require.Equal(t, req.UserID, resp.UserID)

	deleteEvent(event.ID)
}

func TestDeleteEventHandler(t *testing.T) {
	event := addEvent()

	request := httptest.NewRequest("DELETE", "/events/{id}", nil)
	request.SetPathValue("id", event.ID)

	responseWriter := httptest.NewRecorder()

	handler.DeleteEvent(responseWriter, request)

	require.Equal(t, http.StatusOK, responseWriter.Code)
}

func TestGetEventsHandler(t *testing.T) {
	event := addEvent()

	request := httptest.NewRequest("GET", "/events", nil)
	query := request.URL.Query()
	query.Add(Period, Day)
	query.Add(Year, strconv.Itoa(event.Date.Local().Year()))
	query.Add(Month, strconv.Itoa(int(event.Date.Local().Month())))
	query.Add(Day, strconv.Itoa(event.Date.Local().Day()))
	request.URL.RawQuery = query.Encode()

	responseWriter := httptest.NewRecorder()

	handler.GetEvents(responseWriter, request)

	require.Equal(t, http.StatusOK, responseWriter.Code)

	var resp []EventDto
	err := json.NewDecoder(responseWriter.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp, 1)
	require.Equal(t, event.ID, resp[0].ID)
	require.Equal(t, event.Title, resp[0].Title)
	require.Equal(t, event.Date, resp[0].Date)
	require.Equal(t, event.UserID, resp[0].UserID)

	deleteEvent(event.ID)
}

func addEvent() *EventDto {
	req := EventTitleDto{
		Title: "some title",
	}
	body, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	request.Header.Add(UserID, uuid.NewString())
	responseWriter := httptest.NewRecorder()

	handler.AddEvent(responseWriter, request)

	resp := &EventDto{}
	json.NewDecoder(responseWriter.Body).Decode(resp)

	return resp
}

func deleteEvent(id string) {
	request := httptest.NewRequest("DELETE", "/events/{id}", nil)
	request.SetPathValue("id", id)

	responseWriter := httptest.NewRecorder()

	handler.DeleteEvent(responseWriter, request)
}
