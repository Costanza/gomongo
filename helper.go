package gomongo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewHelper[T any](ctx context.Context, client *MongoDB, collection string) (r *Helper[T]) {
	return &Helper[T]{
		client:     client,
		collection: collection,
	}
}

func (r Helper[T]) entityFieldToBSONTag(entity T, field string) (tag string, e error) {
	tag = FieldToBSONTag(entity, field)

	if tag == "" {
		e = fmt.Errorf("could not find struct field %s for entity %s", field, reflect.TypeOf(entity).Name())
	}

	return tag, e
}

func (h Helper[T]) CreateIndex(ctx context.Context, name string, field string, kind string, unique bool) (idxName string, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		model := mongo.IndexModel{

			Keys:    bson.D{{Key: tag, Value: kind}},
			Options: options.Index().SetName(name).SetUnique(unique),
		}

		idxName, e = col.Indexes().CreateOne(ctx, model)
	}

	return
}

func (h *Helper[T]) InsertOne(ctx context.Context, entity T) (id interface{}, e error) {
	col := h.client.Collection(h.collection)

	result, e := col.InsertOne(ctx, entity)

	if e == nil {
		id = result.InsertedID
	}

	return
}

func (h Helper[T]) InsertMany(ctx context.Context, entities []interface{}) (ids []interface{}, e error) {
	col := h.client.Collection(h.collection)

	result, e := col.InsertMany(ctx, entities)

	return result.InsertedIDs, e
}

func (h *Helper[T]) FindOne(ctx context.Context, field string, value interface{}) (result *T, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		e = col.FindOne(ctx, query).Decode(&entity)
	}

	return &entity, e
}

func (h Helper[T]) FindMany(ctx context.Context, field string, value interface{}, sortField string, sortDir SortDirection) (entities []T, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		opts := options.Find()
		if sortField != "" {
			sortTag, e := h.entityFieldToBSONTag(entity, field)
			if e == nil {
				opts.SetSort(bson.D{primitive.E{Key: sortTag, Value: sortDir}})
			}
		}

		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		var cursor *mongo.Cursor
		cursor, e = col.Find(ctx, query, opts)

		if e == nil {
			e = cursor.All(ctx, &entities)
		}
	}

	return
}

func (h Helper[T]) FindByDateRange(ctx context.Context, field string, start time.Time, end time.Time, sortDir SortDirection) (entities []T, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		opts := options.Find().SetSort(bson.D{primitive.E{Key: tag, Value: sortDir}})

		filter := bson.M{
			tag: bson.M{
				"$gte": start,
				"$lte": end,
			},
		}

		var cursor *mongo.Cursor
		cursor, e = col.Find(ctx, filter, opts)
		if e == nil {
			e = cursor.All(ctx, &entities)
		}
	}

	return
}

func (h *Helper[T]) SaveOne(ctx context.Context, entity *T, field string, value interface{}) (count int64, e error) {
	col := h.client.Collection(h.collection)

	var tag string
	tag, e = h.entityFieldToBSONTag(*entity, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}
		opts := options.Replace().SetUpsert(true)

		result, e := col.ReplaceOne(ctx, query, *entity, opts)

		return result.MatchedCount, e
	}

	return
}

func (h Helper[T]) DeleteOne(ctx context.Context, field string, value interface{}) (count int64, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		result, e := col.DeleteOne(ctx, query)

		return result.DeletedCount, e
	}

	return
}

func (h Helper[T]) DeleteMany(ctx context.Context, field string, value interface{}) (count int64, e error) {
	col := h.client.Collection(h.collection)

	var entity T
	var tag string
	tag, e = h.entityFieldToBSONTag(entity, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		result, e := col.DeleteMany(ctx, query)

		return result.DeletedCount, e
	}

	return
}

func (h Helper[T]) FindOneAndIncrementField(ctx context.Context, field string, value string, updateField string, increment int64) (entity *T, e error) {
	col := h.client.Collection(h.collection)

	var ent T
	var tag string
	tag, e = h.entityFieldToBSONTag(ent, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		var uf string
		uf, e = h.entityFieldToBSONTag(ent, updateField)

		if e == nil {
			update := bson.D{
				primitive.E{Key: "$inc", Value: bson.D{
					primitive.E{Key: uf, Value: increment},
				}},
			}

			opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
			e = col.FindOneAndUpdate(ctx, query, update, opts).Decode(&entity)
		}
	}

	return
}

func (h Helper[T]) FindOneAndUpdateField(ctx context.Context, field string, value string, updateField string, updateValue interface{}) (entity *T, e error) {
	col := h.client.Collection(h.collection)

	var ent T
	var tag string
	tag, e = h.entityFieldToBSONTag(ent, field)

	if e == nil {
		query := bson.D{
			primitive.E{Key: tag, Value: value},
		}

		var uf string
		uf, e = h.entityFieldToBSONTag(ent, updateField)

		if e == nil {
			update := bson.D{
				primitive.E{Key: "$inc", Value: bson.D{
					primitive.E{Key: uf, Value: updateValue},
				}},
			}

			opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
			e = col.FindOneAndUpdate(ctx, query, update, opts).Decode(&entity)
		}
	}

	return
}

func (h Helper[T]) TextSearch(ctx context.Context, term string) (entities []T, e error) {
	col := h.client.Collection(h.collection)

	opts := options.Find().SetSort(bson.D{{"score", bson.D{{"$meta", "textScore"}}}})

	query := bson.D{{"$text", bson.D{{"$search", term}}}}
	cursor, e := col.Find(ctx, query, opts)

	if e == nil {
		e = cursor.All(ctx, &entities)
	}

	return
}

func (h Helper[T]) Iterate(ctx context.Context, cb func(ctx context.Context, n *T) (e error)) (e error) {
	col := h.client.Collection(h.collection)

	var cursor *mongo.Cursor
	cursor, e = col.Find(ctx, bson.D{})
	if e == nil {
		for cursor.Next(ctx) {
			var entity T
			if e = cursor.Decode(&entity); e == nil {
				e = cb(ctx, &entity)
			}
		}

	}

	return
}
