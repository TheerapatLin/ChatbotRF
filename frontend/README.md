# ChatBot Frontend

Vue 3 + Vite frontend application for the ChatBot project.

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Next generation frontend tooling
- **Pinia** - State management
- **Axios** - HTTP client
- **WebSocket** - Real-time communication

## Setup Instructions

### 1. Install Dependencies
```bash
npm install
```

### 2. Run Development Server
```bash
npm run dev
```
App available at `http://localhost:5173`

### 3. Build for Production
```bash
npm run build
```

## Project Structure

```
frontend/
├── src/
│   ├── api/                 # API service layer
│   │   ├── axios.js        # Axios configuration
│   │   ├── personaService.js
│   │   ├── chatService.js
│   │   ├── fileService.js
│   │   └── audioService.js
│   ├── stores/             # Pinia stores
│   │   ├── chat.js         # Chat state management
│   │   ├── persona.js      # Persona state management
│   │   ├── audio.js        # Audio state management
│   │   └── file.js         # File state management
│   ├── services/           # Business logic services
│   │   └── websocketService.js
│   └── main.js
└── vite.config.js
```

## Features Completed (Task 1)

- ✅ Vue 3 + Vite project setup
- ✅ Axios configuration with interceptors
- ✅ Pinia state management setup
- ✅ API service layer (persona, chat, file, audio)
- ✅ Pinia stores (chat, persona, audio, file)
- ✅ WebSocket service with auto-reconnect
- ✅ Environment configuration (.env files)
- ✅ Vite proxy setup for API calls

## Next Steps

Continue with **Task 2: Create Persona Management Components**
