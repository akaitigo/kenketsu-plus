"use client";

import { DonationForm } from "@/components/DonationForm";
import { NextDonationDate } from "@/components/NextDonationDate";
import { api } from "@/lib/api";
import { type Donation, type Gender, parseGender } from "@/types";
import { useCallback, useEffect, useState } from "react";

export default function DonationsPage() {
	const [donations, setDonations] = useState<Donation[]>([]);
	const [gender, setGender] = useState<Gender>("male");

	const fetchDonations = useCallback(async () => {
		const data = await api.get<Donation[]>("/api/donations");
		setDonations(data);
	}, []);

	useEffect(() => {
		void fetchDonations();
	}, [fetchDonations]);

	const handleCreated = () => {
		void fetchDonations();
	};

	return (
		<main style={{ maxWidth: 800, margin: "0 auto", padding: 20 }}>
			<h1>献血記録</h1>

			<section>
				<h2>次回献血可能日</h2>
				<label htmlFor="gender-select">性別: </label>
				<select id="gender-select" value={gender} onChange={(e) => setGender(parseGender(e.target.value))}>
					<option value="male">男性</option>
					<option value="female">女性</option>
				</select>
				<NextDonationDate gender={gender} />
			</section>

			<section>
				<h2>新規記録</h2>
				<DonationForm onCreated={handleCreated} />
			</section>

			<section>
				<h2>献血履歴</h2>
				{donations.length === 0 ? (
					<p>記録がありません</p>
				) : (
					<table style={{ width: "100%", borderCollapse: "collapse" }}>
						<thead>
							<tr>
								<th style={{ textAlign: "left", borderBottom: "1px solid #ccc", padding: 8 }}>日付</th>
								<th style={{ textAlign: "left", borderBottom: "1px solid #ccc", padding: 8 }}>血液型</th>
								<th style={{ textAlign: "left", borderBottom: "1px solid #ccc", padding: 8 }}>種別</th>
								<th style={{ textAlign: "left", borderBottom: "1px solid #ccc", padding: 8 }}>量(ml)</th>
							</tr>
						</thead>
						<tbody>
							{donations.map((d) => (
								<tr key={d.id}>
									<td style={{ padding: 8 }}>{new Date(d.donatedAt).toLocaleDateString("ja-JP")}</td>
									<td style={{ padding: 8 }}>{d.bloodType}</td>
									<td style={{ padding: 8 }}>{d.donationType}</td>
									<td style={{ padding: 8 }}>{d.volumeMl}</td>
								</tr>
							))}
						</tbody>
					</table>
				)}
			</section>
		</main>
	);
}
