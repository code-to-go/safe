export namespace chat {
	
	export class Message {
	    id: number;
	    author: string;
	    // Go type: time.Time
	    time: any;
	    content: string;
	    contentType: string;
	    attachments: number[][];
	    signature: number[];
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.author = source["author"];
	        this.time = this.convertValues(source["time"], null);
	        this.content = source["content"];
	        this.contentType = source["contentType"];
	        this.attachments = source["attachments"];
	        this.signature = source["signature"];
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

export namespace library {
	
	export class File {
	    name: string;
	    id: number;
	    // Go type: time.Time
	    modTime: any;
	    size: number;
	    authorId: string;
	    contentType: string;
	    hash: number[];
	    hashChain: number[][];
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	        this.modTime = this.convertValues(source["modTime"], null);
	        this.size = source["size"];
	        this.authorId = source["authorId"];
	        this.contentType = source["contentType"];
	        this.hash = source["hash"];
	        this.hashChain = source["hashChain"];
	        this.tags = source["tags"];
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

export namespace security {
	
	export class Key {
	    pu: number[];
	    pr?: number[];
	
	    static createFrom(source: any = {}) {
	        return new Key(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pu = source["pu"];
	        this.pr = source["pr"];
	    }
	}
	export class Identity {
	    n: string;
	    m: string;
	    s: Key;
	    e: Key;
	
	    static createFrom(source: any = {}) {
	        return new Identity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.n = source["n"];
	        this.m = source["m"];
	        this.s = this.convertValues(source["s"], Key);
	        this.e = this.convertValues(source["e"], Key);
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

