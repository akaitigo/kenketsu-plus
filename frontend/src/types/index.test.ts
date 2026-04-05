import { describe, expect, it } from "vitest";
import { isBloodType, isDonationType, isGender, parseBloodType, parseDonationType, parseGender } from "./index";

describe("isBloodType", () => {
	it("returns true for valid blood types", () => {
		expect(isBloodType("A+")).toBe(true);
		expect(isBloodType("A-")).toBe(true);
		expect(isBloodType("B+")).toBe(true);
		expect(isBloodType("O-")).toBe(true);
		expect(isBloodType("AB+")).toBe(true);
		expect(isBloodType("AB-")).toBe(true);
	});

	it("returns false for invalid blood types", () => {
		expect(isBloodType("X+")).toBe(false);
		expect(isBloodType("")).toBe(false);
		expect(isBloodType("A")).toBe(false);
	});
});

describe("isDonationType", () => {
	it("returns true for valid donation types", () => {
		expect(isDonationType("whole_400")).toBe(true);
		expect(isDonationType("whole_200")).toBe(true);
		expect(isDonationType("component")).toBe(true);
	});

	it("returns false for invalid donation types", () => {
		expect(isDonationType("unknown")).toBe(false);
		expect(isDonationType("")).toBe(false);
	});
});

describe("isGender", () => {
	it("returns true for valid genders", () => {
		expect(isGender("male")).toBe(true);
		expect(isGender("female")).toBe(true);
	});

	it("returns false for invalid genders", () => {
		expect(isGender("other")).toBe(false);
		expect(isGender("")).toBe(false);
	});
});

describe("parseBloodType", () => {
	it("returns the value for valid blood types", () => {
		expect(parseBloodType("A+")).toBe("A+");
		expect(parseBloodType("O-")).toBe("O-");
	});

	it("throws for invalid blood types", () => {
		expect(() => parseBloodType("X+")).toThrow("Invalid blood type: X+");
	});
});

describe("parseDonationType", () => {
	it("returns the value for valid donation types", () => {
		expect(parseDonationType("whole_400")).toBe("whole_400");
	});

	it("throws for invalid donation types", () => {
		expect(() => parseDonationType("invalid")).toThrow("Invalid donation type: invalid");
	});
});

describe("parseGender", () => {
	it("returns the value for valid genders", () => {
		expect(parseGender("male")).toBe("male");
	});

	it("throws for invalid genders", () => {
		expect(() => parseGender("other")).toThrow("Invalid gender: other");
	});
});
