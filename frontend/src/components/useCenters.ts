"use client";

import { api } from "@/lib/api";
import type { DonationCenter } from "@/types";
import { useEffect, useState } from "react";

interface UseCentersResult {
	centers: DonationCenter[];
	error: string | null;
	loading: boolean;
}

function extractErrorMessage(err: unknown): string {
	if (err instanceof Error) {
		return err.message;
	}
	return "データの取得に失敗しました";
}

export function useCenters(): UseCentersResult {
	const [centers, setCenters] = useState<DonationCenter[]>([]);
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		let cancelled = false;

		async function fetchCenters() {
			try {
				const data = await api.get<DonationCenter[]>("/api/centers");
				if (!cancelled) {
					setCenters(data ?? []);
					setError(null);
				}
			} catch (err: unknown) {
				if (!cancelled) {
					setError(extractErrorMessage(err));
				}
			} finally {
				if (!cancelled) {
					setLoading(false);
				}
			}
		}

		void fetchCenters();

		return () => {
			cancelled = true;
		};
	}, []);

	return { centers, error, loading };
}
