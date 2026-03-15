import { Routes } from "@angular/router";

const routes: Routes = [
    {
        path: "",
        children: [
            { path: '', loadComponent: () => import('./pages/main/main.component'), pathMatch: "full" },
        ]
    }
];

export default routes;
