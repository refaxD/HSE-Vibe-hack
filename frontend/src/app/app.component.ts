import { Component, effect, ElementRef, signal, ViewChild } from "@angular/core";
import { RouterOutlet } from "@angular/router";
import { BottomBarComponent } from "./core/components/bottom-bar/bottom-bar.component";
import MainComponent from "./features/intro/pages/main/main.component";

@Component({
    selector: 'app-root',
    standalone: true,
	imports: [
		BottomBarComponent,
		RouterOutlet,
		MainComponent,
	],
    styles: `
        .container {
            width:            100vw;
            height:           100vh;
            max-height:       100vh;
            display:          flex;
            position:         relative;
            background-color: var(--background-primary);
            overflow-y:       hidden;
            box-sizing:       border-box;
        }

        .router-container {
            width:      100%;
            height:     fit-content;
            min-height: 100%;
            transition: all .3s linear;
            padding:    0;
            overflow-y: hidden;

            .router-content {
                border-radius: var(--br-16);
                overflow:      hidden;
            }
        }
	`,
    template: `
        <div class="container" #container>
             <intro-main/> 
            <div class="router-container" #router>
                <div class="router-content" #router>
                    <router-outlet/>
                </div>
            </div>
            <BottomBar/>
        </div>
	`
})
export class AppComponent {
	@ViewChild('router') router!: ElementRef;

	static bottomBarState = signal<{ percent: number, state: boolean }>({ percent: 0, state: false });

	private readonly sideOpenedPadding = 10;
	private readonly topOpenedPadding = 30;

	constructor() {
		effect(() => {
			const event = AppComponent.bottomBarState();
			this.handleBottomState(event);
		});
	}

	handleBottomState(event: { percent: number; state: boolean }) {
		if (this.router === undefined) return;
		if (event.state)
			this.router.nativeElement.style.transitionDuration = '.3s';
		else
			this.router.nativeElement.style.transitionDuration = '0s';
		setTimeout(() => {
			this.router.nativeElement.style.padding = `${this.topOpenedPadding * event.percent}px ${this.sideOpenedPadding * event.percent}px`;
		})
	}
}
