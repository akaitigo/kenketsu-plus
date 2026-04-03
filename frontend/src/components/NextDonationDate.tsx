"use client";

import { api } from "@/lib/api";
import type { Gender, NextAvailableResult } from "@/types";
import { useEffect, useState } from "react";

interface NextDonationDateProps {
	gender: Gender;
}

export function NextDonationDate({ gender }: NextDonationDateProps) {
	const [result, setResult] = useState<NextAvailableResult | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		setLoading(true);
		api
			.get<NextAvailableResult>(`/api/donations/next-available?gender=${gender}`)
			.then(setResult)
			.catch(() => setResult(null))
			.finally(() => setLoading(false));
	}, [gender]);

	if (loading) return <p>読み込み中...</p>;
	if (!result) return <p>取得に失敗しました</p>;

	return (
		<div
			style={{
				padding: 16,
				borderRadius: 8,
				backgroundColor: result.canDonateToday ? "#dcfce7" : "#fef3c7",
				marginTop: 8,
			}}
		>
			{result.canDonateToday ? (
				<p style={{ color: "#166534", fontWeight: "bold" }}>本日から献血可能です</p>
			) : (
				<>
					<p style={{ fontWeight: "bold" }}>次回献血可能日: {new Date(result.nextDate).toLocaleDateString("ja-JP")}</p>
					<p>あと {result.daysRemaining} 日</p>
				</>
			)}
			<p style={{ fontSize: 14, color: "#666" }}>{result.reason}</p>
		</div>
	);
}
