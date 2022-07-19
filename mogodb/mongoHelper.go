package mogodb

import (
	"MyGameServer/logger"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoHelper struct {
	Options        *options.ClientOptions
	Client         *mongo.Client
	DBName         string
	CollectionName string
	IsConnected    bool
	Url            string
	Collection     *mongo.Collection
}

const TIME_OUT = 10

// Connect 连接数据库
func (mh *MongoHelper) Connect(url string) error {
	if mh.IsConnected {
		return errors.New("[MongoDB] IsConnected Not Again Connect")
	}
	mh.Url = url
	mh.Options = options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), mh.Options)
	if err != nil {
		logger.PopErrorInfo("[MongoDB] Connect Error:")
		logger.PopError(err)
		return err
	}
	mh.Client = client
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("[MongoDB] Connect Ping Error:", err)
		return err
	}
	fmt.Println("[MongoDB] Connect Success!")
	mh.IsConnected = true
	mh.Collection = mh.Client.Database(mh.DBName).Collection(mh.CollectionName)
	return nil
}

// FindAll 从数据库中查找所有的数据
func (mh *MongoHelper) FindAll() ([]bson.M, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	cur, err := mh.Collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("[MongoDB] FindAll Error:", err)
		return nil, err
	}
	defer cur.Close(ctx)
	retValue := make([]bson.M, 0, 50)
	for cur.Next(ctx) {
		var result bson.D
		err = cur.Decode(&result)
		if err != nil {
			fmt.Println("[MongoDB] Decode Error:", err)
			return nil, err
		}
		retValue = append(retValue, result.Map())
	}
	return retValue, nil
}

// FindMany 从数据库中查找符合条件的一些数据
func (mh *MongoHelper) FindMany(filter bson.D) ([]bson.M, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	cur, err := mh.Collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("[MongoDB] Find Error:", err)
		return nil, err
	}
	defer cur.Close(ctx)
	retValue := make([]bson.M, 0, 50)
	for cur.Next(ctx) {
		var result bson.D
		err = cur.Decode(&result)
		if err != nil {
			fmt.Println("[MongoDB] Decode Error:", err)
			return nil, err
		}
		retValue = append(retValue, result.Map())
	}
	return retValue, nil
}

// InsertOne 向数据库中插入一条数据
func (mh *MongoHelper) InsertOne(document any) (primitive.ObjectID, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return primitive.ObjectID{}, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	result, err := mh.Collection.InsertOne(ctx, document)
	if err != nil {
		fmt.Println("[MongoDB] InsertOne Error: ", err)
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// InsertMany 向数据库中插入多条数据
func (mh *MongoHelper) InsertMany(document []any) ([]primitive.ObjectID, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	result, err := mh.Collection.InsertMany(ctx, document)
	if err != nil {
		fmt.Println("[MongoDB] InsertMany Error: ", err)
		return nil, err
	}
	slice := make([]primitive.ObjectID, 0, len(result.InsertedIDs))
	for i := 0; i < len(result.InsertedIDs); i++ {
		slice = append(slice, result.InsertedIDs[i].(primitive.ObjectID))
	}
	return slice, nil
}

// UpdateByID 根据ID更新数据库中的一条数据
func (mh *MongoHelper) UpdateByID(id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	result, err := mh.Collection.UpdateByID(ctx, id, update)
	if err != nil {
		fmt.Println("[MongoDB] UpdateByID Error: ", err)
		return nil, err
	}
	return result, nil
}

// UpdateOne 更新数据库中的一条数据
func (mh *MongoHelper) UpdateOne(filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	result, err := mh.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MongoDB] UpdateOne Error: ", err)
		return nil, err
	}
	return result, nil
}

// UpdateMany 更新数据库中的多条数据
func (mh *MongoHelper) UpdateMany(filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	result, err := mh.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println("[MongoDB] UpdateMany Error: ", err)
		return nil, err
	}
	return result, nil
}

// DeleteOne 删除数据库中的一条数据
func (mh *MongoHelper) DeleteOne(filter bson.D, isDelete bool) (*mongo.DeleteResult, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	if isDelete {
		ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
		defer cancel()
		result, err := mh.Collection.DeleteOne(ctx, filter)
		if err != nil {
			fmt.Println("[MongoDB] DeleteOne Error: ", err)
			return nil, err
		}
		return result, nil
	} else {
		update := bson.D{}
		result, err := mh.UpdateOne(filter, update)
		if err != nil {
			fmt.Println("[MongoDB] DeleteOne IsDelete Error: ", err)
			return nil, err
		}
		return &mongo.DeleteResult{DeletedCount: result.ModifiedCount}, nil
	}
}

// DeleteMany 删除数据库中的多条数据
func (mh *MongoHelper) DeleteMany(filter bson.D, isDelete bool) (*mongo.DeleteResult, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	if isDelete {
		ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
		defer cancel()
		result, err := mh.Collection.DeleteMany(ctx, filter)
		if err != nil {
			fmt.Println("[MongoDB] DeleteMany Error: ", err)
			return nil, err
		}
		return result, nil
	} else {
		update := bson.D{}
		result, err := mh.UpdateMany(filter, update)
		if err != nil {
			fmt.Println("[MongoDB] DeleteMany IsDelete Error: ", err)
			return nil, err
		}
		return &mongo.DeleteResult{DeletedCount: result.ModifiedCount}, nil
	}
}

func (mh *MongoHelper) TestFind() ([]bson.M, error) {
	if !mh.IsConnected {
		connectError := mh.Connect(mh.Url)
		if connectError != nil {
			return nil, connectError
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT*time.Second)
	defer cancel()
	if mh.Collection == nil {
		mh.Collection = mh.Client.Database(mh.DBName).Collection(mh.CollectionName)
	}
	cur, err := mh.Collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("[MongoDB] FindAll Error:", err)
		return nil, err
	}
	defer cur.Close(ctx)
	retValue := make([]bson.M, 0, 50)
	for cur.Next(ctx) {
		var result bson.D
		err = cur.Decode(&result)
		if err != nil {
			fmt.Println("[MongoDB] Decode Error:", err)
			return nil, err
		}
		retValue = append(retValue, result.Map())
	}
	return retValue, nil
}

func NewMongoHelper(dbName, collectionName string) *MongoHelper {
	return &MongoHelper{
		DBName:         dbName,
		CollectionName: collectionName,
		IsConnected:    false,
	}
}
