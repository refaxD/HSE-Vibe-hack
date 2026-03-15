import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from "@angular/router";
import { MapInteractionsService } from "../services/map-interactions.service";

@Injectable({
	providedIn: "root",
})
export class IntroGuard implements CanActivate {

	constructor(
		private mapInteractionsService: MapInteractionsService,
		private router: Router,
	) {}

	canActivate(
		route: ActivatedRouteSnapshot,
		state: RouterStateSnapshot,
	): boolean {
		const userLocation = this.mapInteractionsService.userLocation();
		if (userLocation !== null) {
			this.router.navigate(["home"]);
			return false;
		}
		return true;
	}

}
