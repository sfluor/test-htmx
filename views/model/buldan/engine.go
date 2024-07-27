package buldan

import (
	"fmt"
	"sync"
)

type Engine interface {
	NewGame(admin string) (GameInstance, error)
	JoinGame(ID string, playerID string) (GameInstance, <-chan Event, error)
}

type Event struct {
	raw string
}

type engine struct {
	sync.Mutex

	db DB
	// ID -> Player
	channels map[string]map[Player]chan Event
}

// JoinGame implements Engine.
func (e *engine) JoinGame(ID string, name string) (GameInstance, <-chan Event, error) {
	e.Lock()
	defer e.Unlock()

	newPlayer := player(name)

	if _, ok := e.channels[ID]; !ok {
		return GameInstance{}, nil, fmt.Errorf("Unknown game: %s", ID)
	}

	if _, ok := e.channels[ID][newPlayer]; ok {
		return GameInstance{}, nil, fmt.Errorf("Player already exists and is connected: %+v", newPlayer)
	}

	game, err := e.db.JoinGame(ID, newPlayer)
	if err != nil {
		return GameInstance{}, nil, fmt.Errorf("failed to join game: %s", err)
	}

	e.broadcast(ID, newPlayer, Event{raw: "New player joined: " + name})

	ch := make(chan Event, 10)
	e.channels[ID] = map[Player]chan Event{
		newPlayer: ch,
	}

	return game, ch, nil
}

func (e *engine) broadcast(ID string, from Player, event Event) {
	for player, ch := range e.channels[ID] {
		if player.ID != from.ID {
			ch <- event
		}
	}
}

// NewGame implements Engine.
func (e *engine) NewGame(admin string) (GameInstance, error) {
	e.Lock()
	defer e.Unlock()
	player := player(admin)
	game, err := e.db.CreateLobby(player)
	if err != nil {
		return game, err
	}

	e.channels[game.ID] = map[Player]chan Event{ }

	return game, nil
}

func NewDefaultEngine() Engine {
	return NewEngine(NewInMemory())
}

func NewEngine(db DB) Engine {
	return &engine{
		db:       db,
		channels: make(map[string]map[Player]chan Event),
	}
}

func player(name string) Player {
	// TODO: assign ID properly
	return Player{
		ID:   name,
		Name: name,
	}
}
