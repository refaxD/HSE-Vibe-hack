import { Directive, effect, ElementRef, inject, Input, OnDestroy, signal, WritableSignal } from "@angular/core";
import { AppComponent } from "../app/app.component";

@Directive({
	selector: "[moveable]",
	standalone: true,
})
export class MoveableDirective implements OnDestroy {
	elementRef = inject(ElementRef<HTMLElement>);

	element!: HTMLElement;

	private readonly FAST_DURATION: number = 500;
	private readonly MIN_DURATION: number = 100;

	/** Maximum height of the container */
	@Input({ required: false, alias: "maxHeight" })
	MAX_HEIGHT = window.innerHeight;

	/** Minimum height of the container */
	@Input({ required: false, alias: "minHeight" })
	MIN_HEIGHT: number = 0;

	/** Function that will be call after Hide function finishes */
	@Input({ required: false, alias: "callback" })
	callback: any = () => { return 0; };

	@Input({ required: false, alias: "valid" })
	valid: WritableSignal<boolean> = signal(true);

	@Input({ required: false, alias: "isPopup" })
	popup: boolean = true;

	@Input({ required: false, alias: "state" })
	stateSignal: WritableSignal<boolean> = signal<boolean>(false);

	private localStateSignal: WritableSignal<boolean> = signal<boolean>(false);

	private timeouts: any[] = [];

	constructor() {
		let lastCall = 0;

		this.element = this.elementRef.nativeElement;
		this.element.addEventListener("touchstart",
			(ev: TouchEvent) => {
				const now = Date.now();
				if (now - lastCall > 100) {
					lastCall = now;
					this.handleTouchStart(this.element, ev);
				}
			},
			{ passive: true },
		);

		effect(() => {
			const state = this.stateSignal();
			if (this.ifLocalStateChange())
				return;
			this.localStateSignal.set(state);
			if (state) this.Open(this.element);
			else this.Hide(this.element);
		});
	}

	ngOnDestroy() {
		this.element.removeEventListener("touchstart", (ev: TouchEvent) => this.handleTouchStart(this.element, ev));
	}

	ifLocalStateChange() {
		return this.localStateSignal() === this.stateSignal();
	}

	changeState(state: boolean) {
		this.localStateSignal.set(state);
		this.stateSignal.set(state);
	}

