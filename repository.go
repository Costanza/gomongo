package gomongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository[T any] struct {
	dbClient   *MongoDB
	dbName     string
	collection string
}

func NewRepository[T any](ctx context.Context, server string, user string, password string, port int, dbName string, coll string) (r *Repository[T], e error) {
	r = new(Repository[T])

	r.dbClient, e = NewMongoDB(ctx, server, user, password, port)
	r.dbName = dbName
	r.collection = coll

	return
}

func (r Repository[T]) GetByID(ctx context.Context, ID primitive.ObjectID) (entity *T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.D{
		primitive.E{Key: "_id", Value: ID},
	}

	e = col.FindOne(ctx, filter).Decode(&entity)

	return
}

func (r Repository[T]) GetOne(ctx context.Context, key string, value interface{}) (entity *T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.D{
		primitive.E{Key: key, Value: value},
	}

	e = col.FindOne(ctx, filter).Decode(&entity)

	return
}

func (r Repository[T]) GetByDateRange(ctx context.Context, key string, start time.Time, end time.Time) (entities []T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.M{
		key: bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, e := col.Find(ctx, filter)
	if e == nil {
		e = cursor.All(ctx, &entities)
	}

	return
}

func (r Repository[T]) Search(ctx context.Context, term string) (entities []T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.D{{"$text", bson.D{{"$search", term}}}}
	cursor, e := col.Find(ctx, filter)

	if e == nil {
		e = cursor.All(ctx, &entities)
	}

	return
}

func (r Repository[T]) FindOneAndIncrementField(ctx context.Context, key string, value string, updateKey string, increment int64) (entity *T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.D{
		primitive.E{Key: key, Value: value},
	}

	update := bson.D{
		primitive.E{Key: "$inc", Value: bson.D{
			primitive.E{Key: updateKey, Value: increment},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	opts.SetUpsert(true)

	e = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)

	return
}

func (r Repository[T]) FindOneAndUpdateField(ctx context.Context, key string, value string, updateKey string, updateValue interface{}) (entity *T, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	filter := bson.D{
		primitive.E{Key: key, Value: value},
	}

	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: updateKey, Value: updateValue},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	opts.SetUpsert(true)

	e = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)

	return
}

func (r Repository[T]) Iterate(ctx context.Context, cb func(ctx context.Context, n *T) (e error)) (e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		e = err
	} else {
		for cursor.Next(ctx) {
			var entity T
			if err := cursor.Decode(&entity); err != nil {
				e = err
			} else {
				e = cb(ctx, &entity)
			}

		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return
}

func (r Repository[T]) Save(ctx context.Context, entity *T, key string, value interface{}) (e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	query := bson.M{
		key: value,
	}

	opts := options.Replace().SetUpsert(true)

	_, e = col.ReplaceOne(ctx, query, *entity, opts)

	return
}

func (r Repository[T]) DeleteOne(ctx context.Context, key string, value interface{}) (count int64, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	query := bson.M{
		key: value,
	}

	result, e := col.DeleteOne(ctx, query)

	return result.DeletedCount, e
}

func (r Repository[T]) DeleteMany(ctx context.Context, key string, value interface{}) (count int64, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	query := bson.M{
		key: value,
	}

	result, e := col.DeleteMany(ctx, query)

	return result.DeletedCount, e
}

func (r Repository[T]) CreateIndex(ctx context.Context, name string, field string, kind string, unique bool) (idxName string, e error) {
	col := r.dbClient.Collection(r.dbName, r.collection)

	model := mongo.IndexModel{

		Keys:    bson.D{{Key: field, Value: kind}},
		Options: options.Index().SetName(name).SetUnique(unique),
	}

	idxName, e = col.Indexes().CreateOne(ctx, model)

	return
}
