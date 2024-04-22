export namespace config {
	
	export class Peer {
	    name: string;
	    protocol: string;
	    port: number;
	    addr: string;
	    uuid: string;
	
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
	    }
	}

}

