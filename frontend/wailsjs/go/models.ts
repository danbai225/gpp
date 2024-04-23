export namespace config {
	
	export class Peer {
	    name: string;
	    protocol: string;
	    port: number;
	    addr: string;
	    uuid: string;
	    ping: number;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.protocol = source["protocol"];
	        this.port = source["port"];
	        this.addr = source["addr"];
	        this.uuid = source["uuid"];
	        this.ping = source["ping"];
	    }
	}

}

export namespace data {
	
	export class Status {
	    running: boolean;
	    game_peer?: config.Peer;
	    http_peer?: config.Peer;
	    up: number;
	    down: number;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.game_peer = this.convertValues(source["game_peer"], config.Peer);
	        this.http_peer = this.convertValues(source["http_peer"], config.Peer);
	        this.up = source["up"];
	        this.down = source["down"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

}

