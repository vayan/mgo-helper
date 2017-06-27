package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/transcovo/mgo-helper/mongo"
	mgo "gopkg.in/mgo.v2"
)

type testDocument struct {
	Key   string `bson:"key"`
	Value int    `bson:"value"`
}

/*
TestBaseModel_EnsureIndexes tests that the indexes are correctly created
*/
func TestBaseModel_EnsureIndexes(t *testing.T) {
	config := mongo.Configuration{
		URL:           "mongodb://localhost:27017/some-test-db",
		UseSSL:        false,
		SSLCert:       []byte{},
		PingFrequency: 100,
	}

	db, teardown := mongo.InitMongoFromConfig(config)
	defer teardown()

	// initialize the collection..
	db.C("some-collection").Insert(&testDocument{})

	model := &BaseModel{
		DB:             db,
		CollectionName: "some-collection",
		Indexes: []*mgo.Index{
			{
				Unique: true,
				Name:   "test_1",
				Key:    []string{"first_key"},
			},
		},
	}

	SetupIndexes([]BaseModelInterface{model}, false)

	indexes, err := db.C("some-collection").Indexes()
	assert.Nil(t, err)
	assert.Equal(t, []mgo.Index{
		{Key: []string{"_id"}, Name: "_id_"},
		{Key: []string{"first_key"}, Name: "test_1", Unique: true},
	}, indexes)
}

type mockIndexInserter struct{}

func (mock *mockIndexInserter) EnsureIndex(mgo.Index) error {
	return errors.New("some random error")
}

/*
TestEnsureIndexes_Error tests that if it fails during index insertions, it panics
*/
func TestEnsureIndexes_Error(t *testing.T) {
	failingInserter := &mockIndexInserter{}
	assert.Panics(t, func() {
		ensureIndexes(failingInserter, []*mgo.Index{
			{Key: []string{"some key"}},
		})
	})
}
