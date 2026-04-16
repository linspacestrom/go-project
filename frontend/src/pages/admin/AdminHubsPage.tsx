import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { PageTitle } from "@/shared/ui/PageTitle";

export function AdminHubsPage() {
  const [cityId, setCityId] = useState("");
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [updateId, setUpdateId] = useState("");
  const [updateName, setUpdateName] = useState("");
  const [updateAddress, setUpdateAddress] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: () => mainApi.createHub({ city_id: cityId, name, address }),
    onSuccess: () => {
      setError(null);
      setSuccess("Хаб создан");
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось создать хаб")
  });

  const updateMutation = useMutation({
    mutationFn: () =>
      mainApi.updateHub(updateId, { name: updateName || undefined, address: updateAddress || undefined }),
    onSuccess: () => {
      setError(null);
      setSuccess("Хаб обновлен");
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось обновить хаб")
  });

  function onCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSuccess(null);
    createMutation.mutate();
  }

  function onUpdate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSuccess(null);
    updateMutation.mutate();
  }

  return (
    <div className="grid">
      <PageTitle title="Admin Hubs" subtitle="Создание и обновление хабов" />
      {error ? <ErrorState message={error} /> : null}
      {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
      <Card title="Создать хаб">
        <form onSubmit={onCreate} className="grid">
          <input
            value={cityId}
            onChange={(event) => setCityId(event.target.value)}
            placeholder="City ID"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Название"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={address}
            onChange={(event) => setAddress(event.target.value)}
            placeholder="Адрес"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <Button type="submit" loading={createMutation.isPending}>
            Create
          </Button>
        </form>
      </Card>
      <Card title="Обновить хаб">
        <form onSubmit={onUpdate} className="grid">
          <input
            value={updateId}
            onChange={(event) => setUpdateId(event.target.value)}
            placeholder="Hub ID"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={updateName}
            onChange={(event) => setUpdateName(event.target.value)}
            placeholder="Новое имя"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
          />
          <input
            value={updateAddress}
            onChange={(event) => setUpdateAddress(event.target.value)}
            placeholder="Новый адрес"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
          />
          <Button type="submit" loading={updateMutation.isPending}>
            Update
          </Button>
        </form>
      </Card>
    </div>
  );
}
