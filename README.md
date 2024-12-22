# Setup

## Prerequisites

- A running instance of Postgres
- Go installed

## Instructions

- Clone this repository

```sh
git clone
```

- Install the dependencies

```sh
go install
```

- Create the config file

```sh
touch ~/.gatorconfig.json
```

- The content of `gatorconfig.json` should be the following

```json
{
  // Connection string to your postgres instance
  "db_url": "postgres://bdibon:postgres:@localhost:5432/gator?sslmode=disable"
  // Later if should contain a user field
}
```

- Build the project

```sh
go build
```

- Register a new user

```sh
gator register <username>
```

Now you can run several commands to manage your RSS feeds.

# Usage

- Add a feed to your collection

```sh
gator follow https://news.ycombinator.com/rss
```

- List the feeds you are following

```sh
gator following
```

- Run an infinite loop to aggregate posts from your feeds every 1s

```sh
gator agg 1s
```

- Browse the posts from the feeds you've subscribed to

```sh
gator browse
```
