package internalgrpc

import (
	"context"
	"net"
	"testing"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type ServerTestSuite struct {
	suite.Suite
	server *grpc.Server
	client *grpc.ClientConn
}

func (s *ServerTestSuite) SetupSuite() {
	addr := "localhost:8090"
	size := 1024 * 1024
	log := logger.NewMock()
	app := app.NewApplication(log, memorystorage.New())
	lis := bufconn.Listen(size)
	s.server = grpc.NewServer()
	server := &Server{
		addr: addr,
		log:  log,
		app:  app,
		srv:  s.server,
	}
	RegisterEventServiceServer(s.server, server)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.T().Log("error running server", err.Error())
		}
	}()

	conn, err := grpc.NewClient(
		addr,
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.T().Log("error opening client", err.Error())
	}
	s.client = conn
}

func (s *ServerTestSuite) TearDownSuite() {
	err := s.client.Close()
	if err != nil {
		s.T().Log("error closing client", err.Error())
	}
	s.server.Stop()
}

func TestServer(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}

func (s *ServerTestSuite) TestAddEvent() {
	req := &CreateEventReq{
		Title:  "some event",
		UserID: uuid.NewString(),
	}
	resp, err := NewEventServiceClient(s.client).AddEvent(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	event := resp.GetEvent()
	s.Require().NotNil(event)
	s.Require().NotEmpty(event.ID)
	s.Require().Equal(req.Title, event.Title)
	s.Require().NotNil(event.Date)
	s.Require().Equal(req.UserID, event.UserID)

	s.deleteEvent(event.ID)
}

func (s *ServerTestSuite) TestUpdateEvent() {
	event := s.addEvent()
	reqEvent := &Event{
		ID:     event.ID,
		Title:  "some event 111",
		Date:   event.Date,
		UserID: event.UserID,
	}
	req := &UpdateEventReq{
		Event: reqEvent,
	}
	resp, err := NewEventServiceClient(s.client).UpdateEvent(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	respEvent := resp.GetEvent()
	s.Require().NotNil(respEvent)
	s.Require().Equal(reqEvent.ID, respEvent.ID)
	s.Require().Equal(reqEvent.Title, respEvent.Title)
	s.Require().Equal(reqEvent.Date.AsTime(), respEvent.Date.AsTime())
	s.Require().Equal(reqEvent.UserID, respEvent.UserID)

	s.deleteEvent(event.ID)
}

func (s *ServerTestSuite) TestDeleteEvent() {
	event := s.addEvent()
	req := &DeleteEventReq{
		ID: event.ID,
	}
	resp, err := NewEventServiceClient(s.client).DeleteEvent(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *ServerTestSuite) TestGetEvents() {
	event := s.addEvent()
	req := &GetEventsReq{
		Period: GetEventsReq_PERIOD_DAY,
		Year:   int32(event.Date.AsTime().Local().Year()),
		Month:  int32(event.Date.AsTime().Local().Month()),
		Day:    int32(event.Date.AsTime().Local().Day()),
	}
	resp, err := NewEventServiceClient(s.client).GetEvents(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	events := resp.GetEvents()
	s.Require().Len(events, 1)
	s.Require().Equal(event.ID, events[0].ID)
	s.Require().Equal(event.Title, events[0].Title)
	s.Require().Equal(event.Date, events[0].Date)
	s.Require().Equal(event.UserID, events[0].UserID)

	s.deleteEvent(event.ID)
}

func (s *ServerTestSuite) addEvent() *Event {
	req := &CreateEventReq{
		Title:  "some event",
		UserID: uuid.NewString(),
	}
	resp, err := NewEventServiceClient(s.client).AddEvent(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	return resp.GetEvent()
}

func (s *ServerTestSuite) deleteEvent(id string) {
	req := &DeleteEventReq{
		ID: id,
	}
	resp, err := NewEventServiceClient(s.client).DeleteEvent(context.Background(), req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}
