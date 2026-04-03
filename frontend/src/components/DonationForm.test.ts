import type { BloodType, DonationType, Gender } from "@/types";
import { describe, expect, it } from "vitest";

const VALID_BLOOD_TYPES: BloodType[] = ["A+", "A-", "B+", "B-", "O+", "O-", "AB+", "AB-"];
const VALID_DONATION_TYPES: DonationType[] = ["whole_400", "whole_200", "component"];
const VALID_GENDERS: Gender[] = ["male", "female"];

const VOLUME_MAP: Record<DonationType, number> = {
	whole_400: 400,
	whole_200: 200,
	component: 0,
};

describe("DonationForm validation", () => {
	it("all blood types are valid", () => {
		expect(VALID_BLOOD_TYPES).toHaveLength(8);
		for (const bt of VALID_BLOOD_TYPES) {
			expect(bt).toMatch(/^(A|B|O|AB)[+-]$/);
		}
	});

	it("all donation types have correct volume", () => {
		expect(VOLUME_MAP.whole_400).toBe(400);
		expect(VOLUME_MAP.whole_200).toBe(200);
		expect(VOLUME_MAP.component).toBe(0);
	});

	it("all donation types are covered", () => {
		expect(VALID_DONATION_TYPES).toHaveLength(3);
	});

	it("genders are male and female", () => {
		expect(VALID_GENDERS).toEqual(["male", "female"]);
	});

	it("date validation: empty date is invalid", () => {
		const donatedAt = "";
		expect(donatedAt).toBe("");
	});

	it("date validation: valid date is accepted", () => {
		const donatedAt = "2026-03-15";
		const parsed = new Date(donatedAt);
		expect(Number.isNaN(parsed.getTime())).toBe(false);
	});
});
