package models

import (
	"github.com/sirupsen/logrus"
	"github.com/transcovo/go-chpr-logger"
	"gopkg.in/mgo.v2"
)

/*
BaseModelInterface interface  defines the methods all models must implement.
*/
type BaseModelInterface interface {
	EnsureIndexes(bool)
	GetCollection() *mgo.Collection
}

/*
BaseModel struct declares the required methods to implement the BaseModelInterface interface.
*/
type BaseModel struct {
	CollectionName string
	DB             *mgo.Database
	Indexes        []*mgo.Index
}

/*
EnsureIndexes creates the model indexes in the database.
*/
func (model *BaseModel) EnsureIndexes(background bool) {
	for _, index := range model.Indexes {
		index.Background = background
	}
	ensureIndexes(model.GetCollection(), model.Indexes)
}

/*
GetCollection returns the model collection.
*/
func (model *BaseModel) GetCollection() *mgo.Collection {
	return model.DB.C(model.CollectionName)
}

type indexInserter interface {
	EnsureIndex(mgo.Index) error
}

/*
ensureIndexes sets in mongo a list of indexes in a given collection.
*/
func ensureIndexes(dbCollection indexInserter, indexes []*mgo.Index) {
	for _, index := range indexes {
		err := dbCollection.EnsureIndex(*index)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":   err,
				"index": index,
			}).Error("[models.ensureIndexes] Error creating index")
			panic(err)
		}
	}
}

/*
EnsureIndexes sets in mongo all the indexes for all the given models.
*/
func EnsureIndexes(models []BaseModelInterface, background bool) {
	for _, model := range models {
		model.EnsureIndexes(background)
	}
}
