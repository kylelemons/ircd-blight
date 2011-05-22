package user

import (
	"os"
	"sync"
	"kevlar/ircd/parser"
)

var (
	// Always lock this before locking a user mutex
	userMutex = new(sync.RWMutex)
	userMap   = make(map[string]*User)
	userIDs   = make(chan string)

	// userNicks[nick] = uid
	userNicks = make(map[string]string)
)

var (
	// The server prefix for all user IDs.  Set this before calling
	// NextUserID.
	UserIDPrefix = "000"
)

type userType int

const (
	Unregistered userType = iota
	RegisteredAsUser
	RegisteredAsServer
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

// Store the user information and keep it synchronized across possible
// multiple accesses.
type User struct {
	mutex *sync.RWMutex
	id    string
	user  string
	nick  string
	name  string
	utyp  userType
}

// Get the user ID.
func (u *User) ID() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.id
}

// Get the user nick.
func (u *User) Nick() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.nick
}

// Get the username.
func (u *User) User() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.user
}

// Get the user's long name.
func (u *User) Name() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.name
}

// Get the user's registration type.
func (u *User) Type() userType {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.utyp
}

// Atomically get all of the user's information
func (u *User) Info() (nick, user, name string, regType userType) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.nick, u.user, u.name, u.utyp
}

// Set the user's nick
func (u *User) SetNick(nick string) os.Error {
	if !parser.ValidNick(nick) {
		return parser.NewNumeric(parser.ERR_ERRONEUSNICKNAME, nick)
	}

	lownick := parser.ToLower(nick)

	userMutex.Lock()
	defer userMutex.Unlock()

	if id, used := userNicks[lownick]; used {
		if id == u.ID() {
			return nil
		}
		return parser.NewNumeric(parser.ERR_NICKNAMEINUSE, nick)
	}
	userNicks[lownick] = u.ID()

	lownick = parser.ToLower(u.nick)
	userNicks[lownick] = "", false

	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.nick = nick
	return nil
}

func (u *User) SetUser(user, name string) os.Error {
	if !parser.ValidNick(user) {
		// BUG(kevlar): Document this behavior
		return parser.NewNumeric(parser.ERR_NEEDMOREPARAMS)
	}

	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.user, u.name = user, name
	return nil
}

// Set the user's type
func (u *User) SetType(newType userType) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.utyp = newType
}

// Get the next available unique ID.
func NextUserID() string {
	return UserIDPrefix + <-userIDs
}

func GetInfo(id string) (nick, user, name string, regType userType, ok bool) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	var u *User
	if u, ok = userMap[id]; !ok {
		return
	}

	nick, user, name, regType = u.Info()
	return
}

// Get the User structure for the given ID.  If it does not exist, it is
// created.
func Get(id string) *User {
	userMutex.Lock()
	defer userMutex.Unlock()

	// Database lookup?
	if u, ok := userMap[id]; ok {
		return u
	}

	u := &User{
		mutex: new(sync.RWMutex),
		id:    id,
		nick:  "*",
	}

	userMap[id] = u
	return u
}

func init() {
	go genUserIDs()
}
