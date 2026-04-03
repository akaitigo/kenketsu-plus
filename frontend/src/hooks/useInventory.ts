"use client";

import { api } from "@/lib/api";
import type { BloodInventory } from "@/types";
import { useEffect, useState } from "react";

interface UseInventoryResult {
	inventories: BloodInventory[];
	error: string | null;
	loading: boolean;
}

function extractErrorMessage(err: unknown): string {
	if (err instanceof Error) {
		return err.message;
	}
	return "データの取得に失敗しました";
}

async function fetchInventoryData(): Promise<BloodInventory[]> {
	const data = await api.get<BloodInventory[]>("/api/inventory");
	return data ?? [];
}

export function useInventory(): UseInventoryResult {
	const [inventories, setInventories] = useState<BloodInventory[]>([]);
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		let cancelled = false;

		fetchInventoryData()
			.then((data) => {
				if (!cancelled) {
					setInventories(data);
					setError(null);
				}
			})
			.catch((err: unknown) => {
				if (!cancelled) {
					setError(extractErrorMessage(err));
				}
			})
			.finally(() => {
				if (!cancelled) {
					setLoading(false);
				}
			});

		return () => {
			cancelled = true;
		};
	}, []);

	return { inventories, error, loading };
}
