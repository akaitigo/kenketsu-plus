import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
	title: "Kenketsu-Plus",
	description: "献血施設・登録された空き枠の確認、献血記録、登録された在庫情報の通知",
};

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="ja">
			<body>{children}</body>
		</html>
	);
}
