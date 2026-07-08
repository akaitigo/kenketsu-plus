// @vitest-environment jsdom
import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it } from "vitest";
import { NotificationToggle } from "./NotificationToggle";

afterEach(() => {
	cleanup();
});

// stubPushSupport enables the push-capable code path. jsdom lacks these APIs and
// the component only checks for their presence (`in`), so an empty object suffices.
function stubPushSupport() {
	Object.defineProperty(navigator, "serviceWorker", { value: {}, configurable: true });
	Object.defineProperty(window, "PushManager", { value: {}, configurable: true });
}

describe("NotificationToggle", () => {
	// This case runs before stubPushSupport is applied, so jsdom's default (no
	// serviceWorker / PushManager) drives the unsupported branch.
	it("shows the unsupported message when push APIs are unavailable", () => {
		render(<NotificationToggle vapidPublicKey="test-key" />);

		expect(screen.getByTestId("notification-unsupported")).toBeDefined();
	});

	it("renders the toggle with an aria-label describing the enable action", () => {
		stubPushSupport();
		render(<NotificationToggle vapidPublicKey="test-key" />);

		const button = screen.getByTestId("notification-toggle");
		expect(button.getAttribute("aria-label")).toBe("在庫逼迫通知をオンにする");
		expect(button.getAttribute("aria-pressed")).toBe("false");
		expect(button.getAttribute("aria-busy")).toBe("false");
	});

	it("exposes the aria-label as the button's accessible name", () => {
		stubPushSupport();
		render(<NotificationToggle vapidPublicKey="test-key" />);

		expect(screen.getByRole("button", { name: "在庫逼迫通知をオンにする" })).toBeDefined();
	});
});
