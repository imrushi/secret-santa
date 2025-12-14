# 🎅 Real-time Secret Santa Application

A festive, real-time web application to organize Secret Santa gift exchanges with friends and family. Built with **Svelte 5** for a reactive UI and **Golang** for high-performance WebSocket communication.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Svelte](https://img.shields.io/badge/Svelte-5-orange.svg)

## ✨ Features

- **Real-time Lobby:** Watch friends join the room instantly without refreshing.
- **Zero Authentication:** No emails or passwords required—just enter a name and play.
- **Secret Matching Algorithm:** Ensures a perfect "Derangement" (no one draws themselves) using a Hamiltonian path logic.
- **Festive UI:** "Snow-capped" Material Design interface with falling snow animations.
- **Avatar Selection:** Choose from 15 custom Christmas avatars.
- **Host Controls:** Only the room creator can start the draw.

## 🛠️ Tech Stack

### Frontend

- **Framework:** Svelte 5 (using Runes for state management)
- **Language:** TypeScript
- **Build Tool:** Vite
- **Styling:** Custom CSS with CSS Variables (Snow & Glassmorphism effects)

### Backend

- **Language:** Go (Golang)
- **Communication:** WebSockets (`gorilla/websocket`)
- **Router:** Standard `net/http` ServeMux
- **Storage:** In-memory maps (fast and simple for temporary lobbies)

---

## 🚀 Getting Started

Follow these instructions to run the project locally.

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.18 or higher)
- [Node.js](https://nodejs.org/) (version 18 or higher)

### 1. Backend Setup (Go)

1.  Navigate to the backend folder:
    ```bash
    cd santa-backend
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Start the server:
    ```bash
    go run main.go
    ```
    _The server will start on `http://localhost:8007`_

### 2. Frontend Setup (Svelte)

1.  Open a new terminal and navigate to the frontend folder:
    ```bash
    cd secret-santa-ui
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the development server:
    ```bash
    npm run dev
    ```
4.  Open your browser to the URL shown (usually `http://localhost:5173`).

---

## 🎮 How to Play

1.  **Create a Room:**

    - Open the app.
    - Select your festive avatar (🎅, 🦌, ☃️, etc.).
    - Enter your Name and click **"Create Room"**.
    - Share the 4-digit **Room Code** (e.g., `XJ9Z`) with friends.

2.  **Join a Room:**

    - Friends enter their Name and the **Room Code**.
    - They click **"Join Room"** and appear in the lobby instantly.

3.  **The Draw:**
    - Once everyone has joined (minimum 3 participants), the **Host** clicks the **"Draw Names"** button.
    - Everyone's screen instantly flips to reveal their secret target! 🎁

---

## 📂 Project Structure

```text
.
├── santa-backend/          # Golang Server
│   ├── go.mod
│   ├── go.sum
│   └── main.go             # Entry point, WebSocket handler, Game Logic
│
└── secret-santa-ui/        # Svelte Frontend
    ├── src/
    │   ├── lib/
    │   │   └── Snowfall.svelte  # Background animation component
    │   ├── App.svelte      # Main Application View
    │   ├── types.ts        # TypeScript Interfaces
    │   └── app.css         # Global styles & Snow themes
    └── package.json
```
