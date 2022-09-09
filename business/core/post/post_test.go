package post_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dudakovict/social-network/business/core/post"
	"github.com/dudakovict/social-network/business/data/post/dbschema"
	"github.com/dudakovict/social-network/business/data/post/dbtest"
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

func TestPost(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpost")
	t.Cleanup(teardown)

	core := post.NewCore(log, db, nil)

	t.Log("Given the need to work with Post records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Post.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

			np := post.NewPost{
				Title:       "New Song",
				Description: "Check out my new song!",
				UserID:      "45b5fbd3-755f-4379-8f07-a58d4a30fa2f",
			}

			p, err := core.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create post.", dbtest.Success, testID)

			_, err = core.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create post.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, p.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve post by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve post by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(p, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same post. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same post.", dbtest.Success, testID)

			upd := post.UpdatePost{
				Title:       dbtest.StringPointer("New Album"),
				Description: dbtest.StringPointer("I've just released a new album!"),
			}

			if err := core.Update(ctx, p.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update post.", dbtest.Success, testID)

			saved, err = core.QueryByID(ctx, p.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve post by ID : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve post by ID.", dbtest.Success, testID)

			if saved.Title != *upd.Title {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Title.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Title)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Title)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Title.", dbtest.Success, testID)
			}

			if saved.Description != *upd.Description {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Description.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Description)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Description)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Description.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, p.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete post.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, p.ID)
			if !errors.Is(err, post.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve post.", dbtest.Success, testID)

			_, err = core.QueryByUserID(ctx, p.UserID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve posts by UserID : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve posts by UserID.", dbtest.Success, testID)
		}
	}
}

func TestPagingPost(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	post := post.NewCore(log, db, nil)

	t.Log("Given the need to page through Post records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 posts.", testID)
		{
			ctx := context.Background()

			posts1, err := post.Query(ctx, 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve posts for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve posts for page 1.", dbtest.Success, testID)

			if len(posts1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single post.", dbtest.Success, testID)

			posts2, err := post.Query(ctx, 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve posts for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve posts for page 2.", dbtest.Success, testID)

			if len(posts2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single post : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single post.", dbtest.Success, testID)

			if posts1[0].ID == posts2[0].ID {
				t.Logf("\t\tTest %d:\tPost1: %v", testID, posts1[0].ID)
				t.Logf("\t\tTest %d:\tPost2: %v", testID, posts2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different posts : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different posts.", dbtest.Success, testID)
		}
	}
}
