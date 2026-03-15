import { booleanAttribute, Component, effect, ElementRef, signal, ViewChild } from "@angular/core";
import { EventPlace, Place } from "../../../../generated";
import { MoveableDirective } from "../../../../libs/moveable.directive";
import MapPointModel from "../../../core/models/map-point.model";
import { BypassHtmlSanitizerPipe } from "../../../core/pipes/doom-sanitizer.pipe";
import { MapInteractionsService } from "../../../core/services/map-interactions.service";
import { isInstanceOfEventPlace, isInstanceOfPlace } from "../../../core/services/utils.service";

@Component({
	selector: "MapPointInformation",
	standalone: true,
	imports: [
		BypassHtmlSanitizerPipe,
		MoveableDirective,
	],
	styleUrls: ["./map-point-information.component.css"],
	template: `
        <div style="height: 100vh;" moveable [isPopup]="false" [state]="state" [maxHeight]="maxHeight" class="container" #container>
            <div class="content">
                <!-- container with the arrow -->
                <div class="flex flex-row justify-center w-full">
                    <svg width="30" height="3" fill="#6F6F6F">
                        <path d="M0 1.5C0 0.671573 0.671573 0 1.5 0H28.5C29.3284 0 30 0.671573 30 1.5V1.5C30 2.32843 29.3284 3 28.5 3H1.5C0.671573 3 0 2.32843 0 1.5V1.5Z"></path>
                    </svg>
                </div>
                <!-- container with the information about the point -->
                <div class="router">
                    <!-- point name and address -->
                    @if (point !== null) {
                        <div class="main-header">
                            <p class="name allow-selection">{{ point.name }}</p>
                            <div>
                                <svg (click)="close()" width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M17.707 16.293a1 1 0 0 1-1.414 1.414L12 13.414l-4.293 4.293a1 1 0 0 1-1.414-1.414L10.586 12 6.293 7.707a1 1 0 0 1 1.414-1.414L12 10.586l4.293-4.293a1 1 0 1 1 1.414 1.414L13.414 12l4.293 4.293z"></path>
                                </svg>
                            </div>
                        </div>
                        @if (isInstanceOfPlace(point)) {
                            <p class="address allow-selection">{{ point.address }}</p>
                        }
                    }
                    <!-- router routes -->
                    @if (isInstanceOfPlace(point)) {
                        <div class="routes">
                            <p (click)="route = 'Обзор'" [class.chosen]="route === 'Обзор'">Обзор</p>
                            <p (click)="route = 'Особенности'" [class.chosen]="route === 'Особенности'">Особенности</p>
                        </div>
                    }
                    <!-- router content -->
                    <div class="router-content">
                        @if (route === 'Обзор') {
                            @if (isInstanceOfPlace(point)) {
                                @if (point.email !== null || point.website !== null || point.phone !== null) {
                                    <div class="block">
                                        <div class="header-container">
                                            <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                <path d="M20.186 19.742c1.15-1.15.883-2.424.404-2.707-.336-.198-4.749-2.684-4.749-2.684-.344-.216-.686-.106-.893.142l-.005-.004-1.626 1.625a.674.674 0 0 1-.824.1 14.052 14.052 0 0 1-2.632-2.075 14.054 14.054 0 0 1-2.074-2.632.674.674 0 0 1 .1-.824L9.51 9.057l-.004-.005c.243-.203.361-.544.143-.893 0 0-2.487-4.413-2.685-4.75-.283-.478-1.556-.745-2.707.405-2.566 2.568-1.081 8.207 3.32 12.608 4.398 4.399 10.04 5.887 12.608 3.32z"></path>
                                            </svg>
                                            <p class="header">Контакты</p>
                                        </div>
                                        <div class="block-info">
                                            @if (point.phone !== null) {
                                                <a [href]="getPhoneNumber(point.phone!)" class="phone">+7 {{
                                                        point.phone
                                                    }}</a>
                                            }
                                            @if (point.email !== null) {
                                                <a [href]="'mailto:' + point.email!" class="email">{{ point.email }}</a>
                                            }
                                            @if (point.website !== null) {
                                                <a [href]="'https://' + point.website!" class="website">{{
                                                        point.website
                                                    }}</a>
                                            }
                                        </div>
                                    </div>
                                }

                                @if (point.schedule !== null || point.data?.time !== undefined) {
                                    <div class="block">
                                        <div class="header-container">
                                            <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                <path fill-rule="evenodd" clip-rule="evenodd" d="M12 20a8 8 0 1 1 0-16 8 8 0 0 1 0 16zm.523-11.775a.25.25 0 0 0-.25-.225h-.545a.25.25 0 0 0-.25.227l-.476 4.203c-.028.433.22.856.645 1.022l3.205 1.21a.3.3 0 0 0 .374-.14l.286-.537a.3.3 0 0 0-.1-.39l-2.417-1.714-.472-3.654v-.002z"></path>
                                            </svg>
                                            <p class="header">
                                                Часы работы
                                            </p>
                                        </div>
                                        <div class="block-info">
                                            @if (checkIfOpenedToday(point)) {
                                                @if (point.schedule !== null) {
                                                    @for (day of point.schedule; track point) {
                                                        <div class="flex flex-row justify-between items-center">
                                                            <p style="text-transform: capitalize">{{
                                                                    day.DayOfWeek
                                                                }}</p>
                                                            <p>{{ day.Hours.split('-').join(' - ') }}</p>
                                                        </div>
                                                    }
                                                } @else {
                                                    <p>открыто</p>
                                                }
                                            } @else {
                                                <p class="text-red-700">Сегодня не работает</p>
                                            }
                                        </div>
                                    </div>
                                }
                            } @else if (isInstanceOfEventPlace(point)) {
                            } @else if (point !== null) {

                            }
                        } @else if (route === 'Особенности') {
                            @if (isInstanceOfPlace(point)) {
                                <div class="flex flex-row items-center flex-wrap" style="gap: 10px;">
                                    @for (key of getKeys(point.data!); track point) {
                                        <div class="flex flex-row items-center feature-svg-container" style="gap: 5px;">
                                            <div [innerHtml]="pointDataSvg[key] | bypassHtmlSanitizer"></div>
                                            <p>{{ pointDataNames[key] }}</p>
                                        </div>
                                    }
                                </div>
                            }
                        }
                    </div>
                </div>
                <!-- container with buttons -->
                <div class="buttons-container">
                    <div class="scroll-content">
                        <button (click)="AddPointToPath()" class="solid-btn">Добавить</button>
                        <button class="void-btn">
                            <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 24 24" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                                <path d="M12 5l0 14"/>
                                <path d="M5 12l14 0"/>
                            </svg>
                        </button>
                        <button class="void-btn">
                            <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                                <path d="M19.5 12.572l-7.5 7.428l-7.5 -7.428a5 5 0 1 1 7.5 -6.566a5 5 0 1 1 7.5 6.572"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>
	`,
})
export class MapPointInformationComponent {
	/** Point */
	point: Place | EventPlace | MapPointModel | null = null;

