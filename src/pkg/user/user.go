package user

import (
	"sync"
)

var (
	userMutex = new(sync.Mutex)
	userMap   = make(map[string]*User)
	userIDs   = make(chan string)
)

var (
	UserIDPrefix = "000"
)

func genUserIDs() {
	chars := []byte{'A', 'A', 'A', 'A', 'A', 'A'}
	for {
		userIDs <- string(chars)
		for i := len(chars) - 1; i >= 0; i-- {
			switch chars[i] {
			case '9':
				chars[i] = 'A'
				continue
			case 'Z':
				chars[i] = '0'
			default:
				chars[i]++
			}
			break
		}
	}
}

type User struct {
	mutex *sync.RWMutex
	id    string
	User  string
	Nick  string
	Name  string
}

func NextUserID() string {
	return UserIDPrefix + <-userIDs
}

func NewUser(id string) *User {
	userMutex.Lock()
	defer userMutex.Unlock()

	// Database lookup
	if u, ok := userMap[id]; ok {
		return u
	}

	u := &User{
		mutex: new(sync.RWMutex),
		id:    id,
	}

	userMap[id] = u
	return u
}

func (u *User) Login(user, nick, name string) {
	u.User, u.Nick, u.Name = user, nick, name
}

func init() {
	go genUserIDs()
}
