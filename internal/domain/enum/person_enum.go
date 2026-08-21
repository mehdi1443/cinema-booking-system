package enum

import "strings"

type MoviePersonRole string

const (
	RoleActor    MoviePersonRole = "ACTOR"
	RoleDirector MoviePersonRole = "DIRECTOR"
	RoleWriter   MoviePersonRole = "WRITER"
)

func (r MoviePersonRole) String() string {
	return string(r)
}

func AllMoviePersonRoles() []MoviePersonRole {
	return []MoviePersonRole{
		RoleActor,
		RoleDirector,
		RoleWriter,
	}
}

func (r MoviePersonRole) IsValid() bool {
	for _, role := range AllMoviePersonRoles() {
		if r == role {
			return true
		}
	}
	return false
}

func ParseMoviePersonRole(str string) (MoviePersonRole, bool) {
	parsed := MoviePersonRole(strings.ToUpper(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}
