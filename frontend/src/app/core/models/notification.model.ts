export interface NotificationModel {
	id?: string;
	type: "error" | "warning" | "info";
	message: string;
	timeOut: number;
	callback?: () => void;
}
