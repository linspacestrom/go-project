import { FormEvent, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function AdminCitiesPage() {
  const queryClient = useQueryClient();
  const listQuery = useQuery({
    queryKey: ["admin-cities"],
    queryFn: () => mainApi.listCities({ limit: 100, offset: 0 })
  });

  const [createName, setCreateName] = useState("");
  const [updateId, setUpdateId] = useState("");
  const [updateName, setUpdateName] = useState("");
  const [error, setError] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: () => mainApi.createCity({ name: createName }),
    onSuccess: () => {
      setError(null);
      setCreateName("");
      queryClient.invalidateQueries({ queryKey: ["admin-cities"] });
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось создать город")
  });
  const updateMutation = useMutation({
    mutationFn: () => mainApi.updateCity(updateId, { name: updateName || undefined }),
    onSuccess: () => {
      setError(null);
      setUpdateId("");
      setUpdateName("");
      queryClient.invalidateQueries({ queryKey: ["admin-cities"] });
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось обновить город")
  });

  function onCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createMutation.mutate();
  }

  function onUpdate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    updateMutation.mutate();
  }

  return (
    <div className="grid">
      <PageTitle title="Admin Cities" subtitle="Создание и редактирование городов" />
      {error ? <ErrorState message={error} /> : null}
      <Card title="Создать город">
        <form onSubmit={onCreate} style={{ display: "flex", gap: 8 }}>
          <input
            value={createName}
            onChange={(event) => setCreateName(event.target.value)}
            placeholder="Название города"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px", flex: 1 }}
            required
          />
          <Button type="submit" loading={createMutation.isPending}>
            Create
          </Button>
        </form>
      </Card>
      <Card title="Обновить город">
        <form onSubmit={onUpdate} style={{ display: "grid", gap: 8 }}>
          <input
            value={updateId}
            onChange={(event) => setUpdateId(event.target.value)}
            placeholder="City ID"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={updateName}
            onChange={(event) => setUpdateName(event.target.value)}
            placeholder="Новое имя (optional)"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
          />
          <Button type="submit" loading={updateMutation.isPending}>
            Update
          </Button>
        </form>
      </Card>
      <Card title="Список городов">
        {listQuery.isLoading ? <Loader /> : null}
        {listQuery.isError ? <ErrorState message="Не удалось загрузить города" /> : null}
        {listQuery.data && listQuery.data.items.length === 0 ? <EmptyState title="Города не найдены" /> : null}
        <div className="grid">
          {listQuery.data?.items.map((city) => (
            <div key={city.id} style={{ borderBottom: "1px solid #e2e8f0", paddingBottom: 8 }}>
              {city.name} ({city.id}) | active: {String(city.is_active)}
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
