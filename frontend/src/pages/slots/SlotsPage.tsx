import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { SlotFilters } from "@/features/slots/filter/ui/SlotFilters";
import { SlotFilter, mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";
import { formatDateTime } from "@/shared/lib/date";

export function SlotsPage() {
  const [filter, setFilter] = useState<SlotFilter>({
    limit: 20,
    offset: 0
  });

  const query = useQuery({
    queryKey: ["slots", filter],
    queryFn: () => mainApi.listSlots(filter)
  });

  return (
    <div className="grid">
      <PageTitle title="Слоты" subtitle="Календарь доступных слотов" />
      <SlotFilters onApply={setFilter} />
      {query.isLoading ? <Loader /> : null}
      {query.isError ? <ErrorState message="Не удалось загрузить слоты" /> : null}
      {query.data && query.data.items.length === 0 ? <EmptyState title="Слоты не найдены" /> : null}
      <div className="card-grid">
        {query.data?.items.map((slot) => (
          <Card key={slot.id} title={`Слот ${slot.id.slice(0, 8)}...`}>
            <div className="grid">
              <div>Mentor: {slot.mentor_id}</div>
              <div>City: {slot.city_id}</div>
              <div>Room: {slot.room_id ?? "—"}</div>
              <div>Type: {slot.type}</div>
              <div>Start: {formatDateTime(slot.start_at)}</div>
              <div>End: {formatDateTime(slot.end_at)}</div>
              <div>
                <Badge tone={slot.status === "active" ? "success" : "warning"}>{slot.status}</Badge>
              </div>
              <Link to={`/slots/${slot.id}`}>Открыть детали</Link>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
