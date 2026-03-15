import {
	AfterViewInit,
	Component,
	effect,
	ElementRef, Input,
	QueryList,
	signal,
	ViewChild,
	ViewChildren, WritableSignal,
} from "@angular/core";
import { NavigationEnd, Router, RouterLink, RouterLinkActive } from "@angular/router";
import { filter } from "rxjs";
import { EventPlace, Place } from "../../../../generated";
import { MoveableDirective } from "../../../../libs/moveable.directive";
import { MapClickEntityModel } from "../../models/map-click-entity.model";
import { PlaceSearchPipe } from "../../pipes/place-search.pipe";
import { MapInteractionsService } from "../../services/map-interactions.service";
import { NotificationsService } from "../../services/notifications.service";
import { isInstanceOfEventPlace, isInstanceOfPlace } from "../../services/utils.service";
import { SlideCategoriesComponent } from "../slide-categories.component";

@Component({
	selector: "BottomBarContent",
	standalone: true,
	imports: [
		PlaceSearchPipe,
	],
	styleUrl: "./bottom-bar-content.component.css",
	template: `
        <div class="container">
			@if (this.searchState() && searchArray().length > 0) {
                @for (item of searchArray() | placeSearch:this.searchText; track [searchArray(), searchText]) {
                    <div class="card" (click)="chosenPoint.set(item)">
                        <p class="name">{{ item.name }}</p>
                        @if (item.address !== null) {
                            <p class="address">{{ item.address }}</p>
                        }
                    </div>
                }
			} @else {
				<div class="scroll-container">
					<div class="scroll-content">
	
						<div class="place">
							<img src="intro/background.avif" alt="">
							<div class="marker">
								<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="var(--red-50)">
									<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
									<path d="M18.364 4.636a9 9 0 0 1 .203 12.519l-.203 .21l-4.243 4.242a3 3 0 0 1 -4.097 .135l-.144 -.135l-4.244 -4.243a9 9 0 0 1 12.728 -12.728zm-6.364 3.364a3 3 0 1 0 0 6a3 3 0 0 0 0 -6z"/>
								</svg>
							</div>
							<div class="content">
								<p class="title">Mukebu, Norway</p>
								<p class="text">Горы в дальней части норвегии, круто, стильно, моложежно</p>
								<div class="stars">
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
								</div>
							</div>
						</div>
						
						<div class="place">
							<img src="intro/background.avif" alt="">
							<div class="marker">
								<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="var(--red-50)">
									<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
									<path d="M18.364 4.636a9 9 0 0 1 .203 12.519l-.203 .21l-4.243 4.242a3 3 0 0 1 -4.097 .135l-.144 -.135l-4.244 -4.243a9 9 0 0 1 12.728 -12.728zm-6.364 3.364a3 3 0 1 0 0 6a3 3 0 0 0 0 -6z"/>
								</svg>
							</div>
							<div class="content">
								<p class="title">Mukebu, Norway</p>
								<p class="text">Горы в дальней части норвегии, круто, стильно, моложежно</p>
								<div class="stars">
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
									<svg xmlns="http://www.w3.org/2000/svg" width="14" viewBox="0 0 24 24" fill="var(--yellow-30)">
										<path stroke="none" d="M0 0h24v24H0z" fill="none"/>
										<path d="M8.243 7.34l-6.38 .925l-.113 .023a1 1 0 0 0 -.44 1.684l4.622 4.499l-1.09 6.355l-.013 .11a1 1 0 0 0 1.464 .944l5.706 -3l5.693 3l.1 .046a1 1 0 0 0 1.352 -1.1l-1.091 -6.355l4.624 -4.5l.078 -.085a1 1 0 0 0 -.633 -1.62l-6.38 -.926l-2.852 -5.78a1 1 0 0 0 -1.794 0l-2.853 5.78z"/>
									</svg>
								</div>
							</div>
						</div>
	
					</div>
				</div>
            }
        </div>
	`,
})
export class BottomBarContentComponent {
	@Input() searchState!: WritableSignal<boolean>;
	@Input() searchArray!: WritableSignal<Array<any>>;
	@Input() searchText!: string;
	@Input() chosenPoint!: WritableSignal<EventPlace | Place | null>;
}
