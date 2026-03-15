import { Component, ElementRef, ViewChild } from "@angular/core";
import { environment } from "../../../../environments/environment";
import { EventPlace, Place } from "../../../../generated";
import { MoveableDirective } from "../../../../libs/moveable.directive";
import MapPointModel from "../../../core/models/map-point.model";
import { MapInteractionsService } from "../../../core/services/map-interactions.service";
import { NotificationsService } from "../../../core/services/notifications.service";
import { isInstanceOfEventPlace, isInstanceOfPlace } from "../../../core/services/utils.service";

@Component({
	selector: "PathInformation",
	imports: [
		MoveableDirective,
	],
	styleUrl: "./path-information.component.css",
	template: `
		<div moveable [isPopup]="false" [callback]="closeCallback.bind({}, this.notificationsService, this.notificationCallback.bind({}, this.mapInteractionsService, this.path))" class="container" #container>
			<div class="content">
				<!-- container with the arrow -->
				<div class="arrow-container">
					<svg fill="#6F6F6F">
						<path d="M0 1.5C0 0.671573 0.671573 0 1.5 0H28.5C29.3284 0 30 0.671573 30 1.5V1.5C30 2.32843 29.3284 3 28.5 3H1.5C0.671573 3 0 2.32843 0 1.5V1.5Z"></path>
					</svg>
				</div>
				<!-- container with a short path -->
				<div class="small-info">
					<div class="content">
						<div>
							<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
								<path d="M14 4.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"></path>
								<path d="M14.836 15.133a.66.66 0 0 1 .11.217l1.4 4.734a.716.716 0 0 1-.516.889.76.76 0 0 1-.909-.447l-1.672-4.413-2.815-2.707a1.013 1.013 0 0 1-.29-.93l.6-3.292-1.352.385-1.191 2.664a.635.635 0 0 1-.82.32.593.593 0 0 1-.34-.766L8.242 8.78a.626.626 0 0 1 .308-.332l.077-.037 3.008-1.478.021-.008.015-.005a.935.935 0 0 1 .6-.089c.269.05.5.199.615.425.114.226.154.641.154.641.046.256.062.508.078.76a12.356 12.356 0 0 1 .039 1.42c-.021.72-.07 1.434-.17 2.141l-.053.418 1.833 2.407.068.091z"></path>
								<path d="M14.12 9.253l2.643 2.276c.27.206.316.583.103.843a.644.644 0 0 1-.858.114l-2.128-1.482.023-.129.028-.158.074-.407.026-.145.025-.142.017-.23.025-.29.01-.12.011-.13z"></path>
								<path d="M10.164 14.399c.209.2.441.405.674.609.371.326.751.64 1.085.895-.142.287-.379.76-.379.76l-2.539 3.992a.768.768 0 0 1-1.031.24.707.707 0 0 1-.278-.943l2.117-3.99.029-.277c.023-.225.046-.451.077-.675l.11-.736.065.06.07.065z"></path>
							</svg>
						</div>
					</div>

					<div>
						<svg (click)="close()" class="hide-icon" width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
							<path d="M17.707 16.293a1 1 0 0 1-1.414 1.414L12 13.414l-4.293 4.293a1 1 0 0 1-1.414-1.414L10.586 12 6.293 7.707a1 1 0 0 1 1.414-1.414L12 10.586l4.293-4.293a1 1 0 1 1 1.414 1.414L13.414 12l4.293 4.293z"></path>
						</svg>
					</div>
				</div>
				<hr #separator>
				<!-- container with a full information about the path -->
				<div class="large-info" #largeInfo>
					@for (point of path; track path; let index = $index) {
						<div class="card" (click)="openPointInformation(point)">
							@if (isInstanceOfEventPlace(point)) {

							} @else if (isInstanceOfPlace(point)) {
								<div class="flex flex-row justify-between gap-1 items-center">
									<p class="font-semibold">{{ point.name }}</p>
									<p class="font-normal text-base text-text-secondary">{{
											this.arrivalTimeToPoints[index]
										}}</p>
								</div>
								<p class="text-text-secondary">{{ point.address }}</p>
							} @else {
								<div class="flex flex-row justify-between gap-1 items-center">
									<p class="font-semibold">{{ point.name }}</p>
									<p class="font-normal text-base text-text-secondary">{{
											this.arrivalTimeToPoints[index]
										}}</p>
								</div>
							}
						</div>

						@if (index !== path!.length - 1) {
							<div class="way-block">
								<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
									<path d="M14 4.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"></path>
									<path d="M14.836 15.133a.66.66 0 0 1 .11.217l1.4 4.734a.716.716 0 0 1-.516.889.76.76 0 0 1-.909-.447l-1.672-4.413-2.815-2.707a1.013 1.013 0 0 1-.29-.93l.6-3.292-1.352.385-1.191 2.664a.635.635 0 0 1-.82.32.593.593 0 0 1-.34-.766L8.242 8.78a.626.626 0 0 1 .308-.332l.077-.037 3.008-1.478.021-.008.015-.005a.935.935 0 0 1 .6-.089c.269.05.5.199.615.425.114.226.154.641.154.641.046.256.062.508.078.76a12.356 12.356 0 0 1 .039 1.42c-.021.72-.07 1.434-.17 2.141l-.053.418 1.833 2.407.068.091z"></path>
									<path d="M14.12 9.253l2.643 2.276c.27.206.316.583.103.843a.644.644 0 0 1-.858.114l-2.128-1.482.023-.129.028-.158.074-.407.026-.145.025-.142.017-.23.025-.29.01-.12.011-.13z"></path>
									<path d="M10.164 14.399c.209.2.441.405.674.609.371.326.751.64 1.085.895-.142.287-.379.76-.379.76l-2.539 3.992a.768.768 0 0 1-1.031.24.707.707 0 0 1-.278-.943l2.117-3.99.029-.277c.023-.225.046-.451.077-.675l.11-.736.065.06.07.065z"></path>
								</svg>
								<p>{{ this.travelLengthBetweenPoints[index] }}</p>
								<p>•</p>
								<p>{{ this.travelTimeBetweenPoints[index] }}</p>
							</div>
						}
					}
				</div>
			</div>
		</div>
	`,
})
export class PathInformationComponent {
	/** path that should be displayed */
	path: Array<Place | EventPlace | MapPointModel> | null = null;

