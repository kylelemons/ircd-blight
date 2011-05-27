package user

import (
	"kevlar/ircd/parser"
	"os"
	"strings"
	"sync"
	"time"
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
	mutex  *sync.RWMutex
	ts     int64
	id     string
	user   string
	pass   string
	nick   string
	name   string
	sver   int
	spfx   string
	capab  []string
	server string
	hops   int
	utyp   userType
}

// Get the user ID.
func (u *User) ID() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	if u.utyp == RegisteredAsServer {
		return u.spfx
	}
	return u.id
}

// Get the user nick.
func (u *User) Nick() string {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.nick
}

// Get the username (immutable).
func (u *User) User() string {
	return u.user
}

// Get the user's long name (immutable).
func (u *User) Name() string {
	return u.name
}

// Get the user's registration type (immutable).
func (u *User) Type() userType {
	return u.utyp
}

// Atomically get all of the user's information.
func (u *User) Info() (nick, user, name string, regType userType) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.nick, u.user, u.name, u.utyp
}

// Atomically get all of the server's information.
func (u *User) ServerInfo() (sid, server, pass string, capab []string) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.spfx, u.server, u.pass, u.capab
}

// Set the user's nick.
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
	u.ts = time.Nanoseconds()
	return nil
}

// Set the user and gecos (immutable once set).
func (u *User) SetUser(user, name string) os.Error {
	if len(u.user) > 0 {
		return parser.NewNumeric(parser.ERR_ALREADYREGISTRED)
	}

	if !parser.ValidNick(user) || len(name) == 0 {
		// BUG(kevlar): Document this behavior
		return parser.NewNumeric(parser.ERR_NEEDMOREPARAMS)
	}

	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.user, u.name = user, name
	u.ts = time.Nanoseconds()
	return nil
}

func (u *User) SetPassServer(password, ts, prefix string) os.Error {
	if len(u.user) > 0 {
		return parser.NewNumeric(parser.ERR_ALREADYREGISTRED)
	}

	if len(password) == 0 {
		return os.NewError("Zero-length password")
	}

	if ts != "6" {
		return os.NewError("TS " + ts + " is unsupported")
	}

	if !parser.ValidServerPrefix(prefix) {
		return os.NewError("SID " + prefix + " is invalid")
	}

	u.pass, u.sver, u.spfx = password, 6, prefix
	u.ts = time.Nanoseconds()
	return nil
}

func (u *User) SetCapab(capab string) os.Error {
	if len(u.user) > 0 {
		return parser.NewNumeric(parser.ERR_ALREADYREGISTRED)
	}

	if !strings.Contains(capab, "QS") {
		return os.NewError("QS CAPAB missing")
	}

	if !strings.Contains(capab, "ENCAP") {
		return os.NewError("ENCAP CAPAB missing")
	}

	u.capab = strings.Fields(capab)
	u.ts = time.Nanoseconds()
	return nil
}

func (u *User) SetServer(serv, hops string) os.Error {
	if len(u.user) > 0 {
		return parser.NewNumeric(parser.ERR_ALREADYREGISTRED)
	}

	if len(serv) == 0 {
		return os.NewError("Zero-length server name")
	}

	if hops != "1" {
		return os.NewError("Hops = " + hops + " is unsupported")
	}

	u.server, u.hops = serv, 1
	u.ts = time.Nanoseconds()
	return nil
}

// Set the user's type (immutable once set).
func (u *User) SetType(newType userType) os.Error {
	if u.utyp != Unregistered {
		return parser.NewNumeric(parser.ERR_ALREADYREGISTRED)
	}
	u.utyp = newType
	u.ts = time.Nanoseconds()
	return nil
}

// Get the next available unique ID.
func NextUserID() string {
	return UserIDPrefix + <-userIDs
}

// Atomically retrieve user information.
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

// Get the ID for a particular nick.
func GetID(nick string) (id string, err os.Error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	lownick := parser.ToLower(nick)

	if _, ok := userMap[nick]; ok {
		return nick, nil
	}

	var ok bool
	if id, ok = userNicks[lownick]; !ok {
		err = parser.NewNumeric(parser.ERR_NOSUCHNICK, nick)
	}
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

// Delete the user record.
func Delete(id string) {
	userMutex.Lock()
	defer userMutex.Unlock()

	// Database lookup?
	if u, ok := userMap[id]; ok {
		nick := u.Nick()
		userNicks[nick] = "", false
		userMap[id] = nil, false
	}
}

func init() {
	go genUserIDs()
}
