import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { loadAuthSession } from "@/shared/api/auth-storage";
import { notificationApi } from "@/shared/api/notification-api";

export function useNotificationsRealtime() {
  const queryClient = useQueryClient();

  useEffect(() => {
    const session = loadAuthSession();
    if (!session?.accessToken) {
      return;
    }

    let intervalId: number | null = null;
    let source: EventSource | null = null;
    const invalidate = () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications-unread"] });
    };

    try {
      source = new EventSource(notificationApi.createRealtimeUrl(session.accessToken));
      source.onmessage = () => invalidate();
      source.onerror = () => {
        source?.close();
        source = null;
        intervalId = window.setInterval(invalidate, 15_000);
      };
    } catch {
      intervalId = window.setInterval(invalidate, 15_000);
    }

    return () => {
      source?.close();
      if (intervalId !== null) {
        window.clearInterval(intervalId);
      }
    };
  }, [queryClient]);
}
