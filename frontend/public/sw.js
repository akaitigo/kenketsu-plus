self.addEventListener("push", (event) => {
	// M-1: guard against malformed push data
	let data = {};
	try {
		data = event.data ? event.data.json() : {};
	} catch (_) {
		data = { title: "Kenketsu-Plus", body: "通知があります" };
	}

	const title = data.title || "Kenketsu-Plus";
	const options = {
		body: data.body || "通知があります",
		icon: "/icon-192.png",
		badge: "/badge-72.png",
		data: {
			url: data.url || "/",
		},
	};

	event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener("notificationclick", (event) => {
	event.notification.close();
	const url = event.notification.data?.url || "/";
	event.waitUntil(clients.openWindow(url));
});
