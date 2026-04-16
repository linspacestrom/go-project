import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notificationApi } from "@/shared/api/notification-api";
import { Badge } from "@/shared/ui/Badge";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";
import { formatDateTime } from "@/shared/lib/date";

export function NotificationsPage() {
  const queryClient = useQueryClient();
  const listQuery = useQuery({
    queryKey: ["notifications"],
    queryFn: () => notificationApi.listNotifications({ limit: 100, offset: 0 })
  });

  const markReadMutation = useMutation({
    mutationFn: (id: string) => notificationApi.markRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications-unread"] });
    }
  });

  const markAllMutation = useMutation({
    mutationFn: () => notificationApi.markReadAll(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications-unread"] });
    }
  });

  return (
    <div className="grid">
      <PageTitle title="Уведомления" subtitle="Список уведомлений и статус прочтения" />
      <div>
        <Button type="button" onClick={() => markAllMutation.mutate()} loading={markAllMutation.isPending}>
          Пометить все прочитанными
        </Button>
      </div>
      <Card>
        {listQuery.isLoading ? <Loader /> : null}
        {listQuery.isError ? <ErrorState message="Не удалось загрузить уведомления" /> : null}
        {listQuery.data && listQuery.data.items.length === 0 ? <EmptyState title="Уведомлений нет" /> : null}
        <div className="grid">
          {listQuery.data?.items.map((item) => (
            <Card key={item.id} title={item.title}>
              <div className="grid">
                <div>{item.body}</div>
                <div>
                  Status:{" "}
                  <Badge tone={item.is_read ? "neutral" : "warning"}>{item.is_read ? "read" : "unread"}</Badge>
                </div>
                <div>Type: {item.type}</div>
                <div>Channel: {item.delivery_channel}</div>
                <div>Created: {formatDateTime(item.created_at)}</div>
                {!item.is_read ? (
                  <Button
                    type="button"
                    variant="secondary"
                    loading={markReadMutation.isPending}
                    onClick={() => markReadMutation.mutate(item.id)}
                  >
                    Mark as read
                  </Button>
                ) : null}
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}
