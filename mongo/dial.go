/*
Package mongo implements the required plumbing to connect to mongo using ssl.
*/
package mongo

import (
	"crypto/tls"
	"crypto/x509"
	"net"

	"github.com/transcovo/go-chpr-logger"
	"github.com/transcovo/mgo-helper/utils"
	"gopkg.in/mgo.v2"
)

var tlsDialer = tls.Dial

var sslSessionFactory = mgo.DialWithInfo

func makeMgoDialServer(tlsConfig *tls.Config) func(addr *mgo.ServerAddr) (net.Conn, error) {
	return func(addr *mgo.ServerAddr) (net.Conn, error) {
		logger.Info("[makeMgoDialServer] Calling Dial")
		conn, err := tlsDialer("tcp", addr.String(), tlsConfig)
		if err != nil {
			logger.WithField("err", err).Error("[makeMgoDialServer] Error while dialing")
		}
		return conn, err
	}
}

/*
DialWithSSL connects to a mongo database with SSL using the server public certificate passed with the ca argument.

It panics if anything goes wrong.
*/
func DialWithSSL(mongoURL string, ca []byte) *mgo.Session {
	roots := x509.NewCertPool()

	roots.AppendCertsFromPEM(ca)

	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = roots

	dialInfo, err := mgo.ParseURL(mongoURL)
	utils.PanicIfError(err)

	dialInfo.DialServer = makeMgoDialServer(tlsConfig)
	//Here is the session you are looking for. Up to you from here ;)
	logger.Info("[DialWithSSL] Calling DialWithInfo")
	session, err := sslSessionFactory(dialInfo)
	utils.PanicIfError(err)

	return session
}

/*
DialWithoutSSL connects to a mongodb database without using SSL.

It panics if anything goes wrong.
*/
func DialWithoutSSL(mongoURL string) *mgo.Session {
	session, err := mgo.Dial(mongoURL)
	utils.PanicIfError(err)

	return session
}

/*
Dial connects to amongodb database with our without using SSL, depending on the value of the "useSSL" param.

It panics if anything goes wrong.
*/
func Dial(useSSL bool, mongoURL string, ca []byte) *mgo.Session {
	if useSSL {
		logger.Info("[Dial] Connecting with SSL")
		return DialWithSSL(mongoURL, ca)
	}
	logger.Info("[Dial] Connecting without SSL")
	return DialWithoutSSL(mongoURL)
}
