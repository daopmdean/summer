package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/daopmdean/summer/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *Instance) Create(ctx context.Context, ent interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	obj, err := ConvertToBson(ent)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: " + err.Error())
	}

	obj["created_time"] = time.Now()

	result, err := m.col.InsertOne(ctx, obj)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: " + err.Error())
	}

	obj["_id"] = result.InsertedID
	ent, err = m.convertToObj(obj)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: " + err.Error())
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

func (m *Instance) CreateMany(ctx context.Context, ents interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	list, err := InterfaceSlice(ents)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid slice input")
	}

	now := time.Now()

	var bsonList []interface{}
	for _, item := range list {
		bi, err := ConvertToBson(item)
		if err != nil {
			return common.BuildMongoErr("Mongodb err: invalid bson object")
		}

		bi["created_time"] = now
		bsonList = append(bsonList, bi)
	}

	result, err := m.col.InsertMany(ctx, bsonList)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: create many failed")
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Create many %s (s) successfullly", m.ColName),
		Data:    result.InsertedIDs,
	}
}

func (m *Instance) Query(
	ctx context.Context,
	filter interface{},
	skip, limit int64,
	sortFields *bson.M,
) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	if skip < 0 {
		skip = 0
	}

	if limit <= 0 {
		limit = 50
	}

	opt := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	if sortFields != nil {
		opt.Sort = sortFields
	}

	converted, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	cur, err := m.col.Find(ctx, converted, opt)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: query failed with err" + err.Error())
	}
	if cur.Err() != nil {
		return common.BuildMongoErr("Mongodb err: query failed with cur err" + cur.Err().Error())
	}

	list := m.newObjectSlice(int(limit))
	err = cur.All(ctx, list)
	if err != nil {
		return &common.Response{
			Status:  common.ResponseStatus.NotFound,
			Message: fmt.Sprintf("Not found any match %s", m.ColName),
		}
	}

	err = cur.Close(ctx)
	if err != nil {
		return common.BuildMongoErr("Mongodb err close cursor: " + err.Error())
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Query %s success", m.ColName),
		Data:    list,
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
