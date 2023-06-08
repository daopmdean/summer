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

func (m *Instance) GetClient() *mongo.Client {
	return m.db.Client()
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

	if obj["created_time"] == nil {
		obj["created_time"] = time.Now()
	}

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

func (m *Instance) QueryWithOpt(
	ctx context.Context,
	filter interface{},
	opt *options.FindOptions,
) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
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

	list := m.newObjectSlice(cur.RemainingBatchLength())
	err = cur.All(ctx, &list)
	if err != nil {
		return common.BuildMongoErr("Mongodb err decode cursor: " + err.Error())
	}

	err = cur.Close(ctx)
	if err != nil {
		return common.BuildMongoErr("Mongodb err close cursor: " + err.Error())
	}

	if reflect.ValueOf(list).Len() == 0 {
		return common.BuildQueryNotFound("Not found any match " + m.ColName)
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Query %s success", m.ColName),
		Data:    list,
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

	list := m.newObjectSlice(cur.RemainingBatchLength())
	err = cur.All(ctx, &list)
	if err != nil {
		return &common.Response{
			Status: common.ResponseStatus.NotFound,
			Error: &common.ErrorResponse{
				ErrorCode:    "COLLECTION_NOT_FOUND",
				ErrorMessage: fmt.Sprintf("Not found any match %s", m.ColName),
			},
		}
	}

	err = cur.Close(ctx)
	if err != nil {
		return common.BuildMongoErr("Mongodb err close cursor: " + err.Error())
	}

	if reflect.ValueOf(list).Len() == 0 {
		return common.BuildQueryNotFound("Not found any match " + m.ColName)
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Query %s success", m.ColName),
		Data:    list,
	}
}

func (m *Instance) QueryOne(ctx context.Context, filter interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	converted, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	result := m.col.FindOne(ctx, converted)
	if result.Err() == mongo.ErrNoDocuments {
		return common.BuildQueryNotFound("Not found any match " + m.ColName)
	} else if result.Err() != nil {
		return common.BuildMongoErr("Mongodb err: query failed with err " + result.Err().Error())
	}

	return m.parseSingleResult(result, "query one")
}

func (m *Instance) Update(ctx context.Context,
	filter interface{}, updater interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	convertedFilter, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	convertedUpdater, err := ConvertToBson(updater)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid updater input")
	}
	delete(convertedUpdater, "created_time")
	convertedUpdater["last_update_time"] = time.Now()

	result, err := m.col.UpdateMany(ctx, convertedFilter, bson.M{
		"$set": convertedUpdater,
	})
	if err != nil {
		return common.BuildMongoErr("Mongodb err: update failed with err" + err.Error())
	}

	if result.MatchedCount == 0 {
		return &common.Response{
			Status:  common.ResponseStatus.NotFound,
			Message: fmt.Sprintf("Not found any match %s", m.ColName),
		}
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Update %s success", m.ColName),
		Data:    []int64{result.MatchedCount},
	}
}

func (m *Instance) UpdateOne(ctx context.Context,
	filter interface{}, updater interface{},
	opts ...*options.FindOneAndUpdateOptions) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	convertedFilter, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	convertedUpdater, err := ConvertToBson(updater)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid updater input")
	}
	delete(convertedUpdater, "created_time")
	convertedUpdater["last_update_time"] = time.Now()

	if len(opts) == 0 {
		after := options.After
		opts = append(opts, &options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		})
	}

	result := m.col.FindOneAndUpdate(ctx, convertedFilter, bson.M{
		"$set": convertedUpdater,
	}, opts...)
	if result.Err() != nil {
		return common.BuildMongoErr("Mongodb err: update failed with err" + result.Err().Error())
	}

	return m.parseSingleResult(result, "update one")
}

func (m *Instance) Upsert(ctx context.Context,
	filter interface{}, updater interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	convertedFilter, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	convertedUpdater, err := ConvertToBson(updater)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid updater input")
	}
	delete(convertedUpdater, "created_time")
	convertedUpdater["last_update_time"] = time.Now()

	result := m.col.FindOneAndUpdate(ctx, convertedFilter, bson.M{
		"$set": convertedUpdater,
		"$setOnInsert": bson.M{
			"created_time": time.Now(),
		},
	}, &options.FindOneAndUpdateOptions{
		Upsert:         &[]bool{true}[0],
		ReturnDocument: &[]options.ReturnDocument{options.After}[0],
	})
	if result.Err() != nil {
		return common.BuildMongoErr("Mongodb err: upsert failed with err" + result.Err().Error())
	}

	return m.parseSingleResult(result, "upsert")
}

