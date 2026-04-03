export type CenterStatus = "open" | "closed" | "full";

export interface DonationCenter {
	id: string;
	name: string;
	address: string;
	lat: number;
	lng: number;
	capacity: number;
	availableSlots: number;
	status: CenterStatus;
	createdAt: string;
	updatedAt: string;
}

export type BloodType = "A+" | "A-" | "B+" | "B-" | "O+" | "O-" | "AB+" | "AB-";

export type DonationType = "whole_400" | "whole_200" | "component";

export type Gender = "male" | "female";

export interface Donation {
	id: string;
	bloodType: BloodType;
	donationType: DonationType;
	gender: Gender;
	donatedAt: string;
	volumeMl: number;
	memo: string;
	createdAt: string;
}

export type InventoryLevel = "critical" | "low" | "normal" | "sufficient";

export interface BloodInventory {
	id: string;
	bloodType: BloodType;
	level: InventoryLevel;
	updatedAt: string;
}

export interface PushSubscription {
	id: string;
	endpoint: string;
	p256dh: string;
	auth: string;
	createdAt: string;
}

export interface NextAvailableResult {
	nextDate: string;
	daysRemaining: number;
	canDonateToday: boolean;
	reason: string;
}
