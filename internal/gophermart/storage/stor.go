package storage

import "time"

type Storage struct {
	userStor map[int]User
}

type User struct {
	Id            int
	Username      string
	Password_hash string
	Created       time.Time
}

func New() (*Storage, error) {
	var stor Storage
	stor.userStor = make(map[int]User)
	return &stor, nil
}

func (stor *Storage) NewUser(user User) error {
	for _, u := range stor.userStor {
		if u.Username == user.Username {
			return ErrUserExist
		}
	}
	user.Id = len(stor.userStor) + 1
	stor.userStor[user.Id] = user
	return nil
}

func (stor *Storage) GetUser(uname string) (*User, error) {
	for _, u := range stor.userStor {
		if u.Username == uname {
			return &u, nil
		}
	}
	return nil, ErrUserNotExist
}
