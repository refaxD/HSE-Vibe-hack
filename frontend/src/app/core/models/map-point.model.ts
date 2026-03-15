import { Category } from "../../../generated";

export default interface MapPointModel {
    name: string;
	categories: Category[];
	latitude: number;
	longitude: number;
    isLiked: boolean;
}
