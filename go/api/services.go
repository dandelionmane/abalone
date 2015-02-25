package api

import (
	"io"

	"github.com/jinzhu/gorm"
)

type Services struct {
	Games   GamesService
	Matches MatchesService
	Players PlayersService
	Users   UsersService

	DB *gorm.DB
}

type GamesService interface {
	List() ([]Game, error)
	ListDetailled() ([]GameWithDetails, error)
}

type PlayersService interface {
	Upload(userID int64, p Player, executable io.Reader) (*Player, error)
	Create(userID int64, p Player) (*Player, error)
	List() ([]Player, error)
	Delete(id int64) error
}

type UsersService interface {
	Create(User) (*User, error)
	List() ([]User, error)
	Delete(id int64) error
}

type MatchesService interface {
	// Run creates a match and schedules it for execution.
	Run(playerID1, playerID2 int64) (*Match, error)

	// Create creates a match.
	// TODO Create(playerID1, playerID2 int64) (*Match, error)
}
