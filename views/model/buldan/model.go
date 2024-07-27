package buldan

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	ID   string
	Name string
}

type GameStatus int8

const DefaultRounds = 5

const (
	GameStatusLobby GameStatus = iota
	GameStatusRunning
	GameStatusFinished
)

type GameInstance struct {
	GameSettings
	ID                 string
	StartTime          time.Time
	Status             GameStatus
	Admin              Player
	Players            []Player
	CurrentLetter      rune
	CurrentPlayerIndex int
}

type GameSettings struct {
	rounds int
}

type Guess struct {
	GameID string
	Guess  string
	Player Player
}

type GuessResult struct {
	Correct       bool
	NextPlayer    Player
	NextLetter    rune
	RoundFinished bool
	GameFinished  bool
}

type DB interface {
	Get(ID string) (GameInstance, bool, error)

	CreateLobby(admin Player) (GameInstance, error)

	UpdateSettings(ID string, settings GameSettings) error

	JoinGame(ID string, player Player) (GameInstance, error)

	StartGame(ID string) error

	// Returns the next player that should play, along with the letter
	// or if the game is finished, returns false.
	Guess(ID string, player Player, guess string) (GuessResult, error)

	OpenLobbies() ([]GameInstance, error)
}

type inMemoryDB struct {
	instances     map[string]*GameInstance
	guessesByGame map[string][]Guess
}

var _ DB = &inMemoryDB{}

func NewInMemory() DB {
	return &inMemoryDB{
		instances:     make(map[string]*GameInstance),
		guessesByGame: make(map[string][]Guess),
	}
}


// Get implements DB.
func (i *inMemoryDB) Get(ID string) (GameInstance, bool, error) {
    game, ok := i.instances[ID]

    if !ok {
        return GameInstance{}, false, nil
    }

    return *game, ok, nil
}


func (i *inMemoryDB) nextGameID() string {
	id := ""

	// find first unused ID
	for {
		if _, ok := i.instances[id]; !ok && id != "" {
			return id
		}

		raw := make([]byte, 6)
		if _, err := crand.Read(raw); err != nil {
			panic(err)
		}
		id = fmt.Sprintf("%X", raw)
	}
}

// CreateLobby implements DB.
func (i *inMemoryDB) CreateLobby(admin Player) (GameInstance, error) {
	game := GameInstance{
		GameSettings: GameSettings{
			rounds: DefaultRounds,
		},
		ID:        i.nextGameID(),
		StartTime: time.Now(),
		Status:    GameStatusLobby,
		Admin:     admin,
		Players:   []Player{},
	}

	i.instances[game.ID] = &game
	i.guessesByGame[game.ID] = make([]Guess, 0)
	return game, nil
}

func (i *inMemoryDB) advancePlayer(gi *GameInstance) Player {
	gi.CurrentPlayerIndex++
	if gi.CurrentPlayerIndex >= len(gi.Players) {
		gi.CurrentPlayerIndex = 0
	}

	return gi.Players[gi.CurrentPlayerIndex]
}

// Guess implements DB.
func (i *inMemoryDB) Guess(ID string, player Player, guess string) (GuessResult, error) {
	result := GuessResult{}
	game, ok := i.instances[ID]
	if !ok {
		return result, fmt.Errorf("No game found with ID: %s", ID)
	}

	expectedPlayer := game.Players[game.CurrentPlayerIndex]
	if expectedPlayer.ID != player.ID {
		return result, fmt.Errorf("Expected player %s to play but it was %s", expectedPlayer.ID, player.ID)
	}

	i.guessesByGame[game.ID] = append(i.guessesByGame[game.ID], Guess{
		GameID: game.ID,
		Guess:  guess,
		Player: player,
	})

	// TODO: implement whether guess is correct or not
	// TODO: implement if we guessed everything or not

	return GuessResult{
		Correct:       true,
		NextPlayer:    i.advancePlayer(game),
		NextLetter:    game.CurrentLetter,
		RoundFinished: false,
		GameFinished:  false,
	}, nil
}

// JoinGame implements DB.
func (i *inMemoryDB) JoinGame(ID string, player Player) (GameInstance, error) {
	// TODO: check duplicate IDs
	game, ok := i.instances[ID]
	if !ok {
		return GameInstance{}, fmt.Errorf("No game found with ID: %s", ID)
	}

	for _, p := range game.Players {
		if p.ID == player.ID {
			return GameInstance{},fmt.Errorf("A player is already registered with this ID: %s", p.ID)
		}
	}

	game.Players = append(game.Players, player)
	return *game, nil
}

// OpenLobbies implements DB.
func (i *inMemoryDB) OpenLobbies() ([]GameInstance, error) {
	lobbies := make([]GameInstance, 0)

	for _, game := range i.instances {
		if game.Status == GameStatusLobby {
			lobbies = append(lobbies, *game)
		}
	}
	return lobbies, nil
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// StartGame implements DB.
func (i *inMemoryDB) StartGame(ID string) error {
	game, ok := i.instances[ID]
	if !ok {
		return fmt.Errorf("No game found with ID: %s", ID)
	}

	game.Status = GameStatusRunning
	game.CurrentLetter = rune(letters[rand.Intn(len(letters))])

	return nil
}

// UpdateSettings implements DB.
func (i *inMemoryDB) UpdateSettings(ID string, settings GameSettings) error {
	game, ok := i.instances[ID]
	if !ok {
		return fmt.Errorf("No game found with ID: %s", ID)
	}

	game.GameSettings = settings
	return nil
}