	/** Container element */
	@ViewChild("container") container!: ElementRef<HTMLDivElement>;
	/** Element with full information about the points on the path */
	@ViewChild("largeInfo") largeInfoContainer!: ElementRef<HTMLDivElement>;
	/** Element with separates largeInfo from top part */
	@ViewChild("separator") separatorElement!: ElementRef<HTMLHRElement>;

	/** travel length between points of the path */
	travelLengthBetweenPoints: string[] = [];
	/** travel time between points of the path */
	travelTimeBetweenPoints: string[] = [];
	/** arrival time into a point */
	arrivalTimeToPoints: string[] = [];

	/** Initial height of the container */
	private readonly INITIAL_HEIGHT = window.innerHeight * 0.7;
	/** Maximum height of the container */
	private readonly MAX_HEIGHT = window.innerHeight;
	/** Minimum height of the container */
	private MIN_HEIGHT: number = 0;
	/** Minimum deltaY to interact with the container */
	private MIN_INTERACTION_DELTA: number = 40;

	protected readonly isInstanceOfEventPlace = isInstanceOfEventPlace;
	protected readonly isInstanceOfPlace = isInstanceOfPlace;

	constructor(
		public notificationsService: NotificationsService,
		public mapInteractionsService: MapInteractionsService,
	) {
		/** Listener for path Points change */
		this.mapInteractionsService.pathPoints.subscribe(path => {
			this.travelTimeBetweenPoints.length = 0;
			this.travelLengthBetweenPoints.length = 0;
			this.arrivalTimeToPoints.length = 0;

			if (path === null) {
				if (this.mapInteractionsService.pathInformationState.value !== -1)
					this.mapInteractionsService.pathInformationState.next(-1);
			} else {
				if (this.mapInteractionsService.pathInformationState.value !== 1)
					this.mapInteractionsService.pathInformationState.next(1);
			}
			this.path = path;
			if (this.path !== null)
				this.getPathInfo().then();
		});
		/** Listener for path information state change */
		this.mapInteractionsService.pathInformationState.subscribe(state => {
			if (state == -1)
				this.Hide();
			else if (state == 1)
				this.Show();
		});
	}

