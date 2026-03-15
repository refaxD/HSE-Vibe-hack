import { Component, effect, ElementRef, signal, ViewChild } from "@angular/core";
import { Router } from "@angular/router";
import { MoveableDirective } from "../../../../../libs/moveable.directive";
import { MapInteractionsService } from "../../../../core/services/map-interactions.service";

@Component({
	selector: "intro-main",
	standalone: true,
	imports: [
		MoveableDirective,
	],
	templateUrl: "./main.component.html",
	styleUrl: "./main.component.css",
})
export default class MainComponent {
	@ViewChild("circle") circle!: ElementRef<HTMLDivElement>;

	disabled = signal(false);

	state = signal<boolean>(false);

	finish = signal(false);

	constructor(
		private router: Router,
		private mapInteractionsService: MapInteractionsService,
	) {
		effect(() => {
			const location = this.mapInteractionsService.userLocation();
			this.disabled.set(location !== null);
		});
		effect(() => {
			const state = this.state();
			if (!state) return;
			setTimeout(() => {
				this.finish.set(true);
			}, 300);
		});
	}
}
