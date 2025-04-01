export namespace common {
	
	export class Value {
	    Format: string;
	    Data: number[];
	
	    static createFrom(source: any = {}) {
	        return new Value(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Format = source["Format"];
	        this.Data = source["Data"];
	    }
	}
	export class Type {
	    Text: string;
	
	    static createFrom(source: any = {}) {
	        return new Type(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Text = source["Text"];
	    }
	}
	export class ItemWithID {
	    // Go type: Type
	    Type: any;
	    Text: string;
	    Values: Value[];
	    ID: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemWithID(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = this.convertValues(source["Type"], null);
	        this.Text = source["Text"];
	        this.Values = this.convertValues(source["Values"], Value);
	        this.ID = source["ID"];
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

}

