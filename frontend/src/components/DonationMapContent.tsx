"use client";

import type { DonationCenter } from "@/types";
import "leaflet/dist/leaflet.css";
import { type Icon, icon } from "leaflet";
import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";

const TOKYO_CENTER: [number, number] = [35.6812, 139.7671];
const DEFAULT_ZOOM = 12;

const defaultIcon: Icon = icon({
	iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
	iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
	shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
	iconSize: [25, 41],
	iconAnchor: [12, 41],
	popupAnchor: [1, -34],
	shadowSize: [41, 41],
});

function statusLabel(status: DonationCenter["status"]): string {
	switch (status) {
		case "open":
			return "受付中";
		case "closed":
			return "受付終了";
		case "full":
			return "満員";
		default:
			return "不明";
	}
}

function statusColor(status: DonationCenter["status"]): string {
	switch (status) {
		case "open":
			return "#16a34a";
		case "closed":
			return "#dc2626";
		case "full":
			return "#ea580c";
		default:
			return "#6b7280";
	}
}

interface DonationMapContentProps {
	centers: DonationCenter[];
}

export default function DonationMapContent({ centers }: DonationMapContentProps) {
	return (
		<MapContainer
			center={TOKYO_CENTER}
			zoom={DEFAULT_ZOOM}
			scrollWheelZoom={true}
			style={{ height: "100%", width: "100%" }}
		>
			<TileLayer
				attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
				url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
			/>
			{centers.map((center) => (
				<Marker key={center.id} position={[center.lat, center.lng]} icon={defaultIcon}>
					<Popup>
						<div style={{ minWidth: 180 }}>
							<h3 style={{ margin: "0 0 4px 0", fontSize: 14, fontWeight: 600 }}>{center.name}</h3>
							<p style={{ margin: "0 0 4px 0", fontSize: 12, color: "#4b5563" }}>{center.address}</p>
							<p style={{ margin: "0 0 4px 0", fontSize: 12 }}>
								<span
									style={{
										color: statusColor(center.status),
										fontWeight: 600,
									}}
								>
									{statusLabel(center.status)}
								</span>
							</p>
							<p style={{ margin: 0, fontSize: 12, color: "#6b7280" }}>
								空き: {center.availableSlots} / {center.capacity}
							</p>
						</div>
					</Popup>
				</Marker>
			))}
		</MapContainer>
	);
}
