package parser

import (
	"strings"
	"unicode/utf8"
)

func isletter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
func isdigit(r rune) bool {
	return (r >= '0' && r <= '9')
}
func isspecial(r rune) bool {
	return (r >= '[' && r <= '`') || (r >= '{' && r <= '}')
}

// TODO(kevlar): Tests
func ToLower(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= '^' {
			return r - 'A' + 'a'
		}
		return r
	}, str)
}

func ToUpper(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= '~' {
			return r - 'a' + 'A'
		}
		return r
	}, str)
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
	first, _ := utf8.DecodeRuneInString(pfx)
	if !isdigit(first) {
		return false
	}
	for _, r := range pfx {
		if !isletter(r) && !isdigit(r) {
			return false
		}
	}
	return true
}

func ValidNick(str string) bool {
	if len(str) == 0 {
		return false
	}
	first, _ := utf8.DecodeRuneInString(str)
	if isdigit(first) || first == '-' {
		return false
	}
	for _, r := range str {
		if !isletter(r) && !isdigit(r) && !isspecial(r) && r != '-' {
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
	for _, r := range str {
		switch r {
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
	return strings.Map(func(r rune) rune {
		if r >= '!' && r <= '~' && r != ':' {
			return r
		}
		return -1
	}, str)
}
