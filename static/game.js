// Game state
let ws = null;
let gameId = null;
let playerId = null;
let playerName = null;
let currentState = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Artemis Expedition loaded');
});

// Join or create game
async function joinGame() {
    playerName = document.getElementById('player-name').value.trim();
    gameId = document.getElementById('game-id').value.trim();

    if (!playerName) {
        alert('Please enter your name!');
        return;
    }

    // Generate player ID
    playerId = generateUUID();

    // If no game ID provided, create new game
    if (!gameId) {
        try {
            const response = await fetch('/api/create', { method: 'POST' });
            const data = await response.json();
            gameId = data.gameId;
            console.log('Created new game:', gameId);
        } catch (error) {
            console.error('Error creating game:', error);
            alert('Failed to create game. Please try again.');
            return;
        }
    }

    // Connect via WebSocket
    connectWebSocket();
}

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('✅ Connected to server');
        updateConnectionStatus('🟢 CONNECTED');

        // Send initial connection message
        ws.send(JSON.stringify({
            type: 'connect',
            gameId: gameId,
            playerId: playerId,
            name: playerName
        }));

        // Show game screen
        document.getElementById('welcome-screen').style.display = 'none';
        document.getElementById('game-screen').style.display = 'block';

        // Show game ID for sharing
        addLogEntry('system', `✅ Connected to mission ${gameId}`);
        addLogEntry('system', `Share this Game ID with your crew: ${gameId}`);

        // Show start button if we're the first player
        document.getElementById('start-btn').style.display = 'block';
    };

    ws.onclose = () => {
        console.log('❌ Disconnected from server');
        updateConnectionStatus('🔴 DISCONNECTED');
        addLogEntry('error', '❌ Connection lost. Please refresh the page.');
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        addLogEntry('error', '❌ Connection error occurred.');
    };

    ws.onmessage = (event) => {
        try {
            const message = JSON.parse(event.data);
            handleMessage(message);
        } catch (error) {
            console.error('Error parsing message:', error);
        }
    };
}

function handleMessage(message) {
    console.log('📨 Received:', message);

    switch (message.type) {
        case 'state':
            updateGameState(message.data);
            break;
        case 'system':
            const systemData = JSON.parse(message.data);
            addLogEntry('system', systemData.text);
            break;
        case 'roll':
            const rollData = JSON.parse(message.data);
            addLogEntry('roll', rollData.text);
            break;
        default:
            console.log('Unknown message type:', message.type);
    }
}

function updateGameState(state) {
    currentState = state;
    console.log('🎮 State updated:', state);

    // Update round and phase
    document.getElementById('round-number').textContent = state.round || 0;
    document.getElementById('phase-name').textContent = state.phase || 'waiting';
    document.getElementById('game-status').textContent = (state.phase || 'WAITING').toUpperCase();

    // Update players
    updatePlayerList(state.players);

    // Update resources
    updateResources(state.resources);

    // Update ship status
    document.getElementById('ship-condition').textContent = state.shipCondition + '%';
    document.getElementById('mission-progress').textContent = state.missionProgress + '%';

    // Update story
    if (state.story) {
        document.getElementById('story-text').innerHTML = formatStory(state.story);
        // Scroll to top of story
        document.getElementById('story-text').scrollTop = 0;
    }

    // Update choices
    updateChoices(state.choices);

    // Hide start button if game started
    if (state.phase !== 'waiting') {
        document.getElementById('start-btn').style.display = 'none';
    }

    // Check for game end
    if (state.phase === 'ended') {
        addLogEntry('system', '🎬 Mission Complete! Check the Mission Log for your ending.');
        document.getElementById('choices-container').innerHTML = '<p class="text-muted">Mission concluded.</p>';
    }
}

function updatePlayerList(players) {
    const container = document.getElementById('player-list');

    if (!players || Object.keys(players).length === 0) {
        container.innerHTML = '<p class="text-muted">Waiting for crew members...</p>';
        return;
    }

    let html = '';
    for (const [id, player] of Object.entries(players)) {
        const isCurrentPlayer = id === playerId;
        const currentClass = player.hasVoted ? 'player-item' : 'player-item';

        html += `
            <div class="${currentClass} ${isCurrentPlayer ? 'slide-in' : ''}">
                <div class="player-name">${player.name} ${isCurrentPlayer ? '(You)' : ''}</div>
                <div class="player-role">${player.role}</div>
                <div class="player-stats">
                    <span>❤️ ${player.health}%</span>
                    <span>${player.hasVoted ? '✅ Voted' : '⏳ Waiting'}</span>
                </div>
                <div class="player-skills">
                    ${Object.entries(player.skills || {}).map(([skill, value]) =>
                        `<span>${getSkillIcon(skill)} ${skill}: ${value}</span>`
                    ).join('')}
                </div>
            </div>
        `;
    }

    container.innerHTML = html;
}

