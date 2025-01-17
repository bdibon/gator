package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bdibon/gator/internal/config"
	"github.com/bdibon/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	dbQueries := database.New(db)

	programState := state{&cfg, dbQueries}

	cmds := commands{handlers: map[string]func(s *state, cmd command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerList)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("missing argument: command name")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{cmdName, cmdArgs}

	err = cmds.run(&programState, cmd)
	if err != nil {
		log.Fatalf("error running %s: %v", cmd.name, err)
	}
}
