package mongo

import (
	"testing"

	mgo "github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

/*
TestInitMongoFromConfig_Success tests that InitMongoFromConfig returns a working db
*/
func TestInitMongoFromConfig_Success(t *testing.T) {
	mongoConfig := Configuration{
		PingFrequency: 100,
		SSLCert:       []byte{},
		UseSSL:        false,
		URL:           "mongodb://localhost:27017/some-test-db",
	}
	db, teardown := InitMongoFromConfig(mongoConfig)
	defer teardown()
	assert.IsType(t, &mgo.Database{}, db)
	assert.Equal(t, "some-test-db", db.Name)
	pingErr := db.Session.Ping()
	assert.Nil(t, pingErr)
}
