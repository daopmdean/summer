package mongodb

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Instance struct {
	ColName     string
	DBName      string
	TemplateObj interface{}

	db  *mongo.Database
	col *mongo.Collection
}

func (m *Instance) SetDB(database *mongo.Database) {
	m.db = database
	m.col = m.db.Collection(m.ColName)
	m.DBName = database.Name()
}

func (m *Instance) convertToObj(b bson.M) (interface{}, error) {
	obj := m.newObject()
	if b == nil {
		return obj, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bytes, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (m *Instance) newObject() interface{} {
	t := reflect.TypeOf(m.TemplateObj)
	v := reflect.New(t)
	return v.Interface()
}

func (m *Instance) newObjectSlice(limit int) interface{} {
	t := reflect.TypeOf(m.TemplateObj)
	v := reflect.MakeSlice(reflect.SliceOf(t), 0, limit)
	return v.Interface()
}
