package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: agg <time_between_reqs>")
	}

	requestedTime := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(requestedTime)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", requestedTime)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s.db)
	}
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("illegal argument, usage: addfeed <name of the feed> <url of feed>")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
			Url:       url,
			UserID:    currentUser.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println(feed)
	cmd.args = cmd.args[1:]
	return handlerFollow(s, cmd, currentUser)
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("illegal argument, usage: feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from db: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%s @%s created by %s\n", feed.FeedName, feed.FeedUrl, feed.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: follow <feed_url>")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		url,
	)
	if err != nil {
		return fmt.Errorf("error getting feed info: %w", err)
	}

	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    currentUser.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("illegal argument, usage: following")
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		currentUser.ID,
	)
	if err != nil {
		return fmt.Errorf("error getting feeds followed by user: %w", err)
	}

	for _, feedFollow := range feedFollows {
		fmt.Printf(" - %s\n", feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: follow <feed_url>")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		url,
	)
	if err != nil {
		return fmt.Errorf("error getting feed info: %w", err)
	}

	if err = s.db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			UserID: currentUser.ID,
			FeedID: feed.ID,
		},
	); err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 0 || len(cmd.args) != 1 {
		return errors.New("illegal argument, usage: browse [limit]")
	}

	var limit int32
	if len(cmd.args) == 1 {
		intInput, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(intInput)
	}
	limit = 2

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: currentUser.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return fmt.Errorf("error getting posts for user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf(" - %s\n", post.Title)
	}

	return nil
}
