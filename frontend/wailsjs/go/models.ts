export namespace playlist {
	
	export class Channel {
	    name: string;
	    url: string;
	    logo: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new Channel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.logo = source["logo"];
	        this.category = source["category"];
	    }
	}

}

