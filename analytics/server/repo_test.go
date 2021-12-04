package server_test

import (
	"context"
	"encoding/json"
	"teknologi-umum-bot/analytics/server"
	"testing"
	"time"
)

func TestGetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c, err := db.Connx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO analytics
			(user_id, username, display_name, counter, created_at, joined_at, updated_at)
			VALUES
			($1, $2, $3, $4, $5, $6, $7)`,
		1,
		"user1",
		"User 1",
		1,
		time.Now(),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	deps := &server.Dependency{
		DB:     db,
		Memory: memory,
	}

	res, err := deps.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}

	var user []server.User
	err = json.Unmarshal(res, &user)
	if err != nil {
		t.Error(err)
	}

	if len(user) != 1 {
		t.Errorf("Expected 1 user, got %d", len(user))
	}

	if user[0].UserID != 1 {
		t.Error("user id should be 1, got:", user[0].UserID)
	}

	if user[0].Username != "user1" {
		t.Error("username should be user1, got:", user[0].Username)
	}

	if user[0].DisplayName != "User 1" {
		t.Error("display name should be User 1, got:", user[0].DisplayName)
	}

	if user[0].Counter != 1 {
		t.Error("counter should be 1, got:", user[0].Counter)
	}

	// try to get it again
	res2, err := deps.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}

	if string(res) != string(res2) {
		t.Errorf("result should be the same\n\nres1: %s\n\nres2: %s", string(res), string(res2))
	}
}

func TestGetTotal(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// create a dummy user struct slice
	users := []server.User{
		{UserID: 1, Username: "user1", DisplayName: "User 1", Counter: 1},
		{UserID: 2, Username: "user2", DisplayName: "User 2", Counter: 2},
		{UserID: 3, Username: "user3", DisplayName: "User 3", Counter: 3},
	}

	// convert users slice to single slice with no keys, just values.
	var usersSlice []interface{}
	for _, v := range users {
		usersSlice = append(usersSlice, v.UserID)
		usersSlice = append(usersSlice, v.Username)
		usersSlice = append(usersSlice, v.DisplayName)
		usersSlice = append(usersSlice, v.Counter)
		usersSlice = append(usersSlice, v.CreatedAt)
		usersSlice = append(usersSlice, v.JoinedAt)
		usersSlice = append(usersSlice, v.UpdatedAt)
	}

	c, err := db.Connx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO analytics
			(user_id, username, display_name, counter, created_at, joined_at, updated_at)
			VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14),
			($15, $16, $17, $18, $19, $20, $21)`,
		usersSlice...,
	)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	deps := &server.Dependency{
		DB:     db,
		Memory: memory,
	}

	data, err := deps.GetTotal(ctx)
	if err != nil {
		t.Error(err)
	}

	if string(data) != "3" {
		t.Errorf("Expected 3, got %s", data)
	}
}
