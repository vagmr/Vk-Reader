export namespace main {
	
	export class Book {
	    id: string;
	    title: string;
	    author: string;
	    cover: string;
	    wordCount: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Book(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.author = source["author"];
	        this.cover = source["cover"];
	        this.wordCount = source["wordCount"];
	        this.status = source["status"];
	    }
	}
	export class Chapter {
	    title: string;
	    content: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new Chapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.content = source["content"];
	        this.id = source["id"];
	    }
	}
	export class BookDetail {
	    title: string;
	    chapters: {[key: string]: Chapter};
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new BookDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.chapters = this.convertValues(source["chapters"], Chapter, true);
	        this.status = source["status"];
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
	
	export class LogConfig {
	    enabled: boolean;
	    level: string;
	
	    static createFrom(source: any = {}) {
	        return new LogConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.level = source["level"];
	    }
	}
	export class RetryConfig {
	    maxRetries: number;
	    retryDelay: number;
	
	    static createFrom(source: any = {}) {
	        return new RetryConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxRetries = source["maxRetries"];
	        this.retryDelay = source["retryDelay"];
	    }
	}
	export class DelayConfig {
	    min: number;
	    max: number;
	
	    static createFrom(source: any = {}) {
	        return new DelayConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.min = source["min"];
	        this.max = source["max"];
	    }
	}
	export class ToValues {
	    txt: string;
	    epub: string;
	    html: string;
	
	    static createFrom(source: any = {}) {
	        return new ToValues(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.txt = source["txt"];
	        this.epub = source["epub"];
	        this.html = source["html"];
	    }
	}
	export class Config {
	    savePath: string;
	    saveFormat: string;
	    threadCount: number;
	    convertType: ToValues;
	    delay: DelayConfig;
	    retryConfig: RetryConfig;
	    logConfig: LogConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.savePath = source["savePath"];
	        this.saveFormat = source["saveFormat"];
	        this.threadCount = source["threadCount"];
	        this.convertType = this.convertValues(source["convertType"], ToValues);
	        this.delay = this.convertValues(source["delay"], DelayConfig);
	        this.retryConfig = this.convertValues(source["retryConfig"], RetryConfig);
	        this.logConfig = this.convertValues(source["logConfig"], LogConfig);
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
	
	export class DownloadProgress {
	    total: number;
	    current: number;
	    chapterTitle: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.current = source["current"];
	        this.chapterTitle = source["chapterTitle"];
	        this.status = source["status"];
	    }
	}
	
	

}