	/** Function that fetches information about the path */
	async getPathInfo() {
		/** Function that rounds number for `count` digits after the point */
		function roundToOne(num: number, count: number = 1) {
			return Number(num.toFixed(count));
		}

		/** Function that converts date to string */
		function dateToString(date: Date) {
			return `${ date.getHours().toString().padStart(2, "0") }:${ date.getMinutes().toString().padStart(2, "0") }`;
		}

		if (this.path === null) return;

		// getting and pushing string with path start time
		let currentTime = new Date();
		this.arrivalTimeToPoints.push(dateToString(currentTime));

		// iterating through all the segments of the path
		for (let i = 1; i < this.path.length; i++) {
			const pointA = this.path[i - 1]; // start point of the segment
			const pointB = this.path[i]; // finish point of the segment
			// coordinates of the segment start and finish
			const coordinates: number[][] = [
				[pointA.latitude, pointA.longitude],
				[pointB.latitude, pointB.longitude],
			];

			// fetches info from OSM
			const response = await fetch(`${ environment.RouteUrl }/v2/directions/foot-walking/geojson`, {
				method: "POST",
				headers: {
					"Authorization": environment.RouteKey,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ "coordinates": coordinates }),
			});
			const data = await response.json();
			const segments = data.features[0].properties.segments[0];
			const distance: number = segments.distance; // the length of the path in meters
			const duration: number = segments.duration; // the duration of the path in seconds

			if (distance > 1000) // we should display it in kilometers
				this.travelLengthBetweenPoints.push(`${ roundToOne(distance / 1000) } км`);
			else // we should display it in meters
				this.travelLengthBetweenPoints.push(`${ roundToOne(distance) } м`);

			// duration in minutes
			const minutes_duration: number = Math.max(1, roundToOne(duration / 60, 0));
			if (minutes_duration > 60) { // we should display duration in hours
				const hours_duration: number = roundToOne(minutes_duration / 60, 0);
				currentTime.setHours(currentTime.getHours() + hours_duration);
				this.travelTimeBetweenPoints.push(`${ hours_duration } ч`);
				this.arrivalTimeToPoints.push(dateToString(currentTime));
			} else { // we should display duration in minutes
				currentTime.setMinutes(currentTime.getMinutes() + minutes_duration);
				this.travelTimeBetweenPoints.push(`${ minutes_duration } мин`);
				this.arrivalTimeToPoints.push(dateToString(currentTime));
			}
		}
	}

	/** Opens information about the point on the path */
	openPointInformation(point: Place | EventPlace | MapPointModel) {
		this.mapInteractionsService.chosenMapPoint.next(point);
	}

	closeCallback(notificationsService: NotificationsService, callback: any) {
		notificationsService.addNotification({
			timeOut: 5000,
			message: "Маршрут закрыт",
			type: "info",
			callback: callback,
		});
	}

	notificationCallback(service: MapInteractionsService, path: any) {
		service.pathPoints.next(path);
	}

	/** Closes path point */
	close() {
		this.mapInteractionsService.pathInformationState.next(-1);

		this.closeCallback(this.notificationsService, this.notificationCallback.bind({}, this.mapInteractionsService, this.path));

		setTimeout(() => {
			this.mapInteractionsService.pathPoints.next(null);
		}, 400);
	}

	/** Function that hides container element */
	Hide() {
		if (this.container === undefined) return;
		this.container.nativeElement.style.transitionDuration = ".3s";
		setTimeout(() => {
			this.container.nativeElement.style.maxHeight = "0px";
		}, 100);
		setTimeout(() => {
			this.container.nativeElement.style.display = "none";
			this.container.nativeElement.style.transitionDuration = "0s";
			this.mapInteractionsService.pathInformationState.next(-2);
		}, 400);
	}

	/** Function that shows container element */
	private Show() {
		this.container.nativeElement.style.display = "flex";
		this.container.nativeElement.style.transitionDuration = ".3s";
		setTimeout(() => {
			this.container.nativeElement.style.maxHeight = `${ this.INITIAL_HEIGHT }px`;
		}, 100);
		setTimeout(() => {
			this.mapInteractionsService.pathInformationState.next(2);
			this.container.nativeElement.style.transitionDuration = "0s";
		}, 400);
	}
}
