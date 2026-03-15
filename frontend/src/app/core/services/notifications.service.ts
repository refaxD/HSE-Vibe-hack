import { Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { v4 as uuidv4 } from "uuid";
import { NotificationModel } from "../models/notification.model";

@Injectable({
	providedIn: "root",
})
export class NotificationsService {
	queue = new BehaviorSubject<Array<NotificationModel>>([]);

	constructor() { }

	addNotification(notification: NotificationModel) {
		notification.id = uuidv4();
		const queue = this.queue.value;
		queue.push(notification);
		this.queue.next(queue);
	}

	popBack() {
		const queue = this.queue.value;
		if (queue.length === 0) return;
		queue.pop();
		this.queue.next(queue);
	}
}
