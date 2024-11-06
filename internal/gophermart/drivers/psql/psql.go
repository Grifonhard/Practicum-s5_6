package psql

type DB struct {
}

func New() (*DB, error) {
	var db DB
	return &db, nil
}

func (db *DB) push() {
}