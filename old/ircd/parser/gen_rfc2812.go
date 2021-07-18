// Automatically generated from ../../doc/IRC-RFC2812.txt ../../doc/IRC-CustomNumerics.txt
// DO NOT EDIT

package parser

// IRC Numerics
const (
	RPL_WELCOME           = "001"
	RPL_YOURHOST          = "002"
	RPL_CREATED           = "003"
	RPL_MYINFO            = "004"
	RPL_BOUNCE            = "005"
	RPL_TRACELINK         = "200"
	RPL_TRACECONNECTING   = "201"
	RPL_TRACEHANDSHAKE    = "202"
	RPL_TRACEUNKNOWN      = "203"
	RPL_TRACEOPERATOR     = "204"
	RPL_TRACEUSER         = "205"
	RPL_TRACESERVER       = "206"
	RPL_TRACESERVICE      = "207"
	RPL_TRACENEWTYPE      = "208"
	RPL_TRACECLASS        = "209"
	RPL_STATSLINKINFO     = "211"
	RPL_STATSCOMMANDS     = "212"
	RPL_ENDOFSTATS        = "219"
	RPL_UMODEIS           = "221"
	RPL_SERVLIST          = "234"
	RPL_SERVLISTEND       = "235"
	RPL_STATSUPTIME       = "242"
	RPL_STATSOLINE        = "243"
	RPL_ISUPPORT          = "250"
	RPL_LUSERCLIENT       = "251"
	RPL_LUSEROP           = "252"
	RPL_LUSERUNKNOWN      = "253"
	RPL_LUSERCHANNELS     = "254"
	RPL_LUSERME           = "255"
	RPL_ADMINME           = "256"
	RPL_ADMINEMAIL        = "259"
	RPL_TRACELOG          = "261"
	RPL_TRACEEND          = "262"
	RPL_TRYAGAIN          = "263"
	RPL_AWAY              = "301"
	RPL_USERHOST          = "302"
	RPL_ISON              = "303"
	RPL_UNAWAY            = "305"
	RPL_NOWAWAY           = "306"
	RPL_WHOISUSER         = "311"
	RPL_WHOISSERVER       = "312"
	RPL_WHOISOPERATOR     = "313"
	RPL_WHOWASUSER        = "314"
	RPL_ENDOFWHO          = "315"
	RPL_WHOISIDLE         = "317"
	RPL_ENDOFWHOIS        = "318"
	RPL_WHOISCHANNELS     = "319"
	RPL_LIST              = "322"
	RPL_LISTEND           = "323"
	RPL_CHANNELMODEIS     = "324"
	RPL_UNIQOPIS          = "325"
	RPL_NOTOPIC           = "331"
	RPL_TOPIC             = "332"
	RPL_INVITING          = "341"
	RPL_SUMMONING         = "342"
	RPL_INVITELIST        = "346"
	RPL_ENDOFINVITELIST   = "347"
	RPL_EXCEPTLIST        = "348"
	RPL_ENDOFEXCEPTLIST   = "349"
	RPL_VERSION           = "351"
	RPL_WHOREPLY          = "352"
	RPL_NAMREPLY          = "353"
	RPL_LINKS             = "364"
	RPL_ENDOFLINKS        = "365"
	RPL_ENDOFNAMES        = "366"
	RPL_BANLIST           = "367"
	RPL_ENDOFBANLIST      = "368"
	RPL_ENDOFWHOWAS       = "369"
	RPL_INFO              = "371"
	RPL_MOTD              = "372"
	RPL_ENDOFINFO         = "374"
	RPL_MOTDSTART         = "375"
	RPL_ENDOFMOTD         = "376"
	RPL_YOUREOPER         = "381"
	RPL_REHASHING         = "382"
	RPL_YOURESERVICE      = "383"
	RPL_TIME              = "391"
	RPL_USERSSTART        = "392"
	RPL_USERS             = "393"
	RPL_ENDOFUSERS        = "394"
	RPL_NOUSERS           = "395"
	ERR_NOSUCHNICK        = "401"
	ERR_NOSUCHSERVER      = "402"
	ERR_NOSUCHCHANNEL     = "403"
	ERR_CANNOTSENDTOCHAN  = "404"
	ERR_TOOMANYCHANNELS   = "405"
	ERR_WASNOSUCHNICK     = "406"
	ERR_TOOMANYTARGETS    = "407"
	ERR_NOSUCHSERVICE     = "408"
	ERR_NOORIGIN          = "409"
	ERR_NORECIPIENT       = "411"
	ERR_NOTEXTTOSEND      = "412"
	ERR_NOTOPLEVEL        = "413"
	ERR_WILDTOPLEVEL      = "414"
	ERR_BADMASK           = "415"
	ERR_UNKNOWNCOMMAND    = "421"
	ERR_NOMOTD            = "422"
	ERR_NOADMININFO       = "423"
	ERR_FILEERROR         = "424"
	ERR_NONICKNAMEGIVEN   = "431"
	ERR_ERRONEUSNICKNAME  = "432"
	ERR_NICKNAMEINUSE     = "433"
	ERR_NICKCOLLISION     = "436"
	ERR_UNAVAILRESOURCE   = "437"
	ERR_USERNOTINCHANNEL  = "441"
	ERR_NOTONCHANNEL      = "442"
	ERR_USERONCHANNEL     = "443"
	ERR_NOLOGIN           = "444"
	ERR_SUMMONDISABLED    = "445"
	ERR_USERSDISABLED     = "446"
	ERR_NOTREGISTERED     = "451"
	ERR_NEEDMOREPARAMS    = "461"
	ERR_ALREADYREGISTRED  = "462"
	ERR_NOPERMFORHOST     = "463"
	ERR_PASSWDMISMATCH    = "464"
	ERR_YOUREBANNEDCREEP  = "465"
	ERR_KEYSET            = "467"
	ERR_CHANNELISFULL     = "471"
	ERR_UNKNOWNMODE       = "472"
	ERR_INVITEONLYCHAN    = "473"
	ERR_BANNEDFROMCHAN    = "474"
	ERR_BADCHANNELKEY     = "475"
	ERR_BADCHANMASK       = "476"
	ERR_NOCHANMODES       = "477"
	ERR_BANLISTFULL       = "478"
	ERR_NOPRIVILEGES      = "481"
	ERR_CHANOPRIVSNEEDED  = "482"
	ERR_CANTKILLSERVER    = "483"
	ERR_RESTRICTED        = "484"
	ERR_UNIQOPPRIVSNEEDED = "485"
	ERR_NOOPERHOST        = "491"
	ERR_UMODEUNKNOWNFLAG  = "501"
	ERR_USERSDONTMATCH    = "502"
	RPL_CUSTOM            = "999"
)

