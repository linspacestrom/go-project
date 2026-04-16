import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { mainApi } from "@/shared/api/main-api";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function RoomsPage() {
  const [params, setParams] = useSearchParams();
  const cityId = params.get("cityId") ?? "";

  const query = useQuery({
    queryKey: ["rooms", cityId],
    enabled: Boolean(cityId),
    queryFn: () => mainApi.listRoomsByCity(cityId, { limit: 50, offset: 0 })
  });

  return (
    <div className="grid">
      <PageTitle title="Комнаты" subtitle="Комнаты хабов выбранного города" />
      <Card title="Параметры">
        <label className="grid" style={{ maxWidth: 500 }}>
          <span>City ID</span>
          <input
            value={cityId}
            onChange={(event) => setParams({ cityId: event.target.value })}
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
          />
        </label>
      </Card>
      {!cityId ? <EmptyState title="Укажи cityId для загрузки комнат" /> : null}
      {query.isLoading ? <Loader /> : null}
      {query.isError ? <ErrorState message="Не удалось загрузить комнаты" /> : null}
      {query.data && query.data.items.length === 0 ? <EmptyState title="Комнаты не найдены" /> : null}
      <div className="card-grid">
        {query.data?.items.map((room) => (
          <Card key={room.id} title={room.name}>
            <div className="grid">
              <div>ID: {room.id}</div>
              <div>Hub ID: {room.hub_id}</div>
              <div>Capacity: {room.capacity}</div>
              <div>Тип: {room.room_type ?? "—"}</div>
              <div>Активна: {room.is_active ? "Да" : "Нет"}</div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
