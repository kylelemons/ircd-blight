package parser

const (
	// Client commands
	CMD_NICK   = "NICK"
	CMD_USER   = "USER"
	CMD_SERVER = "SERVER"
	CMD_CAPAB  = "CAPAB"
	CMD_PASS   = "PASS"
	CMD_ERROR  = "ERROR"
	CMD_QUIT   = "QUIT"
	CMD_PING   = "PING"
	CMD_PONG   = "PONG"

	CMD_OPER = "OPER"
	CMD_MODE = "MODE"

	CMD_JOIN  = "JOIN"
	CMD_PART  = "PART"
	CMD_WHO   = "WHO"
	CMD_TOPIC = "TOPIC"
	CMD_NAMES = "NAMES"

	CMD_WALLOPS = "WALLOPS"
	CMD_PRIVMSG = "PRIVMSG"
	CMD_NOTICE  = "NOTICE"

	// Server commands
	CMD_SJOIN = "SJOIN"
	CMD_SID   = "SID"
	CMD_SQUIT = "SQUIT"
	CMD_UID   = "UID"
	CMD_EUID  = "EUID"
	CMD_ENCAP = "ENCAP"
	CMD_BMASK = "BMASK"
	CMD_TB    = "TB"

	// Internal commands
	INT_DELUSER = "deluser" // Delete all UIDs in DestIDs
)
