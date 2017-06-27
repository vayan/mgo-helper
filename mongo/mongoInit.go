package mongo

import (
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
InitMongoFromConfig setups a DB instance based on the viper config.
It also sets a session ping checker to crash the process if the ping fails.
Aim: let the server/worker restart in case of a mongo stepdown for instance, so it can recover properly.

*/
func InitMongoFromConfig(config Configuration) (*mgo.Database, func()) {
	session := Dial(
		config.UseSSL,
		config.URL,
		config.SSLCert,
	)
	pingSession(session, config.PingFrequency)
	db := session.DB("") // use database name from URL
	teardown := session.Close
	return db, teardown
}

func pingSession(session *mgo.Session, frequency time.Duration) {
	frequency = time.Millisecond * frequency
	checker := CreateMongoSessionPinger(session, frequency)
	checker.Listen()
}
