"use client";

import type { BloodInventory } from "@/types";

interface InventoryAlertProps {
	inventories: BloodInventory[];
}

export default function InventoryAlert({ inventories }: InventoryAlertProps) {
	const urgentItems = inventories.filter((inv) => inv.level === "critical" || inv.level === "low");

	if (urgentItems.length === 0) {
		return null;
	}

	const criticalItems = urgentItems.filter((inv) => inv.level === "critical");
	const lowItems = urgentItems.filter((inv) => inv.level === "low");

	return (
		<div
			data-testid="inventory-alert"
			role="alert"
			style={{
				padding: "16px",
				borderRadius: "8px",
				backgroundColor: criticalItems.length > 0 ? "#fef2f2" : "#fff7ed",
				border: `1px solid ${criticalItems.length > 0 ? "#fecaca" : "#fed7aa"}`,
				marginBottom: "16px",
			}}
		>
			<p
				style={{
					fontWeight: 700,
					fontSize: "16px",
					color: criticalItems.length > 0 ? "#dc2626" : "#f97316",
					margin: "0 0 8px 0",
				}}
			>
				{criticalItems.length > 0 ? "血液在庫が逼迫しています" : "血液在庫が不足しています"}
			</p>
			{criticalItems.length > 0 && (
				<p style={{ margin: "0 0 4px 0", color: "#dc2626" }}>
					逼迫: {criticalItems.map((inv) => inv.bloodType).join(", ")}
				</p>
			)}
			{lowItems.length > 0 && (
				<p style={{ margin: "0", color: "#f97316" }}>不足: {lowItems.map((inv) => inv.bloodType).join(", ")}</p>
			)}
		</div>
	);
}
