import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { notificationApi } from "@/shared/api/notification-api";

export function NotificationBell() {
  const unreadQuery = useQuery({
    queryKey: ["notifications-unread"],
    queryFn: () => notificationApi.listUnreadNotifications()
  });

  const unreadCount = unreadQuery.data?.items.length ?? 0;

  return (
    <Link
      to="/notifications"
      style={{
        position: "relative",
        display: "inline-flex",
        border: "1px solid #94a3b8",
        borderRadius: 10,
        padding: "8px 12px",
        background: "white"
      }}
    >
      Уведомления
      {unreadCount > 0 ? (
        <span
          style={{
            position: "absolute",
            top: -6,
            right: -6,
            minWidth: 20,
            height: 20,
            borderRadius: 999,
            background: "#dc2626",
            color: "white",
            fontSize: 12,
            display: "grid",
            placeItems: "center",
            padding: "0 6px"
          }}
        >
          {unreadCount}
        </span>
      ) : null}
    </Link>
  );
}
