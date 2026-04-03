"use client";

import type { BloodInventory, InventoryLevel } from "@/types";

const LEVEL_COLORS: Record<InventoryLevel, string> = {
	critical: "#dc2626",
	low: "#f97316",
	normal: "#22c55e",
	sufficient: "#3b82f6",
};

const LEVEL_HEIGHTS: Record<InventoryLevel, number> = {
	critical: 25,
	low: 50,
	normal: 75,
	sufficient: 100,
};

const LEVEL_LABELS: Record<InventoryLevel, string> = {
	critical: "逼迫",
	low: "不足",
	normal: "通常",
	sufficient: "十分",
};

interface InventoryChartProps {
	inventories: BloodInventory[];
}

export default function InventoryChart({ inventories }: InventoryChartProps) {
	if (inventories.length === 0) {
		return <p data-testid="inventory-empty">在庫データがありません</p>;
	}

	return (
		<div
			data-testid="inventory-chart"
			style={{
				display: "flex",
				gap: "12px",
				alignItems: "flex-end",
				justifyContent: "center",
				padding: "24px 16px",
				minHeight: "280px",
			}}
		>
			{inventories.map((inv) => {
				const color = LEVEL_COLORS[inv.level];
				const heightPct = LEVEL_HEIGHTS[inv.level];
				const label = LEVEL_LABELS[inv.level];

				return (
					<div
						key={inv.id}
						style={{
							display: "flex",
							flexDirection: "column",
							alignItems: "center",
							flex: "1 1 0",
							maxWidth: "80px",
						}}
					>
						<span
							style={{
								fontSize: "12px",
								color: color,
								fontWeight: 600,
								marginBottom: "4px",
							}}
						>
							{label}
						</span>
						<div
							data-testid={`bar-${inv.bloodType}`}
							style={{
								width: "100%",
								height: `${heightPct * 2}px`,
								backgroundColor: color,
								borderRadius: "4px 4px 0 0",
								transition: "height 0.3s ease",
								minHeight: "20px",
							}}
							role="meter"
							aria-label={`${inv.bloodType}: ${label}`}
							aria-valuenow={heightPct}
							aria-valuemin={0}
							aria-valuemax={100}
						/>
						<span
							style={{
								marginTop: "8px",
								fontWeight: 700,
								fontSize: "14px",
							}}
						>
							{inv.bloodType}
						</span>
					</div>
				);
			})}
		</div>
	);
}
