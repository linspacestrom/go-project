import { useQuery } from "@tanstack/react-query";
import { CreateMentorRequestForm } from "@/features/mentor-requests/create/ui/CreateMentorRequestForm";
import { mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function MentorRequestsPage() {
  const query = useQuery({
    queryKey: ["mentor-requests"],
    queryFn: () => mainApi.listMentorRequests({ limit: 50, offset: 0 })
  });

  return (
    <div className="grid">
      <PageTitle title="Mentor Requests" subtitle="Создание и отслеживание запросов ментору" />
      <CreateMentorRequestForm />
      <Card title="Список запросов">
        {query.isLoading ? <Loader /> : null}
        {query.isError ? <ErrorState message="Не удалось загрузить запросы" /> : null}
        {query.data && query.data.items.length === 0 ? <EmptyState title="Запросы отсутствуют" /> : null}
        <div className="grid">
          {query.data?.items.map((request) => (
            <Card key={request.id} title={`Request ${request.id.slice(0, 8)}...`}>
              <div className="grid">
                <div>Type: {request.request_type}</div>
                <div>Status: <Badge>{request.status}</Badge></div>
                <div>Mentee: {request.mentee_id}</div>
                <div>Mentor: {request.mentor_id ?? "—"}</div>
                <div>Skill: {request.skill_id ?? "—"}</div>
                <div>City: {request.city_id}</div>
                <div>Comment: {request.comment ?? "—"}</div>
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}