// NumericName maps IRC numerics to their human-readable names.
var NumericName = map[string]string{
	ERR_ALREADYREGISTRED:  "ERR_ALREADYREGISTRED",
	ERR_BADCHANMASK:       "ERR_BADCHANMASK",
	ERR_BADCHANNELKEY:     "ERR_BADCHANNELKEY",
	ERR_BADMASK:           "ERR_BADMASK",
	ERR_BANLISTFULL:       "ERR_BANLISTFULL",
	ERR_BANNEDFROMCHAN:    "ERR_BANNEDFROMCHAN",
	ERR_CANNOTSENDTOCHAN:  "ERR_CANNOTSENDTOCHAN",
	ERR_CANTKILLSERVER:    "ERR_CANTKILLSERVER",
	ERR_CHANNELISFULL:     "ERR_CHANNELISFULL",
	ERR_CHANOPRIVSNEEDED:  "ERR_CHANOPRIVSNEEDED",
	ERR_ERRONEUSNICKNAME:  "ERR_ERRONEUSNICKNAME",
	ERR_FILEERROR:         "ERR_FILEERROR",
	ERR_INVITEONLYCHAN:    "ERR_INVITEONLYCHAN",
	ERR_KEYSET:            "ERR_KEYSET",
	ERR_NEEDMOREPARAMS:    "ERR_NEEDMOREPARAMS",
	ERR_NICKCOLLISION:     "ERR_NICKCOLLISION",
	ERR_NICKNAMEINUSE:     "ERR_NICKNAMEINUSE",
	ERR_NOADMININFO:       "ERR_NOADMININFO",
	ERR_NOCHANMODES:       "ERR_NOCHANMODES",
	ERR_NOLOGIN:           "ERR_NOLOGIN",
	ERR_NOMOTD:            "ERR_NOMOTD",
	ERR_NONICKNAMEGIVEN:   "ERR_NONICKNAMEGIVEN",
	ERR_NOOPERHOST:        "ERR_NOOPERHOST",
	ERR_NOORIGIN:          "ERR_NOORIGIN",
	ERR_NOPERMFORHOST:     "ERR_NOPERMFORHOST",
	ERR_NOPRIVILEGES:      "ERR_NOPRIVILEGES",
	ERR_NORECIPIENT:       "ERR_NORECIPIENT",
	ERR_NOSUCHCHANNEL:     "ERR_NOSUCHCHANNEL",
	ERR_NOSUCHNICK:        "ERR_NOSUCHNICK",
	ERR_NOSUCHSERVER:      "ERR_NOSUCHSERVER",
	ERR_NOSUCHSERVICE:     "ERR_NOSUCHSERVICE",
	ERR_NOTEXTTOSEND:      "ERR_NOTEXTTOSEND",
	ERR_NOTONCHANNEL:      "ERR_NOTONCHANNEL",
	ERR_NOTOPLEVEL:        "ERR_NOTOPLEVEL",
	ERR_NOTREGISTERED:     "ERR_NOTREGISTERED",
	ERR_PASSWDMISMATCH:    "ERR_PASSWDMISMATCH",
	ERR_RESTRICTED:        "ERR_RESTRICTED",
	ERR_SUMMONDISABLED:    "ERR_SUMMONDISABLED",
	ERR_TOOMANYCHANNELS:   "ERR_TOOMANYCHANNELS",
	ERR_TOOMANYTARGETS:    "ERR_TOOMANYTARGETS",
	ERR_UMODEUNKNOWNFLAG:  "ERR_UMODEUNKNOWNFLAG",
	ERR_UNAVAILRESOURCE:   "ERR_UNAVAILRESOURCE",
	ERR_UNIQOPPRIVSNEEDED: "ERR_UNIQOPPRIVSNEEDED",
	ERR_UNKNOWNCOMMAND:    "ERR_UNKNOWNCOMMAND",
	ERR_UNKNOWNMODE:       "ERR_UNKNOWNMODE",
	ERR_USERNOTINCHANNEL:  "ERR_USERNOTINCHANNEL",
	ERR_USERONCHANNEL:     "ERR_USERONCHANNEL",
	ERR_USERSDISABLED:     "ERR_USERSDISABLED",
	ERR_USERSDONTMATCH:    "ERR_USERSDONTMATCH",
	ERR_WASNOSUCHNICK:     "ERR_WASNOSUCHNICK",
	ERR_WILDTOPLEVEL:      "ERR_WILDTOPLEVEL",
	ERR_YOUREBANNEDCREEP:  "ERR_YOUREBANNEDCREEP",
	RPL_ADMINEMAIL:        "RPL_ADMINEMAIL",
	RPL_ADMINME:           "RPL_ADMINME",
	RPL_AWAY:              "RPL_AWAY",
	RPL_BANLIST:           "RPL_BANLIST",
	RPL_BOUNCE:            "RPL_BOUNCE",
	RPL_CHANNELMODEIS:     "RPL_CHANNELMODEIS",
	RPL_CREATED:           "RPL_CREATED",
	RPL_CUSTOM:            "RPL_CUSTOM",
	RPL_ENDOFBANLIST:      "RPL_ENDOFBANLIST",
	RPL_ENDOFEXCEPTLIST:   "RPL_ENDOFEXCEPTLIST",
	RPL_ENDOFINFO:         "RPL_ENDOFINFO",
	RPL_ENDOFINVITELIST:   "RPL_ENDOFINVITELIST",
	RPL_ENDOFLINKS:        "RPL_ENDOFLINKS",
	RPL_ENDOFMOTD:         "RPL_ENDOFMOTD",
	RPL_ENDOFNAMES:        "RPL_ENDOFNAMES",
	RPL_ENDOFSTATS:        "RPL_ENDOFSTATS",
	RPL_ENDOFUSERS:        "RPL_ENDOFUSERS",
	RPL_ENDOFWHO:          "RPL_ENDOFWHO",
	RPL_ENDOFWHOIS:        "RPL_ENDOFWHOIS",
	RPL_ENDOFWHOWAS:       "RPL_ENDOFWHOWAS",
	RPL_EXCEPTLIST:        "RPL_EXCEPTLIST",
	RPL_INFO:              "RPL_INFO",
	RPL_INVITELIST:        "RPL_INVITELIST",
	RPL_INVITING:          "RPL_INVITING",
	RPL_ISON:              "RPL_ISON",
	RPL_ISUPPORT:          "RPL_ISUPPORT",
	RPL_LINKS:             "RPL_LINKS",
	RPL_LIST:              "RPL_LIST",
	RPL_LISTEND:           "RPL_LISTEND",
	RPL_LUSERCHANNELS:     "RPL_LUSERCHANNELS",
	RPL_LUSERCLIENT:       "RPL_LUSERCLIENT",
	RPL_LUSERME:           "RPL_LUSERME",
	RPL_LUSEROP:           "RPL_LUSEROP",
	RPL_LUSERUNKNOWN:      "RPL_LUSERUNKNOWN",
	RPL_MOTD:              "RPL_MOTD",
	RPL_MOTDSTART:         "RPL_MOTDSTART",
	RPL_MYINFO:            "RPL_MYINFO",
	RPL_NAMREPLY:          "RPL_NAMREPLY",
	RPL_NOTOPIC:           "RPL_NOTOPIC",
	RPL_NOUSERS:           "RPL_NOUSERS",
	RPL_NOWAWAY:           "RPL_NOWAWAY",
	RPL_REHASHING:         "RPL_REHASHING",
	RPL_SERVLIST:          "RPL_SERVLIST",
	RPL_SERVLISTEND:       "RPL_SERVLISTEND",
	RPL_STATSCOMMANDS:     "RPL_STATSCOMMANDS",
	RPL_STATSLINKINFO:     "RPL_STATSLINKINFO",
	RPL_STATSOLINE:        "RPL_STATSOLINE",
	RPL_STATSUPTIME:       "RPL_STATSUPTIME",
	RPL_SUMMONING:         "RPL_SUMMONING",
	RPL_TIME:              "RPL_TIME",
	RPL_TOPIC:             "RPL_TOPIC",
	RPL_TRACECLASS:        "RPL_TRACECLASS",
	RPL_TRACECONNECTING:   "RPL_TRACECONNECTING",
	RPL_TRACEEND:          "RPL_TRACEEND",
	RPL_TRACEHANDSHAKE:    "RPL_TRACEHANDSHAKE",
	RPL_TRACELINK:         "RPL_TRACELINK",
	RPL_TRACELOG:          "RPL_TRACELOG",
	RPL_TRACENEWTYPE:      "RPL_TRACENEWTYPE",
	RPL_TRACEOPERATOR:     "RPL_TRACEOPERATOR",
	RPL_TRACESERVER:       "RPL_TRACESERVER",
	RPL_TRACESERVICE:      "RPL_TRACESERVICE",
	RPL_TRACEUNKNOWN:      "RPL_TRACEUNKNOWN",
	RPL_TRACEUSER:         "RPL_TRACEUSER",
	RPL_TRYAGAIN:          "RPL_TRYAGAIN",
	RPL_UMODEIS:           "RPL_UMODEIS",
	RPL_UNAWAY:            "RPL_UNAWAY",
	RPL_UNIQOPIS:          "RPL_UNIQOPIS",
	RPL_USERHOST:          "RPL_USERHOST",
	RPL_USERS:             "RPL_USERS",
	RPL_USERSSTART:        "RPL_USERSSTART",
	RPL_VERSION:           "RPL_VERSION",
	RPL_WELCOME:           "RPL_WELCOME",
	RPL_WHOISCHANNELS:     "RPL_WHOISCHANNELS",
	RPL_WHOISIDLE:         "RPL_WHOISIDLE",
	RPL_WHOISOPERATOR:     "RPL_WHOISOPERATOR",
	RPL_WHOISSERVER:       "RPL_WHOISSERVER",
	RPL_WHOISUSER:         "RPL_WHOISUSER",
	RPL_WHOREPLY:          "RPL_WHOREPLY",
	RPL_WHOWASUSER:        "RPL_WHOWASUSER",
	RPL_YOUREOPER:         "RPL_YOUREOPER",
	RPL_YOURESERVICE:      "RPL_YOURESERVICE",
	RPL_YOURHOST:          "RPL_YOURHOST",
}

