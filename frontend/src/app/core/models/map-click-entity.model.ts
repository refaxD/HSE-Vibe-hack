export interface MapClickEntityModel {
    geometry: MapClickEntityGeometryModel,
    id: string;
    properties: MapClickEntityPropertiesModel;
	_props?: any;
	element?: any;
}

export interface MapClickEntityGeometryModel {
    type: string;
    coordinates: number[];
}

export interface MapClickEntityPropertiesModel {
    name: string | undefined;
}
