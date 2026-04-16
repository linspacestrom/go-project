import { FormEvent, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { Select } from "@/shared/ui/Select";

export function CreateBookingForm() {
  const queryClient = useQueryClient();
  const [bookingType, setBookingType] = useState<"room_only" | "room_with_mentor" | "mentor_call" | "event_seat">(
    "room_only"
  );
  const [slotId, setSlotId] = useState("");
  const [roomId, setRoomId] = useState("");
  const [startAt, setStartAt] = useState("");
  const [endAt, setEndAt] = useState("");
  const [meetingUrl, setMeetingUrl] = useState("");
  const [seatNumber, setSeatNumber] = useState("");
  const [idempotencyKey, setIdempotencyKey] = useState("");
  const [error, setError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () =>
      mainApi.createBooking({
        booking_type: bookingType,
        slot_id: slotId || undefined,
        room_id: roomId || undefined,
        start_at: new Date(startAt).toISOString(),
        end_at: new Date(endAt).toISOString(),
        meeting_url: meetingUrl || undefined,
        seat_number: seatNumber ? Number(seatNumber) : undefined,
        idempotency_key: idempotencyKey || undefined
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bookings"] });
      setError(null);
    },
    onError: (e) => {
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось создать бронирование");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!startAt || !endAt || new Date(startAt) >= new Date(endAt)) {
      setError("Проверьте интервал времени");
      return;
    }
    mutation.mutate();
  }

  return (
    <Card title="Создать бронирование">
      <form onSubmit={onSubmit} className="grid">
        <Select
          label="Тип бронирования"
          value={bookingType}
          onChange={(event) =>
            setBookingType(
              event.target.value as "room_only" | "room_with_mentor" | "mentor_call" | "event_seat"
            )
          }
          options={[
            { label: "Room only", value: "room_only" },
            { label: "Room with mentor", value: "room_with_mentor" },
            { label: "Mentor call", value: "mentor_call" },
            { label: "Event seat", value: "event_seat" }
          ]}
        />
        <Input label="Slot ID (optional)" value={slotId} onChange={(event) => setSlotId(event.target.value)} />
        <Input label="Room ID (optional)" value={roomId} onChange={(event) => setRoomId(event.target.value)} />
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
          label="Meeting URL (optional)"
          value={meetingUrl}
          onChange={(event) => setMeetingUrl(event.target.value)}
        />
        <Input
          label="Seat Number (optional)"
          type="number"
          value={seatNumber}
          onChange={(event) => setSeatNumber(event.target.value)}
        />
        <Input
          label="Idempotency Key (optional)"
          value={idempotencyKey}
          onChange={(event) => setIdempotencyKey(event.target.value)}
        />
        {error ? <ErrorState message={error} /> : null}
        <Button type="submit" loading={mutation.isPending}>
          Создать
        </Button>
      </form>
    </Card>
  );
}
