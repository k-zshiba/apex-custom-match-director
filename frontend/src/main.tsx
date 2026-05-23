import React, { useCallback, useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { backend } from './api';
import type { MatchState, PublicSettings, ReceiverStatus, Team } from './types';
import './styles.css';

const initialSettings: PublicSettings = {
  liveApiPort: 7777,
  teamCount: 4,
  teamSizeCap: 5,
  discordWebhookUrl: '',
  autoPostDiscord: true
};

const initialState: MatchState = {
  matchId: '',
  players: [],
  rankings: [],
  teams: [],
  unknownEvents: 0,
  completed: false,
  posted: false
};

function normalizeState(nextState: MatchState): MatchState {
  return {
    ...nextState,
    players: nextState.players ?? [],
    rankings: nextState.rankings ?? [],
    teams: (nextState.teams ?? []).map((team) => ({
      ...team,
      players: team.players ?? []
    }))
  };
}

function App() {
  const [settings, setSettings] = useState<PublicSettings>(initialSettings);
  const [state, setState] = useState<MatchState>(initialState);
  const [status, setStatus] = useState<ReceiverStatus>({ running: false, address: '' });
  const [message, setMessage] = useState('');
  const [fatalError, setFatalError] = useState('');

  const refresh = useCallback(async () => {
    try {
      const api = backend();
      const [nextSettings, nextState, nextStatus] = await Promise.all([
        api.GetSettings(),
        api.GetMatchState(),
        api.ReceiverStatus()
      ]);
      setSettings(nextSettings);
      setState(normalizeState(nextState));
      setStatus(nextStatus);
      setFatalError('');
    } catch (error) {
      setFatalError(error instanceof Error ? error.message : String(error));
    }
  }, []);

  useEffect(() => {
    void refresh();
    const timer = window.setInterval(() => void refresh(), 1500);
    return () => window.clearInterval(timer);
  }, [refresh]);

  async function saveSettings() {
    if (!settings) {
      return;
    }
    await backend().SaveSettings(settings);
    setMessage('Settings saved');
    await refresh();
  }

  async function toggleReceiver() {
    if (status.running) {
      await backend().StopReceiver();
      setMessage('Receiver stopped');
    } else {
      await backend().StartReceiver();
      setMessage('Receiver started');
    }
    await refresh();
  }

  async function generateTeams() {
    const teams: Team[] = await backend().GenerateTeams();
    setState((current) => ({ ...current, teams: teams.map((team) => ({ ...team, players: team.players ?? [] })) }));
    setMessage('Teams generated');
  }

  return (
    <main className="shell">
      <header className="topbar">
        <div>
          <h1>Apex Custom Match Director</h1>
          <p>LiveAPI receiver, ranking, team balancing, and Discord posting.</p>
        </div>
        <button className={status.running ? 'danger' : 'primary'} onClick={() => void toggleReceiver()}>
          {status.running ? 'Stop Receiver' : 'Start Receiver'}
        </button>
      </header>

      <section className="status">
        <div>
          <span>Receiver</span>
          <strong>{status.running ? status.address : 'Stopped'}</strong>
        </div>
        <div>
          <span>Match</span>
          <strong>{state.completed ? 'Complete' : 'Running'}</strong>
        </div>
        <div>
          <span>Discord</span>
          <strong>{state.posted ? 'Posted' : settings.autoPostDiscord ? 'Auto' : 'Off'}</strong>
        </div>
        <div>
          <span>Unknown events</span>
          <strong>{state.unknownEvents}</strong>
        </div>
      </section>

      {message && <p className="notice">{message}</p>}
      {fatalError && <p className="error">{fatalError}</p>}
      {state.lastError && <p className="error">{state.lastError}</p>}

      <div className="grid">
        <section>
          <h2>Settings</h2>
          <label>
            LiveAPI port
            <input
              type="number"
              value={settings.liveApiPort}
              onChange={(event) => setSettings({ ...settings, liveApiPort: Number(event.target.value) })}
            />
          </label>
          <label>
            Team count
            <input
              min={2}
              max={4}
              type="number"
              value={settings.teamCount}
              onChange={(event) => setSettings({ ...settings, teamCount: Number(event.target.value) })}
            />
          </label>
          <label>
            Team size cap
            <input
              min={1}
              type="number"
              value={settings.teamSizeCap}
              onChange={(event) => setSettings({ ...settings, teamSizeCap: Number(event.target.value) })}
            />
          </label>
          <label>
            Discord webhook URL
            <input
              type="password"
              value={settings.discordWebhookUrl}
              onChange={(event) => setSettings({ ...settings, discordWebhookUrl: event.target.value })}
            />
          </label>
          <label className="check">
            <input
              type="checkbox"
              checked={settings.autoPostDiscord}
              onChange={(event) => setSettings({ ...settings, autoPostDiscord: event.target.checked })}
            />
            Auto-post completed matches
          </label>
          <button className="primary" onClick={() => void saveSettings()}>Save Settings</button>
        </section>

        <section>
          <h2>Rankings</h2>
          <div className="table">
            <div className="row head">
              <span>#</span><span>Player</span><span>K</span><span>Dmg</span><span>Score</span>
            </div>
            {state.rankings.map((player) => (
              <div className="row" key={player.playerId}>
                <span>{player.rank}</span>
                <span>{player.playerName}</span>
                <span>{player.kills}</span>
                <span>{player.damage}</span>
                <span>{player.score}</span>
              </div>
            ))}
            {state.rankings.length === 0 && <p className="empty">Waiting for LiveAPI events.</p>}
          </div>
        </section>
      </div>

      <section>
        <div className="sectionHeader">
          <h2>Next Teams</h2>
          <button onClick={() => void generateTeams()}>Generate Teams</button>
        </div>
        <div className="teams">
          {state.teams.map((team) => (
            <article key={team.name}>
              <h3>{team.name}</h3>
              <p>Average {team.averageScore.toFixed(1)}</p>
              <ul>
                {team.players.map((player) => (
                  <li key={player.playerId}>{player.playerName}<span>{player.kills} kills</span></li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
