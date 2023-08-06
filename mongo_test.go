package gomongo

import (
	"context"
	"testing"
)

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
