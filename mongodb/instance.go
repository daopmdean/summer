package mongodb

import (
	"context"
	"reflect"
	"time"

	"github.com/daopmdean/summer/common"
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

func (m *Instance) Create(ent interface{}) *common.Response {
	if m.col == nil {
		return common.ResponseError("Mongodb err: Collection is nil " + m.ColName)
	}

	obj, err := ConvertToBson(ent)
	if err != nil {
		return common.ResponseError("Mongodb err: " + err.Error())
	}

	obj["created_time"] = time.Now()

	result, err := m.col.InsertOne(context.TODO(), obj)
	if err != nil {
		return common.ResponseError("Mongodb err: " + err.Error())
	}

	obj["_id"] = result.InsertedID
	ent, err = m.convertToObj(obj)
	if err != nil {
		return common.ResponseError("Mongodb err: " + err.Error())
	}

	slice := m.newObjectSlice(1)
	sliceValue := reflect.Append(
		reflect.ValueOf(slice),
		reflect.Indirect(reflect.ValueOf(ent)),
	)

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: "Create " + m.ColName + " success",
		Data:    sliceValue.Interface(),
	}
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
