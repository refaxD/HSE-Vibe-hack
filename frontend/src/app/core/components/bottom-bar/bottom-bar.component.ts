import {
	AfterViewInit,
	Component,
	effect,
	ElementRef,
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
import { BottomBarContentComponent } from "./bottom-bar-content.component";

/**
 * places Caret to the end of the editor
 *
 * @param {HTMLDivElement} editor editor in which caret should be placed
 * @param {number} pos place where cater should be places
 */
function placeCaretAt(editor: HTMLDivElement, pos: number) {
	const range = document.createRange();
	const sel = window.getSelection()!;

	const children = editor.childNodes;
	let [index, position] = [0, pos];
	while (index < children.length) {
		//@ts-ignore
		const length = ( children[index].nodeValue || children[index].innerText ).length;
		if (position <= length) break;
		position -= length;
		index++;
	}
	if (children[index] === undefined)
		return;
	else if (children[index].nodeValue === null)
		range.setStart(children[index].childNodes[0], position);
	else
		range.setStart(children[index], position);
	range.collapse(true);

	sel.removeAllRanges();
	sel.addRange(range);
}

/** Returns caret position in editor */
function getCaretPosition(editor: HTMLDivElement) {
	let caretOffset = 0;
	const selection = window.getSelection()!;
	if (selection.rangeCount > 0) {
		const range = selection.getRangeAt(0);
		const preCaretRange = range.cloneRange();
		preCaretRange.selectNodeContents(editor);
		preCaretRange.setEnd(range.endContainer, range.endOffset);
		caretOffset = preCaretRange.toString().length;
	}
	return caretOffset;
}

/** swaps two elements in array */
function swapElements(arr: any[], index1: number, index2: number): void {
	[arr[index1], arr[index2]] = [arr[index2], arr[index1]];
}

@Component({
	selector: "BottomBar",
	standalone: true,
	imports: [
		SlideCategoriesComponent,
		PlaceSearchPipe,
		MoveableDirective,
		BottomBarContentComponent,
	],
	styleUrl: "./bottom-bar.component.css",
	template: `
        <div [minHeight]="MIN_HEIGHT" [state]="bottomBarState" [maxHeight]="window.innerHeight - 67" moveable class="container" #container>
            <div class="flex flex-col gap-3">
                <div class="map-addon" #mapAddonContainer>
                    <div class="fold-mark-container">
                        <svg class="fold-mark">
                            <path d="" fill="#6F6F6F" #foldMark/>
                        </svg>
                    </div>

                    <div class="content">
                        <div class="prompt-input-container">
                            <div class="prompt-input-scroll" #promptInputScroll>
                                <div dir="auto" tabindex="0" class="input allow-selection" (click)="handleClickPromptInput()" (keydown)="handleKeyDownPromptInput($event)" (keyup)="handleCaretPositionChanges($event)" (input)="handlePromptInput()" (click)="handleCaretPositionChanges($event)" role="textbox" contenteditable="true" #promptInput></div>
                                <span class="message-placeholder" #promptPlaceholder>Prompt...</span>
                            </div>
                        </div>

                        <CategoriesSlider/>
                    </div>
                </div>

				<div class="main">
                    <BottomBarContent [searchArray]="searchArray" [searchState]="searchState" [searchText]="searchText" [chosenPoint]="chosenPoint"/>
				</div>
            </div>
        </div>

        <div class="user-position" #openMenu>
            <p>{{ this.mapInteractionService.userLocation()?.city }}</p>
        </div>

        <div class="background" #openMenu></div>

        <button (click)="stickUserPosition()" class="stick-position-marker" [style.background-color]="mapInteractionService.mapScrolled() === -1 ? 'var(--blue-60)' : 'var(--neutral-1)'" #stickPositionButton>
            <svg width="30" height="30" viewBox="0 0 30 30" xmlns="http://www.w3.org/2000/svg">
                <path d="M19.37 15.7965C19.2479 15.7479 19.1174 15.7239 18.9859 15.7257C18.8545 15.7276 18.7247 15.7553 18.604 15.8073C18.4833 15.8593 18.374 15.9345 18.2824 16.0288C18.1908 16.123 18.1186 16.2344 18.07 16.3565C18.0214 16.4786 17.9974 16.6091 17.9992 16.7406C18.001 16.872 18.0288 17.0018 18.0808 17.1225C18.1328 17.2432 18.208 17.3525 18.3023 17.4441C18.3965 17.5358 18.5079 17.6079 18.63 17.6565C20.09 18.2365 21 19.1365 21 20.0065C21 21.4265 18.54 23.0065 15 23.0065C11.46 23.0065 9 21.4265 9 20.0065C9 19.1365 9.91 18.2365 11.37 17.6565C11.6166 17.5584 11.8142 17.3663 11.9192 17.1225C12.0243 16.8787 12.0281 16.6032 11.93 16.3565C11.8319 16.1099 11.6398 15.9123 11.396 15.8073C11.1522 15.7023 10.8767 15.6984 10.63 15.7965C8.36 16.6965 7 18.2665 7 20.0065C7 22.8065 10.51 25.0065 15 25.0065C19.49 25.0065 23 22.8065 23 20.0065C23 18.2665 21.64 16.6965 19.37 15.7965ZM14 12.8665V20.0065C14 20.2717 14.1054 20.5261 14.2929 20.7136C14.4804 20.9012 14.7348 21.0065 15 21.0065C15.2652 21.0065 15.5196 20.9012 15.7071 20.7136C15.8946 20.5261 16 20.2717 16 20.0065V12.8665C16.9427 12.6231 17.7642 12.0443 18.3106 11.2385C18.857 10.4327 19.0908 9.45533 18.9681 8.48951C18.8454 7.5237 18.3747 6.63578 17.6442 5.99219C16.9137 5.3486 15.9736 4.99353 15 4.99353C14.0264 4.99353 13.0863 5.3486 12.3558 5.99219C11.6253 6.63578 11.1546 7.5237 11.0319 8.48951C10.9092 9.45533 11.143 10.4327 11.6894 11.2385C12.2358 12.0443 13.0573 12.6231 14 12.8665ZM15 7.00651C15.3956 7.00651 15.7822 7.12381 16.1111 7.34357C16.44 7.56334 16.6964 7.87569 16.8478 8.24115C16.9991 8.6066 17.0387 9.00873 16.9616 9.39669C16.8844 9.78466 16.6939 10.141 16.4142 10.4207C16.1345 10.7004 15.7781 10.8909 15.3902 10.9681C15.0022 11.0453 14.6001 11.0057 14.2346 10.8543C13.8692 10.7029 13.5568 10.4466 13.3371 10.1177C13.1173 9.78876 13 9.40208 13 9.00651C13 8.47608 13.2107 7.96737 13.5858 7.5923C13.9609 7.21723 14.4696 7.00651 15 7.00651Z"/>
            </svg>
        </button>
	`,
})
export class BottomBarComponent implements AfterViewInit {
	@ViewChildren("openMenu") openMenuItems!: QueryList<ElementRef<HTMLDivElement>>; // elements that should be visible when menu openes

	@ViewChild("container") container!: ElementRef<HTMLDivElement>; // Main container of all elements

	@ViewChild("mainContainer") mainContainer!: ElementRef<HTMLDivElement>; // container with main elements
	@ViewChild("mapAddonContainer") mainAddonContainer!: ElementRef<HTMLDivElement>; // container with prompt input
	@ViewChild("foldMark") foldMark!: ElementRef<SVGPathElement>; // mark for folding the container when we are on the /home route

	@ViewChild("stickPositionButton") stickPositionButton!: ElementRef<HTMLButtonElement>; // button that sticks our position to geolocation

	@ViewChild("promptInput") promptInput!: ElementRef<HTMLDivElement>; // prompt input element
	@ViewChild("promptPlaceholder") promptPlaceholder!: ElementRef<HTMLSpanElement>; // prompt input placeholder
	@ViewChild("promptInputScroll") promptInputScroll!: ElementRef<HTMLDivElement>; // prompt input scroll container

	/** Array with search values */
	searchArray = signal<Array<any>>([]);
	/** state for bottom bar content */
	searchState = signal<boolean>(false);
	/** chosen point */
	chosenPoint = signal<EventPlace | Place | null>(null);
	/** Search text for searching pipe */
	searchText: string = "";
	/** Bool param if last typed button to prompt input is valid */
	typedKey: string = "";
	/** Prompt text */
	promptText: string = "";

	leftIndex: number[] = [];
	rightIndex: number[] = [];
	chosenElements: Array<EventPlace | Place | null> = [];

	searchBoxIndex: number = -1;

	/** bottom bar state */
	public readonly bottomBarState = signal<boolean>(false);
	private bottomBarTimeouts: any[] = [];

	public readonly MIN_HEIGHT = 233;

	constructor(
		private router: Router,
		private readonly notificationsService: NotificationsService,
		public readonly mapInteractionService: MapInteractionsService,
	) {
		// effect for mapInteractionService.mainContainerState
		effect(() => {
			//if (this.mapInteractionService.mainContainerState() === -1) this.hideMainContainer();
			//else if (this.mapInteractionService.mainContainerState() === 1) this.showMainContainer();
		});
		// effect for mapInteractionService.mapScrolled
		effect(() => {
			if (this.mapInteractionService.mapScrolled() !== 0)
				this.mapInteractionService.mainContainerState.set(-1 * this.mapInteractionService.mapScrolled());
		});
		// effect for closing bottom bar event
		effect(() => {
			const state = this.bottomBarState();
			if (this.container === undefined) return;

			this.bottomBarTimeouts.forEach(t => clearTimeout(t));
			this.bottomBarTimeouts = [];

			if (state) {
				this.stickPositionButton.nativeElement.style.opacity = "0";
				for (let i = 0; i < this.openMenuItems.length; i++)
					this.openMenuItems.get(i)!.nativeElement.style.display = "flex";
				this.bottomBarTimeouts.push(
					setTimeout(() => {
						this.bottomBarTimeouts = this.bottomBarTimeouts.splice(0,1);
						for (let i = 0; i < this.openMenuItems.length; i++)
							this.openMenuItems.get(i)!.nativeElement.style.opacity = "1";
					}),
				);
				this.bottomBarTimeouts.push(
					setTimeout(() => {
						this.bottomBarTimeouts = this.bottomBarTimeouts.splice(0,1);
						this.stickPositionButton.nativeElement.style.display = "none";
					}, 300),
				);
			}
			else {
				this.stickPositionButton.nativeElement.style.display = "flex";
				for (let i = 0; i < this.openMenuItems.length; i++)
					this.openMenuItems.get(i)!.nativeElement.style.opacity = "0";
				this.bottomBarTimeouts.push(
					setTimeout(() => {
						this.bottomBarTimeouts = this.bottomBarTimeouts.splice(0,1);
						this.stickPositionButton.nativeElement.style.opacity = "1";
					}),
				);
				this.bottomBarTimeouts.push(
					setTimeout(() => {
						this.bottomBarTimeouts = this.bottomBarTimeouts.splice(0,1);
						for (let i = 0; i < this.openMenuItems.length; i++)
							this.openMenuItems.get(i)!.nativeElement.style.display = "none";
					}, 300),
				);
			}
		});
		// effect for chosen point
		effect(() => {
			const point = this.chosenPoint();
			if (point !== null) {
				this.chosenPoint.set(null);
				this.handlePromptSearchBoxClick(point);
				this.bottomBarState.set(false);
			}
		});
		// Listener for path information container opening
		this.mapInteractionService.pathInformationState.subscribe(state => {
			if (state === 1) {
				this.closeMainContainer();
			}
			else if (state === -2 && this.mapInteractionService.chosenMapPoint.value === null) {
				this.openMainContainer();
			}
		});
	}

	ngAfterViewInit() {
		/** setting up router listener */
		this.router.events.pipe(
			filter((event: any) => event instanceof NavigationEnd),
		).subscribe(() => {
			this.routeChangeHandler();
		});
	}

	closeMainContainer() {
		this.container.nativeElement.style.transitionDuration = ".3s";
		setTimeout(() => {
			this.container.nativeElement.style.transform = "translateY(100%)";
		}, 100);
		setTimeout(() => {
			this.container.nativeElement.style.display = "none";
			this.container.nativeElement.style.transitionDuration = "0s";
		}, 400);
	}

	openMainContainer() {
		this.container.nativeElement.style.transitionDuration = ".3s";
		this.container.nativeElement.style.display = "flex";
		setTimeout(() => {
			this.container.nativeElement.style.transform = "translateY(0%)";
		}, 100);
		setTimeout(() => {
			this.container.nativeElement.style.transitionDuration = "0s";
			this.stickPositionButton.nativeElement.style.top = `${ window.innerHeight - this.MIN_HEIGHT - 43 }px`;
		}, 400);
	}

	stickUserPosition() {
		this.mapInteractionService.chosenMapPoint.next(null);
		this.mapInteractionService.mapScrolled.set(-1)
	}

	/** listens to route change */
	routeChangeHandler() {
		/* Showing map addon if we are in the home page */
		if (this.router.url.split("/").length === 1 || this.router.url.split("/")[1] === "") {
			this.mainAddonContainer.nativeElement.style.display = "flex";
			setTimeout(() => this.mainAddonContainer.nativeElement.style.maxHeight = "100vh");
			setTimeout(() => this.stickPositionButton.nativeElement.style.opacity = "1", 300);
		}
		else {
			this.mainAddonContainer.nativeElement.style.maxHeight = "0";
			this.stickPositionButton.nativeElement.style.opacity = "0";
			setTimeout(() => this.mainAddonContainer.nativeElement.style.display = "none", 300);
		}
	}

	/** shows the search box */
	showSearchBox() {
		this.searchState.set(true);
	}

	/** hides search box */
	hideSearchBox() {
		this.searchState.set(false);
	}

	/** shows prompt placeholder */
	showPromptPlaceholder() {
		this.promptPlaceholder.nativeElement.style.opacity = "1";
		this.promptPlaceholder.nativeElement.style.transform = "translateX(0)";
	}

	/** hides prompt placeholder */
	hidePromptPlaceholder() {
		this.promptPlaceholder.nativeElement.style.opacity = "0";
		this.promptPlaceholder.nativeElement.style.transform = "translateX(50px)";
	}

	/** changes editor input size after typing */
	changeInputSize() {
		// check if we should add new line on the div and add to scroll-container size
		const current_input_height: number = this.promptInput.nativeElement.getBoundingClientRect().height!;
		this.promptInputScroll.nativeElement.style.height = `clamp(var(--base-height), ${ current_input_height }px, var(--base-height-max))`;
	}

	/** Handler for prompt search box click event */
	handlePromptSearchBoxClick(target: Place | EventPlace) {
		const editor = this.promptInput.nativeElement;
		let editorText = editor.innerText;
		let index = this.searchBoxIndex;
		let newCaretPosition;

		const currentValue = editorText;
		editorText = currentValue.slice(0, this.leftIndex[index] + 1);
		editorText += target.name;
		if (currentValue.slice(this.rightIndex[index] + 1).length === 0 || currentValue.slice(this.rightIndex[index] + 1)[0] !== " ")
			editorText += " ";
		editorText += currentValue.slice(this.rightIndex[index] + 1);

		this.chosenElements[index] = target;
		const length = target.name?.length!;
		this.rightIndex[index] = this.leftIndex[index] + length;
		newCaretPosition = this.rightIndex[index] + 2;


		// highlighting the segments in the line with spans
		const highlightedString = this.highlightEditorText(editorText);

		// updating prompt text and editor
		this.promptText = editorText;
		editor.innerHTML = highlightedString;
		placeCaretAt(editor, newCaretPosition);

		// hiding search box and updating input size
		this.hideSearchBox();
		this.changeInputSize();

		// if we are not choosing Embedding we should show it on the map
		if (isInstanceOfPlace(target) || isInstanceOfEventPlace(target)) {
			const selectedPoint: MapClickEntityModel = {
				id: "1",
				geometry: {
					type: "Point",
					coordinates: [target.latitude, target.longitude],
				},
				properties: {
					name: target.name,
				},
			};
			this.mapInteractionService.selectedPointOnPromptInput.next(selectedPoint);
			this.mapInteractionService.mapScrolled.set(1);
		}
	}

	/** Handler for prompt input click */
	handleClickPromptInput() {
		//this.bottomBarState.set(true);
	}

	/** Handler for prompt input key button down event */
	handleKeyDownPromptInput(event: any) {
		this.typedKey = event.key;
		//this.bottomBarState.set(true);
	}

	/** Handler for caret position change */
	handleCaretPositionChanges(event: any) {
		if (event.key === undefined || ( event.key.slice(0, 5) === "Arrow" )) {
			const editor = this.promptInput.nativeElement;
			const editorText = editor.innerText;
			const caretPosition = getCaretPosition(editor);
			let index = -1;
			for (let i = 0; i < this.leftIndex.length; i++)
				if (this.leftIndex[i] + 1 <= caretPosition && caretPosition <= this.rightIndex[i] + 1)
					index = i;

			if (index !== -1 && this.chosenElements[index] === null) {
				if (editorText.charAt(this.leftIndex[index]) === "#")
					this.searchArray = this.mapInteractionService.places;
				else
					this.searchArray = this.mapInteractionService.places; // change to event places
				this.searchText = editorText.slice(this.leftIndex[index] + 1, this.rightIndex[index] + 1);
				this.searchBoxIndex = index;
				this.showSearchBox();
			}
			else
				this.hideSearchBox();
		}
	}

	/** Handler for prompt input event */
	handlePromptInput() {
		const editor = this.promptInput.nativeElement;
		const caretPosition = getCaretPosition(editor);
		let newCaretPosition = caretPosition;
		let editorText = editor.innerText;

		if (this.typedKey === "Backspace" || this.typedKey === "Delete") {
			// we are looking for a segment that will be affected by the current change
			let index = -1;
			for (let i = 0; i < this.leftIndex.length; i++)
				if (this.leftIndex[i] + ( this.typedKey === "Backspace" ? 0 : 0 ) <= caretPosition && caretPosition <= this.rightIndex[i])
					index = i;

			// removing this segment
			if (index !== -1) {
				const currentValue = editorText;
				editorText = currentValue.slice(0, this.leftIndex[index]);
				editorText += currentValue.slice(this.rightIndex[index]);
				this.hideSearchBox();
				newCaretPosition = this.leftIndex[index];
				this.leftIndex.splice(index, 1);
				this.rightIndex.splice(index, 1);
				this.chosenElements.splice(index, 1);
			}
		}
		else if (this.typedKey === "Enter") {
			editorText = this.promptText;

			// checking that none of chosenElements is invalid
			let isValid = true;
			for (const chosenElement of this.chosenElements)
				if (chosenElement === null) {
					this.notificationsService.addNotification({
						message: "выбрано некорректное место или событие",
						type: "error",
						timeOut: 10 * 1000,
					});
					isValid = false;
				}

			if (isValid) {
				let promptText: string = "";
				let lastIndex = 0;

				// adding segments
				for (let index = 0; index < this.leftIndex.length; index++) {
					// adding normal string segment
					promptText += editorText.slice(lastIndex, this.leftIndex[index]);

					// adding highlighted string
					const item = this.chosenElements[index];
					if (isInstanceOfPlace(item))
						promptText += `<#fixed:${ item._id }:${ item.name }>`;
					else if (isInstanceOfEventPlace(item))
						promptText += `<@event:${ item._id }>`;
					else
						promptText += `<#embedding:${ item }>`;

					// updating last placed index
					lastIndex = this.rightIndex[index] + 1;
				}
				promptText += editorText.slice(lastIndex);

				this.mapInteractionService.generatePath(promptText);
			}
		}
		else if (this.typedKey === "#" || this.typedKey === "@") {
			// we are looking for a segment that will be affected by the current change
			this.bottomBarState.set(true);
			let index = -1;
			for (let i = 0; i < this.leftIndex.length; i++)
				if (this.leftIndex[i] < caretPosition && caretPosition <= this.rightIndex[i] + 2)
					index = i;

			if (index !== -1) {
				editorText = this.promptText;
				newCaretPosition = caretPosition - 1;
			}
			else {

				this.leftIndex.push(caretPosition - 1);
				this.rightIndex.push(caretPosition - 1);
				this.chosenElements.push(null);

				let lastIndex = this.leftIndex.length - 1;
				while (lastIndex > 0 && this.leftIndex[lastIndex] < this.leftIndex[lastIndex - 1]) {
					swapElements(this.leftIndex, lastIndex, lastIndex - 1);
					swapElements(this.rightIndex, lastIndex, lastIndex - 1);
					swapElements(this.chosenElements, lastIndex, lastIndex - 1);
					lastIndex--;
				}

				if (editorText.charAt(caretPosition - 1) === "#")
					this.searchArray = this.mapInteractionService.places;
				else
					this.searchArray = this.mapInteractionService.places; // TODO change to event places

				this.searchText = "";
				this.searchBoxIndex = this.leftIndex.length - 1;
				this.showSearchBox();
			}
		}
		else {
			// we are looking for a segment that will be affected by the current change
			let index = -1;
			for (let i = 0; i < this.leftIndex.length; i++)
				if (this.leftIndex[i] < caretPosition && caretPosition <= this.rightIndex[i] + 2)
					index = i;

			// updating this segment
			if (index !== -1) {
				this.rightIndex[index]++;
				this.searchText = editorText.slice(this.leftIndex[index], this.rightIndex[index] + 1);
			}
		}

		// updating position of left and right indexes
		let lastIndex = 0;
		for (let index = 0; index < editorText.length; index++) {
			if (editorText.charAt(index) === "#" || editorText.charAt(index) === "@") {
				const delta = index - this.leftIndex[lastIndex];
				this.leftIndex[lastIndex] += delta;
				this.rightIndex[lastIndex] += delta;
				lastIndex++;
			}
		}

		// removing all invalid indexes
		for (let index = 0; index < this.leftIndex.length; index++) {
			if (editorText.length <= this.leftIndex[index] || ( editorText.charAt(this.leftIndex[index]) !== "#" && editorText.charAt(this.leftIndex[index]) !== "@" )) {
				this.leftIndex.splice(index, 1);
				this.rightIndex.splice(index, 1);
				this.chosenElements.splice(index, 1);
				index--;
			}
		}

		// highlighting the segments in the line with spans
		const highlightedString = this.highlightEditorText(editorText);

		// showing placeholder
		if (editorText === "\n" || editorText === "")
			this.showPromptPlaceholder();
		else
			this.hidePromptPlaceholder();

		this.promptText = editorText;
		editor.innerHTML = highlightedString;
		placeCaretAt(editor, newCaretPosition);

		this.changeInputSize();
	}

	/** highlights EVENTS and PLACES in editor text */
	highlightEditorText(editorText: string) {
		let highlightedString: string = "";
		let lastIndex = 0;
		// highlighting segments
		for (let index = 0; index < this.leftIndex.length; index++) {
			// adding normal string segment
			highlightedString += editorText.slice(lastIndex, this.leftIndex[index]);

			// adding highlighted string
			const hashStyles = "border-radius: var(--br-100); background: rgba(67, 70, 218, 0.5); color: var(--text-primary); padding: 0 10px;";
			const atStyles = "border-radius: var(--br-100); background: rgba(172, 0, 230, 0.5); color: var(--text-primary); padding: 0 10px;";
			if (editorText.charAt(this.leftIndex[index]) === "#")
				highlightedString += `<span style="${ hashStyles }">${ editorText.slice(this.leftIndex[index], this.rightIndex[index] + 1) }</span>`;
			else
				highlightedString += `<span style="${ atStyles }">${ editorText.slice(this.leftIndex[index], this.rightIndex[index] + 1) }</span>`;

			// updating last placed index
			lastIndex = this.rightIndex[index] + 1;
		}
		highlightedString += editorText.slice(lastIndex);
		// returns highlighted result
		return highlightedString;
	}

	protected readonly window = window;
}
