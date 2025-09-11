package main

import (
	"context"
	"fmt"

	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(
			context.Background(),
			s.cfg.CurrentUsername,
		)
		if err != nil {
			return fmt.Errorf("error getting current user info: %w", err)
		}

		return handler(s, cmd, currentUser)
	}
}
