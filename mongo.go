package gomongo

import (
	"context"
	"fmt"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB is a wrapper around a mongo connection
type MongoDB struct {
	connStr string
	dbName  string
	client  *mongo.Client
}

// NewMongoDB creates a new Mongo clent
func NewMongoDB(ctx context.Context, server string, user string, password string, port int, dbName string) (db *MongoDB, e error) {
	e = nil

	db = new(MongoDB)
	db.dbName = dbName
	db.connStr = "mongodb"
	if port == 0 {
		db.connStr += "+srv"
	}
	db.connStr += "://"

	if user != "" {
		db.connStr += url.PathEscape(user)

		if password != "" {
			db.connStr += ":" + url.PathEscape(password) + "@"
		}
	}

	db.connStr += server
	if port != 0 {
		db.connStr = fmt.Sprintf("%s:%d/", db.connStr, port)
	}

	clientOptions := options.Client().ApplyURI(db.connStr)

	db.client, e = mongo.Connect(ctx, clientOptions)

	if e == nil {
		e = db.client.Ping(ctx, nil)
	} else {
		fmt.Printf("error connecting to mongodb\n")
	}

	return
}

func (m *MongoDB) Database() (s string) {
	return m.dbName
}

func (m *MongoDB) Collection(collectionName string) (col *mongo.Collection) {
	col = m.client.Database(m.dbName).Collection(collectionName)

	return
}

func (m *MongoDB) ConnectionString() (s string) {
	return m.connStr
}

func (m *MongoDB) Client() (client *mongo.Client) {
	client = m.client

	return
}

func (m *MongoDB) Disconnect(ctx context.Context) (e error) {
	e = m.client.Disconnect(ctx)

	return
}
