import { HttpClient } from "@angular/common/http";
import { Injectable, signal } from "@angular/core";
import { Geolocation, Position } from "@capacitor/geolocation";
import { BehaviorSubject, firstValueFrom } from "rxjs";
import { Category, EventPlace, Place } from "../../../generated";
import GeolocationModel from "../models/geolocation.model";
import { LocationModel } from "../models/location.model";
import { MapClickEntityModel } from "../models/map-click-entity.model";
import MapPointModel from "../models/map-point.model";
import { ApiService } from "./api.service";

@Injectable({
	providedIn: "root",
})
export class MapInteractionsService {
	/** all categories for the project */
	categories = signal<Category[]>([]);

	/** Array with all places */
	public places = signal<Array<Place>>([]);

	/** points of the generated path that should be displayed */
	pathPoints = new BehaviorSubject<Array<Place | EventPlace | MapPointModel> | null>(null);

	/** point that was selected during typing prompt, point should be displayed on the map */
	selectedPointOnPromptInput = new BehaviorSubject<MapClickEntityModel | null>(null);

	/** variable that responsible for map scrolling action
	 * ** -1 ** - map never been scrolled
	 * ** 0 ** - map was scrolled
	 * ** 1 ** - map is being scrolled */
	mapScrolled = signal<number>(-1);

	/** variable that responsible for mainContainer state <br>
	 * ** -1 ** - hidden
	 * ** 0 ** - custom height
	 * ** 1 ** - shown */
	mainContainerState = signal<number>(-1);

	/** all chosen categories that should be represented on the card */
	chosenCategories = signal<Category[]>([]);

	/** chosen point on the map or **null** if user hadn't chosen anything */
	chosenMapPoint = new BehaviorSubject<Place | EventPlace | MapPointModel | null>(null);

	/** variable that responsible for the last position of the user marker <br> **null** - map was not loaded yet and last position doesn't exist <br> **GeolocationModel** - prev position of the user marker */
	userPosition = signal<GeolocationModel>({ longitude: 55, latitude: 37 });

	/** last user location */
	userLocation = signal<LocationModel | null>(null);

	/** variable that responsible for path information container start
	 * -2 - closed
	 * -1 - should close
	 * 1 - should open
	 * 2 - opened
	 */
	pathInformationState = new BehaviorSubject<number>(2);

	constructor(
		private readonly apiService: ApiService,
		private readonly http: HttpClient,
	) {
		Geolocation.watchPosition({}, (pos: Position | null) => {
			if (pos === null) return;
			this.userPosition.set({ longitude: pos.coords.longitude, latitude: pos.coords.latitude });
			this.updateCity().then();
		});
		this.apiService.getAllCategories().then(categories => this.categories.set(categories));
		this.apiService.getAllPlaces().then(places => {
			this.places.set(places);
			//this.pathPoints.next([this.places()[0], this.places()[1], this.places()[2]]);
			//this.pathInformationState.next(1);
		});
	}

	/**
	 * Function that checkes if category was chosen
	 * @param {Category} category category that should be checked for being chosen
	 * @description function that checks `category` for being in chosen list
	 */
	checkCategoryChosen(category: Category) {
		return this.chosenCategories().indexOf(category) !== -1;
	}

	/** Function that updates the user city */
	async updateCity() {
		let bdcApi = "https://api.bigdatacloud.net/data/reverse-geocode-client";
		bdcApi = bdcApi
			+ "?latitude=" + this.userPosition()?.latitude
			+ "&longitude=" + this.userPosition()?.longitude
			+ "&localityLanguage=ru";

		const resp = await firstValueFrom(this.http.get(bdcApi));
		this.userLocation.set(resp as LocationModel);
	}

	/**
	 * Setting new value to pathPoints
	 * @param {Array<Place | EventPlace | MapPointModel> | null} path new value of path
	 */
	setNewPath(path: Array<Place | EventPlace | MapPointModel> | null) {
		this.pathPoints.next(path);
	}

	async generatePath(promptText: string) {
		const position = this.userPosition();
		const path = await this.apiService.createPath(promptText, [37.626225, 55.753236]); // TODO set normal coordinates
		this.pathPoints.next(path);
	}
}
