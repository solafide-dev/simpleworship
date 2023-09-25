export namespace main {

	export class DisplayServerData {
		type: string;
		meta: { [key: string]: any };

		static createFrom(source: any = {}) {
			return new DisplayServerData(source);
		}

		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.type = source["type"];
			this.meta = source["meta"];
		}
	}

}

