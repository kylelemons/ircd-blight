package parser

import (
	"strings"
)

func isletter(rune int) bool {
	return (rune >= 'a' && rune <= 'z') || (rune >= 'A' && rune <= 'Z')
}
func isdigit(rune int) bool {
	return (rune >= '0' && rune <= '9')
}
func isspecial(rune int) bool {
	return (rune >= '[' && rune <= '`') || (rune >= '{' && rune <= '}')
}

// TODO(kevlar): Tests
func ToLower(str string) string {
	return strings.Map(func(rune int) int {
		if rune >= 'A' && rune <= '^' {
			return rune - 'A' + 'a'
		}
		return rune
	},str)
}

func ToUpper(str string) string {
	return strings.Map(func(rune int) int {
		if rune >= 'a' && rune <= '~' {
			return rune - 'a' + 'A'
		}
		return rune
	},str)
}

func ValidServerName(str string) bool {
	if len(str) == 0 {
		return false
	}
	dot := 0
	for _, rune := range str {
		if !isletter(rune) && !isdigit(rune) && !isspecial(rune) && rune != '-' && rune != '.' {
			return false
		}
		if rune == '.' {
			dot++
		}
	}
	return dot > 0
}

func ValidServerPrefix(pfx string) bool {
	if len(pfx) != 3 {
		return false
	}
	if !isdigit(int(pfx[0])) {
		return false
	}
	for _, rune := range pfx {
		if !isletter(rune) && !isdigit(rune) {
			return false
		}
	}
	return true
}

func ValidNick(str string) bool {
	if len(str) == 0 {
		return false
	}
	if isdigit(int(str[0])) || str[0] == '-' {
		return false
	}
	for _, rune := range str {
		if !isletter(rune) && !isdigit(rune) && !isspecial(rune) && rune != '-' {
			return false
		}
	}
	return true
}

func ValidChannel(str string) bool {
	if len(str) == 0 {
		return false
	}
	if str[0] != '#' {
		return false
	}
	for _, rune := range str {
		switch rune {
		case 0x00:
			return false
		case 0x07:
			return false
		case '\r':
			return false
		case '\n':
			return false
		case ' ':
			return false
		case ',':
			return false
		case ':':
			return false
		}
	}
	return true
}

func StripUnsafe(str string) string {
	return strings.Map(func(rune int) int {
		if rune >= '!' && rune <= '~' && rune != ':' {
			return rune
		}
		return -1
	},str)
}
