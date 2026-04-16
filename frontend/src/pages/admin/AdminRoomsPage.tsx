import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { PageTitle } from "@/shared/ui/PageTitle";

export function AdminRoomsPage() {
  const [hubId, setHubId] = useState("");
  const [name, setName] = useState("");
  const [capacity, setCapacity] = useState("1");
  const [roomId, setRoomId] = useState("");
  const [updateName, setUpdateName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: () => mainApi.createRoom({ hub_id: hubId, name, capacity: Number(capacity) }),
    onSuccess: () => {
      setError(null);
      setSuccess("Комната создана");
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось создать комнату")
  });

  const updateMutation = useMutation({
    mutationFn: () => mainApi.updateRoom(roomId, { name: updateName || undefined }),
    onSuccess: () => {
      setError(null);
      setSuccess("Комната обновлена");
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось обновить комнату")
  });

  const deleteMutation = useMutation({
    mutationFn: () => mainApi.deleteRoom(roomId),
    onSuccess: () => {
      setError(null);
      setSuccess("Комната удалена");
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : "Не удалось удалить комнату")
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
      <PageTitle title="Admin Rooms" subtitle="Управление жизненным циклом комнат" />
      {error ? <ErrorState message={error} /> : null}
      {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
      <Card title="Создать комнату">
        <form onSubmit={onCreate} className="grid">
          <input
            value={hubId}
            onChange={(event) => setHubId(event.target.value)}
            placeholder="Hub ID"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Название комнаты"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={capacity}
            onChange={(event) => setCapacity(event.target.value)}
            placeholder="Capacity"
            type="number"
            min={1}
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <Button type="submit" loading={createMutation.isPending}>
            Create
          </Button>
        </form>
      </Card>
      <Card title="Обновить / удалить комнату">
        <form onSubmit={onUpdate} className="grid">
          <input
            value={roomId}
            onChange={(event) => setRoomId(event.target.value)}
            placeholder="Room ID"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
            required
          />
          <input
            value={updateName}
            onChange={(event) => setUpdateName(event.target.value)}
            placeholder="Новое имя (optional)"
            style={{ border: "1px solid #94a3b8", borderRadius: 10, padding: "10px 12px" }}
          />
          <div style={{ display: "flex", gap: 8 }}>
            <Button type="submit" loading={updateMutation.isPending}>
              Update
            </Button>
            <Button
              type="button"
              variant="danger"
              loading={deleteMutation.isPending}
              onClick={() => deleteMutation.mutate()}
            >
              Delete
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
