import { AfterViewInit, Component, effect, ElementRef, OnDestroy, ViewChild } from "@angular/core";
import { Router } from "@angular/router";
import { environment } from "../../../../../environments/environment";
import { EventPlace, Place } from "../../../../../generated";
import { ApiService } from "../../../../core/services/api.service";
import {
	MapPointInformationComponent,
} from "../../../../pop-up-menus/components/map-point-information/map-point-information.component";
import { NotificationsComponent } from "../../../../pop-up-menus/components/notifications/notifications.component";
import { PathInformationComponent } from "../../../../pop-up-menus/components/path-information/path-information.component";
import { MapClickEntityModel } from "../../../../core/models/map-click-entity.model";
import { MapClickEventModel } from "../../../../core/models/map-click-event.model";
import MapPointModel from "../../../../core/models/map-point.model";
import { MapInteractionsService } from "../../../../core/services/map-interactions.service";

declare var ymaps3: any;

@Component({
	selector: "home-main",
	standalone: true,
	imports: [
		MapPointInformationComponent,
		PathInformationComponent,
		NotificationsComponent,
	],
	templateUrl: "./main.component.html",
	styleUrl: "./main.component.css",
})
export default class MainComponent implements AfterViewInit, OnDestroy {
	/** Yandex map container */
	map: any = null;
	/** user marker that is displayed on the map */
	userMarker: any = null;
	/** point that was chosen during prompt typing */
	displayedPromptPoint: any = null;
	/** path points */
	pathPoints: any[] | null = null;
	/** path lines */
	pathLines: any[] | null = null;
	/** Yandex map HTML DOOM container */
	@ViewChild("ymaps") mapElement!: ElementRef<HTMLDivElement>;

	/** default location setting to the map */
	private readonly DEFAULT_LOCATION_SETTINGS = {
		duration: 500,
		easing: "ease-in-out",
	};

	constructor(
		private apiService: ApiService,
		private router: Router,
		private mapInteractionService: MapInteractionsService,
	) {
		/** Adding Listener for MapScrolling changes */
		effect(() => {
			if (this.mapInteractionService.mapScrolled() === -1 && this.map !== null)
				this.stickPosition();
		});
		/** Adding Listener for user geolocation changes */
		effect(() => {
			if (this.mapInteractionService.userPosition() && this.map !== null)
				this.geolocationChange();
		});
		/** Adding Listener for prompt input point chosen state */
		this.mapInteractionService.selectedPointOnPromptInput.subscribe(() => {
			if (this.map !== null)
				this.displaySelectedPromptInputPoint(this.mapInteractionService.selectedPointOnPromptInput.value);
		});
		/** Adding Listener for route points change */
		this.mapInteractionService.pathPoints.subscribe(path => {
			if (this.map !== null)
				this.displayPathPoints().then();
		});
		/** Listener for chosen point changes */
		this.mapInteractionService.chosenMapPoint.subscribe(point => {
			if (point === null) return;
			const coords = [point.latitude, point.longitude];
			// moves map camera to marker
			this.map.update({
				location: {
					center: coords,
					zoom: 16,
					...this.DEFAULT_LOCATION_SETTINGS,
				},
			});
		});
	}

	ngAfterViewInit() {
		/** when ymaps ready we create the map */
		this.initMap().then();
		/** getting points from route */
		this.getPointsFromRoute();
	}

	ngOnDestroy() {
		this.map.destroy();
	}

	/** Gets path points from route url */
	private async getPointsFromRoute() {
		let route: string[] | string = this.router.url.split('?')[1];
		if (route === undefined) return;
		route = route.split('&');
		if (route.length < 2) return;
		const points = route.map(item => item.split('='));
		let places = [];
		for (const point of points)
			places.push(await (this.apiService.getPlaceById(point[1])));
		if (places.length > 0)
			this.mapInteractionService.pathPoints.next(places);
	}

	/** Function that centers map according to `lastPosition` */
	public stickPosition() {
		if (this.mapInteractionService.userPosition() === null) return;
		const location: number[] = [this.mapInteractionService.userPosition()?.longitude!, this.mapInteractionService.userPosition()?.latitude!];
		this.map.update({
			location: {
				center: location,
				zoom: 16,
				...this.DEFAULT_LOCATION_SETTINGS,
			},
		});
	}

