package game

import (
	"math/rand"
	"time"
)

// MonsterInputType exposes monster input to UI2D
type MonsterInputType int

const (
	NoInput MonsterInputType = iota
	KeyPress
)

// Monster is an enemy entity
type Monster struct {
	Character
	Typ MonsterInputType
}

// NewRat spawns a slow monster
// Why a map? Can iterate over maps fast, and access values by key
//   level.monsters[pos]
//   for key, value := range level.Monster { }
func NewRat(p Pos) *Monster {
	return &Monster{
		Character: Character{
			Entity: Entity{
				Pos:  p,
				Name: "Rat",
				Rune: 'R',
			},
			Hitpoints:    5,
			Strength:     5,
			Speed:        1.5,
			ActionPoints: 0.0,
			SightRange:   10.0,
			Items:        []*Item{NewCredits(Pos{})},
			PatternRNG:   rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}
}

// NewSpider spawns a fast monster
func NewSpider(p Pos) *Monster {
	return &Monster{
		Character: Character{
			Entity: Entity{
				Pos:  p,
				Name: "Spider",
				Rune: 'S',
			},
			Hitpoints:    100,
			Strength:     0,
			Speed:        1.0,
			ActionPoints: 0.0,
			SightRange:   10.0,
			PatternRNG:   rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}
}

// Update searches for player position
func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
	if len(positions) == 0 {
		// Nothing we can do, pass turn
		m.Pass()
		return
	}
	moveIndex := 1 // Move 1 position closer if we have a path, and we're not on top of the player (>1)
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

// Autoplay will automatically play the provided burst
func (m *Monster) Autoplay(level *Level) {
	// Pause before playing to simulate a real player
	amt := time.Duration(50) // 400-500ms
	time.Sleep(time.Millisecond * amt)
	if level.LastEvent == Attack {
		for {
			// Wait random interval
			amt := time.Duration(100 + rand.Intn(600)) // 100-700ms
			time.Sleep(time.Millisecond * amt)
			// Play note
			if len(m.Burst.Notes) > 0 {
				m.Typ = KeyPress
				m.Burst.Notes = m.Burst.Notes[1:]
				if len(m.Burst.Notes) == 0 {
					level.LastEvent = Damage
					level.ResolveDamage()
					return
				}
			}
		}
	}
}

// Pass prevents monsters from building up large sums of action points
func (m *Monster) Pass() {
	m.ActionPoints -= m.Speed
}

// Kill drops monster's items onto its current position
func (m *Monster) Kill(level *Level) {
	// Remove a monster from the map when it is dead.
	// It is safe to delete from a map while iterative over it. (cool!)
	delete(level.Monsters, m.Pos)

	groundItems := level.Items[m.Pos]
	for _, item := range m.Items {
		item.Pos = m.Pos
		groundItems = append(groundItems, item)
	}
	// TODO(max): will overwrite items on that tile
	level.Items[m.Pos] = groundItems
}

// Move moves towards the player position
func (m *Monster) Move(to Pos, level *Level) {
	if level.LastEvent != Attack && level.LastEvent != Damage {
		_, exists := level.Monsters[to] // Is there something at the position we want to move to?
		if !exists && to != level.Player.Pos {
			delete(level.Monsters, m.Pos) // Delete current, add new
			level.Monsters[to] = m
			m.Pos = to
			return
		}
		// If there is another monster in the way, don't attack the player
		if to == level.Player.Pos {
			level.Attack(&m.Character, &level.Player.Character)
			// Run blocking events in a seperate goroutine
			go func() {
				m.Autoplay(level)
			}()
		}
	}
}
