const defaultMainApiUrl = "http://localhost:8080";
const defaultNotificationApiUrl = "http://localhost:8081";

export const env = {
  mainApiUrl: import.meta.env.VITE_MAIN_API_URL ?? defaultMainApiUrl,
  notificationApiUrl: import.meta.env.VITE_NOTIFICATION_API_URL ?? defaultNotificationApiUrl,
  notificationRealtimeUrl:
    import.meta.env.VITE_NOTIFICATION_REALTIME_URL ??
    `${import.meta.env.VITE_NOTIFICATION_API_URL ?? defaultNotificationApiUrl}/api/v1/notifications/realtime`
};