	/** Function that creates Yandex map */
	private async initMap(): Promise<void> {
		await ymaps3.ready;

		const { YMap, YMapDefaultSchemeLayer, YMapListener, YMapDefaultFeaturesLayer } = ymaps3;

		let center: number[] = [this.mapInteractionService.userPosition()?.latitude!, this.mapInteractionService.userPosition()?.longitude!];

		this.map = new YMap(this.mapElement.nativeElement, {
			location: {
				center: center,
				zoom: 10,
			},
		});

		this.map.addChild(new YMapDefaultSchemeLayer());
		this.map.addChild(new YMapDefaultFeaturesLayer());
		this.map.addChild(new YMapListener({
			onTouchStart: () => this.handleMapMoveStart(),
			onMouseDown: () => this.handleMapMoveStart(),
			onTouchMove: () => this.handleMapMoveStart(),
			onMouseMove: () => this.handleMapMoveStart(),
			onClick: (e: MapClickEventModel | undefined) => this.handleMapClick(e),
		}));

		if (this.mapInteractionService.userPosition() !== null)
			this.geolocationChange();
		await this.displayPathPoints();
	}

	/** Handler for map move event begin */
	private handleMapMoveStart() {
		this.mapInteractionService.mapScrolled.set(1);
		setTimeout(() => this.mapInteractionService.mapScrolled.set(0), 100);
	}

	/**
	 * Function that displays point that was chosen during prompt input
	 * @param {MapClickEntityModel | null} point point that should be displayed
	 */
	private displaySelectedPromptInputPoint(point: MapClickEntityModel | null) {
		if (point === null) {
			if (this.displayedPromptPoint !== null)
				this.map.remove(this.displayedPromptPoint);
			this.displayedPromptPoint = null;
		} else {
			const { YMapMarker } = ymaps3; // getting point interface from ymap
			const coords = point.geometry.coordinates; // point coordinates
			// updating map camera view
			this.map.update({
				location: {
					center: coords,
					zoom: 16,
					...this.DEFAULT_LOCATION_SETTINGS,
				},
			});

			if (this.displayedPromptPoint === null) { // if point don't exist on the map with now we should create it
				// creating marker element for the point
				const markerElement = document.createElement("div");
				markerElement.className = "selected-prompt-place-map-marker";
				markerElement.innerText = "1";
				// adding point to the map
				this.displayedPromptPoint = new YMapMarker({
					coordinates: coords,
				}, markerElement);
				this.map.addChild(this.displayedPromptPoint);
			} else { // if point already exists on the map we just update it position
				this.displayedPromptPoint.update({ coordinates: coords });
			}
		}
	}

	/**
	 * Handler for map click event
	 * @param {MapClickEventModel | undefined} event map click event
	 */
	private async handleMapClick(event: MapClickEventModel | undefined) {
		if (event?.type === "marker") { // we clicked on the marker on the map
			const element = event.entity.element;
			const text = element.innerText;
			if (text === "Я") return;
			else {
				const index = Number(text) - 1;
				const point = this.mapInteractionService.pathPoints.value![index];
				this.mapInteractionService.chosenMapPoint.next(point);
			}
		} else if (event && event.entity.properties.name) {
			this.mapInteractionService.chosenMapPoint.next({
				name: event.entity.properties.name!,
				categories: [],
				latitude: event.entity.geometry.coordinates[0],
				longitude: event.entity.geometry.coordinates[1],
				isLiked: false,
			});
		} else {
			this.mapInteractionService.chosenMapPoint.next(null);
		}
	}

