import { Link, useLocation } from "react-router-dom";
import { useAuth } from "@/features/auth/model/auth-context";

type NavItem = {
  to: string;
  label: string;
  roles?: Array<"student" | "mentor" | "admin">;
};

const links: NavItem[] = [
  { to: "/", label: "Dashboard" },
  { to: "/me", label: "Профиль" },
  { to: "/cities", label: "Города" },
  { to: "/hubs", label: "Хабы" },
  { to: "/rooms", label: "Комнаты" },
  { to: "/slots", label: "Слоты" },
  { to: "/bookings", label: "Бронирования" },
  { to: "/mentor-requests", label: "Mentor Requests" },
  { to: "/notifications", label: "Уведомления" },
  { to: "/mentor/skills", label: "Mentor Skills", roles: ["mentor"] },
  { to: "/mentor/slots", label: "Mentor Slots", roles: ["mentor"] },
  { to: "/mentor/requests", label: "Mentor Queue", roles: ["mentor"] },
  { to: "/admin/register-mentor", label: "Admin: Register Mentor", roles: ["admin"] },
  { to: "/admin/cities", label: "Admin: Cities", roles: ["admin"] },
  { to: "/admin/hubs", label: "Admin: Hubs", roles: ["admin"] },
  { to: "/admin/rooms", label: "Admin: Rooms", roles: ["admin"] },
  { to: "/admin/users/city", label: "Admin: User City", roles: ["admin"] },
  { to: "/admin/analytics", label: "Admin: Analytics", roles: ["admin"] }
];

export function Sidebar() {
  const location = useLocation();
  const { role } = useAuth();

  return (
    <aside
      style={{
        width: 270,
        borderRight: "1px solid #cbd5e1",
        padding: 14,
        background: "#0f172a",
        color: "#e2e8f0"
      }}
    >
      <h2 style={{ margin: "6px 6px 14px", fontSize: 20 }}>Student & T</h2>
      <nav style={{ display: "grid", gap: 6 }}>
        {links
          .filter((item) => !item.roles || (role ? item.roles.includes(role) : false))
          .map((item) => {
            const active = location.pathname === item.to;
            return (
              <Link
                key={item.to}
                to={item.to}
                style={{
                  borderRadius: 10,
                  padding: "9px 10px",
                  background: active ? "#1d4ed8" : "transparent",
                  color: "#e2e8f0"
                }}
              >
                {item.label}
              </Link>
            );
          })}
      </nav>
    </aside>
  );
}
