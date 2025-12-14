<script lang="ts">
  import Snowfall from "./lib/Snowfall.svelte";
  import type { AppState, Participant } from "./types";

  // --- STATE ---
  // --- SEPARATE STATES FOR BETTER REACTIVITY ---
  let userName = $state("");
  let roomCode = $state("");
  let currentView = $state("landing"); // 'landing' | 'lobby' | 'results'
  const christmasEmojis = [
    "🎅",
    "🎄",
    "🎁",
    "☃️",
    "❄️",
    "🦌",
    "🤶",
    "🔔",
    "🌟",
    "🍪",
    "🥛",
    "🧦",
    "🕯️",
    "🛷",
    "🎀",
  ];
  let selectedAvatar = $state(
    christmasEmojis[Math.floor(Math.random() * christmasEmojis.length)]
  );

  // This is the important one:
  let participants = $state<Participant[]>([]);
  let myTarget = $state<Participant | null>(null);

  let amIHost = $derived(
    participants.find((p) => p.name === userName)?.isHost ?? false
  );
  let hostName = $derived(participants.find((p) => p.isHost)?.name ?? "Host");
  // Derived state
  let isDrawButtonDisabled = $derived(participants.length < 2);

  let socket: WebSocket | null = null;

  function connectWebSocket(action: "create" | "join") {
    // 1. Log that we are trying to connect
    console.log(
      `Connecting to WS... Action: ${action}, Room: ${roomCode}, Name: ${userName}`
    );

    const baseUrl = import.meta.env.VITE_WS_URL || "ws://localhost:8007/ws";

    const wsUrl = `${baseUrl}?room=${roomCode}&name=${userName}&action=${action}&avatar=${encodeURIComponent(
      selectedAvatar
    )}`;
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log("✅ WebSocket Connected!");
      currentView = "lobby";
    };

    socket.onmessage = (event) => {
      // 2. Log every message received from Go
      // console.log("📩 Message received:", event.data);

      try {
        const msg = JSON.parse(event.data);

        if (msg.type === "UPDATE_PARTICIPANTS") {
          console.log("👥 Updating participants list to:", msg.payload);
          // 3. Direct assignment triggers the UI update
          participants = msg.payload;
        } else if (msg.type === "MATCH_RESULT") {
          const targetName = msg.payload;
          const targetObj = participants.find((p) => p.name === targetName);
          myTarget = targetObj || {
            name: targetName,
            avatar: "🎁",
            isHost: false,
          };
          currentView = "results";
        } else if (msg.type === "ERROR") {
          alert("Server Error: " + msg.payload);
          currentView = "landing";
        }
      } catch (err) {
        console.error("❌ Error parsing JSON:", err);
      }
    };

    socket.onerror = (error) => {
      console.error("❌ WebSocket Error:", error);
      alert("Connection failed! Is the Go server running?");
    };

    socket.onclose = () => {
      console.log("⚠️ WebSocket Disconnected");
    };
  }

  function createRoom() {
    if (!userName) return alert("Santa needs a name!");
    // Generate random 4-letter code
    roomCode = Math.random().toString(36).substring(2, 6).toUpperCase();
    // Clear list initially
    participants = [];
    connectWebSocket("create");
  }

  function joinRoom() {
    if (!userName || !roomCode) return alert("Check inputs!");
    participants = [];
    connectWebSocket("join");
  }

  function startGame() {
    if (socket && participants.length >= 2) {
      socket.send(JSON.stringify({ type: "START_GAME", payload: null }));
    } else {
      alert("Need at least 3 people to start!");
    }
  }
</script>

<Snowfall />

