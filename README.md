# Blog Aggregator (gator CLI)

This is a command-line blog aggregator written in Go. It allows you to subscribe to RSS feeds, manage your subscriptions, and fetch new posts into a PostgreSQL database.

---

## Prerequisites

To use this project, you’ll need the following installed:

- [Go](https://go.dev/dl/) (1.20 or higher recommended)
- [PostgreSQL](https://www.postgresql.org/download/) (running locally or accessible via connection string)
- Git (for cloning the repository)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/Sanghun1Adam1Park/blog-aggregator-go.git
cd blog-aggregator-go
```

Install PostgreSQL:

For macOS:
```
brew install postgresql
```

For Linux:
```
sudo apt-get install postgresql
```

For Windows:
```
choco install postgresql
```

Run the Migration:

```
goose -dir sql/schema postgres "DBURL" up
sqlc generate
```

---

## Usage

Build the project:

```
go build -o gator
```

Run the application:

```
./gator
```

### Commands

- `login <username>`: Log in to the application.
- `register <username>`: Register a new user.
- `users`: List all users.
- `agg`: Aggregate new posts from subscribed feeds.
- `addfeed <url> <name>`: Add a new RSS feed to your subscriptions.
- `feeds`: List all subscribed feeds.
- `follow <url>`: Subscribe to a new RSS feed.
- `unfollow <url>`: Unsubscribe from a feed.
- `following`: List all feeds you are following.
- `browse`: Browse all posts.

---

## Roadmap

* Add sorting and filtering options to the browse command
* Add pagination to the browse command
* Add concurrency to the agg command so that it can fetch more frequently
* Add a search command that allows for fuzzy searching of posts
* Add bookmarking or liking posts
* Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
* Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
* Write a service manager that keeps the agg command running in the background and restarts it if it crashes
