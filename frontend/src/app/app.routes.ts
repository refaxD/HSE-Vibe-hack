import { Routes } from "@angular/router";
import { HomeGuard } from "./core/guards/home.guard";
import { IntroGuard } from "./core/guards/intro.guard";

export const routes: Routes = [
	{
		path: "",
		loadChildren: () => import("./features/home/home.routes"),
		//canActivate: [HomeGuard],
	},
	{
		path: "**",
		pathMatch: "full",
		redirectTo: "",
	},
];