func (m *Instance) Delete(ctx context.Context, filter interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	converted, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	if len(converted) == 0 {
		return common.BuildMongoErr("Mongodb err: empty filter, delete all is not allowed")
	}

	result, err := m.col.DeleteMany(ctx, converted)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: delete failed with err" + err.Error())
	}

	if result.DeletedCount == 0 {
		return &common.Response{
			Status:  common.ResponseStatus.NotFound,
			Message: fmt.Sprintf("Not found any match %s", m.ColName),
		}
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Delete %s success", m.ColName),
		Data:    []int64{result.DeletedCount},
	}
}

func (m *Instance) Count(ctx context.Context, filter interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	converted, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	if len(converted) == 0 {
		return m.DocCount(ctx)
	}

	count, err := m.col.CountDocuments(ctx, converted)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: count failed with err" + err.Error())
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Count %s success", m.ColName),
		Data:    []int64{count},
	}
}

func (m *Instance) DocCount(ctx context.Context) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	count, err := m.col.EstimatedDocumentCount(ctx)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: count failed with err" + err.Error())
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Count %s success", m.ColName),
		Data:    []int64{count},
	}
}

func (m *Instance) IncreValue(ctx context.Context, filter interface{},
	fieldName string, value int) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	convertedFilter, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	if len(fieldName) == 0 {
		return common.BuildMongoErr("Mongodb err: empty field name")
	}

	updater := bson.M{
		"$inc": bson.D{
			{Key: fieldName, Value: value},
		},
	}

	result := m.col.FindOneAndUpdate(ctx, convertedFilter, updater,
		&options.FindOneAndUpdateOptions{
			Upsert:         &[]bool{true}[0],
			ReturnDocument: &[]options.ReturnDocument{options.After}[0],
		},
	)
	if result.Err() != nil {
		return common.BuildMongoErr("Mongodb err: update failed with err" + result.Err().Error())
	}

	return m.parseSingleResult(result, "incre value")
}

func (m *Instance) Distinct(ctx context.Context, fieldName string,
	filter interface{}, opts ...*options.DistinctOptions) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	convertedFilter, err := ConvertToBson(filter)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: invalid filter input")
	}

	result, err := m.col.Distinct(ctx, fieldName, convertedFilter, opts...)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: distinct failed with err" + err.Error())
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Distinct %s success", m.ColName),
		Data:    result,
	}
}

func (m *Instance) Aggregate(ctx context.Context,
	pipeline interface{}, result interface{}) *common.Response {
	if m.col == nil {
		return common.BuildMongoErr("Mongodb err: Collection is nil " + m.ColName)
	}

	cur, err := m.col.Aggregate(ctx, pipeline)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: aggregate failed with err" + err.Error())
	}

	err = cur.All(ctx, result)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: aggregate failed with err" + err.Error())
	}

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("Aggregate %s success", m.ColName),
	}
}

func (m *Instance) parseSingleResult(result *mongo.SingleResult, action string) *common.Response {
	obj := m.newObject()
	err := result.Decode(obj)
	if err != nil {
		return common.BuildMongoErr("Mongodb err: " + action + " failed with err" + err.Error())
	}

	list := m.newObjectSlice(1)
	listValue := reflect.Append(reflect.ValueOf(list),
		reflect.Indirect(reflect.ValueOf(obj)))

	return &common.Response{
		Status:  common.ResponseStatus.Success,
		Message: fmt.Sprintf("%s %s success", action, m.ColName),
		Data:    listValue.Interface(),
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
	return reflect.
		New(t).
		Interface()
}

func (m *Instance) newObjectSlice(limit int) interface{} {
	t := reflect.TypeOf(m.TemplateObj)
	return reflect.
		MakeSlice(reflect.SliceOf(t), 0, limit).
		Interface()
}
