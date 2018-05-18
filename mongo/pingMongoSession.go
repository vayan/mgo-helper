package mongo

import (
	"io"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
	"github.com/transcovo/go-chpr-logger"
	"github.com/transcovo/go-chpr-metrics"
)

/*
Listener is the interface for the only exposed part of the sessionChecker that is public:
its Listen function.
*/
type Listener interface {
	Listen()
}

type pinger interface {
	Ping() error
}

type errorHandler interface {
	onError(error, string)
}

/*
processQuitter implements the errorHandler interface
*/
type processQuitter struct {
	logger *logrus.Logger
}

/*
onError logs and exits the process with status code 1 (see logger.Fatal doc)
*/
func (loggerWrapper *processQuitter) onError(err error, message string) {
	loggerWrapper.logger.WithError(err).Fatal(message)
}

type sessionChecker struct {
	session      pinger
	trigger      <-chan time.Time
	errorHandler errorHandler
}

func (checker *sessionChecker) Listen() {
	logger.Info("[mongo.PingMongoSession] Starting ping")
	go func() {
		defer func() {
			// This defer is only for testing
			// For production, an Error log warns that the go routine failed and hence the go routine failed
			if r := recover(); r != nil {
				logger.WithField("err", r).Error("[mongo.PingMongoSession] Recovered from failing Ping gorountine")
			}
		}()
		for range checker.trigger {
			checker.checkConnection()
		}
	}()
}

func (checker *sessionChecker) checkConnection() {
	err := checker.session.Ping()
	if err == nil {
		metrics.Increment("mongo.ping_session.ok")
		return
	}
	checker.crashOnError(err)
}

func (checker *sessionChecker) crashOnError(err error) {
	if err == io.EOF {
		metrics.Increment("mongo.ping_session.eof_err")
		checker.errorHandler.onError(err, "[mongo.PingMongoSession] EOF DB Error, crashing")
	} else {
		metrics.Increment("mongo.ping_session.unknown_err")
		checker.errorHandler.onError(err, "[mongo.PingMongoSession] Unknown DB error, crashing")
	}
}

/*
CreateMongoSessionPinger creates a Listener based on a mgo session that:
- regularly pings the DB session
- in case of an EOF/Unknown error, it exits the process so the server can restart and recover properly
*/
func CreateMongoSessionPinger(session *mgo.Session, pingInterval time.Duration) Listener {
	ticker := time.NewTicker(pingInterval)
	handler := &processQuitter{logger.GetLogger()}
	checker := &sessionChecker{
		session:      session,
		trigger:      ticker.C,
		errorHandler: handler,
	}
	return checker
}
