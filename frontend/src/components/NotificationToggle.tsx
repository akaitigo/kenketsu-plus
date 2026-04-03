"use client";

import { useEffect, useState } from "react";

interface NotificationToggleProps {
	vapidPublicKey: string;
}

export function NotificationToggle({ vapidPublicKey }: NotificationToggleProps) {
	const [supported, setSupported] = useState(false);
	const [permission, setPermission] = useState<NotificationPermission>("default");
	const [subscriptionId, setSubscriptionId] = useState<string | null>(null);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		setSupported("serviceWorker" in navigator && "PushManager" in window);
		if ("Notification" in window) {
			setPermission(Notification.permission);
		}
	}, []);

	const handleToggle = async () => {
		setLoading(true);
		try {
			if (subscriptionId) {
				const { unsubscribeFromPush } = await import("@/lib/push");
				await unsubscribeFromPush(subscriptionId);
				setSubscriptionId(null);
			} else {
				const { subscribeToPush } = await import("@/lib/push");
				const sub = await subscribeToPush(vapidPublicKey);
				if (sub) {
					setSubscriptionId(sub.id);
					setPermission(Notification.permission);
				}
			}
		} catch {
			// Notification permission denied or API error
		} finally {
			setLoading(false);
		}
	};

	if (!supported) {
		return <p data-testid="notification-unsupported">このブラウザはプッシュ通知に対応していません</p>;
	}

	if (permission === "denied") {
		return <p data-testid="notification-denied">通知がブロックされています。ブラウザの設定から許可してください。</p>;
	}

	return (
		<button
			type="button"
			onClick={handleToggle}
			disabled={loading}
			data-testid="notification-toggle"
			style={{
				padding: "8px 16px",
				borderRadius: 6,
				border: "none",
				backgroundColor: subscriptionId ? "#ef4444" : "#3b82f6",
				color: "white",
				cursor: loading ? "wait" : "pointer",
			}}
		>
			{loading ? "処理中..." : subscriptionId ? "通知をOFF" : "通知をON"}
		</button>
	);
}
