package mongo

import (
	"context"
	"time"

	mgo "gopkg.in/mgo.v2"
)

/*
Configuration defines the require fields in order to correctly set up a mongo connection
*/
type Configuration struct {
	URL           string
	UseSSL        bool
	SSLCert       []byte
	PingFrequency time.Duration
}

/*
InitMongoFromConfig setups a DB instance based on the config.
It also sets a session ping checker to crash the process if the ping fails.
Aim: let the server/worker restart in case of a mongo stepdown for instance, so it can recover properly.
*/
func InitMongoFromConfig(config Configuration) (*mgo.Database, func()) {
	ctx := context.Background()
	return initWithContext(ctx, config)
}

func initWithContext(ctx context.Context, config Configuration) (*mgo.Database, func()) {
	cancelCtx, cancel := context.WithCancel(ctx)
	session := dial(
		config.UseSSL,
		config.URL,
		config.SSLCert,
	)
	pingSession(cancelCtx, session, config.PingFrequency)
	db := session.DB("") // use database name from URL
	teardown := func() {
		cancel()
		session.Close()
	}
	return db, teardown
}

func pingSession(ctx context.Context, session *mgo.Session, frequency time.Duration) {
	frequency = time.Millisecond * frequency
	checker := createMongoSessionPinger(session, frequency)
	checker.Listen(ctx)
}
