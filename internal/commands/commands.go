package commands

import (
	"context"
	"fmt"
	"os"
	"time"

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
