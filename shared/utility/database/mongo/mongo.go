package mongo

import (
	context2 "context"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMainMongo(context ctx.BackgroundContext) *mongo.Client {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		mongoDial := configuration.GetConfig(context).Database.MongoDial

		c, err := mongo.NewClient(options.Client().ApplyURI(mongoDial))
		if err != nil {
			return c, err
		}
		err = c.Connect(context2.Background())
		return c, err
	}).(*mongo.Client)
}

func GetMainMongoDatabase(context ctx.BackgroundContext) *mongo.Database {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		mongoDb := configuration.GetConfig(context).Database.MongoDb

		return GetMainMongo(context).Database(mongoDb), nil
	}).(*mongo.Database)
}

func GetMainMongoGridFs(context ctx.BackgroundContext) *gridfs.Bucket {
	type mainMongoGridFs struct{}
	return context.Persist(mainMongoGridFs{}, func() (interface{}, error) {
		return gridfs.NewBucket(GetMainMongoDatabase(context))
	}).(*gridfs.Bucket)
}

func MustObjectIdFromHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return id
}
