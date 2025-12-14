// src/types.ts

/**
 * Defines a participant in the Secret Santa room.
 */
export interface Participant {
  name: string;
  isHost: boolean;
  avatar: string;
  // In a real app, we'd add an ID here too
  // id: string;
}

/**
 * Defines the entire application state.
 */
export interface AppState {
  currentView: "landing" | "lobby" | "results";
  userName: string;
  roomCode: string;
  participants: Participant[];
  myTarget: string | null; // The name of the person this user must gift.
}

/**
 * Defines the structure for messages sent from the Go Backend.
 */
export interface ServerMessage {
  type: "JOIN" | "PARTICIPANTS_UPDATE" | "MATCH_RESULT" | "ERROR";
  payload: any;
}
