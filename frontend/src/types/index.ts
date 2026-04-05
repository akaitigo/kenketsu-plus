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

// Type guard functions for validated select inputs

const VALID_BLOOD_TYPES: ReadonlySet<string> = new Set<string>(["A+", "A-", "B+", "B-", "O+", "O-", "AB+", "AB-"]);

const VALID_DONATION_TYPES: ReadonlySet<string> = new Set<string>(["whole_400", "whole_200", "component"]);

const VALID_GENDERS: ReadonlySet<string> = new Set<string>(["male", "female"]);

export function isBloodType(value: string): value is BloodType {
	return VALID_BLOOD_TYPES.has(value);
}

export function isDonationType(value: string): value is DonationType {
	return VALID_DONATION_TYPES.has(value);
}

export function isGender(value: string): value is Gender {
	return VALID_GENDERS.has(value);
}

export function parseBloodType(value: string): BloodType {
	if (!isBloodType(value)) {
		throw new Error(`Invalid blood type: ${value}`);
	}
	return value;
}

export function parseDonationType(value: string): DonationType {
	if (!isDonationType(value)) {
		throw new Error(`Invalid donation type: ${value}`);
	}
	return value;
}

export function parseGender(value: string): Gender {
	if (!isGender(value)) {
		throw new Error(`Invalid gender: ${value}`);
	}
	return value;
}
