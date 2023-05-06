package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strconv"
	"time"
)

type MilliTime time.Time

// UnmarshalBSONValue UnmarshalBSON bson转go对象
func (th *MilliTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	milli, _, ok := bsoncore.ReadDateTime(data)
	if ok {
		*th = MilliTime(time.UnixMilli(milli))
	}

	return nil
}

// MarshalBSONValue go对象转bson
func (th MilliTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(th))
}

func (th *MilliTime) UnmarshalJSON(data []byte) error {
	milli, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*th = MilliTime(time.UnixMilli(milli))
	return nil
}

func (th MilliTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(th).UnixMilli())
}

func NewMilliTime() MilliTime {
	return MilliTime(time.Now())
}

func MilliTimeFromTime(time time.Time) MilliTime {
	return MilliTime(time)
}

func NewMilliTimeFromTimePtr(time *time.Time) *MilliTime {
	return (*MilliTime)(time)
}

//func (th MilliTime) Add(t time.Duration) MilliTime {
//	return MilliTime(time.Time(th).Add(t))
//}
