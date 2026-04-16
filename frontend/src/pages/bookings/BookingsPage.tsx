import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CreateBookingForm } from "@/features/bookings/create/ui/CreateBookingForm";
import { mainApi } from "@/shared/api/main-api";
import { Badge } from "@/shared/ui/Badge";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";
import { formatDateTime } from "@/shared/lib/date";

export function BookingsPage() {
  const queryClient = useQueryClient();
  const listQuery = useQuery({
    queryKey: ["bookings"],
    queryFn: () => mainApi.listBookings({ limit: 50, offset: 0 })
  });
  const cancelMutation = useMutation({
    mutationFn: (id: string) => mainApi.cancelBooking(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["bookings"] })
  });

  return (
    <div className="grid">
      <PageTitle title="Бронирования" subtitle="Создание, просмотр и отмена" />
      <CreateBookingForm />
      <Card title="Список моих броней">
        {listQuery.isLoading ? <Loader /> : null}
        {listQuery.isError ? <ErrorState message="Не удалось загрузить бронирования" /> : null}
        {listQuery.data && listQuery.data.items.length === 0 ? <EmptyState title="Бронирования отсутствуют" /> : null}
        <div className="grid">
          {listQuery.data?.items.map((booking) => (
            <Card key={booking.id} title={`Бронь ${booking.id.slice(0, 8)}...`}>
              <div className="grid">
                <div>Type: {booking.booking_type}</div>
                <div>Status: <Badge tone={booking.status === "active" ? "success" : "warning"}>{booking.status}</Badge></div>
                <div>Slot ID: {booking.slot_id ?? "—"}</div>
                <div>Room ID: {booking.room_id ?? "—"}</div>
                <div>Start: {formatDateTime(booking.start_at)}</div>
                <div>End: {formatDateTime(booking.end_at)}</div>
                <Button
                  type="button"
                  variant="danger"
                  loading={cancelMutation.isPending}
                  onClick={() => cancelMutation.mutate(booking.id)}
                >
                  Отменить
                </Button>
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}
