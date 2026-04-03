// @vitest-environment jsdom
import type { BloodInventory } from "@/types";
import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it } from "vitest";
import InventoryAlert from "./InventoryAlert";

afterEach(() => {
	cleanup();
});

const createInventory = (bloodType: string, level: BloodInventory["level"]): BloodInventory => ({
	id: `inv-${bloodType}`,
	bloodType: bloodType as BloodInventory["bloodType"],
	level,
	updatedAt: "2026-04-04T00:00:00Z",
});

describe("InventoryAlert", () => {
	it("renders nothing when all levels are normal or sufficient", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "normal"),
			createInventory("B+", "sufficient"),
			createInventory("O+", "normal"),
		];

		const { container } = render(<InventoryAlert inventories={inventories} />);
		expect(container.innerHTML).toBe("");
	});

	it("renders nothing for empty inventories", () => {
		const { container } = render(<InventoryAlert inventories={[]} />);
		expect(container.innerHTML).toBe("");
	});

	it("shows critical alert when critical items exist", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "critical"),
			createInventory("B+", "normal"),
			createInventory("O+", "low"),
		];

		render(<InventoryAlert inventories={inventories} />);

		const alert = screen.getByTestId("inventory-alert");
		expect(alert).toBeDefined();
		expect(screen.getByText("血液在庫が逼迫しています")).toBeDefined();
		expect(screen.getByText("逼迫: A+")).toBeDefined();
		expect(screen.getByText("不足: O+")).toBeDefined();
	});

	it("shows low alert when only low items exist (no critical)", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "low"),
			createInventory("B+", "normal"),
			createInventory("O+", "low"),
		];

		render(<InventoryAlert inventories={inventories} />);

		const alert = screen.getByTestId("inventory-alert");
		expect(alert).toBeDefined();
		expect(screen.getByText("血液在庫が不足しています")).toBeDefined();
		expect(screen.getByText("不足: A+, O+")).toBeDefined();
	});

	it("shows multiple critical blood types", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "critical"),
			createInventory("B+", "critical"),
			createInventory("O+", "normal"),
		];

		render(<InventoryAlert inventories={inventories} />);

		expect(screen.getByText("逼迫: A+, B+")).toBeDefined();
	});

	it("has alert role for accessibility", () => {
		const inventories: BloodInventory[] = [createInventory("A+", "critical")];

		render(<InventoryAlert inventories={inventories} />);

		const alert = screen.getByRole("alert");
		expect(alert).toBeDefined();
	});
});