	/** Handles container touch start event */
	handleTouchStart(container: HTMLElement, event: TouchEvent) {
		if (container.scrollTop !== 0) return;
		if (!this.valid()) return;
		event.stopPropagation();
		this.stopAnimation(container);

		const startTouchY = event.touches[0].clientY;
		const timeStart = new Date(); // touch start
		let lastY = startTouchY;
		let deltaWithOutScroll = 0; // height change without scroll effecting
		let direction = "";
		let moving = false;

		/** handles touch move event */
		const handleTouchMove = (ev: TouchEvent) => {
			if (moving) return;
			moving = true;
			const currentY = ev.touches[0].clientY; // touch position
			const currentHeight = container.getBoundingClientRect().height; // container height
			const currentScroll = container.scrollTop; // container scroll

			let height = 0; // new container height
			let delta = lastY - currentY; // delta touch

			if (delta > 0) { // scrolling to top ⬆
				/* we should first add maxHeight of the container
				* then if it is 100vh we should scroll the content */
				direction = "top";
				let addToHeight = Math.min(delta, this.MAX_HEIGHT - currentHeight);

				if (addToHeight > 0) { // we should add to maxHeight of the container
					height = currentHeight + addToHeight;
					container.style.maxHeight = `${ height }px`;
					delta -= addToHeight;
					deltaWithOutScroll += addToHeight;
				}

				if (delta > 0 && this.popup) { // we should scroll the content
					this.changeState(true);
					container.style.overflowY = "scroll";
					container.scrollTo({ top: currentScroll + delta });
				}
			}
			else if (delta < 0) { // scrolling to bottom ⬇
				/* here first we should scroll content till the top
				* then we decrease container maxHeight */
				direction = "bottom";
				delta = Math.abs(delta);
				let addToScroll = Math.min(delta, currentScroll);

				if (addToScroll > 0) { // we should scroll the content
					container.scrollTo({ top: currentScroll - addToScroll });
					delta -= addToScroll;
				}

				if (delta > 0) { // we should add to maxHeight of the container
					container.style.overflowY = "hidden";
					height = Math.max(this.MIN_HEIGHT, currentHeight - delta);
					container.style.maxHeight = `${ height }px`;
					deltaWithOutScroll += delta;
				}
			}

			if (height !== 0 && this.popup) {
				const percent = Math.sqrt(( height - this.MIN_HEIGHT ) / ( this.MAX_HEIGHT - this.MIN_HEIGHT ));
				AppComponent.bottomBarState.set({ percent: percent, state: false });
			}

			lastY = currentY;
			moving = false;
		};
		/** handles touch end event */
		const handleTouchEnd = (ev: TouchEvent) => {
			const timeEnd = new Date(); // touch finish
			const duration = timeEnd.getTime() - timeStart.getTime(); // touch duration

			const currentHeight = container.getBoundingClientRect().height;
			const validHeight = this.MAX_HEIGHT - this.MIN_HEIGHT;
			const min_interaction_height = validHeight / 2 + this.MIN_HEIGHT;
			container.style.transitionDuration = ".3s";

			if (duration < this.MIN_DURATION || deltaWithOutScroll === 0) {
				if (currentHeight < ( this.MAX_HEIGHT - this.MIN_HEIGHT ) / 2 + this.MIN_HEIGHT)
					this.Hide(container);
				else
					this.Open(container);
				return;
			}

			setTimeout(() => {
				// setting new height to an object
				if (direction === "bottom") { // to bottom
					if (container.scrollTop !== 0)
						return;

					if (duration < this.FAST_DURATION)
						this.Hide(container);
					else if (currentHeight <= min_interaction_height)
						this.Hide(container);
					else
						this.Open(container);
				}
				else { // to top
					if (duration < this.FAST_DURATION)
						this.Open(container);
					else if (currentHeight >= min_interaction_height)
						this.Open(container);
					else
						this.Hide(container);
				}
			});

			setTimeout(() => {
				container.style.transitionDuration = "0s";
				if (container.getBoundingClientRect().height === 0)
					this.Hide(container);
			}, 300);

			container.removeEventListener("touchmove", handleTouchMove);
			container.removeEventListener("touchend", handleTouchEnd);
		};

		container.addEventListener("touchmove", handleTouchMove);
		container.addEventListener("touchend", handleTouchEnd);
	}

	private stopAnimation(container: HTMLElement) {
		if (container === undefined) return;

		this.timeouts.forEach(t => {
			clearTimeout(t);
		});
		this.timeouts = [];

		const height = container.getBoundingClientRect().height;
		container.style.maxHeight = `${ height }px`;
		setTimeout(() => {
			container.style.transitionDuration = "0s";
		});

		if (this.popup) {
			const percent = Math.sqrt(( height - this.MIN_HEIGHT ) / ( this.MAX_HEIGHT - this.MIN_HEIGHT ));
			AppComponent.bottomBarState.set({ percent: percent, state: false });
		}
	}

	private Open(container: HTMLElement) {
		if (container === undefined) return;

		this.changeState(true);
		if (this.popup) {
			AppComponent.bottomBarState.set({ percent: 1, state: true });
			container.style.overflowY = "scroll";
		}

		container.style.transitionDuration = ".3s";
		if (this.MIN_HEIGHT === 0)
			container.style.display = "flex";
		this.timeouts.push(
			setTimeout(() => {
				container.style.maxHeight = `${ this.MAX_HEIGHT }px`;
				//this.timeouts = this.timeouts.splice(0, 1);
			}, 0),
		);
		this.timeouts.push(
			setTimeout(() => {
				container.style.transitionDuration = "0s";
				//this.timeouts = this.timeouts.splice(0, 1);
			}, 300),
		);
	}

	/** Hides MapPointInformation content */
	private Hide(container: HTMLElement) {
		if (container === undefined) return;

		this.changeState(false);
		if (this.popup) {
			AppComponent.bottomBarState.set({ percent: 0, state: true });
			container.style.overflowY = "hidden";
		}

		container.style.transitionDuration = ".3s";
		this.timeouts.push(
			setTimeout(() => {
				container.style.maxHeight = `${ this.MIN_HEIGHT }px`;
				//this.timeouts = this.timeouts.splice(0, 1);
			}, 0),
		);
		this.timeouts.push(
			setTimeout(() => {
				if (this.MIN_HEIGHT === 0)
					container.style.display = "none";
				container.style.transitionDuration = "0s";
				//this.timeouts = this.timeouts.splice(0, 1);
			}, 300),
		);
	}
}