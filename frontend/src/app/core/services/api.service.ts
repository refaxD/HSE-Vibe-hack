import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { firstValueFrom } from "rxjs";
import { environment } from "../../../environments/environment";
import { BackendApiService, Category, Configuration, EventPlace, Place } from "../../../generated";

@Injectable({
	providedIn: "root",
})
export class ApiService {
	public api: BackendApiService;

	constructor(
		private http: HttpClient,
	) {
		this.api = new BackendApiService(this.http, environment.apiBaseUrl, new Configuration());
	}

	/**
	 * Function that gets all places from backend
	 * @returns {Promise<Place[]>} list of places
	 */
	async getAllPlaces(): Promise<Place[]> {
		try {
			return await firstValueFrom(this.api.placesGetAll());
		} catch (error) {
			throw ( error );
		}
	}

	/**
	 * Function that gets all categories from backend
	 * @returns {Promise<Category[]>} list of categories
	 */
	async getAllCategories(): Promise<Category[]> {
		try {
			return await firstValueFrom(this.api.categoriesGetAll());
		} catch (error) {
			throw ( error );
		}
	}

	async getPlaceById(id: string): Promise<Place | EventPlace> {
		try {
			return await firstValueFrom(this.api.placesGetById(id));
		} catch (error) {
			try {
				return await firstValueFrom(this.api.eventPlacesGetById(id));
			} catch (error) {
				throw ( 'Not found' );
			}
		}
	}

	async createPath(prompt: string, startPosition: number[]): Promise<Array<Place | EventPlace>> {
		const testPrompt: Array<any> = [
			{ type: "category", categories: { "Кафе": 90, "Ресторан": 60, "Бар": 50 } },
			{ type: "fixed", coords: [37.626225, 55.753236], name: "«Галерея-мастерская Варшавка»", time: 90 },
			{ type: "category", categories: { "Парк": 76, "Каток": 45, "Ботанический сад": 30, "Пикник": 67 } },
		]; // TODO remove after testing
		try {
			return await firstValueFrom(this.api.pathCreate({ prompt: prompt, startPosition: startPosition }));
		} catch (error) {
			throw ( error );
		}
	}
}
