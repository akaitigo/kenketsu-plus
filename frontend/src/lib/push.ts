import type { PushSubscription as PushSub } from "@/types";
import { api } from "./api";

function urlBase64ToUint8Array(base64String: string): Uint8Array {
	const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
	const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
	const rawData = atob(base64);
	const outputArray = new Uint8Array(rawData.length);
	for (let i = 0; i < rawData.length; ++i) {
		outputArray[i] = rawData.charCodeAt(i);
	}
	return outputArray;
}

export async function registerServiceWorker(): Promise<ServiceWorkerRegistration | null> {
	if (!("serviceWorker" in navigator)) {
		return null;
	}
	return navigator.serviceWorker.register("/sw.js");
}

export async function subscribeToPush(vapidPublicKey: string): Promise<PushSub | null> {
	const registration = await registerServiceWorker();
	if (!registration) return null;

	const keyArray = urlBase64ToUint8Array(vapidPublicKey);
	const subscription = await registration.pushManager.subscribe({
		userVisibleOnly: true,
		applicationServerKey: keyArray.buffer as ArrayBuffer,
	});

	const json = subscription.toJSON();
	const result = await api.post<PushSub>("/api/subscriptions", {
		endpoint: json.endpoint,
		p256dh: json.keys?.p256dh ?? "",
		auth: json.keys?.auth ?? "",
	});

	return result;
}

export async function unsubscribeFromPush(subscriptionId: string): Promise<void> {
	await api.delete(`/api/subscriptions/${subscriptionId}`);

	const registration = await navigator.serviceWorker.ready;
	const subscription = await registration.pushManager.getSubscription();
	if (subscription) {
		await subscription.unsubscribe();
	}
}
