package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/FooWho/gator/internal/config"
	"github.com/FooWho/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandler map[string]func(*state, command) error
}

type state struct {
	config *config.Config
	db     *database.Queries
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login requires username as argument")
	}
	username := cmd.args[0]
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("call to GetUser() failed for user: %s\n", username)
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Printf("Username set - %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register requires username as argument")
	}
	user, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("unable to create user: %s", cmd.args[0])
	}
	err = handlerLogin(s, command{name: "login", args: cmd.args})
	if err != nil {
		return err
	}
	fmt.Printf("created user: %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("reset does not take arguments")
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("could not TRUNCATE TABLE users")
	}
	fmt.Printf("TRUNCATED TABLE users done\n")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("users does not take arguments")
	}
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("users could not get users: %v\n", err)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.config.CurrentUserName {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Print("\n")
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("agg needs one argument, time_between_reqs")
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not parse time duration: %s", cmd.args[0])
	}
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs.String())

	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Println("Tick at ", t.Format("15:04:05.000"))
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("scrapeFeed() returned %v", err)
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed requires two arguments")
	}

	feed, err := s.db.AddFeed(context.Background(),
		database.AddFeedParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    user.ID,
		})
	if err != nil {
		return fmt.Errorf("unable to add feed %s", cmd.args[0])
	}
	fmt.Printf("Feed created: \n")
	fmt.Printf("    ID: %s\n", feed.ID)
	fmt.Printf("    CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("    UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("    Name: %s\n", feed.Name)
	fmt.Printf("    Url: %s\n", feed.Url)
	fmt.Printf("    UserID: %s\n", feed.UserID)

	feed_follow, err := s.db.FollowFeed(context.Background(),
		database.FollowFeedParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})

	fmt.Printf("User %s is now following %s\n", feed_follow.UserName, feed_follow.FeedName)
	fmt.Printf("    ID: %s\n", feed_follow.ID)
	fmt.Printf("    CreatedAt: %v\n", feed_follow.CreatedAt)
	fmt.Printf("    UpdatedAt: %v\n", feed_follow.UpdatedAt)
	fmt.Printf("    UserID: %s\n", feed_follow.UserID)
	fmt.Printf("    FeedID: %s\n", feed_follow.FeedID)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("feeds takes no arguments")
	}
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return errors.New("unable to get feeds in handlerListFeeds()")
	}
	for _, feed := range feeds {
		userName, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("unable to obtain user for feed %s", feed.Name)
		}
		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("    ID: %s\n", feed.ID)
		fmt.Printf("    CreatedAt: %v\n", feed.CreatedAt)
		fmt.Printf("    UpdatedAt: %v", feed.UpdatedAt)
		fmt.Printf("    Url: %s\n", feed.Url)
		fmt.Printf("    User: %s\n", userName.Name)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow requires a URL argument")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to obtain feed %s in handlerFollowFeed", cmd.args[0])
	}
	feed_follow, err := s.db.FollowFeed(context.Background(),
		database.FollowFeedParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	fmt.Printf("User %s is now following %s\n", feed_follow.UserName, feed_follow.FeedName)
	fmt.Printf("    ID: %s\n", feed_follow.ID)
	fmt.Printf("    CreatedAt: %v\n", feed_follow.CreatedAt)
	fmt.Printf("    UpdatedAt: %v\n", feed_follow.UpdatedAt)
	fmt.Printf("    UserID: %s\n", feed_follow.UserID)
	fmt.Printf("    FeedID: %s\n", feed_follow.FeedID)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("following takes no arguments")
	}
	feeds_followed, err := s.db.GetFeedFollowsByUserID(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to obtain feeds for %s in handlerFollowing", s.config.CurrentUserName)
	}
	fmt.Printf("%s follows:\n", user.Name)
	for _, feed := range feeds_followed {
		fmt.Printf("   Feed Name: %s\n", feed.Name)
		fmt.Printf("   URL: %s\n", feed.Url)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("unfollow requires the url as an argument")
	}

	feed_unfollowed, err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{Url: cmd.args[0], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("unable to unfollow feed %s for %s", cmd.args[0], user.ID)
	}
	fmt.Printf("%s unfollowed:\n", user.Name)
	fmt.Printf("    %v\n", feed_unfollowed)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	switch len(cmd.args) {
	case 0:
		limit = 2
	case 1:
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("browse accepts an optional integer argument for the limit of posts - got: %s", cmd.args[0])
		}
	default:
		return errors.New("browse accepts a single optional argument")
	}
	params := database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)}
	posts, err := s.db.GetPostsForUser(context.Background(), params)

	for _, post := range posts {
		fmt.Printf("        Title: %s\n", post.Title)
		fmt.Printf("        Link: %s\n", post.Url)
		fmt.Printf("        Date: %s\n", post.PublishedAt)
		fmt.Printf("        Description: %s\n\n", post.Description)
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if s.config == nil {
		return errors.New("state not initialized")
	}
	f, ok := c.cmdHandler[cmd.name]
	if !ok {
		return fmt.Errorf("no handler for command: %s", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.cmdHandler == nil {
		c.cmdHandler = make(map[string]func(*state, command) error)
	}
	c.cmdHandler[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to login %s in middlewareLoggedIn", s.config.CurrentUserName)
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New("unable to obtain next feed in handlerAgg()")
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch %s", feed.Name)
	}
	feed, err = s.db.MarkFeedFetched(context.Background(),
		database.MarkFeedFetchedParams{
			ID: feed.ID,
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("unable to mark feed %s as fetched", feed.Name)
	}
	fmt.Printf("Fetched %s:\n", rssFeed.Channel.Title)
	fmt.Printf("    Title: %s\n", rssFeed.Channel.Title)
	fmt.Printf("    Description: %s\n", rssFeed.Channel.Description)
	fmt.Printf("    Url: %s\n", rssFeed.Channel.Link)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("        Title: %s\n", item.Title)
		fmt.Printf("        Link: %s\n", item.Link)
		fmt.Printf("        Date: %s\n", item.PubDate)
		fmt.Printf("        Description: %s\n\n", item.Description)
		nts := sql.NullTime{}
		ts, err := time.Parse("2006-01-02", item.PubDate)
		if err != nil {
			nts.Valid = false
		} else {
			nts.Time = ts
			nts.Valid = true
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: nts,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("Got err: %v\n", err)
		} else {
			fmt.Printf("Posted: %s\n", post.Title)
		}
	}
	return nil
}
