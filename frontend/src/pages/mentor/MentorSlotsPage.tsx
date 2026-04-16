import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CreateMentorSlotForm } from "@/features/mentor/slots/ui/CreateMentorSlotForm";
import { mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";
import { formatDateTime } from "@/shared/lib/date";

export function MentorSlotsPage() {
  const queryClient = useQueryClient();
  const listQuery = useQuery({
    queryKey: ["mentor-slots"],
    queryFn: () => mainApi.listMentorSlots()
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => mainApi.deleteMentorSlot(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["mentor-slots"] })
  });

  return (
    <div className="grid">
      <PageTitle title="Mentor Slots" subtitle="Создание и управление слотами ментора" />
      <CreateMentorSlotForm />
      <Card title="Мои слоты">
        {listQuery.isLoading ? <Loader /> : null}
        {listQuery.isError ? <ErrorState message="Не удалось загрузить mentor slots" /> : null}
        {listQuery.data && listQuery.data.items.length === 0 ? <EmptyState title="Слотов пока нет" /> : null}
        <div className="grid">
          {listQuery.data?.items.map((slot) => (
            <Card key={slot.id} title={`Slot ${slot.id.slice(0, 8)}...`}>
              <div className="grid">
                <div>Type: {slot.type}</div>
                <div>City: {slot.city_id}</div>
                <div>Room: {slot.room_id ?? "—"}</div>
                <div>Start: {formatDateTime(slot.start_at)}</div>
                <div>End: {formatDateTime(slot.end_at)}</div>
                <div>
                  Status:{" "}
                  <Badge tone={slot.status === "active" ? "success" : "warning"}>{slot.status}</Badge>
                </div>
                <Button
                  type="button"
                  variant="danger"
                  loading={deleteMutation.isPending}
                  onClick={() => deleteMutation.mutate(slot.id)}
                >
                  Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}
