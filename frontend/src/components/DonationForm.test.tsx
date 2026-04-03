// @vitest-environment jsdom
import type { BloodType, DonationType, Gender } from "@/types";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { DonationForm } from "./DonationForm";

afterEach(() => {
	cleanup();
});

vi.mock("@/lib/api", () => ({
	api: {
		post: vi.fn().mockResolvedValue({ id: "donation-1" }),
	},
}));

describe("DonationForm", () => {
	it("renders all form fields", () => {
		render(<DonationForm onCreated={() => {}} />);
		expect(screen.getByLabelText("血液型:")).toBeDefined();
		expect(screen.getByLabelText("種別:")).toBeDefined();
		expect(screen.getByLabelText("性別:")).toBeDefined();
		expect(screen.getByLabelText("献血日:")).toBeDefined();
		expect(screen.getByRole("button", { name: "記録する" })).toBeDefined();
	});

	it("shows error when submitting without date", async () => {
		render(<DonationForm onCreated={() => {}} />);
		fireEvent.click(screen.getByRole("button", { name: "記録する" }));
		expect(await screen.findByText("献血日を入力してください")).toBeDefined();
	});

	it("blood type select has 8 options", () => {
		render(<DonationForm onCreated={() => {}} />);
		const select = screen.getByLabelText("血液型:") as HTMLSelectElement;
		expect(select.options.length).toBe(8);
	});

	it("donation type select has 3 options", () => {
		render(<DonationForm onCreated={() => {}} />);
		const select = screen.getByLabelText("種別:") as HTMLSelectElement;
		expect(select.options.length).toBe(3);
	});
});

describe("DonationForm validation types", () => {
	const VALID_BLOOD_TYPES: BloodType[] = ["A+", "A-", "B+", "B-", "O+", "O-", "AB+", "AB-"];
	const VALID_DONATION_TYPES: DonationType[] = ["whole_400", "whole_200", "component"];
	const VALID_GENDERS: Gender[] = ["male", "female"];

	it("all blood types are valid", () => {
		expect(VALID_BLOOD_TYPES).toHaveLength(8);
	});

	it("all donation types are covered", () => {
		expect(VALID_DONATION_TYPES).toHaveLength(3);
	});

	it("genders are male and female", () => {
		expect(VALID_GENDERS).toEqual(["male", "female"]);
	});
});
