"use client";

import { api } from "@/lib/api";
import type { BloodType, Donation, DonationType, Gender } from "@/types";
import { type FormEvent, useState } from "react";

interface DonationFormProps {
	onCreated: () => void;
}

const BLOOD_TYPES: BloodType[] = ["A+", "A-", "B+", "B-", "O+", "O-", "AB+", "AB-"];
const DONATION_TYPES: { value: DonationType; label: string }[] = [
	{ value: "whole_400", label: "全血400ml" },
	{ value: "whole_200", label: "全血200ml" },
	{ value: "component", label: "成分献血" },
];

export function DonationForm({ onCreated }: DonationFormProps) {
	const [bloodType, setBloodType] = useState<BloodType>("A+");
	const [donationType, setDonationType] = useState<DonationType>("whole_400");
	const [gender, setGender] = useState<Gender>("male");
	const [donatedAt, setDonatedAt] = useState("");
	const [error, setError] = useState("");
	const [submitting, setSubmitting] = useState(false);

	const volumeMap: Record<DonationType, number> = {
		whole_400: 400,
		whole_200: 200,
		component: 0,
	};

	const handleSubmit = async (e: FormEvent) => {
		e.preventDefault();
		setError("");

		if (!donatedAt) {
			setError("献血日を入力してください");
			return;
		}

		setSubmitting(true);
		try {
			await api.post<Donation>("/api/donations", {
				bloodType,
				donationType,
				gender,
				donatedAt: new Date(donatedAt).toISOString(),
				volumeMl: volumeMap[donationType],
			});
			setDonatedAt("");
			onCreated();
		} catch {
			setError("登録に失敗しました");
		} finally {
			setSubmitting(false);
		}
	};

	return (
		<form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: 12, maxWidth: 400 }}>
			<div>
				<label htmlFor="blood-type">血液型: </label>
				<select id="blood-type" value={bloodType} onChange={(e) => setBloodType(e.target.value as BloodType)}>
					{BLOOD_TYPES.map((bt) => (
						<option key={bt} value={bt}>
							{bt}
						</option>
					))}
				</select>
			</div>
			<div>
				<label htmlFor="donation-type">種別: </label>
				<select
					id="donation-type"
					value={donationType}
					onChange={(e) => setDonationType(e.target.value as DonationType)}
				>
					{DONATION_TYPES.map((dt) => (
						<option key={dt.value} value={dt.value}>
							{dt.label}
						</option>
					))}
				</select>
			</div>
			<div>
				<label htmlFor="gender">性別: </label>
				<select id="gender" value={gender} onChange={(e) => setGender(e.target.value as Gender)}>
					<option value="male">男性</option>
					<option value="female">女性</option>
				</select>
			</div>
			<div>
				<label htmlFor="donated-at">献血日: </label>
				<input id="donated-at" type="date" value={donatedAt} onChange={(e) => setDonatedAt(e.target.value)} />
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<button type="submit" disabled={submitting}>
				{submitting ? "登録中..." : "記録する"}
			</button>
		</form>
	);
}
