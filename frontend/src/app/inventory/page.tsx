"use client";

import InventoryAlert from "@/components/InventoryAlert";
import InventoryChart from "@/components/InventoryChart";
import { NotificationToggle } from "@/components/NotificationToggle";
import { useInventory } from "@/hooks/useInventory";

const VAPID_PUBLIC_KEY = process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY ?? "";

export default function InventoryPage() {
	const { inventories, error, loading } = useInventory();

	return (
		<main style={{ maxWidth: "800px", margin: "0 auto", padding: "24px 16px" }}>
			<div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 24 }}>
				<h1 style={{ fontSize: "24px", fontWeight: 700 }}>血液型別在庫ダッシュボード</h1>
				{VAPID_PUBLIC_KEY && <NotificationToggle vapidPublicKey={VAPID_PUBLIC_KEY} />}
			</div>

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
