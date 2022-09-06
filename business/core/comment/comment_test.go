package comment_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dudakovict/social-network/business/core/comment"
	"github.com/dudakovict/social-network/business/data/comment/dbschema"
	"github.com/dudakovict/social-network/business/data/comment/dbtest"
	"github.com/dudakovict/social-network/foundation/docker"
	"github.com/google/go-cmp/cmp"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestComment(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testcomment")
	t.Cleanup(teardown)

	core := comment.NewCore(log, db)

	t.Log("Given the need to work with Comment records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Comment.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

			nc := comment.NewComment{
				Description: "Check out my new song!",
				PostID:      "45b5fbd3-755f-4379-8f07-a58d4a30fa2f",
				UserID:      "45b5fbd3-755f-4379-8f07-a58d4a30fa2f",
			}

			c, err := core.Create(ctx, nc, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create comment.", dbtest.Success, testID)

			_, err = core.Create(ctx, nc, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create comment.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, c.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comment by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comment by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(c, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same comment. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same comment.", dbtest.Success, testID)

			upd := comment.UpdateComment{
				Description: dbtest.StringPointer("I've just released a new album!"),
			}

			if err := core.Update(ctx, c.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update comment.", dbtest.Success, testID)

			saved, err = core.QueryByID(ctx, c.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comment by ID : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comment by ID.", dbtest.Success, testID)

			if saved.Description != *upd.Description {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Description.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Description)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Description)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Description.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, c.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete comment.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, c.ID)
			if !errors.Is(err, comment.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve comment.", dbtest.Success, testID)

			_, err = core.QueryByUserID(ctx, c.UserID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comments by UserID : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comments by UserID.", dbtest.Success, testID)

			_, err = core.QueryByPostID(ctx, c.PostID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comments by PostID : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comments by PostID.", dbtest.Success, testID)
		}
	}
}

func TestPagingComment(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	comment := comment.NewCore(log, db)

	t.Log("Given the need to page through Comment records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 comments.", testID)
		{
			ctx := context.Background()

			comments1, err := comment.Query(ctx, 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comments for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comments for page 1.", dbtest.Success, testID)

			if len(comments1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single comment.", dbtest.Success, testID)

			comments2, err := comment.Query(ctx, 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve comments for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve comments for page 2.", dbtest.Success, testID)

			if len(comments2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single comment : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single comment.", dbtest.Success, testID)

			if comments1[0].ID == comments2[0].ID {
				t.Logf("\t\tTest %d:\tComment1: %v", testID, comments1[0].ID)
				t.Logf("\t\tTest %d:\tComment2: %v", testID, comments2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different comments : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different comments.", dbtest.Success, testID)
		}
	}
}
