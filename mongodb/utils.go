package mongodb

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func ConvertToBson(ent interface{}) (bson.M, error) {
	if ent == nil {
		return bson.M{}, nil
	}

	bytes, err := bson.Marshal(ent)
	if err != nil {
		return nil, err
	}

	obj := bson.M{}
	err = bson.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("InterfaceSlice() given none-slice type")
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result, nil
}
