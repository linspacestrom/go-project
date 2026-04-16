import { FormEvent, useState } from "react";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { Input } from "@/shared/ui/Input";
import { Select } from "@/shared/ui/Select";
import { SlotFilter } from "@/shared/api/main-api";

type Props = {
  onApply: (filter: SlotFilter) => void;
};

export function SlotFilters({ onApply }: Props) {
  const [filter, setFilter] = useState<SlotFilter>({
    limit: 20,
    offset: 0
  });

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onApply(filter);
  }

  return (
    <Card title="Фильтры слотов">
      <form onSubmit={submit} className="grid">
        <div className="grid" style={{ gridTemplateColumns: "repeat(3, minmax(120px, 1fr))" }}>
          <Input
            label="City ID"
            value={filter.city_id ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, city_id: event.target.value }))}
          />
          <Input
            label="Hub ID"
            value={filter.hub_id ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, hub_id: event.target.value }))}
          />
          <Input
            label="Room ID"
            value={filter.room_id ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, room_id: event.target.value }))}
          />
          <Input
            label="Mentor ID"
            value={filter.mentor_id ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, mentor_id: event.target.value }))}
          />
          <Select
            label="Тип"
            value={filter.type ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, type: event.target.value }))}
            options={[
              { label: "Все", value: "" },
              { label: "Online", value: "online" },
              { label: "Offline", value: "offline" }
            ]}
          />
          <Input
            label="Статус"
            value={filter.status ?? ""}
            onChange={(event) => setFilter((prev) => ({ ...prev, status: event.target.value }))}
          />
          <Input
            label="Start From"
            type="datetime-local"
            value={filter.start_from ?? ""}
            onChange={(event) =>
              setFilter((prev) => ({ ...prev, start_from: event.target.value ? new Date(event.target.value).toISOString() : "" }))
            }
          />
          <Input
            label="End To"
            type="datetime-local"
            value={filter.end_to ?? ""}
            onChange={(event) =>
              setFilter((prev) => ({ ...prev, end_to: event.target.value ? new Date(event.target.value).toISOString() : "" }))
            }
          />
          <Input
            label="Limit"
            type="number"
            value={String(filter.limit ?? 20)}
            onChange={(event) => setFilter((prev) => ({ ...prev, limit: Number(event.target.value) }))}
          />
        </div>
        <Button type="submit">Применить фильтры</Button>
      </form>
    </Card>
  );
}
