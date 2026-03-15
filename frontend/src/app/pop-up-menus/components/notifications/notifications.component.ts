import { Component, ElementRef, ViewChild, ViewEncapsulation } from "@angular/core";
import { NotificationModel } from "../../../core/models/notification.model";
import { NotificationsService } from "../../../core/services/notifications.service";

@Component({
	selector: "Notifications",
	imports: [],
	styleUrl: "./notifications.component.css",
	encapsulation: ViewEncapsulation.None,
	template: `
		<div class="notifications-container" #container>
			<!-- 			<div class="notification"> -->
			<!-- 				<p>Я крутое уведомление</p> -->
			<!-- 				<div class="callback-action"> -->
			<!-- 					<p>Отмена</p> -->
			<!-- 				</div> -->
			<!-- 			</div> -->
		</div>
	`,
})
export class NotificationsComponent {
	@ViewChild("container") container!: ElementRef<HTMLDivElement>;

	constructor(
		private notificationService: NotificationsService,
	) {
		this.notificationService.queue.subscribe(queue => {
			if (queue.length === 0) return;
			this.addNotification(queue[queue.length - 1]);
			this.notificationService.popBack();
		});
	}

	addNotification(notification: NotificationModel) {
		let touchPosition: null | number[] = null;

		const containerTouchStart = (event: TouchEvent) => {
			container.style.opacity = ".5";
			touchPosition = [event.touches[0].clientX, event.touches[0].clientY];
			addEventListener("touchend", touchEnd);
			addEventListener("touchmove", touchMove);
		};
		const touchMove = (event: TouchEvent) => {
			touchPosition = [event.touches[0].clientX, event.touches[0].clientY];
		};
		const touchEnd = () => {
			container.style.opacity = "1";
			removeEventListener("touchend", touchEnd);
			removeEventListener("touchmove", touchMove);

			const bounds = container.getBoundingClientRect();
			if (touchPosition === null) return;
			if (bounds.top <= touchPosition[1] && touchPosition[1] <= bounds.bottom && bounds.left <= touchPosition[0] && touchPosition[0] <= bounds.right)
				this.removeNotification(notification);
		};

		const container = document.createElement("div");
		container.ontouchstart = containerTouchStart;

		if (notification.type === "error") {
			container.style.background = "rgb(255,0,0)";
		}

		const message = document.createElement("p");
		message.innerText = notification.message;
		container.appendChild(message);

		const cancel = document.createElement("p");
		cancel.innerText = "Отмена";
		cancel.onclick = this.clickHandler.bind({}, this, notification);

		if (notification.callback !== undefined)
			container.appendChild(cancel);

		const progressBar = document.createElement("div");
		container.appendChild(progressBar);

		container.setAttribute("notification-id", notification.id!);

		this.container.nativeElement.appendChild(container);

		setTimeout(() => {
			container.style.transform = "translateX(0%)";
			container.style.opacity = "1";
			progressBar.style.transitionDuration = `${ notification.timeOut }ms`;
			setTimeout(() => {
				progressBar.style.width = "100%";
			});
		});

		setTimeout(() => this.removeNotification(notification), notification.timeOut);
	}

	async removeNotification(notification: NotificationModel) {
		const container = this.container.nativeElement;
		for (let i = 0; i < container.children.length; i++) {
			const child = container.children.item(i)!;
			if (child.getAttribute("notification-id") === notification.id) {
				//@ts-ignore
				child.style.transform = "translateX(100%)";
				//@ts-ignore
				child.style.opacity = "0";
				setTimeout(() => {
					container.removeChild(child);
				}, 300);
			}
		}
	}

	clickHandler(that: any, notification: NotificationModel) {
		that.removeNotification(notification).then(() => {
			if (notification.callback !== undefined)
				notification.callback();
		});
	}
}
