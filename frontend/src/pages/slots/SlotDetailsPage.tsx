import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";
import { formatDateTime } from "@/shared/lib/date";

export function SlotDetailsPage() {
  const { id = "" } = useParams();
  const query = useQuery({
    queryKey: ["slot", id],
    enabled: Boolean(id),
    queryFn: () => mainApi.getSlotById(id)
  });

  return (
    <div className="grid">
      <PageTitle title="Детали слота" subtitle={id} />
      {query.isLoading ? <Loader /> : null}
      {query.isError ? <ErrorState message="Не удалось загрузить слот" /> : null}
      {query.data ? (
        <Card>
          <div className="grid" style={{ gridTemplateColumns: "repeat(2, minmax(160px, 1fr))" }}>
            <div>ID: {query.data.id}</div>
            <div>Mentor ID: {query.data.mentor_id}</div>
            <div>Room ID: {query.data.room_id ?? "—"}</div>
            <div>City ID: {query.data.city_id}</div>
            <div>Type: {query.data.type}</div>
            <div>Capacity: {query.data.capacity ?? "—"}</div>
            <div>Start: {formatDateTime(query.data.start_at)}</div>
            <div>End: {formatDateTime(query.data.end_at)}</div>
            <div>
              Status:{" "}
              <Badge tone={query.data.status === "active" ? "success" : "warning"}>{query.data.status}</Badge>
            </div>
            <div>Meeting URL: {query.data.meeting_url ?? "—"}</div>
          </div>
        </Card>
      ) : null}
    </div>
  );
}
