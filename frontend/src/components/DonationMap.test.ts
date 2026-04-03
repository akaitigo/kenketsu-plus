import type { DonationCenter } from "@/types";
import { beforeEach, describe, expect, it, vi } from "vitest";

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
	mockFetch.mockReset();
});

const sampleCenter: DonationCenter = {
	id: "center-1",
	name: "Shibuya Center",
	address: "Shibuya, Tokyo",
	lat: 35.658,
	lng: 139.7016,
	capacity: 50,
	availableSlots: 10,
	status: "open",
	createdAt: "2026-01-01T00:00:00Z",
	updatedAt: "2026-01-01T00:00:00Z",
};

describe("DonationMap data fetching", () => {
	it("fetches centers from the API", async () => {
		const centers: DonationCenter[] = [sampleCenter];

		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve(centers),
		});

		const { api } = await import("@/lib/api");
		const result = await api.get<DonationCenter[]>("/api/centers");

		expect(result).toEqual(centers);
		expect(mockFetch).toHaveBeenCalledWith(
			"http://localhost:8080/api/centers",
			expect.objectContaining({
				headers: { "Content-Type": "application/json" },
			}),
		);
	});

	it("handles empty center list", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve([]),
		});

		const { api } = await import("@/lib/api");
		const result = await api.get<DonationCenter[]>("/api/centers");

		expect(result).toEqual([]);
	});

	it("handles API error", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 500,
			text: () => Promise.resolve("Internal Server Error"),
		});

		const { api, ApiError } = await import("@/lib/api");

		await expect(api.get<DonationCenter[]>("/api/centers")).rejects.toThrow(ApiError);
	});

	it("DonationCenter type has required fields", () => {
		expect(sampleCenter.id).toBe("center-1");
		expect(sampleCenter.name).toBe("Shibuya Center");
		expect(sampleCenter.lat).toBe(35.658);
		expect(sampleCenter.lng).toBe(139.7016);
		expect(sampleCenter.status).toBe("open");
		expect(sampleCenter.availableSlots).toBe(10);
		expect(sampleCenter.capacity).toBe(50);
	});

	it("fetches centers with distance filter", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve([sampleCenter]),
		});

		const { api } = await import("@/lib/api");
		const result = await api.get<DonationCenter[]>("/api/centers?lat=35.66&lng=139.70&radius=5");

		expect(result).toHaveLength(1);
		expect(mockFetch).toHaveBeenCalledWith(
			"http://localhost:8080/api/centers?lat=35.66&lng=139.70&radius=5",
			expect.objectContaining({
				headers: { "Content-Type": "application/json" },
			}),
		);
	});
});
