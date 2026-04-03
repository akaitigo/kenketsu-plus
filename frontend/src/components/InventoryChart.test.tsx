// @vitest-environment jsdom
import type { BloodInventory } from "@/types";
import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it } from "vitest";
import InventoryChart from "./InventoryChart";

afterEach(() => {
	cleanup();
});

const createInventory = (bloodType: string, level: BloodInventory["level"]): BloodInventory => ({
	id: `inv-${bloodType}`,
	bloodType: bloodType as BloodInventory["bloodType"],
	level,
	updatedAt: "2026-04-04T00:00:00Z",
});

describe("InventoryChart", () => {
	it("renders empty message when no inventories", () => {
		render(<InventoryChart inventories={[]} />);
		expect(screen.getByTestId("inventory-empty")).toBeDefined();
		expect(screen.getByText("在庫データがありません")).toBeDefined();
	});

	it("renders bars for all blood types", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "normal"),
			createInventory("A-", "critical"),
			createInventory("B+", "low"),
			createInventory("B-", "sufficient"),
			createInventory("O+", "normal"),
			createInventory("O-", "normal"),
			createInventory("AB+", "normal"),
			createInventory("AB-", "normal"),
		];

		render(<InventoryChart inventories={inventories} />);

		expect(screen.getByTestId("inventory-chart")).toBeDefined();

		for (const inv of inventories) {
			expect(screen.getByTestId(`bar-${inv.bloodType}`)).toBeDefined();
		}
	});

	it("applies correct colors based on inventory level", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "critical"),
			createInventory("B+", "low"),
			createInventory("O+", "normal"),
			createInventory("AB+", "sufficient"),
		];

		render(<InventoryChart inventories={inventories} />);

		const criticalBar = screen.getByTestId("bar-A+");
		expect(criticalBar.style.backgroundColor).toBe("rgb(220, 38, 38)");

		const lowBar = screen.getByTestId("bar-B+");
		expect(lowBar.style.backgroundColor).toBe("rgb(249, 115, 22)");

		const normalBar = screen.getByTestId("bar-O+");
		expect(normalBar.style.backgroundColor).toBe("rgb(34, 197, 94)");

		const sufficientBar = screen.getByTestId("bar-AB+");
		expect(sufficientBar.style.backgroundColor).toBe("rgb(59, 130, 246)");
	});

	it("displays level labels", () => {
		const inventories: BloodInventory[] = [
			createInventory("A+", "critical"),
			createInventory("B+", "low"),
			createInventory("O+", "normal"),
			createInventory("AB+", "sufficient"),
		];

		render(<InventoryChart inventories={inventories} />);

		expect(screen.getByText("逼迫")).toBeDefined();
		expect(screen.getByText("不足")).toBeDefined();
		expect(screen.getByText("通常")).toBeDefined();
		expect(screen.getByText("十分")).toBeDefined();
	});

	it("displays blood type labels", () => {
		const inventories: BloodInventory[] = [createInventory("A+", "normal"), createInventory("O-", "low")];

		render(<InventoryChart inventories={inventories} />);

		expect(screen.getByText("A+")).toBeDefined();
		expect(screen.getByText("O-")).toBeDefined();
	});
});
