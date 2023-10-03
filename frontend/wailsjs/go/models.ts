export namespace main {
	
	export class DisplayServerData {
	    type: string;
	    meta: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new DisplayServerData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.meta = source["meta"];
	    }
	}
	export class Meta {
	    title: string;
	    artist: string;
	
	    static createFrom(source: any = {}) {
	        return new Meta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.artist = source["artist"];
	    }
	}
	export class ServiceItemMeta {
	    songId?: string;
	    verseReference?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceItemMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.songId = source["songId"];
	        this.verseReference = source["verseReference"];
	    }
	}
	export class ServiceItem {
	    title: string;
	    type: string;
	    meta?: ServiceItemMeta;
	
	    static createFrom(source: any = {}) {
	        return new ServiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.type = source["type"];
	        this.meta = this.convertValues(source["meta"], ServiceItemMeta);
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
	export class OrderOfService {
	    swdt: string;
	    id: string;
	    title: string;
	    date: string;
	    items: ServiceItem[];
	
	    static createFrom(source: any = {}) {
	        return new OrderOfService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.swdt = source["swdt"];
	        this.id = source["id"];
	        this.title = source["title"];
	        this.date = source["date"];
	        this.items = this.convertValues(source["items"], ServiceItem);
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
	
	
	export class Slide {
	    section: string;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new Slide(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.section = source["section"];
	        this.text = source["text"];
	    }
	}
	export class SongPart {
	    id: string;
	    lines: string[];
	
	    static createFrom(source: any = {}) {
	        return new SongPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.lines = source["lines"];
	    }
	}
	export class Song {
	    swdt: string;
	    id: string;
	    title: string;
	    attribution: string;
	    license: string;
	    notes: string;
	    parts: SongPart[];
	    order: string[];
	
	    static createFrom(source: any = {}) {
	        return new Song(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.swdt = source["swdt"];
	        this.id = source["id"];
	        this.title = source["title"];
	        this.attribution = source["attribution"];
	        this.license = source["license"];
	        this.notes = source["notes"];
	        this.parts = this.convertValues(source["parts"], SongPart);
	        this.order = source["order"];
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
	
	export class SongSlide {
	    meta: Meta;
	    slides: Slide[];
	
	    static createFrom(source: any = {}) {
	        return new SongSlide(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], Meta);
	        this.slides = this.convertValues(source["slides"], Slide);
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

