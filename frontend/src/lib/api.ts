const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

class ApiError extends Error {
	constructor(
		public status: number,
		message: string,
	) {
		super(message);
		this.name = "ApiError";
	}
}

// M-2: accept optional runtime validator to avoid unsafe `as T` casts
async function request<T>(path: string, options?: RequestInit, validate?: (data: unknown) => T): Promise<T> {
	const url = `${API_BASE}${path}`;
	const res = await fetch(url, {
		headers: { "Content-Type": "application/json" },
		...options,
	});

	if (!res.ok) {
		const body = await res.text();
		throw new ApiError(res.status, body);
	}

	const json: unknown = await res.json();
	if (validate) {
		return validate(json);
	}
	// Fallback: callers that don't provide a validator accept the risk
	return json as T;
}

export const api = {
	get: <T>(path: string, validate?: (data: unknown) => T) => request<T>(path, undefined, validate),
	post: <T>(path: string, body: unknown, validate?: (data: unknown) => T) =>
		request<T>(path, { method: "POST", body: JSON.stringify(body) }, validate),
	put: <T>(path: string, body: unknown, validate?: (data: unknown) => T) =>
		request<T>(path, { method: "PUT", body: JSON.stringify(body) }, validate),
	delete: (path: string) => request<void>(path, { method: "DELETE" }),
};

export { ApiError };
