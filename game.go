package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID       string
	Name     string
	Conn     *websocket.Conn
	Game     *Game
	Role     string
	Health   int
	Skills   map[string]int
	IsAlive  bool
	HasVoted bool
}

type Game struct {
	ID              string
	Players         map[string]*Player
	State           *GameState
	Register        chan *Player
	Unregister      chan *Player
	Broadcast       chan Message
	mu              sync.RWMutex
	CurrentRound    int
	CurrentTurn     int
	GameStarted     bool
	StoryEngine     *StoryEngine
}

type GameState struct {
	Round           int                 `json:"round"`
	Turn            int                 `json:"turn"`
	CurrentPlayer   string              `json:"currentPlayer"`
	Phase           string              `json:"phase"`
	Story           string              `json:"story"`
	Choices         []Choice            `json:"choices"`
	Players         map[string]PlayerState `json:"players"`
	Resources       Resources           `json:"resources"`
	ShipCondition   int                 `json:"shipCondition"`
	MissionProgress int                 `json:"missionProgress"`
	TimeRemaining   int                 `json:"timeRemaining"`
}

type PlayerState struct {
	Name     string         `json:"name"`
	Role     string         `json:"role"`
	Health   int            `json:"health"`
	Skills   map[string]int `json:"skills"`
	IsAlive  bool           `json:"isAlive"`
	HasVoted bool           `json:"hasVoted"`
}

type Choice struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	SkillCheck  string `json:"skillCheck"`
	Difficulty  int    `json:"difficulty"`
	Description string `json:"description"`
}

type Resources struct {
	Oxygen int `json:"oxygen"`
	Fuel   int `json:"fuel"`
	Power  int `json:"power"`
	Food   int `json:"food"`
}

type Message struct {
	Type     string          `json:"type"`
	PlayerID string          `json:"playerId"`
	Data     json.RawMessage `json:"data"`
}

func NewGame(id string) *Game {
	return &Game{
		ID:          id,
		Players:     make(map[string]*Player),
		State:       &GameState{
			Round:         0,
			Phase:         "waiting",
			Players:       make(map[string]PlayerState),
			Resources:     Resources{Oxygen: 100, Fuel: 100, Power: 100, Food: 100},
			ShipCondition: 100,
			TimeRemaining: 100,
		},
		Register:    make(chan *Player),
		Unregister:  make(chan *Player),
		Broadcast:   make(chan Message),
		StoryEngine: NewStoryEngine(),
	}
}

func (g *Game) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case player := <-g.Register:
			g.mu.Lock()
			g.Players[player.ID] = player
			player.Health = 100
			player.IsAlive = true
			player.Skills = map[string]int{
				"engineering": rand.Intn(5) + 1,
				"piloting":    rand.Intn(5) + 1,
				"science":     rand.Intn(5) + 1,
				"medical":     rand.Intn(5) + 1,
				"leadership":  rand.Intn(5) + 1,
			}

			// Assign roles
			roles := []string{"Captain", "Engineer", "Pilot", "Scientist", "Medic", "Navigator", "Communications", "Security", "Geologist", "Technician"}
			if len(g.Players) <= len(roles) {
				player.Role = roles[len(g.Players)-1]
			} else {
				player.Role = "Crew Member"
			}

			g.State.Players[player.ID] = PlayerState{
				Name:    player.Name,
				Role:    player.Role,
				Health:  player.Health,
				Skills:  player.Skills,
				IsAlive: player.IsAlive,
			}
			g.mu.Unlock()

			g.broadcastState()
			log.Printf("Player %s joined game %s as %s", player.Name, g.ID, player.Role)

		case player := <-g.Unregister:
			g.mu.Lock()
			if _, ok := g.Players[player.ID]; ok {
				delete(g.Players, player.ID)
				delete(g.State.Players, player.ID)
			}
			g.mu.Unlock()
			g.broadcastState()

		case msg := <-g.Broadcast:
			g.handleMessage(msg)

		case <-ticker.C:
			// Game tick for time-based events
			if g.GameStarted && g.State.Phase == "action" {
				// Could add time pressure here
			}
		}
	}
}

