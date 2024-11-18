package server

import "flag"

var (
	flagRunAddress  string
	flagDatabaseURI string
)

func parseFlags() {
	flag.StringVar(&flagRunAddress, "a", "localhost:8080", "address and port of the service run")
	flag.StringVar(&flagDatabaseURI, "d", "postgres://postgres:postgres@localhost:5432/postgres", "database connection address")
	flag.Parse()
}
