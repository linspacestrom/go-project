import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { mainApi } from "@/shared/api/main-api";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function HubsPage() {
  const [params, setParams] = useSearchParams();
  const cityId = params.get("cityId") ?? "";

  const query = useQuery({
    queryKey: ["hubs", cityId],
    enabled: Boolean(cityId),
    queryFn: () => mainApi.listHubsByCity(cityId, { limit: 50, offset: 0 })
  });

  return (
    <div className="grid">
      <PageTitle title="Хабы" subtitle="Список хабов выбранного города" />
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
      {!cityId ? <EmptyState title="Укажи cityId для загрузки хабов" /> : null}
      {query.isLoading ? <Loader /> : null}
      {query.isError ? <ErrorState message="Не удалось загрузить хабы" /> : null}
      {query.data && query.data.items.length === 0 ? <EmptyState title="Хабы не найдены" /> : null}
      <div className="card-grid">
        {query.data?.items.map((hub) => (
          <Card key={hub.id} title={hub.name}>
            <div className="grid">
              <div>ID: {hub.id}</div>
              <div>Address: {hub.address}</div>
              <div>City ID: {hub.city_id}</div>
              <div>Активен: {hub.is_active ? "Да" : "Нет"}</div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
