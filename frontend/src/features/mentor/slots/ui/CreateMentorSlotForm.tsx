import { FormEvent, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { Select } from "@/shared/ui/Select";

export function CreateMentorSlotForm() {
  const queryClient = useQueryClient();
  const [cityId, setCityId] = useState("");
  const [roomId, setRoomId] = useState("");
  const [type, setType] = useState<"online" | "offline">("online");
  const [startAt, setStartAt] = useState("");
  const [endAt, setEndAt] = useState("");
  const [meetingUrl, setMeetingUrl] = useState("");
  const [capacity, setCapacity] = useState("");
  const [error, setError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () =>
      mainApi.createMentorSlot({
        city_id: cityId,
        room_id: roomId || undefined,
        type,
        start_at: new Date(startAt).toISOString(),
        end_at: new Date(endAt).toISOString(),
        meeting_url: meetingUrl || undefined,
        capacity: capacity ? Number(capacity) : undefined
      }),
    onSuccess: () => {
      setError(null);
      queryClient.invalidateQueries({ queryKey: ["mentor-slots"] });
    },
    onError: (e) => {
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось создать слот");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!startAt || !endAt || new Date(startAt) >= new Date(endAt)) {
      setError("Проверьте время начала и окончания");
      return;
    }
    if (type === "offline" && !roomId) {
      setError("Для offline слота обязателен room_id");
      return;
    }
    if (type === "online" && roomId) {
      setError("Для online слота room_id должен быть пустым");
      return;
    }
    if (!meetingUrl) {
      setError("Meeting URL обязателен");
      return;
    }
    mutation.mutate();
  }

  return (
    <Card title="Создать слот ментора">
      <form onSubmit={onSubmit} className="grid">
        <Input label="City ID" value={cityId} onChange={(event) => setCityId(event.target.value)} required />
        <Select
          label="Тип"
          value={type}
          onChange={(event) => setType(event.target.value as "online" | "offline")}
          options={[
            { label: "Online", value: "online" },
            { label: "Offline", value: "offline" }
          ]}
        />
        <Input
          label="Room ID (для offline)"
          value={roomId}
          onChange={(event) => setRoomId(event.target.value)}
        />
        <Input
          label="Start At"
          type="datetime-local"
          value={startAt}
          onChange={(event) => setStartAt(event.target.value)}
          required
        />
        <Input
          label="End At"
          type="datetime-local"
          value={endAt}
          onChange={(event) => setEndAt(event.target.value)}
          required
        />
        <Input
          label="Meeting URL"
          value={meetingUrl}
          onChange={(event) => setMeetingUrl(event.target.value)}
          required
        />
        <Input
          label="Capacity (optional)"
          type="number"
          value={capacity}
          onChange={(event) => setCapacity(event.target.value)}
        />
        {error ? <ErrorState message={error} /> : null}
        <Button type="submit" loading={mutation.isPending}>
          Создать слот
        </Button>
      </form>
    </Card>
  );
}
