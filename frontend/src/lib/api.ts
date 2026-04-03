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

async function request<T>(path: string, options?: RequestInit): Promise<T> {
	const url = `${API_BASE}${path}`;
	const res = await fetch(url, {
		headers: { "Content-Type": "application/json" },
		...options,
	});

	if (!res.ok) {
		const body = await res.text();
		throw new ApiError(res.status, body);
	}

	return res.json() as Promise<T>;
}

export const api = {
	get: <T>(path: string) => request<T>(path),
	post: <T>(path: string, body: unknown) => request<T>(path, { method: "POST", body: JSON.stringify(body) }),
	put: <T>(path: string, body: unknown) => request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
	delete: (path: string) => request<void>(path, { method: "DELETE" }),
};

export { ApiError };
