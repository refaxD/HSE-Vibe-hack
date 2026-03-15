import { Pipe, PipeTransform } from "@angular/core";
import { Place } from "../../../generated";

@Pipe({
	name: "placeSearch",
})
export class PlaceSearchPipe implements PipeTransform {

	transform(places: Place[], search: string): Place[] {
		function splitByNonLetters(str: string) {
			return str.split(/[^a-zA-Zа-яА-ЯёЁ0-9]+/).filter(Boolean); // Filter removes empty strings
		}

		function checkString(value: string, find: string, maxMistake: number = 0) {
			let differs: number = 0;
			for (let i = 0; i < Math.min(value.length, find.length); i++)
				if (value[i] !== find[i])
					differs++;
			differs += Math.max(0, find.length - value.length);
			if (value.length < 3) return differs === 0;
			if (value.length < 5) return differs === Math.min(maxMistake, 1);
			return differs === Math.min(maxMistake, 2);
		}

		function checkArrays(value: string, search: string, maxMistake: number = 0) {
			const value_arr = splitByNonLetters(value);
			const search_arr = splitByNonLetters(search);

			let right = 0;
			for (const item of value_arr) {
				if (right === search_arr.length)
					return true;
				if (checkString(item, search_arr[right], maxMistake))
					right++;
			}

			return right === search_arr.length;
		}

		let result: Place[] = [];

		search = search.toLowerCase();
		for (let mistake = 0; mistake <= 2; mistake++)
			for (const place of places) {
				if (checkArrays(( place.name || "" ).toLowerCase(), search, mistake))
					result.push(place);

				if (result.length >= 10)
					return result;
			}

		return result;
	}

}
