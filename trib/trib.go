// Package trib defines basic interfaces and constants
// for Tribbler service implementation.
package trib

const (
	MaxUsernameLen = 15   // Maximum length of a username
	MaxTribLen     = 140  // Maximum length of a tribble
	MaxTribFetch   = 100  // Maximum count of tribbles for Home() and Tribs()
	MinListUser    = 20   // Minimum count of users required for ListUsers()
	MaxFollowing   = 2000 // Maximum count of users that one can follow
)

type Log struct {
	Version   uint64    // server clock
	Op string
	Pos int
	Content string
}

type Server interface {
	Hello() error

	SearchFile(user, filename string) ([]string, string, error)

	UpdateFile(filename string, cmd Log) error

	Latest(filename string, version uint64) string

	LogoutUser(username string) error

}

// Checks if a username is a valid one. Returns true if it is.
func IsValidUsername(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > MaxUsernameLen {
		return false
	}

	for i, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}

		if i > 0 && r >= '0' && r <= '9' {
			continue
		}

		return false
	}

	return true
}
