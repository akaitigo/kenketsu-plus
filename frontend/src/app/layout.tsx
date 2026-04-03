import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
	title: "Kenketsu-Plus",
	description: "献血ルーム空き状況・献血記録管理・血液型別在庫通知",
};

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="ja">
			<body>{children}</body>
		</html>
	);
}
