import { MapClickEntityModel } from "./map-click-entity.model";

export interface MapClickEventModel {
    entity: MapClickEntityModel;
    layer: string;
    source: string;
    type: string;
}