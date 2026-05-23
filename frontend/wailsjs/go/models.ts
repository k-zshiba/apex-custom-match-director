export namespace domain {
	
	export class Team {
	    name: string;
	    players: Ranking[];
	    averageScore: number;
	    totalScore: number;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.players = this.convertValues(source["players"], Ranking);
	        this.averageScore = source["averageScore"];
	        this.totalScore = source["totalScore"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Ranking {
	    rank: number;
	    playerId: string;
	    playerName: string;
	    kills: number;
	    deaths: number;
	    damage: number;
	    weaponLevel: number;
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new Ranking(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank = source["rank"];
	        this.playerId = source["playerId"];
	        this.playerName = source["playerName"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.damage = source["damage"];
	        this.weaponLevel = source["weaponLevel"];
	        this.score = source["score"];
	    }
	}
	export class PlayerStats {
	    id: string;
	    name: string;
	    kills: number;
	    deaths: number;
	    damage: number;
	    weaponLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.damage = source["damage"];
	        this.weaponLevel = source["weaponLevel"];
	    }
	}
	export class MatchState {
	    matchId: string;
	    startedAt?: string;
	    completedAt?: string;
	    players: PlayerStats[];
	    rankings: Ranking[];
	    teams: Team[];
	    unknownEvents: number;
	    completed: boolean;
	    posted: boolean;
	    lastError?: string;
	
	    static createFrom(source: any = {}) {
	        return new MatchState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchId = source["matchId"];
	        this.startedAt = source["startedAt"];
	        this.completedAt = source["completedAt"];
	        this.players = this.convertValues(source["players"], PlayerStats);
	        this.rankings = this.convertValues(source["rankings"], Ranking);
	        this.teams = this.convertValues(source["teams"], Team);
	        this.unknownEvents = source["unknownEvents"];
	        this.completed = source["completed"];
	        this.posted = source["posted"];
	        this.lastError = source["lastError"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PublicSettings {
	    liveApiPort: number;
	    teamCount: number;
	    teamSizeCap: number;
	    discordWebhookUrl: string;
	    autoPostDiscord: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PublicSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.liveApiPort = source["liveApiPort"];
	        this.teamCount = source["teamCount"];
	        this.teamSizeCap = source["teamSizeCap"];
	        this.discordWebhookUrl = source["discordWebhookUrl"];
	        this.autoPostDiscord = source["autoPostDiscord"];
	    }
	}
	
	export class ReceiverStatus {
	    running: boolean;
	    address: string;
	
	    static createFrom(source: any = {}) {
	        return new ReceiverStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.address = source["address"];
	    }
	}

}

