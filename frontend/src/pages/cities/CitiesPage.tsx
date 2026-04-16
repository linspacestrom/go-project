import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { mainApi } from "@/shared/api/main-api";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function CitiesPage() {
  const query = useQuery({
    queryKey: ["cities", { limit: 50, offset: 0 }],
    queryFn: () => mainApi.listCities({ limit: 50, offset: 0 })
  });

  return (
    <div className="grid">
      <PageTitle title="Города" subtitle="Справочник городов платформы" />
      {query.isLoading ? <Loader /> : null}
      {query.isError ? <ErrorState message="Не удалось загрузить города" /> : null}
      {query.data && query.data.items.length === 0 ? <EmptyState title="Города не найдены" /> : null}
      <div className="card-grid">
        {query.data?.items.map((city) => (
          <Card key={city.id} title={city.name}>
            <div className="grid">
              <div>ID: {city.id}</div>
              <div>Активен: {city.is_active ? "Да" : "Нет"}</div>
              <div style={{ display: "flex", gap: 8 }}>
                <Link to={`/hubs?cityId=${city.id}`}>Открыть хабы</Link>
                <Link to={`/rooms?cityId=${city.id}`}>Открыть комнаты</Link>
              </div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