// NumericText maps IRC numerics to their text descriptions.
var NumericText = map[string]string{
	ERR_ALREADYREGISTRED:  `Unauthorized command (already registered)`,
	ERR_BADCHANMASK:       `<channel> :Bad Channel Mask`,
	ERR_BADCHANNELKEY:     `<channel> :Cannot join channel (+k)`,
	ERR_BADMASK:           `<mask> :Bad Server/host mask`,
	ERR_BANLISTFULL:       `<channel> <char> :Channel list is full`,
	ERR_BANNEDFROMCHAN:    `<channel> :Cannot join channel (+b)`,
	ERR_CANNOTSENDTOCHAN:  `<channel name> :Cannot send to channel`,
	ERR_CANTKILLSERVER:    `You can't kill a server!`,
	ERR_CHANNELISFULL:     `<channel> :Cannot join channel (+l)`,
	ERR_CHANOPRIVSNEEDED:  `<channel> :You're not channel operator`,
	ERR_ERRONEUSNICKNAME:  `<nick> :Erroneous nickname`,
	ERR_FILEERROR:         `File error doing <file op> on <file>`,
	ERR_INVITEONLYCHAN:    `<channel> :Cannot join channel (+i)`,
	ERR_KEYSET:            `<channel> :Channel key already set`,
	ERR_NEEDMOREPARAMS:    `<command> :Not enough parameters`,
	ERR_NICKCOLLISION:     `<nick> :Nickname collision KILL from <user>@<host>`,
	ERR_NICKNAMEINUSE:     `<nick> :Nickname is already in use`,
	ERR_NOADMININFO:       `<server> :No administrative info available`,
	ERR_NOCHANMODES:       `<channel> :Channel doesn't support modes`,
	ERR_NOLOGIN:           `<user> :User not logged in`,
	ERR_NOMOTD:            `MOTD File is missing`,
	ERR_NONICKNAMEGIVEN:   `No nickname given`,
	ERR_NOOPERHOST:        `No O-lines for your host`,
	ERR_NOORIGIN:          `No origin specified`,
	ERR_NOPERMFORHOST:     `Your host isn't among the privileged`,
	ERR_NOPRIVILEGES:      `Permission Denied- You're not an IRC operator`,
	ERR_NORECIPIENT:       `No recipient given (<command>)`,
	ERR_NOSUCHCHANNEL:     `<channel name> :No such channel`,
	ERR_NOSUCHNICK:        `<nickname> :No such nick/channel`,
	ERR_NOSUCHSERVER:      `<server name> :No such server`,
	ERR_NOSUCHSERVICE:     `<service name> :No such service`,
	ERR_NOTEXTTOSEND:      `No text to send`,
	ERR_NOTONCHANNEL:      `<channel> :You're not on that channel`,
	ERR_NOTOPLEVEL:        `<mask> :No toplevel domain specified`,
	ERR_NOTREGISTERED:     `You have not registered`,
	ERR_PASSWDMISMATCH:    `Password incorrect`,
	ERR_RESTRICTED:        `Your connection is restricted!`,
	ERR_SUMMONDISABLED:    `SUMMON has been disabled`,
	ERR_TOOMANYCHANNELS:   `<channel name> :You have joined too many channels`,
	ERR_TOOMANYTARGETS:    `<target> :<error code> recipients. <abort message>`,
	ERR_UMODEUNKNOWNFLAG:  `Unknown MODE flag`,
	ERR_UNAVAILRESOURCE:   `<nick/channel> :Nick/channel is temporarily unavailable`,
	ERR_UNIQOPPRIVSNEEDED: `You're not the original channel operator`,
	ERR_UNKNOWNCOMMAND:    `<command> :Unknown command`,
	ERR_UNKNOWNMODE:       `<char> :is unknown mode char to me for <channel>`,
	ERR_USERNOTINCHANNEL:  `<nick> <channel> :They aren't on that channel`,
	ERR_USERONCHANNEL:     `<user> <channel> :is already on channel`,
	ERR_USERSDISABLED:     `USERS has been disabled`,
	ERR_USERSDONTMATCH:    `Cannot change mode for other users`,
	ERR_WASNOSUCHNICK:     `<nickname> :There was no such nickname`,
	ERR_WILDTOPLEVEL:      `<mask> :Wildcard in toplevel domain`,
	ERR_YOUREBANNEDCREEP:  `You are banned from this server`,
	RPL_ADMINEMAIL:        `<admin info>`,
	RPL_ADMINME:           `<server> :Administrative info`,
	RPL_AWAY:              `<nick> :<away message>`,
	RPL_BANLIST:           `<channel> <banmask>`,
	RPL_BOUNCE:            `Try server <server name>, port <port number>`,
	RPL_CHANNELMODEIS:     `<channel> <mode> <mode params>`,
	RPL_CREATED:           `This server was created <date>`,
	RPL_CUSTOM:            `<param> <param> :Custom Numeric`,
	RPL_ENDOFBANLIST:      `<channel> :End of channel ban list`,
	RPL_ENDOFEXCEPTLIST:   `<channel> :End of channel exception list`,
	RPL_ENDOFINFO:         `End of INFO list`,
	RPL_ENDOFINVITELIST:   `<channel> :End of channel invite list`,
	RPL_ENDOFLINKS:        `<mask> :End of LINKS list`,
	RPL_ENDOFMOTD:         `End of MOTD command`,
	RPL_ENDOFNAMES:        `<channel> :End of NAMES list`,
	RPL_ENDOFSTATS:        `<stats letter> :End of STATS report`,
	RPL_ENDOFUSERS:        `End of users`,
	RPL_ENDOFWHO:          `<name> :End of WHO list`,
	RPL_ENDOFWHOIS:        `<nick> :End of WHOIS list`,
	RPL_ENDOFWHOWAS:       `<nick> :End of WHOWAS`,
	RPL_EXCEPTLIST:        `<channel> <exceptionmask>`,
	RPL_INFO:              `<string>`,
	RPL_INVITELIST:        `<channel> <invitemask>`,
	RPL_INVITING:          `<channel> <nick>`,
	RPL_ISON:              `*1<nick> *( " " <nick> )`,
	RPL_ISUPPORT:          `<supported> :are supported by this server`,
	RPL_LINKS:             `<mask> <server> :<hopcount> <server info>`,
	RPL_LIST:              `<channel> <# visible> :<topic>`,
	RPL_LISTEND:           `End of LIST`,
	RPL_LUSERCHANNELS:     `<integer> :channels formed`,
	RPL_LUSERCLIENT:       `There are <integer> users and <integer> services on <integer> servers`,
	RPL_LUSERME:           `I have <integer> clients and <integer> servers`,
	RPL_LUSEROP:           `<integer> :operator(s) online`,
	RPL_LUSERUNKNOWN:      `<integer> :unknown connection(s)`,
	RPL_MOTD:              `- <text>`,
	RPL_MOTDSTART:         `- <server> Message of the day - `,
	RPL_MYINFO:            `<servername> <version> <available user modes> <available channel modes>`,
	RPL_NAMREPLY:          `( "=" / "*" / "@" ) <channel> :[ "@" / "+" ] <nick> *( " " [ "@" / "+" ] <nick> )`,
	RPL_NOTOPIC:           `<channel> :No topic is set`,
	RPL_NOUSERS:           `Nobody logged in`,
	RPL_NOWAWAY:           `You have been marked as being away`,
	RPL_REHASHING:         `<config file> :Rehashing`,
	RPL_SERVLIST:          `<name> <server> <mask> <type> <hopcount> <info>`,
	RPL_SERVLISTEND:       `<mask> <type> :End of service listing`,
	RPL_STATSCOMMANDS:     `<command> <count> <byte count> <remote count>`,
	RPL_STATSLINKINFO:     `<linkname> <sendq> <sent messages> <sent Kbytes> <received messages> <received Kbytes> <time open>`,
	RPL_STATSOLINE:        `O <hostmask> * <name>`,
	RPL_STATSUPTIME:       `Server Up %d days %d:%02d:%02d`,
	RPL_SUMMONING:         `<user> :Summoning user to IRC`,
	RPL_TIME:              `<server> :<string showing server's local time>`,
	RPL_TOPIC:             `<channel> :<topic>`,
	RPL_TRACECLASS:        `Class <class> <count>`,
	RPL_TRACECONNECTING:   `Try. <class> <server>`,
	RPL_TRACEEND:          `<server name> <version & debug level> :End of TRACE`,
	RPL_TRACEHANDSHAKE:    `H.S. <class> <server>`,
	RPL_TRACELINK:         `Link <version & debug level> <destination> <next server> V<protocol version> <link uptime in seconds> <backstream sendq> <upstream sendq>`,
	RPL_TRACELOG:          `File <logfile> <debug level>`,
	RPL_TRACENEWTYPE:      `<newtype> 0 <client name>`,
	RPL_TRACEOPERATOR:     `Oper <class> <nick>`,
	RPL_TRACESERVER:       `Serv <class> <int>S <int>C <server> <nick!user|*!*>@<host|server> V<protocol version>`,
	RPL_TRACESERVICE:      `Service <class> <name> <type> <active type>`,
	RPL_TRACEUNKNOWN:      `???? <class> [<client IP address in dot form>]`,
	RPL_TRACEUSER:         `User <class> <nick>`,
	RPL_TRYAGAIN:          `<command> :Please wait a while and try again.`,
	RPL_UMODEIS:           `<user mode string>`,
	RPL_UNAWAY:            `You are no longer marked as being away`,
	RPL_UNIQOPIS:          `<channel> <nickname>`,
	RPL_USERHOST:          `*1<reply> *( " " <reply> )`,
	RPL_USERS:             `<username> <ttyline> <hostname>`,
	RPL_USERSSTART:        `UserID   Terminal  Host`,
	RPL_VERSION:           `<version>.<debuglevel> <server> :<comments>`,
	RPL_WELCOME:           `Welcome to the Internet Relay Network <nick>!<user>@<host>`,
	RPL_WHOISCHANNELS:     `<nick> :*( ( "@" / "+" ) <channel> " " )`,
	RPL_WHOISIDLE:         `<nick> <integer> :seconds idle`,
	RPL_WHOISOPERATOR:     `<nick> :is an IRC operator`,
	RPL_WHOISSERVER:       `<nick> <server> :<server info>`,
	RPL_WHOISUSER:         `<nick> <user> <host> * :<real name>`,
	RPL_WHOREPLY:          `<channel> <user> <host> <server> <nick> ( "H" / "G" > ["*"] [ ( "@" / "+" ) ] :<hopcount> <real name>`,
	RPL_WHOWASUSER:        `<nick> <user> <host> * :<real name>`,
	RPL_YOUREOPER:         `You are now an IRC operator`,
	RPL_YOURESERVICE:      `You are service <servicename>`,
	RPL_YOURHOST:          `Your host is <servername>, running version <ver>`,
}