	/** Function that displays path on the map */
	private async displayPathPoints() {
		const points = this.mapInteractionService.pathPoints.value;

		if (this.pathPoints !== null) // we should remove displayed path from the map
			this.pathPoints.forEach(point => this.map.removeChild(point));
		if (this.pathLines !== null)
			this.pathLines.forEach(line => this.map.removeChild(line));

		if (points === null) // we have no path to display
			this.pathPoints = this.pathLines = null;
		else { // we have a new path to display
			const { YMapMarker } = ymaps3;
			this.pathPoints = [];
			this.pathLines = [];
			let bounds: number[][] = []; // bounds for camera move

			for (let point_index = 0; point_index < points.length; point_index++) {
				// creating marker element for the point
				const markerElement = document.createElement("div");
				markerElement.className = "selected-prompt-place-map-marker";
				markerElement.innerText = `${ point_index + 1 }`;
				// adding point to the map
				const displayed_point = new YMapMarker({
					coordinates: [points[point_index].latitude, points[point_index].longitude],
				}, markerElement);
				// adding created points to array
				this.pathPoints.push(displayed_point);
				// displaying points on the map
				this.map.addChild(displayed_point);
				// updating bounds with point position
				const paddingBottom = 0.0017, paddingTop = 0.001;
				if (bounds.length === 0) {
					bounds.push([points[point_index].latitude, points[point_index].longitude - paddingBottom]);
					bounds.push([points[point_index].latitude, points[point_index].longitude + paddingTop]);
				} else {
					bounds[0][0] = Math.min(bounds[0][0], points[point_index].latitude);
					bounds[0][1] = Math.min(bounds[0][1], points[point_index].longitude - paddingBottom);
					bounds[1][0] = Math.max(bounds[1][0], points[point_index].latitude);
					bounds[1][1] = Math.max(bounds[1][1], points[point_index].longitude + paddingTop);
				}
			}

			// fetching route
			const route = await ( this.fetchRoute(points) );
			const geometry: number[][] = route.features[0].geometry.coordinates;
			// displaying route on the map
			for (let i = 1; i < geometry.length; i++) {
				const pointA = geometry[i - 1], pointB = geometry[i];
				const line = this.drawLine(pointA, pointB); // creating line between to point on the path
				// adding line to the map
				this.map.addChild(line);
				this.pathLines.push(line);
			}

			// bounding map to show all the created route
			this.mapInteractionService.mapScrolled.set(0);
			this.map.update({
				location: {
					bounds: bounds,
					padding: {
						top: 40,
						left: 40,
						right: 40,
						bottom: 40,
					},
					...this.DEFAULT_LOCATION_SETTINGS,
				},
			});
		}
	}

	/**
	 * Function that fetches route with points from the OSM
	 * @param {Array<Place | EventPlace | MapPointModel>} points points on the route
	 */
	private async fetchRoute(points: Array<Place | EventPlace | MapPointModel>) {
		try {
			// creating array with points coordinates
			const coordinates: number[][] = [];
			for (const point of points)
				coordinates.push([point.latitude, point.longitude]);
			// fetching path from OSM API
			const response = await fetch(`${ environment.RouteUrl }/v2/directions/foot-walking/geojson`, {
				method: "POST",
				headers: {
					"Authorization": environment.RouteKey,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ "coordinates": coordinates }),
			});
			return response.json();
		} catch (error) {
			return error;
		}
	}

	/** Function that draws like between two points */
	private drawLine(pointA: number[], pointB: number[]) {
		const { YMapFeature } = ymaps3;
		// reversing points coordinates if they are in wrong order
		if (pointA[0] > pointA[1])
			pointA.reverse();
		if (pointB[0] > pointB[1])
			pointB.reverse();
		// basic line PROPS
		const FEATURE_PROPS = {
			geometry: {
				type: "LineString",
				coordinates: [pointA, pointB],
			},
			style: {
				simplificationRate: 0,
				stroke: [
					{ color: "rgba(25, 109, 255, 0.6)", width: 10 },
				],
				fill: "rgba(25, 109, 255, 0.6)",
			},
		};
		// returns created line
		return new YMapFeature(FEATURE_PROPS);
	}

	/** Function that makes changes to map and userMarker after the user geolocation changes */
	private geolocationChange() {
		const updateUserMarker = (coords: number[]) => {
			if (this.userMarker === null) {
				const markerElement = document.createElement("div");
				const insideMarkerElement = document.createElement("div");
				insideMarkerElement.className = "inner-circle";
				insideMarkerElement.innerText = "Я";
				markerElement.className = "user-map-marker";
				markerElement.appendChild(insideMarkerElement);

				this.userMarker = new YMapMarker({
					coordinates: coords,
				}, markerElement);
				this.map.addChild(this.userMarker);
			} else
				this.userMarker.update({ coordinates: coords });
		};

		const { geolocation, YMapMarker } = ymaps3;
		// Changing map geolocation according to new positions
		geolocation.getPosition().then((result: any) => {
			if (this.mapInteractionService.mapScrolled() === -1)
				this.map.update({
					location: {
						center: result.coords,
						zoom: 16,
						...this.DEFAULT_LOCATION_SETTINGS,
					},
				});
			updateUserMarker(result.coords);
		});
	}
}
