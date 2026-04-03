import { describe, it, expect, vi, beforeEach } from "vitest";
import { api, ApiError } from "./api";

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
	mockFetch.mockReset();
});

describe("api", () => {
	it("GET request returns parsed JSON", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ status: "ok" }),
		});

		const result = await api.get<{ status: string }>("/health");
		expect(result).toEqual({ status: "ok" });
		expect(mockFetch).toHaveBeenCalledWith(
			"http://localhost:8080/health",
			expect.objectContaining({ headers: { "Content-Type": "application/json" } }),
		);
	});

	it("POST request sends body as JSON", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ id: "1" }),
		});

		const result = await api.post<{ id: string }>("/api/centers", { name: "test" });
		expect(result).toEqual({ id: "1" });
		expect(mockFetch).toHaveBeenCalledWith(
			"http://localhost:8080/api/centers",
			expect.objectContaining({
				method: "POST",
				body: '{"name":"test"}',
			}),
		);
	});

	it("throws ApiError on non-ok response", async () => {
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 404,
			text: () => Promise.resolve("not found"),
		});

		await expect(api.get("/api/unknown")).rejects.toThrow(ApiError);
	});
});
