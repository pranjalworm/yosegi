export namespace main {
	
	export class GenerateResult {
	    outputPath: string;
	    timeTaken: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputPath = source["outputPath"];
	        this.timeTaken = source["timeTaken"];
	    }
	}

}

