"use client";
import dynamic from "next/dynamic";
import { useCenters } from "./useCenters";

const DonationMapContent = dynamic(() => import("./DonationMapContent"), {
	ssr: false,
	loading: () => (
		<div style={{ height: "100%", display: "flex", alignItems: "center", justifyContent: "center" }}>
			<p>マップを読み込んでいます...</p>
		</div>
	),
});

function CenteredMessage({ children }: { children: React.ReactNode }) {
	return (
		<div style={{ height: "100%", display: "flex", alignItems: "center", justifyContent: "center" }}>{children}</div>
	);
}

export default function DonationMap() {
	const { centers, error, loading } = useCenters();

	if (loading) {
		return (
			<CenteredMessage>
				<p>読み込み中...</p>
			</CenteredMessage>
		);
	}

	if (error) {
		return (
			<CenteredMessage>
				<p style={{ color: "#dc2626" }}>{error}</p>
			</CenteredMessage>
		);
	}

	return <DonationMapContent centers={centers} />;
}
