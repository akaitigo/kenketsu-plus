"use client";

import InventoryAlert from "@/components/InventoryAlert";
import InventoryChart from "@/components/InventoryChart";
import { useInventory } from "@/hooks/useInventory";

export default function InventoryPage() {
	const { inventories, error, loading } = useInventory();

	return (
		<main style={{ maxWidth: "800px", margin: "0 auto", padding: "24px 16px" }}>
			<h1 style={{ fontSize: "24px", fontWeight: 700, marginBottom: "24px" }}>血液型別在庫ダッシュボード</h1>

			{loading && <p data-testid="loading">読み込み中...</p>}

			{error && (
				<p data-testid="error" style={{ color: "#dc2626" }}>
					{error}
				</p>
			)}

			{!loading && !error && (
				<>
					<InventoryAlert inventories={inventories} />
					<InventoryChart inventories={inventories} />
				</>
			)}
		</main>
	);
}