func (g *Game) handleMessage(msg Message) {
	switch msg.Type {
	case "start":
		g.startGame()
	case "choice":
		var data struct {
			ChoiceID string `json:"choiceId"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Error unmarshaling choice: %v", err)
			return
		}
		g.handleChoice(msg.PlayerID, data.ChoiceID)
	case "roll":
		var data struct {
			Dice string `json:"dice"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Error unmarshaling roll: %v", err)
			return
		}
		g.handleRoll(msg.PlayerID, data.Dice)
	case "vote":
		var data struct {
			ChoiceID string `json:"choiceId"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Error unmarshaling vote: %v", err)
			return
		}
		g.handleVote(msg.PlayerID, data.ChoiceID)
	}
}

func (g *Game) startGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.GameStarted || len(g.Players) == 0 {
		return
	}

	g.GameStarted = true
	g.CurrentRound = 1
	g.State.Round = 1
	g.State.Phase = "story"

	// Start the story
	scene := g.StoryEngine.GetScene(0)
	g.State.Story = scene.Text
	g.State.Choices = scene.Choices

	g.broadcastState()
	g.broadcastMessage("system", "🚀 The Artemis Expedition begins! May fortune favor the bold...")
}

func (g *Game) handleChoice(playerID, choiceID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.Players[playerID]
	if !exists {
		return
	}

	// Find the choice
	var selectedChoice *Choice
	for i := range g.State.Choices {
		if g.State.Choices[i].ID == choiceID {
			selectedChoice = &g.State.Choices[i]
			break
		}
	}

	if selectedChoice == nil {
		return
	}

	// Perform skill check if needed
	result := g.performSkillCheck(player, selectedChoice)

	// Advance story based on result
	g.advanceStory(choiceID, result)
}

func (g *Game) handleVote(playerID, choiceID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.Players[playerID]
	if !exists || player.HasVoted {
		return
	}

	player.HasVoted = true
	state := g.State.Players[playerID]
	state.HasVoted = true
	g.State.Players[playerID] = state

	// For simplicity, the first vote determines the action
	// In a full implementation, you'd count votes
	var selectedChoice *Choice
	for i := range g.State.Choices {
		if g.State.Choices[i].ID == choiceID {
			selectedChoice = &g.State.Choices[i]
			break
		}
	}

	if selectedChoice != nil {
		g.broadcastMessage("system", player.Name+" voted for: "+selectedChoice.Text)
	}

	// Check if all players voted
	allVoted := true
	for _, p := range g.Players {
		if !p.HasVoted {
			allVoted = false
			break
		}
	}

	if allVoted {
		// Reset votes and process the choice (using first vote for simplicity)
		for _, p := range g.Players {
			p.HasVoted = false
		}
		g.handleChoice(playerID, choiceID)
	}
}

func (g *Game) performSkillCheck(player *Player, choice *Choice) bool {
	if choice.SkillCheck == "" {
		return true
	}

	skillValue := player.Skills[choice.SkillCheck]
	roll := RollDice("1d20")
	total := roll + skillValue

	success := total >= choice.Difficulty

	msg := ""
	if success {
		msg = fmt.Sprintf("✅ %s rolled %d + %d (%s) = %d vs DC %d - SUCCESS!",
			player.Name, roll, skillValue, choice.SkillCheck, total, choice.Difficulty)
	} else {
		msg = fmt.Sprintf("❌ %s rolled %d + %d (%s) = %d vs DC %d - FAILED!",
			player.Name, roll, skillValue, choice.SkillCheck, total, choice.Difficulty)
	}

	g.broadcastMessage("roll", msg)
	return success
}

func (g *Game) handleRoll(playerID, diceType string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.Players[playerID]
	if !exists {
		return
	}

	result := RollDice(diceType)
	g.broadcastMessage("roll", fmt.Sprintf("🎲 %s rolled %s: %d", player.Name, diceType, result))
}

func (g *Game) advanceStory(choiceID string, success bool) {
	// Update game state based on choice
	currentScene := g.StoryEngine.GetCurrentScene()
	nextSceneID := currentScene.GetNextScene(choiceID, success)

	g.StoryEngine.CurrentScene = nextSceneID
	nextScene := g.StoryEngine.GetScene(nextSceneID)

	g.CurrentRound++
	g.State.Round = g.CurrentRound
	g.State.Story = nextScene.Text
	g.State.Choices = nextScene.Choices
	g.State.MissionProgress += 10

	// Apply consequences
	if !success {
		g.State.Resources.Oxygen -= 10
		g.State.ShipCondition -= 5
	} else {
		g.State.MissionProgress += 5
	}

	// Check for ending
	if nextScene.IsEnding {
		g.State.Phase = "ended"
		g.broadcastMessage("system", "🎬 "+nextScene.EndingType)
	}

	// Reset votes
	for pid := range g.State.Players {
		state := g.State.Players[pid]
		state.HasVoted = false
		g.State.Players[pid] = state
	}

	g.broadcastState()
}

func (g *Game) broadcastState() {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, player := range g.Players {
		if err := player.Conn.WriteJSON(Message{
			Type: "state",
			Data: g.marshalState(),
		}); err != nil {
			log.Printf("Error sending state: %v", err)
		}
	}
}

func (g *Game) broadcastMessage(msgType, text string) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	data, _ := json.Marshal(map[string]string{"text": text})
	for _, player := range g.Players {
		if err := player.Conn.WriteJSON(Message{
			Type: msgType,
			Data: data,
		}); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func (g *Game) marshalState() json.RawMessage {
	data, _ := json.Marshal(g.State)
	return data
}
