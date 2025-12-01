export namespace main {
	
	export class DriveInfo {
	    path: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new DriveInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.label = source["label"];
	    }
	}
	export class FilterOptions {
	    fileTypes: string[];
	    validOnly: boolean;
	    invalidOnly: boolean;
	    searchText: string;
	    minSize: number;
	    maxSize: number;
	    sortBy: string;
	    sortDesc: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FilterOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileTypes = source["fileTypes"];
	        this.validOnly = source["validOnly"];
	        this.invalidOnly = source["invalidOnly"];
	        this.searchText = source["searchText"];
	        this.minSize = source["minSize"];
	        this.maxSize = source["maxSize"];
	        this.sortBy = source["sortBy"];
	        this.sortDesc = source["sortDesc"];
	    }
	}
	export class FilterResult {
	    files: scanner.FileInfo[];
	    totalCount: number;
	
	    static createFrom(source: any = {}) {
	        return new FilterResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], scanner.FileInfo);
	        this.totalCount = source["totalCount"];
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

export namespace scanner {
	
	export class FileInfo {
	    path: string;
	    name: string;
	    size: number;
	    // Go type: time
	    modTime: any;
	    extension: string;
	    fileType: string;
	    isValid: boolean;
	    invalidReason?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.modTime = this.convertValues(source["modTime"], null);
	        this.extension = source["extension"];
	        this.fileType = source["fileType"];
	        this.isValid = source["isValid"];
	        this.invalidReason = source["invalidReason"];
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
	export class ExportOptions {
	    destPath: string;
	    files: FileInfo[];
	    keepStructure: boolean;
	    overwrite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.destPath = source["destPath"];
	        this.files = this.convertValues(source["files"], FileInfo);
	        this.keepStructure = source["keepStructure"];
	        this.overwrite = source["overwrite"];
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
	export class ExportProgress {
	    total: number;
	    completed: number;
	    failed: number;
	    current: string;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new ExportProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.completed = source["completed"];
	        this.failed = source["failed"];
	        this.current = source["current"];
	        this.percent = source["percent"];
	    }
	}
	export class ExportResult {
	    success: number;
	    failed: number;
	    failedFiles: string[];
	    skippedFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new ExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.failed = source["failed"];
	        this.failedFiles = source["failedFiles"];
	        this.skippedFiles = source["skippedFiles"];
	    }
	}
	
	export class ScanOptions {
	    rootPath: string;
	    includeTypes: string[];
	    excludePaths: string[];
	    validateFiles: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScanOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rootPath = source["rootPath"];
	        this.includeTypes = source["includeTypes"];
	        this.excludePaths = source["excludePaths"];
	        this.validateFiles = source["validateFiles"];
	    }
	}
	export class ScanResult {
	    files: FileInfo[];
	    totalCount: number;
	    validCount: number;
	    invalidCount: number;
	    scanTime: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], FileInfo);
	        this.totalCount = source["totalCount"];
	        this.validCount = source["validCount"];
	        this.invalidCount = source["invalidCount"];
	        this.scanTime = source["scanTime"];
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

