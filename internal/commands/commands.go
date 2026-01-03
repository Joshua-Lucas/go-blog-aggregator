package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Joshua-Lucas/blog-aggregator/internal/commands/feed"
	"github.com/Joshua-Lucas/blog-aggregator/internal/config"
	"github.com/Joshua-Lucas/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	println("commands")
	println(cmd.Name)
	println(cmd.Args)
	v, ok := c.Handlers[cmd.Name]

	if !ok {
		return fmt.Errorf("There is no handler for the provided Command")
	}

	err := v(s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("Login in handler expects a single argument, the username")
	}

	userName := cmd.Args[0]

	user, err := s.Db.GetUser(context.Background(), userName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = s.Config.SetUser(user.Name)

	if err != nil {
		return err
	}

	println("User has been set")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) <= 0 {
		return fmt.Errorf("Register command expects a single argument, the username")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	}

	user, err := s.Db.CreateUser(context.Background(), newUser)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = s.Config.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User created: %v", user)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	println("Reset successful")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, val := range users {

		if s.Config.CurrentUserName == val.Name {
			fmt.Printf("%s (current)\n", val.Name)
			continue
		}

		fmt.Printf("%s\n", val.Name)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {

	f, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v", f)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {

	if len(cmd.Args) <= 0 {
		return fmt.Errorf("addfeed command expects two arguments, the feed name and feed url")
	}

	if len(cmd.Args) > 2 || len(cmd.Args) == 1 {
		return fmt.Errorf("addfeed command expects two arguments, the feed name and feed url. We detect %v arguments passed", len(cmd.Args))
	}

	userName := s.Config.CurrentUserName

	// Grab current user from database
	user, err := s.Db.GetUser(context.Background(), userName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Need arguments of name and url
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(context.Background(), newFeed)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v", feed)

	return nil
}