	/** Route for private router */
	route: string = "Обзор";
	/** names for point features */
	public readonly pointDataNames: Record<string, string> = {
		"hasFoodPoint": "точка питания",
		"hasChangeRoom": "раздевалка",
		"hasToilet": "туалет",
		"hasWIFI": "интернет",
		"hasWater": "водоём",
		"hasChild": "игровая площадка",
		"hasSport": "спорт площадка",
	};
	/** svg elements for point features */
	public readonly pointDataSvg: Record<string, string> = {
		"hasFoodPoint": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path fill='#C6C6C6FF' d=\"M16.588 5.191l.058 .045l.078 .074l.072 .084l.013 .018a.998 .998 0 0 1 .182 .727l-.022 .111l-.03 .092c-.99 2.725 -.666 5.158 .679 7.706a4 4 0 1 1 -4.613 4.152l-.005 -.2l.005 -.2a4.002 4.002 0 0 1 2.5 -3.511c-.947 -2.03 -1.342 -4.065 -1.052 -6.207c-.166 .077 -.332 .15 -.499 .218l.094 -.064c-2.243 1.47 -3.552 3.004 -3.98 4.57a4.5 4.5 0 1 1 -7.064 3.906l-.004 -.212l.005 -.212a4.5 4.5 0 0 1 5.2 -4.233c.332 -1.073 .945 -2.096 1.83 -3.069c-1.794 -.096 -3.586 -.759 -5.355 -1.986l-.268 -.19l-.051 -.04l-.046 -.04l-.044 -.044l-.04 -.046l-.04 -.05l-.032 -.047l-.035 -.06l-.053 -.11l-.038 -.116l-.023 -.117l-.005 -.042l-.005 -.118l.01 -.118l.023 -.117l.038 -.115l.03 -.066l.023 -.045l.035 -.06l.032 -.046l.04 -.051l.04 -.046l.044 -.044l.046 -.04l.05 -.04c4.018 -2.922 8.16 -2.922 12.177 0z\" /></svg>",
		"hasChangeRoom": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"  stroke=\"currentColor\"  stroke-width=\"2\"  stroke-linecap=\"round\"  stroke-linejoin=\"round\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M12 9l-7.971 4.428a2 2 0 0 0 -1.029 1.749v.823a2 2 0 0 0 2 2h1\" /><path d=\"M18 18h1a2 2 0 0 0 2 -2v-.823a2 2 0 0 0 -1.029 -1.749l-7.971 -4.428c-1.457 -.81 -1.993 -2.333 -2 -4a2 2 0 1 1 4 0\" /><path d=\"M6 16m0 2a2 2 0 0 1 2 -2h8a2 2 0 0 1 2 2v1a2 2 0 0 1 -2 2h-8a2 2 0 0 1 -2 -2z\" /></svg>",
		"hasToilet": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\" ><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path fill='#C6C6C6FF' d=\"M19 4a3 3 0 0 1 3 3v10a3 3 0 0 1 -3 3h-14a3 3 0 0 1 -3 -3v-10a3 3 0 0 1 3 -3zm-7.534 4a1 1 0 0 0 -.963 .917l-.204 2.445l-.405 -.81l-.063 -.11a1 1 0 0 0 -1.725 .11l-.406 .81l-.203 -2.445a1 1 0 0 0 -.963 -.917l-.117 .003a1 1 0 0 0 -.914 1.08l.5 6l.016 .117c.175 .91 1.441 1.115 1.875 .247l1.106 -2.211l1.106 2.211c.452 .904 1.807 .643 1.89 -.364l.5 -6a1 1 0 0 0 -.913 -1.08zm4.034 0a2.5 2.5 0 0 0 -2.5 2.5v3a2.5 2.5 0 1 0 5 0a1 1 0 0 0 -2 0a.5 .5 0 1 1 -1 0v-3a.5 .5 0 1 1 1 0a1 1 0 0 0 2 0a2.5 2.5 0 0 0 -2.5 -2.5\" /></svg>",
		"hasWIFI": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"  stroke=\"currentColor\"  stroke-width=\"2\"  stroke-linecap=\"round\"  stroke-linejoin=\"round\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M12 18l.01 0\" /><path d=\"M9.172 15.172a4 4 0 0 1 5.656 0\" /><path d=\"M6.343 12.343a8 8 0 0 1 11.314 0\" /><path d=\"M3.515 9.515c4.686 -4.687 12.284 -4.687 17 0\" /></svg>",
		"hasWater": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"  stroke=\"currentColor\"  stroke-width=\"2\"  stroke-linecap=\"round\"  stroke-linejoin=\"round\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M3 7c3 -2 6 -2 9 0s6 2 9 0\" /><path d=\"M3 17c3 -2 6 -2 9 0s6 2 9 0\" /><path d=\"M3 12c3 -2 6 -2 9 0s6 2 9 0\" /></svg>",
		"hasChild": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"  stroke=\"currentColor\"  stroke-width=\"2\"  stroke-linecap=\"round\"  stroke-linejoin=\"round\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M3 21v-15l5 -3l5 3v15\" /><path d=\"M8 21v-7\" /><path d=\"M3 14h10\" /><path d=\"M6 10a2 2 0 1 1 4 0\" /><path d=\"M13 13c6 0 3 8 8 8\" /></svg>",
		"hasSport": "<svg  width=\"24\"  height=\"24\"  viewBox=\"0 0 24 24\"  stroke=\"currentColor\"  stroke-width=\"2\"  stroke-linecap=\"round\"  stroke-linejoin=\"round\"><path stroke=\"none\" d=\"M0 0h24v24H0z\" fill=\"none\"/><path d=\"M11 4a1 1 0 1 0 2 0a1 1 0 0 0 -2 0\" /><path d=\"M3 17l5 1l.75 -1.5\" /><path d=\"M14 21v-4l-4 -3l1 -6\" /><path d=\"M6 12v-3l5 -1l3 3l3 1\" /><path d=\"M19.5 20a.5 .5 0 1 0 0 -1a.5 .5 0 0 0 0 1z\" fill=\"currentColor\" /></svg>",
	};

