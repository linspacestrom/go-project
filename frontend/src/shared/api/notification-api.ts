import { Notification } from "@/entities/notification/model/types";
import { loadAuthSession, clearAuthSession } from "@/shared/api/auth-storage";
import { createHttpClient } from "@/shared/api/http-client";
import { ListResponse, PaginationQuery } from "@/shared/api/types";
import { env } from "@/shared/config/env";
import { emitAuthExpired } from "@/shared/lib/auth-events";
import { refreshAccessToken } from "@/shared/api/main-api";

const client = createHttpClient({
  baseUrl: env.notificationApiUrl,
  getAccessToken: () => loadAuthSession()?.accessToken ?? null,
  refreshAccessToken,
  onAuthFailure: () => {
    clearAuthSession();
    emitAuthExpired();
  }
});

export const notificationApi = {
  listNotifications(query: PaginationQuery) {
    return client.request<ListResponse<Notification>>("/api/v1/notifications", { query });
  },
  listUnreadNotifications() {
    return client.request<ListResponse<Notification>>("/api/v1/notifications/unread");
  },
  markRead(id: string) {
    return client.request<void>(`/api/v1/notifications/${id}/read`, {
      method: "PATCH"
    });
  },
  markReadAll() {
    return client.request<void>("/api/v1/notifications/read-all", {
      method: "PATCH"
    });
  },
  createRealtimeUrl(accessToken: string) {
    const url = new URL(env.notificationRealtimeUrl);
    url.searchParams.set("token", accessToken);
    return url.toString();
  }
};
