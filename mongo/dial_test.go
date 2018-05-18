package mongo

import (
	"crypto/tls"
	"testing"

	mgo "github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/transcovo/mgo-helper/utils"
)

func TestDial(t *testing.T) {
	session := Dial(false, "mongodb://localhost:27017/test", []byte{})
	session.Close()
}

func TestDial_SSL(t *testing.T) {
	defer (func() {
		sslSessionFactory = mgo.DialWithInfo
	})()
	sslSessionFactoryCalled := false
	sslSessionFactory = func(info *mgo.DialInfo) (*mgo.Session, error) {
		sslSessionFactoryCalled = true
		return nil, nil
	}

	Dial(true, "mongodb://localhost:27017/test", []byte{})

	assert.True(t, sslSessionFactoryCalled)
}

func TestPanicIfError_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Panic did not occur")
		}
	}()

	utils.PanicIfError(errors.New("An error"))
}

func TestPanicIfError_NoPanic(t *testing.T) {
	utils.PanicIfError(nil)
}

func TestMakeMgoDialServer(t *testing.T) {
	mgoDialServer := makeMgoDialServer(&tls.Config{})
	mgoDialServer(&mgo.ServerAddr{})
}
