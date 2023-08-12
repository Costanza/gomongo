package gomongo

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestData struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Data string             `json:"data" bson:"data"`
	Seq  int64              `json:"seq" bson:"seq"`
}

func createTestDB(ctx context.Context) (db *MongoDB, e error) {
	server := "localhost"
	user := ""
	password := ""
	port := 27017
	name := "MongoTest"

	return NewMongoDB(ctx, server, user, password, port, name)
}

// This all really should be handled by mocks, but local test db is easier for now.
func TestNewMongoDB(t *testing.T) {
	ctx := context.Background()
	db, e := createTestDB(ctx)

	if e != nil {
		t.Errorf("error creating mongodb: %v", e)
	} else {
		connStr := "mongodb://localhost:27017/"
		if db.ConnectionString() != connStr {
			t.Errorf("incorrect connection string want %s have: %s", connStr, db.ConnectionString())
		}
	}
}

func TestFindOneByFieldName(t *testing.T) {
	ctx := context.Background()
	db, e := createTestDB(ctx)

	helper := NewHelper[TestData](context.Background(), db, "testdata")
	if e != nil {
		t.Errorf("error connecting to mongo: %v", e)
	} else {
		td1 := TestData{
			ID:   primitive.NewObjectID(),
			Name: "test1",
			Data: "test that thing",
			Seq:  1,
		}

		_, e = helper.InsertOne(ctx, "testdata", td1)
		if e != nil {
			t.Errorf("error inserting test data: %v", e)
		} else {
			result, e := helper.FindOne(ctx, "testdata", "Name", "test1")

			if e != nil {
				t.Errorf("error finding test data: %v", e)
			} else {
				fmt.Printf("entity: %v\n", result)
				if result.Name != "test1" {
					t.Errorf("findone have: %s want: %s", result.Name, "test1")
				}
			}

			count, e := helper.DeleteOne(ctx, "Name", result.Name)
			if e != nil || count != 1 {
				t.Errorf("error deleting test data: %v", e)
			}
		}
	}
}

// func TestNewMongoDBWithSrv(t *testing.T) {
// 	ctx := context.Background()
// 	server := "test.mongodb.net"
// 	user := "admin"
// 	password := "123456"
// 	port := 0

// 	db, e := NewMongoDB(ctx, server, user, password, port, "MongoTest")
// 	if e != nil {
// 		if e.Error() == "error parsing uri: lookup _mongodb._tcp.test.mongodb.net: dnsquery: DNS name does not exist." {
// 			connStr := "mongodb+srv://admin:123456@test.mongodb.net"
// 			if db.ConnectionString() != connStr {
// 				t.Errorf("incorrect connection string want %s have: %s", connStr, db.ConnectionString())
// 			}
// 		} else {
// 			t.Errorf("unexpected error creating mongodb: %v", e)
// 		}
// 	} else {
// 		t.Errorf("should have returned DNS error and didn't")
// 	}
// }

// func TestFindOneAndIncrement(t *testing.T) {

// 	// Should move to mocks.
// 	ctx := context.Background()
// 	server := "localhost"
// 	user := ""
// 	password := ""
// 	port := 27017

// 	repo, e := NewRepository[TestData](ctx, server, user, password, port, "MongoTest", "testdata")
// 	if e == nil {
// 		td, e := repo.FindOneAndIncrementField(ctx, "name", "test", "seq", 1)
// 		if e == nil {
// 			if td.Seq != 1 {
// 				t.Errorf("findandincrement have: %d want: %d", td.Seq, 1)
// 			} else {
// 				td, e := repo.FindOneAndIncrementField(ctx, "name", "test", "seq", 5)
// 				if e == nil {
// 					if td.Seq != 6 {
// 						t.Errorf("findandincrement have: %d want: %d", td.Seq, 6)
// 					}
// 				} else {
// 					t.Errorf("error finding and incrementing: %v", e)
// 				}
// 			}

// 			count, e := repo.DeleteOne(ctx, "name", td.Name)
// 			if e != nil || count != 1 {
// 				t.Errorf("error deleting test data: %v", e)
// 			}
// 		} else {
// 			t.Errorf("error finding and incrementing: %v", e)
// 		}
// 	} else {
// 		t.Errorf("error creating mongodb: %v", e)
// 	}
// }

// func TestFindOneAndUpdate(t *testing.T) {

// 	// Should move to mocks.
// 	ctx := context.Background()
// 	server := "localhost"
// 	user := ""
// 	password := ""
// 	port := 27017

