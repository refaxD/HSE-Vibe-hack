import { EventPlace, Place } from "../../../generated";

/**
 * Function that checks if object is instance of Place
 * @param {any} obj object which should be checked
 * @return {obj is Place}
 */
export function isInstanceOfPlace(obj: any): obj is Place {
	return (
		( obj !== null && obj !== undefined ) &&
		( typeof obj.type === "number" ) &&
		( typeof obj.longitude === "number" ) &&
		( typeof obj.latitude === "number" )
	);
}

/**
 * Function that checks if object is instance of EventPlace
 * @param {any} obj object which should be checked
 * @return {obj is EventPlace}
 */
export function isInstanceOfEventPlace(obj: any): obj is EventPlace {
	return (
		( obj !== null && obj !== undefined ) &&
		( typeof obj.name === "string" ) &&
		( typeof obj.description === "string" ) &&
		( typeof obj.longitude === "number" ) &&
		( typeof obj.latitude === "number" ) &&
		( typeof obj.start_time === "string" ) &&
		( typeof obj.finish_time === "string" )
	);
}
