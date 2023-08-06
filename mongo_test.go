package gomongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestData struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Seq  int64              `json:"seq" bson:"seq"`
}

func TestNewMongoDB(t *testing.T) {
	ctx := context.Background()
	server := "localhost"
	user := ""
	password := ""
	port := 27017

	db, e := NewMongoDB(ctx, server, user, password, port)
	if e == nil {
		connStr := "mongodb://localhost:27017/"
		if db.ConnectionString() != connStr {
			t.Errorf("incorrect connection string want %s have: %s", connStr, db.ConnectionString())
		}
	} else {
		t.Errorf("error creating mongodb: %v", e)
	}
}

func TestNewMongoDBWithSrv(t *testing.T) {
	ctx := context.Background()
	server := "test.mongodb.net"
	user := "admin"
	password := "123456"
	port := 0

	db, e := NewMongoDB(ctx, server, user, password, port)
	if e != nil {
		if e.Error() == "error parsing uri: lookup _mongodb._tcp.test.mongodb.net: dnsquery: DNS name does not exist." {
			connStr := "mongodb+srv://admin:123456@test.mongodb.net"
			if db.ConnectionString() != connStr {
				t.Errorf("incorrect connection string want %s have: %s", connStr, db.ConnectionString())
			}
		} else {
			t.Errorf("unexpected error creating mongodb: %v", e)
		}
	} else {
		t.Errorf("should have returned DNS error and didn't")
	}
}

func TestFindOneAndIncrement(t *testing.T) {

	// Should move to mocks.
	ctx := context.Background()
	server := "localhost"
	user := ""
	password := ""
	port := 27017

	repo, e := NewRepository[TestData](ctx, server, user, password, port, "MongoTest", "testdata")
	if e == nil {
		td, e := repo.FindOneAndIncrementField(ctx, "name", "test", "seq", 1)
		if e == nil {
			if td.Seq != 1 {
				t.Errorf("findandincrement have: %d want: %d", td.Seq, 1)
			} else {
				td, e := repo.FindOneAndIncrementField(ctx, "name", "test", "seq", 5)
				if e == nil {
					if td.Seq != 6 {
						t.Errorf("findandincrement have: %d want: %d", td.Seq, 6)
					}
				} else {
					t.Errorf("error finding and incrementing: %v", e)
				}
			}

			count, e := repo.DeleteOne(ctx, "name", td.Name)
			if e != nil || count != 1 {
				t.Errorf("error deleting test data: %v", e)
			}
		} else {
			t.Errorf("error finding and incrementing: %v", e)
		}
	} else {
		t.Errorf("error creating mongodb: %v", e)
	}
}