function updateResources(resources) {
    if (!resources) return;

    updateResourceBar('power', resources.power);
    updateResourceBar('oxygen', resources.oxygen);
    updateResourceBar('fuel', resources.fuel);
    updateResourceBar('food', resources.food);
}

function updateResourceBar(resource, value) {
    const bar = document.getElementById(`${resource}-bar`);
    const valueSpan = document.getElementById(`${resource}-value`);

    if (bar && valueSpan) {
        bar.style.width = value + '%';
        valueSpan.textContent = value + '%';

        // Change color based on value
        if (value < 30) {
            bar.classList.add('low');
        } else {
            bar.classList.remove('low');
        }
    }
}

function updateChoices(choices) {
    const container = document.getElementById('choices-container');

    if (!choices || choices.length === 0) {
        container.innerHTML = '<p class="text-muted">No actions available yet.</p>';
        return;
    }

    let html = '';
    choices.forEach(choice => {
        html += `
            <div class="choice-card slide-in" onclick="makeChoice('${choice.id}')">
                <div class="choice-header">
                    <div class="choice-text">${choice.text}</div>
                    ${choice.skillCheck ? `<div class="choice-skill">${getSkillIcon(choice.skillCheck)} ${choice.skillCheck.toUpperCase()}</div>` : ''}
                </div>
                ${choice.description ? `<div class="choice-description">${choice.description}</div>` : ''}
                ${choice.difficulty > 0 ? `<div class="choice-difficulty">Difficulty Check: ${choice.difficulty}</div>` : ''}
            </div>
        `;
    });

    container.innerHTML = html;
}

function formatStory(text) {
    // Add some formatting to the story text
    return text
        .replace(/⚠️/g, '<span style="color: var(--accent-yellow)">⚠️</span>')
        .replace(/🚀/g, '<span style="color: var(--accent-blue)">🚀</span>')
        .replace(/✅/g, '<span style="color: var(--accent-green)">✅</span>')
        .replace(/❌/g, '<span style="color: var(--accent-red)">❌</span>')
        .replace(/\n/g, '<br>');
}

function startGame() {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        addLogEntry('error', '❌ Not connected to server');
        return;
    }

    ws.send(JSON.stringify({
        type: 'start'
    }));

    addLogEntry('system', '🚀 Mission starting...');
}

function makeChoice(choiceId) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        addLogEntry('error', '❌ Not connected to server');
        return;
    }

    // Check if player already voted
    if (currentState && currentState.players[playerId]?.hasVoted) {
        addLogEntry('system', 'ℹ️ You have already voted this round');
        return;
    }

    ws.send(JSON.stringify({
        type: 'vote',
        data: JSON.stringify({ choiceId })
    }));

    addLogEntry('system', `✅ You voted for an action`);
}

function rollDice(diceType) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        addLogEntry('error', '❌ Not connected to server');
        return;
    }

    ws.send(JSON.stringify({
        type: 'roll',
        data: JSON.stringify({ dice: diceType })
    }));

    // Show rolling animation
    const resultDiv = document.getElementById('dice-result');
    resultDiv.textContent = '🎲 Rolling...';
    resultDiv.classList.add('pulse');

    setTimeout(() => {
        resultDiv.classList.remove('pulse');
    }, 1000);
}

function addLogEntry(type, text) {
    const log = document.getElementById('event-log');
    const entry = document.createElement('div');
    entry.className = `log-entry ${type}`;

    const timestamp = new Date().toLocaleTimeString();
    entry.innerHTML = `[${timestamp}] ${text}`;

    log.appendChild(entry);
    log.scrollTop = log.scrollHeight;

    // Also show roll results in dice panel
    if (type === 'roll') {
        const resultDiv = document.getElementById('dice-result');
        resultDiv.textContent = text;
    }
}

function updateConnectionStatus(status) {
    document.getElementById('connection-status').textContent = status;
}

function getSkillIcon(skill) {
    const icons = {
        'engineering': '🔧',
        'piloting': '🎮',
        'science': '🔬',
        'medical': '⚕️',
        'leadership': '⭐',
        'navigation': '🧭',
        'communications': '📡'
    };
    return icons[skill] || '⚙️';
}

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // Press 1-9 to quickly select choices
    if (e.key >= '1' && e.key <= '9') {
        const choiceIndex = parseInt(e.key) - 1;
        if (currentState && currentState.choices && currentState.choices[choiceIndex]) {
            makeChoice(currentState.choices[choiceIndex].id);
        }
    }

    // Press Enter on welcome screen to join
    if (e.key === 'Enter' && document.getElementById('welcome-screen').style.display !== 'none') {
        joinGame();
    }
});
