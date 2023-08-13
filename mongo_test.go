package gomongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestData struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Data string             `json:"data" bson:"data"`
	Seq  int64              `json:"seq" bson:"seq"`
}

func createTestDB(ctx context.Context, t *testing.T) (db *MongoDB, e error) {
	server := "localhost"
	user := ""
	password := ""
	port := 27017
	name := "MongoTest"

	result, err := NewMongoDB(ctx, server, user, password, port, name)
	if err != nil {
		t.Errorf("error creating mongodb: %v", e)
	}

	return result, e
}

// This all really should be handled by mocks, but local test db is easier for now.
func TestNewMongoDB(t *testing.T) {
	ctx := context.Background()
	db, e := createTestDB(ctx, t)

	if e == nil {
		connStr := "mongodb://localhost:27017/"
		if db.ConnectionString() != connStr {
			t.Errorf("incorrect connection string want %s have: %s", connStr, db.ConnectionString())
		}
	}
}

func TestNewMongoDBWithSrv(t *testing.T) {
	ctx := context.Background()
	server := "test.mongodb.net"
	user := "admin"
	password := "123456"
	name := "MongoTest"
	port := 0

	db, e := NewMongoDB(ctx, server, user, password, port, name)
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

func TestFindOneByFieldName(t *testing.T) {
	ctx := context.Background()
	db, e := createTestDB(ctx, t)

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

		_, e = helper.InsertOne(ctx, td1)
		if e != nil {
			t.Errorf("error inserting test data: %v", e)
		} else {
			result, e := helper.FindOne(ctx, "Name", "test1")

			if e != nil {
				t.Errorf("error finding test data: %v", e)
			} else {
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

func TestFindOneAndIncrement(t *testing.T) {

	ctx := context.Background()
	db, e := createTestDB(ctx, t)

	helper := NewHelper[TestData](context.Background(), db, "testdata")
	if e == nil {
		td, e := helper.FindOneAndIncrementField(ctx, "Name", "test", "Seq", 1)
		if e != nil {
			t.Errorf("error finding and incrementing: %v", e)
		} else {
			if td.Seq != 1 {
				t.Errorf("findandincrement have: %d want: %d", td.Seq, 1)
			} else {
				td, e := helper.FindOneAndIncrementField(ctx, "Name", "test", "Seq", 5)
				if e == nil {
					if td.Seq != 6 {
						t.Errorf("findandincrement have: %d want: %d", td.Seq, 6)
					}
				} else {
					t.Errorf("error finding and incrementing: %v", e)
				}
			}

			count, e := helper.DeleteOne(ctx, "Name", td.Name)
			if e != nil || count != 1 {
				t.Errorf("error deleting test data: %v", e)
			}
		}
	}
}

func TestFindOneAndUpdate(t *testing.T) {

	ctx := context.Background()
	db, e := createTestDB(ctx, t)

	update := int64(125)

	helper := NewHelper[TestData](context.Background(), db, "testdata")
	if e == nil {
		td, e := helper.FindOneAndUpdateField(ctx, "Name", "test", "Seq", update)
		if e != nil {
			t.Errorf("error finding and incrementing: %v", e)
		} else {
			if td.Seq != update {
				t.Errorf("findandincrement have: %d want: %d", td.Seq, update)
			}

			count, e := helper.DeleteOne(ctx, "Name", td.Name)
			if e != nil || count != 1 {
				t.Errorf("error deleting test data: %v", e)
			}
		}

	}
}

func TestCreateIndex(t *testing.T) {

	ctx := context.Background()
	db, e := createTestDB(ctx, t)

	update := int64(125)

	helper := NewHelper[TestData](context.Background(), db, "testdata")
	if e == nil {
		td, e := helper.FindOneAndUpdateField(ctx, "Name", "test", "Seq", update)
		if e == nil {
			if td.Seq != update {
				t.Errorf("createindex have: %d want: %d", td.Seq, update)
			}

			idxName, e := helper.CreateIndex(ctx, "my_text_index", "Data", "text", false)
			if e == nil {
				if idxName != "my_text_index" {
					t.Errorf("createindex have: %s want: %s", idxName, "my_text_index")
				}
			} else {
				t.Errorf("error creating index: %v", e)
			}

			count, e := helper.DeleteOne(ctx, "Name", td.Name)
			if e != nil || count != 1 {
				t.Errorf("error deleting test data: %v", e)
			}
		}
	}
}

func TestTextSearch(t *testing.T) {
	ctx := context.Background()
	db, e := createTestDB(ctx, t)

	helper := NewHelper[TestData](context.Background(), db, "testdata")
	if e == nil {
		td1 := TestData{
			ID:   primitive.NewObjectID(),
			Name: "test1",
			Data: "test that thing",
			Seq:  1,
		}
		_, e1 := helper.SaveOne(ctx, &td1, "Name", td1.Name)

		td2 := TestData{
			ID:   primitive.NewObjectID(),
			Name: "test2",
			Data: "test this thing",
			Seq:  1,
		}
		_, e2 := helper.SaveOne(ctx, &td2, "Name", td2.Name)

		td3 := TestData{
			ID:   primitive.NewObjectID(),
			Name: "findme",
			Data: "findme please",
			Seq:  1,
		}
		_, e3 := helper.SaveOne(ctx, &td3, "Name", td3.Name)

		if e1 != nil || e2 != nil || e3 != nil {
			t.Errorf("error saving test data: %v - %v - %v", e1, e2, e3)
		}

		idxName, e := helper.CreateIndex(ctx, "my_text_index", "Data", "text", false)
		if e == nil {
			if idxName == "my_text_index" {
				tds, e := helper.TextSearch(ctx, "findme")
				if e == nil {
					if len(tds) == 1 {
						if tds[0].Name == "findme" {
							tds, e := helper.TextSearch(ctx, "test")
							if e == nil {
								if len(tds) == 2 {
									if tds[0].Name != "test1" && tds[1].Name != "test1" {
										t.Errorf("testsearch have: %s - %s want: %s - %s", tds[0].Name, tds[1].Name, "test1", "test2")
									}
								} else {
									t.Errorf("testsearch have: %d results want: %d", len(tds), 2)
								}
							} else {
								t.Errorf("error searching for term: %s - %v", "test", e)
							}
						} else {
							t.Errorf("testsearch have: %s want: %s", tds[0].Name, "findme")
						}

					} else {
						t.Errorf("testsearch have: %d results want: %d", len(tds), 1)
					}

				} else {
					t.Errorf("error searching for term: %s - %v", "find", e)
				}

			} else {
				t.Errorf("testsearch have: %s want: %s", idxName, "my_text_index")
			}
		} else {
			t.Errorf("error creating index: %v", e)
		}

		count, e := helper.DeleteMany(ctx, "Seq", 1)
		if e != nil || count != 3 {
			t.Errorf("error deleting test data: %v", e)
		}

	}
}
