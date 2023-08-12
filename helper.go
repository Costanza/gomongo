package gomongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SortDirection int

const (
	Ascending  SortDirection = 1
	Descending SortDirection = -1
)

type Helper[T any] struct {
	client     *MongoDB
	collection string
}

func NewHelper[T any](ctx context.Context, client *MongoDB, coll string) (r *Helper[T]) {
	return &Helper[T]{
		client:     client,
		collection: coll,
	}
}

func (h *Helper[T]) InsertOne(ctx context.Context, collection string, entity T) (id interface{}, e error) {
	col := h.client.Collection(collection)

	result, e := col.InsertOne(ctx, entity)

	if e == nil {
		id = result.InsertedID
	}

	return
}

func (h *Helper[T]) FindOne(ctx context.Context, collection string, field string, value interface{}) (result *T, e error) {
	col := h.client.Collection(collection)

	var entity T
	sf := FieldToBSONTag(entity, field)

	if field == "" {
		e = fmt.Errorf("could not find struct field %s for entity %s", field, reflect.TypeOf(entity).Name())
	} else {
		query := bson.D{
			primitive.E{Key: sf, Value: value},
		}

		e = col.FindOne(ctx, query).Decode(&entity)
	}

	return &entity, e
}

func (r Helper[T]) DeleteOne(ctx context.Context, field string, value interface{}) (count int64, e error) {
	col := r.client.Collection(r.collection)

	var entity T
	sf := FieldToBSONTag(entity, field)

	var result *mongo.DeleteResult
	if field == "" {
		e = fmt.Errorf("could not find struct field %s for entity %s", field, reflect.TypeOf(entity).Name())
	} else {
		query := bson.D{
			primitive.E{Key: sf, Value: value},
		}

		result, e = col.DeleteOne(ctx, query)
	}

	return result.DeletedCount, e
}

// func (r Helper[T]) Get(ctx context.Context, key string, value interface{}, sortField string, sortDir SortDirection) (entities []T, e error) {
// 	col := r.client.Collection(r.collection)
// 	opts := options.Find()
// 	if sortField != "" {
// 		opts.SetSort(bson.D{{sortField, sortDir}})
// 	}

// 	filter := bson.D{
// 		primitive.E{Key: key, Value: value},
// 	}

// 	cursor, e := col.Find(ctx, filter, opts)

// 	if e == nil {
// 		e = cursor.All(ctx, &entities)
// 	}

// 	return
// }

// func (r Helper[T]) GetOne(ctx context.Context, key string, value interface{}) (entity *T, e error) {
// 	col := r.client.Collection(r.collection)

// 	filter := bson.D{
// 		primitive.E{Key: key, Value: value},
// 	}

// 	e = col.FindOne(ctx, filter).Decode(&entity)

// 	return
// }

// func (r Helper[T]) GetByDateRange(ctx context.Context, key string, start time.Time, end time.Time, sortDir SortDirection) (entities []T, e error) {
// 	col := r.client.Collection(r.collection)
// 	opts := options.Find().SetSort(bson.D{{key, sortDir}})

// 	filter := bson.M{
// 		key: bson.M{
// 			"$gte": start,
// 			"$lte": end,
// 		},
// 	}

// 	cursor, e := col.Find(ctx, filter, opts)
// 	if e == nil {
// 		e = cursor.All(ctx, &entities)
// 	}

// 	return
// }

// func (r Helper[T]) TextSearch(ctx context.Context, term string) (entities []T, e error) {
// 	col := r.client.Collection(r.collection)
// 	opts := options.Find().SetSort(bson.D{{"score", bson.D{{"$meta", "textScore"}}}})

// 	filter := bson.D{{"$text", bson.D{{"$search", term}}}}
// 	cursor, e := col.Find(ctx, filter, opts)

// 	if e == nil {
// 		e = cursor.All(ctx, &entities)
// 	}

// 	return
// }

// func (r Helper[T]) FindOneAndIncrementField(ctx context.Context, key string, value string, updateKey string, increment int64) (entity *T, e error) {
// 	col := r.client.Collection(r.collection)

// 	filter := bson.D{
// 		primitive.E{Key: key, Value: value},
// 	}

// 	update := bson.D{
// 		primitive.E{Key: "$inc", Value: bson.D{
// 			primitive.E{Key: updateKey, Value: increment},
// 		}},
// 	}

// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	opts.SetUpsert(true)

// 	e = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)

// 	return
// }

// func (r Helper[T]) FindOneAndUpdateField(ctx context.Context, key string, value string, updateKey string, updateValue interface{}) (entity *T, e error) {
// 	col := r.client.Collection(r.collection)

// 	filter := bson.D{
// 		primitive.E{Key: key, Value: value},
// 	}

// 	update := bson.D{
// 		primitive.E{Key: "$set", Value: bson.D{
// 			primitive.E{Key: updateKey, Value: updateValue},
// 		}},
// 	}

// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	opts.SetUpsert(true)

// 	e = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)

// 	return
// }

// func (r Helper[T]) Iterate(ctx context.Context, cb func(ctx context.Context, n *T) (e error)) (e error) {
// 	col := r.client.Collection(r.collection)

// 	cursor, err := col.Find(ctx, bson.D{})
// 	if err != nil {
// 		e = err
// 	} else {
// 		for cursor.Next(ctx) {
// 			var entity T
// 			if err := cursor.Decode(&entity); err != nil {
// 				e = err
// 			} else {
// 				e = cb(ctx, &entity)
// 			}

// 		}
// 		if err := cursor.Err(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	return
// }

// func (r Helper[T]) InsertOne(ctx context.Context, entity T) (id interface{}, e error) {
// 	col := r.client.Collection(r.collection)

// 	result, e := col.InsertOne(ctx, entity)

// 	return result.InsertedID, e
// }

// func (r Helper[T]) InsertMany(ctx context.Context, entities []interface{}) (ids []interface{}, e error) {
// 	col := r.client.Collection(r.collection)

// 	result, e := col.InsertMany(ctx, entities)

// 	return result.InsertedIDs, e
// }

// func (r Helper[T]) Save(ctx context.Context, entity *T, key string, value interface{}) (e error) {
// 	col := r.client.Collection(r.collection)

// 	query := bson.M{
// 		key: value,
// 	}

// 	opts := options.Replace().SetUpsert(true)

// 	_, e = col.ReplaceOne(ctx, query, *entity, opts)

// 	return
// }

// func (r Helper[T]) DeleteMany(ctx context.Context, key string, value interface{}) (count int64, e error) {
// 	col := r.client.Collection(r.collection)

// 	query := bson.M{
// 		key: value,
// 	}

// 	result, e := col.DeleteMany(ctx, query)

// 	return result.DeletedCount, e
// }

// func (r Helper[T]) CreateIndex(ctx context.Context, name string, field string, kind string, unique bool) (idxName string, e error) {
// 	col := r.client.Collection(r.collection)

// 	model := mongo.IndexModel{

// 		Keys:    bson.D{{Key: field, Value: kind}},
// 		Options: options.Index().SetName(name).SetUnique(unique),
// 	}

// 	idxName, e = col.Indexes().CreateOne(ctx, model)

// 	return
// }
