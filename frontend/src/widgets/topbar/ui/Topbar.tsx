import { useNavigate } from "react-router-dom";
import { useAuth } from "@/features/auth/model/auth-context";
import { NotificationBell } from "@/features/notifications/ui/NotificationBell";
import { Button } from "@/shared/ui/Button";

export function Topbar() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  async function onLogout() {
    await logout();
    navigate("/login");
  }

  return (
    <header
      style={{
        borderBottom: "1px solid #cbd5e1",
        background: "white",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        padding: "14px 20px"
      }}
    >
      <div style={{ display: "grid", gap: 2 }}>
        <strong>{user?.full_name ?? "Гость"}</strong>
        <span style={{ color: "#64748b", fontSize: 13 }}>
          Role: {user?.role ?? "anonymous"} | Email: {user?.email ?? "—"}
        </span>
      </div>
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        <NotificationBell />
        <Button type="button" variant="ghost" onClick={onLogout}>
          Выйти
        </Button>
      </div>
    </header>
  );
}