// 	repo, e := NewRepository[TestData](ctx, server, user, password, port, "MongoTest", "testdata")
// 	if e == nil {
// 		td, e := repo.FindOneAndUpdateField(ctx, "name", "test", "seq", 25)
// 		if e == nil {
// 			if td.Seq != 25 {
// 				t.Errorf("findandupdate have: %d want: %d", td.Seq, 25)
// 			}

// 			count, e := repo.DeleteOne(ctx, "name", td.Name)
// 			if e != nil || count != 1 {
// 				t.Errorf("error deleting test data: %v", e)
// 			}
// 		}
// 	} else {
// 		t.Errorf("error creating mongodb: %v", e)
// 	}
// }

// func TestCreateIndex(t *testing.T) {

// 	// Should move to mocks.
// 	ctx := context.Background()
// 	server := "localhost"
// 	user := ""
// 	password := ""
// 	port := 27017

// 	repo, e := NewRepository[TestData](ctx, server, user, password, port, "MongoTest", "testdata")
// 	if e == nil {
// 		td, e := repo.FindOneAndUpdateField(ctx, "name", "test", "data", "helloworld")
// 		if e == nil {
// 			if td.Data != "helloworld" {
// 				t.Errorf("createindex have: %s want: %s", td.Data, "helloworld")
// 			}
// 		}

// 		idxName, e := repo.CreateIndex(ctx, "my_text_index", "data", "text", false)
// 		if e == nil {
// 			if idxName != "my_text_index" {
// 				t.Errorf("createindex have: %s want: %s", idxName, "my_text_index")
// 			}
// 		} else {
// 			t.Errorf("error creating index: %v", e)
// 		}

// 		count, e := repo.DeleteOne(ctx, "name", td.Name)
// 		if e != nil || count != 1 {
// 			t.Errorf("error deleting test data: %v", e)
// 		}

// 	} else {
// 		t.Errorf("error creating mongodb: %v", e)
// 	}
// }

// func TestSearch(t *testing.T) {
// 	// Should move to mocks.
// 	ctx := context.Background()
// 	server := "localhost"
// 	user := ""
// 	password := ""
// 	port := 27017

// 	repo, e := NewRepository[TestData](ctx, server, user, password, port, "MongoTest", "testdata")
// 	if e == nil {
// 		td1 := TestData{
// 			ID:   primitive.NewObjectID(),
// 			Name: "test1",
// 			Data: "test that thing",
// 			Seq:  1,
// 		}
// 		e1 := repo.Save(ctx, &td1, "name", td1.Name)

// 		td2 := TestData{
// 			ID:   primitive.NewObjectID(),
// 			Name: "test2",
// 			Data: "test this thing",
// 			Seq:  1,
// 		}
// 		e2 := repo.Save(ctx, &td2, "name", td2.Name)

// 		td3 := TestData{
// 			ID:   primitive.NewObjectID(),
// 			Name: "findme",
// 			Data: "findme please",
// 			Seq:  1,
// 		}
// 		e3 := repo.Save(ctx, &td3, "name", td3.Name)

// 		if e1 != nil || e2 != nil || e3 != nil {
// 			t.Errorf("error saving test data: %v - %v - %v", e1, e2, e3)
// 		}

// 		idxName, e := repo.CreateIndex(ctx, "my_text_index", "data", "text", false)
// 		if e == nil {
// 			if idxName == "my_text_index" {
// 				tds, e := repo.TextSearch(ctx, "findme")
// 				if e == nil {
// 					if len(tds) == 1 {
// 						if tds[0].Name == "findme" {
// 							tds, e := repo.TextSearch(ctx, "test")
// 							if e == nil {
// 								if len(tds) == 2 {
// 									if tds[0].Name != "test1" && tds[1].Name != "test1" {
// 										t.Errorf("testsearch have: %s - %s want: %s - %s", tds[0].Name, tds[1].Name, "test1", "test2")
// 									}
// 								} else {
// 									t.Errorf("testsearch have: %d results want: %d", len(tds), 2)
// 								}
// 							} else {
// 								t.Errorf("error searching for term: %s - %v", "test", e)
// 							}
// 						} else {
// 							t.Errorf("testsearch have: %s want: %s", tds[0].Name, "findme")
// 						}

// 					} else {
// 						t.Errorf("testsearch have: %d results want: %d", len(tds), 1)
// 					}

// 				} else {
// 					t.Errorf("error searching for term: %s - %v", "find", e)
// 				}

// 			} else {
// 				t.Errorf("testsearch have: %s want: %s", idxName, "my_text_index")
// 			}
// 		} else {
// 			t.Errorf("error creating index: %v", e)
// 		}

// 		count, e := repo.DeleteMany(ctx, "seq", 1)
// 		if e != nil || count != 3 {
// 			t.Errorf("error deleting test data: %v", e)
// 		}

// 	}
// }
