export type PublicSettings = {
  liveApiPort: number;
  teamCount: number;
  teamSizeCap: number;
  discordWebhookUrl: string;
  autoPostDiscord: boolean;
};

export type Ranking = {
  rank: number;
  playerId: string;
  playerName: string;
  kills: number;
  deaths: number;
  damage: number;
  weaponLevel: number;
  score: number;
};

export type Team = {
  name: string;
  players: Ranking[];
  averageScore: number;
  totalScore: number;
};

export type MatchState = {
  matchId: string;
  startedAt?: string;
  completedAt?: string;
  players: Array<{
    id: string;
    name: string;
    kills: number;
    deaths: number;
    damage: number;
    weaponLevel: number;
  }>;
  rankings: Ranking[];
  teams: Team[];
  unknownEvents: number;
  completed: boolean;
  posted: boolean;
  lastError?: string;
};

export type ReceiverStatus = {
  running: boolean;
  address: string;
};
