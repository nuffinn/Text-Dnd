# 🚀 The Artemis Expedition

A multiplayer, turn-based, text-based space adventure RPG inspired by Dungeons & Dragons mechanics. Set in a realistic near-future scenario where your crew must navigate the challenges of humanity's first asteroid mining expedition.

![Space Adventure](https://img.shields.io/badge/Genre-Space%20Adventure-blue)
![Players](https://img.shields.io/badge/Players-1--10-green)
![Turn Based](https://img.shields.io/badge/Gameplay-Turn%20Based-orange)

## 🎮 Features

- **Multiplayer Support**: Play with 1-10 players in a shared adventure
- **Turn-Based Gameplay**: Collaborative decision-making with voting system
- **Dice System**: Classic D&D-style dice mechanics (D6, D10, D20, D100)
- **Branching Storyline**: Multiple paths and endings based on your choices
- **Realistic Space Theme**: Grounded in actual space exploration challenges
- **Beautiful Terminal UI**: Retro-futuristic interface with animated starfield
- **Real-time Updates**: WebSocket-powered synchronization across all players
- **Role-Playing Elements**: Each player gets unique skills and roles

## 🛸 Story

**Year 2045**. You are crew aboard the ARTEMIS, humanity's first commercial asteroid mining vessel. Your target: asteroid 16 Psyche, a metallic asteroid worth $10,000 quadrillion in rare metals.

But when you arrive, you discover something unexpected... an ancient alien probe with a warning message. What you do next will determine not just the success of your mission, but potentially the fate of humanity.

## 🎯 Game Mechanics

### Skills System
Each player is randomly assigned skills in:
- 🔧 **Engineering**: Fix systems, repair equipment
- 🎮 **Piloting**: Navigate ships, control drones
- 🔬 **Science**: Analyze data, solve puzzles
- ⚕️ **Medical**: Treat injuries, maintain crew health
- ⭐ **Leadership**: Make hard decisions, inspire crew

### Dice Rolls
- Skill checks use **D20 + Skill Modifier**
- Success determined by meeting or beating Difficulty Class (DC)
- Manual dice rolling available for dramatic moments

### Resources Management
Monitor and manage:
- 🔋 **Power**: Ship systems and equipment
- 💨 **Oxygen**: Life support for the crew
- ⛽ **Fuel**: Propulsion and maneuvers
- 🍕 **Food**: Crew sustenance
- 🚀 **Ship Integrity**: Overall vessel condition

### Multiple Endings
Your choices lead to different outcomes:
- 🌟 **Legendary Ending**: "The Library of Worlds"
- 🏆 **Best Ending**: "The Seekers' Legacy"
- ⭐ **Optimal Ending**: "The Patient Explorer"
- ✨ **Good Ending**: "Those Who Return"
- 🎯 **Heroic Ending**: "No One Left Behind"
- 🔄 **Survival Ending**: "Against the Odds"
- ⚫ **Tragic Endings**: Various failure scenarios

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Modern web browser with WebSocket support

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/nuffinn/Text-Dnd.git
cd Text-Dnd
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the server**
```bash
go run .
```

4. **Open your browser**
```
http://localhost:8080
```

### For Multiplayer

1. **First player**: Create a new game and note the Game ID
2. **Other players**: Enter the same Game ID when joining
3. **Start the mission**: Once all players are ready, click "START MISSION"
4. **Play together**: Vote on decisions and watch the story unfold

## 🎲 How to Play

### Joining a Game

1. Enter your astronaut name
2. Leave Game ID empty to create a new game, OR enter an existing Game ID to join
3. Click "EMBARK ON MISSION"
4. Share the Game ID with your crew

### During the Game

1. **Read the Mission Log**: Follow the story as it unfolds
2. **Discuss with Crew**: Coordinate with your team (use external chat/voice)
3. **Vote on Actions**: Each player votes on what to do next
4. **Watch Skill Checks**: See if your crew's skills are enough to succeed
5. **Manage Resources**: Keep an eye on ship status
6. **Use Dice Roller**: Roll dice for fun or when prompted

### Tips for Success

- 🤝 **Coordinate**: Discuss decisions with your crew before voting
- 📊 **Watch Resources**: Low resources lead to harder challenges
- 🎲 **Know Your Odds**: Higher skills mean better chances of success
- ⏰ **Time Matters**: Some decisions have time pressure
- 💡 **Think Ahead**: Your choices compound over multiple rounds
- 🌟 **Take Risks**: The best endings often require bold decisions

## 🏗️ Technical Architecture

### Backend (Go)
- **main.go**: HTTP server and WebSocket handler
- **game.go**: Core game logic, state management, and player handling
- **story.go**: Story engine with all scenes and branching paths
- **dice.go**: Dice rolling mechanics

### Frontend
- **index.html**: Game UI structure
- **style.css**: Retro-futuristic terminal styling
- **game.js**: WebSocket client, game state management, UI updates

### Key Technologies
- **Go**: Backend server and game engine
- **WebSockets**: Real-time multiplayer communication
- **Gorilla WebSocket**: WebSocket library for Go
- **Vanilla JS**: Frontend without framework dependencies
- **CSS3**: Modern animations and effects

## 🎨 Customization

### Adding New Story Scenes

Edit `story.go` and add new scenes to the `initializeStory()` function:

```go
e.Scenes[99] = &Scene{
    ID: 99,
    Text: "Your custom story text here...",
    Choices: []Choice{
        {
            ID: "choice1",
            Text: "Option 1",
            SkillCheck: "engineering",
            Difficulty: 15,
        },
    },
    NextSuccess: 100,
    NextFailure: 101,
}
```

### Adjusting Game Balance

In `game.go`, modify:
- Initial resource values
- Skill value ranges
- Resource drain rates
- Difficulty class thresholds

### Styling Changes

Edit `static/style.css` to customize:
- Color scheme (CSS variables in `:root`)
- Layout and spacing
- Animations and effects

## 📋 API Reference

### WebSocket Messages

**Client → Server:**
```json
{
  "type": "start|choice|roll|vote",
  "playerId": "uuid",
  "data": {"..."}
}
```

**Server → Client:**
```json
{
  "type": "state|system|roll",
  "data": {"..."}
}
```

### REST Endpoints

- `POST /api/create` - Create a new game, returns `{ "gameId": "..." }`
- `GET /` - Serve the web interface
- `GET /ws` - WebSocket endpoint

## 🐛 Troubleshooting

### Connection Issues
- Ensure port 8080 is not in use
- Check firewall settings
- Verify WebSocket support in browser

### Game Not Starting
- At least one player must click "START MISSION"
- Ensure all players are connected (green status)
- Check browser console for errors

### Sync Issues
- Refresh the page to reconnect
- Ensure stable internet connection
- Check server logs for WebSocket errors

## 🤝 Contributing

Contributions welcome! Areas for improvement:
- Additional story branches and endings
- New game mechanics (inventory, combat, etc.)
- Mobile-responsive improvements
- AI-powered game master mode
- Persistence/save game functionality
- Character customization

## 📜 License

MIT License - feel free to use this project for learning or as a base for your own games!

## 🌟 Credits

Inspired by:
- **Dungeons & Dragons**: Classic tabletop RPG mechanics
- **Delta-V by Daniel Suarez**: Realistic asteroid mining fiction
- **The Expanse**: Grounded space opera storytelling
- **FTL: Faster Than Light**: Space crew decision-making

## 🎮 Future Enhancements

- [ ] Persistent game saves
- [ ] More story branches (Mars missions, space station scenarios)
- [ ] Character creation and customization
- [ ] Inventory and equipment system
- [ ] Voice chat integration
- [ ] Mobile app version
- [ ] Spectator mode
- [ ] Achievement system
- [ ] Steam/Discord integration

---

**Ready to reach for the stars?** 🚀

*"In space, no one can hear you scream... but they can hear you make tough decisions together."*
