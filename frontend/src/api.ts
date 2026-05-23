import type { MatchState, PublicSettings, ReceiverStatus, Team } from './types';

type Backend = {
  GetSettings(): Promise<PublicSettings>;
  SaveSettings(settings: PublicSettings): Promise<void>;
  GetMatchState(): Promise<MatchState>;
  GenerateTeams(): Promise<Team[]>;
  StartReceiver(): Promise<void>;
  StopReceiver(): Promise<void>;
  ReceiverStatus(): Promise<ReceiverStatus>;
};

declare global {
  interface Window {
    go?: {
      main?: {
        App?: Backend;
      };
    };
  }
}

const fallbackState: MatchState = {
  matchId: '',
  players: [],
  rankings: [],
  teams: [],
  unknownEvents: 0,
  completed: false,
  posted: false
};

const fallbackSettings: PublicSettings = {
  liveApiPort: 7777,
  teamCount: 4,
  teamSizeCap: 5,
  discordWebhookUrl: '',
  autoPostDiscord: true
};

const fallback: Backend = {
  async GetSettings() {
    return fallbackSettings;
  },
  async SaveSettings(settings) {
    Object.assign(fallbackSettings, settings);
  },
  async GetMatchState() {
    return fallbackState;
  },
  async GenerateTeams() {
    return [];
  },
  async StartReceiver() {},
  async StopReceiver() {},
  async ReceiverStatus() {
    return { running: false, address: '' };
  }
};

export function backend(): Backend {
  return window.go?.main?.App ?? fallback;
}
