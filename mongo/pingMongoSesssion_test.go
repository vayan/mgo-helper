package mongo

import (
	"errors"
	"io"
	"testing"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

type mockPinger struct {
	err error
}

func (mock *mockPinger) Ping() error {
	return mock.err
}

type mockErrorHandler struct {
	done    chan bool
	err     error
	message string
}

func (mock *mockErrorHandler) onError(err error, msg string) {
	mock.err = err
	mock.message = msg
	mock.done <- true
}

/*
TestSessionChecker_Listen_Error tests that
the check is properly triggered by the channel and that
the errorHandler is properly called when the ping returns an error
*/
func TestSessionChecker_Listen_Error(t *testing.T) {
	channel := make(chan time.Time, 1)
	pinger := &mockPinger{err: io.EOF}
	handler := &mockErrorHandler{done: make(chan bool, 1)}
	checker := &sessionChecker{
		session:      pinger,
		trigger:      channel,
		errorHandler: handler,
	}
	checker.Listen()
	channel <- time.Unix(0, 0)
	<-handler.done

	assert.Equal(t, "[mongo.PingMongoSession] EOF DB Error, crashing", handler.message)
	assert.Equal(t, errors.New("EOF"), handler.err)
}

/*
TestSessionChecker_CheckConnection_NilError tests that the errorHandler is not called
when the ping returns no error
*/
func TestSessionChecker_CheckConnection_NilError(t *testing.T) {
	channel := make(chan time.Time, 1)
	pinger := &mockPinger{err: nil}
	handler := &mockErrorHandler{done: make(chan bool, 1)}
	checker := &sessionChecker{
		session:      pinger,
		trigger:      channel,
		errorHandler: handler,
	}
	checker.checkConnection()

	assert.Equal(t, "", handler.message)
	assert.Nil(t, handler.err)
}

/*
TestSessionChecker_CheckConnection_EOFError tests that the errorHandler is properly called
when the ping returns an EOF
*/
func TestSessionChecker_CheckConnection_EOFError(t *testing.T) {
	channel := make(chan time.Time, 1)
	pinger := &mockPinger{err: io.EOF}
	handler := &mockErrorHandler{done: make(chan bool, 1)}
	checker := &sessionChecker{
		session:      pinger,
		trigger:      channel,
		errorHandler: handler,
	}
	checker.checkConnection()

	assert.Equal(t, "[mongo.PingMongoSession] EOF DB Error, crashing", handler.message)
	assert.Equal(t, errors.New("EOF"), handler.err)
}

/*
TestSessionChecker_CheckConnection_UnknownError tests that the errorHandler is properly called
when the ping returns an error other than an EOF
*/
func TestSessionChecker_CheckConnection_UnknownError(t *testing.T) {
	channel := make(chan time.Time, 1)
	pinger := &mockPinger{err: errors.New("Unknown")}
	handler := &mockErrorHandler{done: make(chan bool, 1)}
	checker := &sessionChecker{
		session:      pinger,
		trigger:      channel,
		errorHandler: handler,
	}
	checker.checkConnection()

	assert.Equal(t, "[mongo.PingMongoSession] Unknown DB error, crashing", handler.message)
	assert.Equal(t, errors.New("Unknown"), handler.err)
}

/*
TestCreateMongoSessionPinger_Success tests the constructor for the actual
MongoSessionPinger using an actual mgo session.
*/
func TestCreateMongoSessionPinger_Success(t *testing.T) {
	session := &mgo.Session{}
	pinger := CreateMongoSessionPinger(session, time.Second)
	assert.NotNil(t, pinger)
}
