package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func getDBStuff() (context.Context, *sql.DB, error) {
	var connStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	return ctx, db, nil
}

func TestGetPageInfo(t *testing.T) {
	ctx, db, err := getDBStuff()
	if err != nil {
		t.Errorf("Couldn't establish database connection: %s\n", err)
	}
	defer db.Close()

	testUUIDs := make(uuid.UUIDs, 2)
	testUUIDs[0], err = uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	testUUIDs[1], err = uuid.Parse("60b6b10c-db33-4b4c-9dcf-566f5b3c59a4")
	revUUID, err := uuid.Parse("4a76899b-0051-444d-92cb-8017f09f2fea")
	if err != nil {
		t.Errorf("Couldn't parse UUID: %s\n", err)
	}
	newsiesArchTime := time.Date(2025, time.November, 12, 0, 0, 0, 0, time.FixedZone("", 0))

	var tests = []struct {
		ctx context.Context
		db *sql.DB
		uuid uuid.UUID
		answer *PageInfo
	}{
		{ctx, db, testUUIDs[0], &PageInfo{testUUIDs[0], "Dan Boone", &revUUID, nil, nil}},
		{ctx, db, testUUIDs[1], &PageInfo{testUUIDs[1], "Newsies", nil, &newsiesArchTime, nil}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("uuid:%s", tt.uuid)
		t.Run(testname, func(t *testing.T) {
			res, err := GetPageInfo(tt.ctx, tt.db, tt.uuid)
			if err != nil {
				t.Errorf("Couldn't execute function: %s\n", err)
			}
			if !reflect.DeepEqual(res, tt.answer) {
				t.Errorf("Result different from expected.\nResult:\t%s\nExpected:\t%s\n", res, tt.answer)
			}
		})
	}
}

func TestGetPageNameUUIDs(t *testing.T) {
	ctx, db, err := getDBStuff()
	if err != nil {
		t.Errorf("Couldn't establish database connection: %s\n", err)
	}
	defer db.Close()

	testUUIDs := make(uuid.UUIDs, 2)
	testUUIDs[0], err = uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	testUUIDs[1], err = uuid.Parse("60b6b10c-db33-4b4c-9dcf-566f5b3c59a4")
	if err != nil {
		t.Errorf("Couldn't parse UUID: %s\n", err)
	}

	expected := make([]NameUUID, 2)
	expected[0] = NameUUID{"Dan Boone", testUUIDs[0]}
	expected[1] = NameUUID{"Newsies", testUUIDs[1]}

	var tests = []struct {
		name string
		ctx context.Context
		db *sql.DB
		answer []NameUUID
	}{
		{"test_db", ctx, db, expected},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("name:%s", tt.name)
		t.Run(testname, func(t *testing.T) {
			res, err := GetPageNameUUIDs(tt.ctx, tt.db)
			if err != nil {
				t.Errorf("Couldn't execute function: %s\n", err)
			}
			if !reflect.DeepEqual(res, tt.answer) {
				t.Errorf("Result different from expected.\nResult:\t%s\nExpected:\t%s\n", res, tt.answer)
			}
		})
	}
}

func TestGetPageRevisionsInfo(t *testing.T) {
	ctx, db, err := getDBStuff()
	if err != nil {
		t.Errorf("Couldn't establish database connection: %s\n", err)
	}
	defer db.Close()

	pageUUID, err := uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	revUUID, err := uuid.Parse("4a76899b-0051-444d-92cb-8017f09f2fea")
	if err != nil {
		t.Errorf("Couldn't parse UUID: %s\n", err)
	}
	revTime := time.Date(2025, time.December, 18, 20, 18, 9, 115549000, time.FixedZone("", 0))

	expected := make([]RevInfo, 1)
	expected[0] = RevInfo{revUUID, revTime, "1197028"}

	var tests = []struct {
		ctx context.Context
		db *sql.DB
		page_id uuid.UUID
		answer []RevInfo
	}{
		{ctx, db, pageUUID, expected},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("page_uuid:%s", tt.page_id)
		t.Run(testname, func(t *testing.T) {
			res, err := GetPageRevisionsInfo(tt.ctx, tt.db, tt.page_id)
			if err != nil {
				t.Errorf("Couldn't execute function: %s\n", err)
			}
			if !reflect.DeepEqual(res, tt.answer) {
				t.Errorf("Result different from expected.\nResult:\t%s\nExpected:\t%s\n", res, tt.answer)
			}
		})
	}
}