<main>
  <div class="center-wrapper">
    <h1 class="main-title">🎅 Secret Santa</h1>

    {#if currentView === "landing"}
      <div class="card snow-capped">
        <h2 class="card-title">Join the Holiday Fun</h2>

        <div class="avatar-grid">
          {#each christmasEmojis as emoji}
            <button
              class="avatar-btn {selectedAvatar === emoji ? 'selected' : ''}"
              onclick={() => (selectedAvatar = emoji)}
              aria-label="Choose avatar {emoji}"
            >
              {emoji}
            </button>
          {/each}
        </div>

        <div class="form-container">
          <div class="input-group">
            <label for="name">Your Name</label>
            <input
              type="text"
              id="name"
              placeholder="e.g. Elf Brown"
              bind:value={userName}
            />
          </div>

          <div class="divider"><span>THEN</span></div>

          <div class="input-group">
            <label for="room">Room Code</label>
            <input
              type="text"
              id="room"
              placeholder="e.g. XMAS"
              bind:value={roomCode}
            />
          </div>

          <div class="button-group">
            <button class="btn primary snow-capped-btn" onclick={createRoom}>
              Create New Room
            </button>
            <button class="btn secondary" onclick={joinRoom}>
              Join Room
            </button>
          </div>
        </div>
      </div>
    {:else if currentView === "lobby"}
      <div class="card snow-capped">
        <div class="lobby-header">
          <h2 class="card-title">
            Room: <span class="code-badge">{roomCode}</span>
          </h2>
          <p class="subtitle">Waiting for Santa...</p>
        </div>

        <div class="list-container">
          {#each participants as p}
            <div class="participant-row">
              <span class="avatar">{p.avatar}</span>
              <span class="p-name">{p.name} {p.isHost ? "👑" : ""}</span>
            </div>
          {/each}
        </div>

        {#if participants.length === 0}
          <div style="text-align:center; padding: 20px; color: #999;">
            Connecting...
          </div>
        {/if}

        <div class="button-group">
          {#if amIHost}
            <button
              class="btn primary snow-capped-btn"
              onclick={startGame}
              disabled={isDrawButtonDisabled}
            >
              {isDrawButtonDisabled ? "Need more people..." : "🎁 Draw Names"}
            </button>
            <p class="host-hint">You are the Host! Click above when ready.</p>
          {:else}
            <div class="waiting-message">
              <span class="pulse">⏳</span>
              Waiting for<strong>{hostName}</strong>to start the exchange...
            </div>
          {/if}
        </div>
      </div>
    {:else if currentView === "results"}
      <div class="card snow-capped result-card">
        <h2 class="card-title">You are the Secret Santa for...</h2>

        <div class="reveal-box">
          <div class="big-avatar">{myTarget?.avatar}</div>
          <div class="target-name">{myTarget?.name}</div>
        </div>

        <p class="note">Shhh! Keep it a secret until the party!</p>
        <button class="btn secondary" onclick={() => location.reload()}
          >Start Over</button
        >
      </div>
    {/if}
  </div>
</main>

<style>
  /* --- GLOBAL LAYOUT --- */
  :global(body) {
    margin: 0;
    font-family: "Roboto", sans-serif;
    background: #1a237e; /* Night Blue */
    color: #333;
  }

  main {
    min-height: 100vh;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    z-index: 10; /* Above snow */
    padding: 20px;
    box-sizing: border-box;
  }
  .host-hint {
    font-size: 0.8rem;
    color: #d32f2f;
    text-align: center;
    margin-top: 5px;
  }

  .waiting-message {
    background: #fff3e0; /* Light Orange/Gold */
    color: #e65100;
    padding: 15px;
    border-radius: 8px;
    text-align: center;
    font-weight: 500;
    border: 1px dashed #ffb74d;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
  }

  .pulse {
    display: inline-block;
    animation: pulse-anim 1.5s infinite;
  }

  @keyframes pulse-anim {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.2);
    }
    100% {
      transform: scale(1);
    }
  }

  .avatar-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr); /* 5 per row */
    gap: 10px;
    margin-bottom: 1rem;
    padding: 10px;
    background: #f5f5f5;
    border-radius: 12px;
  }

  .avatar-btn {
    background: white;
    border: 2px solid transparent;
    font-size: 1.5rem;
    padding: 5px;
    border-radius: 8px;
    cursor: pointer;
    transition:
      transform 0.1s,
      background 0.2s;
  }

  .avatar-btn:hover {
    background: #e3f2fd;
    transform: scale(1.1);
  }

  .avatar-btn.selected {
    background: #e8f5e9;
    border-color: #388e3c; /* Selected Green Border */
    transform: scale(1.1);
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  }

  /* Adjustments for existing classes */
  .avatar {
    font-size: 1.8rem;
  } /* Make lobby avatars bigger */

  .center-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    max-width: 450px; /* Limits width of card on big screens */
    gap: 20px;
  }

  .main-title {
    color: white;
    font-size: 3rem;
    margin: 0;
    text-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    text-align: center;
  }

  /* --- CARD STYLING --- */
  .card {
    background: white;
    width: 100%;
    padding: 2rem;
    border-radius: 16px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    position: relative;
    box-sizing: border-box;
  }

  .card-title {
    margin: 0;
    color: #b71c1c; /* Santa Red */
    text-align: center;
    font-size: 1.5rem;
  }

  .subtitle {
    text-align: center;
    color: #666;
    margin: 5px 0 0 0;
    font-size: 0.9rem;
  }

  /* --- SNOW TOPPING EFFECT (CSS ONLY) --- */
  .snow-capped::before {
    content: "";
    position: absolute;
    top: -10px;
    left: 0;
    right: 0;
    height: 15px;
    background: white;
    border-radius: 16px 16px 5px 5px;
    /* This creates the jagged snow look using radial gradients */
    background-image: radial-gradient(
      circle at 10px 10px,
      white 10px,
      transparent 11px
    );
    background-size: 20px 20px;
    background-repeat: repeat-x;
    background-position: 0 -5px;
    z-index: 1;
  }

  /* --- INPUT FIELDS --- */
  .form-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    text-align: left;
  }

  label {
    font-weight: 700;
    font-size: 0.85rem;
    color: #555;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  input {
    width: 100%;
    padding: 12px 15px;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    font-size: 1rem;
    transition: all 0.2s;
    box-sizing: border-box; /* Critical for padding inside width */
    background: #fdfdfd;
  }

  input:focus {
    border-color: #d32f2f;
    outline: none;
    box-shadow: 0 0 0 4px rgba(211, 47, 47, 0.1);
  }

  /* --- BUTTONS --- */
  .button-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 1rem;
  }

  .btn {
    width: 100%;
    padding: 14px;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: bold;
    cursor: pointer;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    transition: transform 0.1s;
    position: relative;
  }

  .btn:active {
    transform: scale(0.98);
  }

  .btn.primary {
    background: #d32f2f;
    color: white;
    box-shadow: 0 4px 0 #b71c1c; /* 3D effect */
  }

  .btn.primary:disabled {
    background: #ccc;
    box-shadow: none;
    cursor: not-allowed;
  }

  /* Specific snow for button */
  .snow-capped-btn::after {
    content: "❄";
    position: absolute;
    top: -8px;
    right: -5px;
    color: white;
    font-size: 1.2rem;
    transform: rotate(15deg);
  }

  .btn.secondary {
    background: transparent;
    border: 2px solid #388e3c; /* Elf Green */
    color: #388e3c;
  }

  .btn.secondary:hover {
    background: #e8f5e9;
  }

  /* --- DIVIDER --- */
  .divider {
    display: flex;
    align-items: center;
    color: #999;
    font-size: 0.8rem;
    font-weight: bold;
  }
  .divider::before,
  .divider::after {
    content: "";
    flex: 1;
    height: 1px;
    background: #eee;
  }
  .divider span {
    padding: 0 10px;
  }

  /* --- LOBBY LIST --- */
  .code-badge {
    background: #ffebee;
    color: #d32f2f;
    padding: 2px 8px;
    border-radius: 4px;
    font-family: monospace;
    border: 1px dashed #d32f2f;
  }

  .list-container {
    background: #f5f5f5;
    border-radius: 8px;
    padding: 10px;
    max-height: 250px;
    overflow-y: auto;
    border: 1px solid #eee;
  }

  .participant-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px;
    background: white;
    margin-bottom: 8px;
    border-radius: 6px;
    box-shadow: 0 2px 2px rgba(0, 0, 0, 0.05);
  }

  .participant-row:last-child {
    margin-bottom: 0;
  }

  .avatar {
    font-size: 1.2rem;
  }
  .p-name {
    font-weight: 500;
  }

  /* --- RESULT REVEAL --- */
  .reveal-box {
    background: #e8f5e9;
    padding: 2rem;
    border-radius: 50%;
    width: 200px;
    height: 200px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    border: 4px dashed #388e3c;
    animation: popIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .big-avatar {
    font-size: 4rem;
    margin-bottom: 10px;
  }
  .target-name {
    font-size: 1.5rem;
    font-weight: 800;
    color: #1b5e20;
    text-transform: uppercase;
  }

  .note {
    text-align: center;
    font-style: italic;
    color: #777;
    margin-top: 10px;
  }

  @keyframes popIn {
    from {
      transform: scale(0);
      opacity: 0;
    }
    to {
      transform: scale(1);
      opacity: 1;
    }
  }
</style>