	/** Container */
	@ViewChild("container") container!: ElementRef<HTMLDivElement>;

	protected readonly isInstanceOfPlace = isInstanceOfPlace;
	protected readonly isInstanceOfEventPlace = isInstanceOfEventPlace;

	state = signal<boolean>(false);

	protected get maxHeight() {
		return window.innerHeight * 0.7;
	}

	constructor(
		public mapInteractionsService: MapInteractionsService,
	) {
		/** Adding listener for point changes */
		this.mapInteractionsService.chosenMapPoint.subscribe(point => {
			this.point = point;
			this.state.set(this.point !== null);
		});
		/** effect for moveable state */
		effect(() => {
			const state = this.state();

			if (state) {
			} else {
				this.hideCallback(this.mapInteractionsService);
			}
		})
	}

	/** Adds this.point to current path */
	async AddPointToPath() {
		let path = this.mapInteractionsService.pathPoints.value;
		if (path === null) path = [];
		path.push(this.point!);
		this.mapInteractionsService.setNewPath(path);
		this.state.set(false);
	}

	/** Closes map point information container */
	close() {
		this.mapInteractionsService.chosenMapPoint.next(null);
	}

	/** Function that converts phone number into a link with tel */
	getPhoneNumber(str: string) {
		let result = "";
		for (const char of str)
			if ("0" <= char && char <= "9")
				result += char;
		return `tel:+7${ result }`;
	}

	/** Function that checks if point is working today */
	checkIfOpenedToday(point: Place | EventPlace | MapPointModel | null) {
		function convertDate(date: Date) {
			return `${ String(date.getDate()).padStart(2, "0") }:${ String(date.getMonth() + 1).padStart(2, "0") }`;
		}

		if (point === null) return true;
		if (!( "data" in point )) return true;
		if (point.data === null || point.data === undefined) return true;
		if (!( "time" in point.data! )) return true;
		if (point.data.time === undefined || point.data.time === null) return true;
		if (typeof point.data.time !== "string") return true;

		const [start, end] = point.data.time.split("-");
		const current = convertDate(new Date());

		if (current < start) return false;
		if (end < current) return false;
		return true;
	}

	/** Function that gets all keys from data that should be displayed in features */
	getKeys(data: Record<string, any>): string[] {
		const result: string[] = [];
		for (const key of Object.keys(data))
			if (key in this.pointDataNames && data[key] === "да")
				result.push(key);
		return result;
	}

	hideCallback(mapInteractionsService: MapInteractionsService) {
		if (mapInteractionsService.pathPoints.value !== null)
			mapInteractionsService.pathInformationState.next(1);
	}
}
