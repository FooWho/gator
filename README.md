# Gator - A Blog Aggregator

Gator is a guided project from [boot.dev](https://boot.dev).  It is a blog aggregator written in Go with a PostgreSQL backend. If you have Go installed, you can install Gator with the command `go install github.com/FooWho/gator/cmd/gator@latest`. In order to run Gator, you will also need PostgreSQL, and Goose.

Instructions for installing PostgreSQL can be found at the [PostgreSQL website](https://www.postgresql.org/).

Installation of Goose is performed by Go: `go install github.com/pressly/goose/v3/cmd/goose@latest`.

Using psql, create a database named gator and create a username/password for this database.

Create a configuration file in your home directory, `.gatorconfig.json`. The format for this file is:

```json
{
  "db_url": "*postgres://postgres:postgres@localhost:5432/gator*"
}
```

Replace the username/password combo in the above, "postgres:postgres" with whatever username and password you created for your user.

You will need to navigate to the directory where Go installs packages and locate the directory for Gator. From that director, navigate to `sql/schema`.  Run the command `goose postgres "*postgres://postgres:postgres@localhost:5432/gator*" up`. Again, replace the username/password with the username/password for your Gator user. This command will setup the database, once the database has been created.

Gator can be run with the following: `gator \\[command\\] \\[arguments\\]`.

Gator supports the following commands:

- `gator register \\[username\\]` - This command creates a registration for a user.
- `gator login \\[username\\]` - This command designates the currently logged in users.
- `gator users` - This command will list the registered users and show which user is the current user.
- `gator reset` - This command restores the database back to an empty state.
- `gator agg \\[time\\]` - This command will begin aggregating the blogs. The time argument is how long to wait between scrapes of the blogs. Example time formats - `30s`, `15m`, `1h`, etc.
- `gator addfeed \\[feedname\\] \\[feed_url\\]` - This creates an entry for a blog feed named `feedname` located at `feed_url`.
- `gator feeds` - This command will list the feeds that have been registered with Gator.
- `gator follow \\[feed_url\\]` - This command will put the designated feed in the list of feeds followed by the active user.
- `gator following` - This command will list all feeds followed by the current user.
- `gator unfollow \\[feed_url\\]` - This command will unfollow the designated feed for the current user.
- `gator browse \\[optional_limit\\]` - This command will show the `optional_limit` most recent posts from the user's followed feeds. If no limit is supplied, it will default to a limit of 2.

The program is pretty rough around the edges (as you can see from the installation instructions). Lots of useful information is dumped to the terminal while it is running, including any error messages. Of particular note, when using the `agg` command, you will see messages for every fetch from the feeds. This is useful to detect something is wrong and you are potentially hammering your feeds with requests for content.

