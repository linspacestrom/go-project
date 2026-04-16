import { useQuery } from "@tanstack/react-query";
import { MentorRequestActions } from "@/features/mentor/requests/ui/MentorRequestActions";
import { mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function MentorRequestsQueuePage() {
  const query = useQuery({
    queryKey: ["mentor-own-requests"],
    queryFn: () => mainApi.listMentorOwnRequests({ limit: 50, offset: 0 })
  });

  return (
    <div className="grid">
      <PageTitle title="Mentor Queue" subtitle="Обработка запросов, назначенных ментору" />
      <Card title="Запросы">
        {query.isLoading ? <Loader /> : null}
        {query.isError ? <ErrorState message="Не удалось загрузить mentor requests" /> : null}
        {query.data && query.data.items.length === 0 ? <EmptyState title="Запросов пока нет" /> : null}
        <div className="grid">
          {query.data?.items.map((request) => (
            <Card key={request.id} title={`Request ${request.id.slice(0, 8)}...`}>
              <div className="grid">
                <div>Type: {request.request_type}</div>
                <div>Status: <Badge>{request.status}</Badge></div>
                <div>Comment: {request.comment ?? "—"}</div>
                <MentorRequestActions requestId={request.id} />
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}
