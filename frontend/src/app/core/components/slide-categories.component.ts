import { Component, Input } from "@angular/core";
import { DomSanitizer } from "@angular/platform-browser";
import { Category } from "../../../generated";
import { MapInteractionsService } from "../services/map-interactions.service";

@Component({
    selector: 'CategoriesSlider',
    standalone: true,
    imports: [],
    styles: `
        .container {
            display:    block;
            width:      100%;
            overflow-x: scroll;
        }

        .content {
            display:        flex;
            flex-direction: row;
            align-items:    center;
            gap:            20px;
        }

        .card {
            cursor:              pointer;
            width:               max-content;
            padding:             5px 10px;
            display:             flex;
            flex-direction:      row;
            gap:                 5px;
            align-items:         center;
            border-radius:       var(--br-100);
            border:              1px solid var(--text-primary);
            color:               var(--text-primary);
            filter:              brightness(.7);
            transition-duration: .3s;
        }

        .active {
            filter:       brightness(1);
            background:   var(--blue-60);
            border-color: var(--blue-60);
        }
    `,
    template: `
        <div class="container">
            <div class="content"
                 [style.width]="flexWrap ? '100%' : 'fit-content'"
                 [style.flex-wrap]="flexWrap ? 'wrap':'nowrap'"
            >
				@for (category of mapInteractions.categories(); track mapInteractions.categories()) {
                    <div (click)="categoryClickHandler(category)"
                         [class.active]="mapInteractions.checkCategoryChosen(category)"
                         class="card"
                    >
						<div [innerHTML]="sanitizer.bypassSecurityTrustHtml(category.svg)"></div>
                        <p>{{ category.name }}</p>
                    </div>
                }
            </div>
        </div>
    `
})
export class SlideCategoriesComponent {
	/** variable that is responsible for content to be flex wrapped
	 * `true` - we wrap the content
	 * `false` - content is inline */
    @Input({ alias: 'flexWrap' }) flexWrap: boolean = false;

    constructor(
        public mapInteractions: MapInteractionsService,
		public sanitizer: DomSanitizer,
    ) {}

	/**
	 * Function that changes chosenCategory in `map-integration service`
	 * @param {Category} newCategory category to change state
	 */
	categoryClickHandler(newCategory: Category) {
        let currentValue = this.mapInteractions.chosenCategories();
        if (this.mapInteractions.checkCategoryChosen(newCategory))
            currentValue.splice(currentValue.indexOf(newCategory), 1);
        else
            currentValue.push(newCategory);
        this.mapInteractions.chosenCategories.set(currentValue);
    }
}
