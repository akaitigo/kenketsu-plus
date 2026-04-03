import DonationMap from "@/components/DonationMap";
import type { Metadata } from "next";

export const metadata: Metadata = {
	title: "献血ルームマップ | Kenketsu-Plus",
	description: "近くの献血ルームをマップで確認できます",
};

export default function MapPage() {
	return (
		<main style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
			<header
				style={{
					padding: "12px 16px",
					borderBottom: "1px solid #e5e7eb",
					display: "flex",
					alignItems: "center",
					gap: 12,
				}}
			>
				<a href="/" style={{ textDecoration: "none", color: "#6b7280", fontSize: 14 }}>
					&larr; ホーム
				</a>
				<h1 style={{ margin: 0, fontSize: 18, fontWeight: 600 }}>献血ルームマップ</h1>
			</header>
			<div style={{ flex: 1, position: "relative" }}>
				<DonationMap />
			</div>
		</main>
	);
}
