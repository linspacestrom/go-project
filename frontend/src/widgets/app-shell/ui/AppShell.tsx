import { Outlet } from "react-router-dom";
import { useNotificationsRealtime } from "@/features/notifications/model/useNotificationsRealtime";
import { Sidebar } from "@/widgets/sidebar/ui/Sidebar";
import { Topbar } from "@/widgets/topbar/ui/Topbar";

export function AppShell() {
  useNotificationsRealtime();

  return (
    <div style={{ minHeight: "100vh", display: "flex" }}>
      <Sidebar />
      <div style={{ flex: 1, display: "grid", gridTemplateRows: "auto 1fr" }}>
        <Topbar />
        <main className="container">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
