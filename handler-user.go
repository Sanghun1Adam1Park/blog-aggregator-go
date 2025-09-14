package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: login <username>")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(
		context.Background(),
		username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no such user: %s", username)
		}
		return fmt.Errorf("could not fetch user: %w", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("username set to %s\n", s.cfg.CurrentUsername)
	return nil
}

func handlerRegsiter(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: register <username>")
	}

	name := cmd.args[0]
	_, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	if err := handlerLogin(s, cmd); err != nil {
		return err
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("illegal argument, usage: reset")
	}

	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error reseting the users table: %w", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("illegal argument, usage: users")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users from table: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
