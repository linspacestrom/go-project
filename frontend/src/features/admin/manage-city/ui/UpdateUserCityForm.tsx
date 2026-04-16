import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";

export function UpdateUserCityForm() {
  const [userId, setUserId] = useState("");
  const [cityId, setCityId] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () => mainApi.updateUserCity(userId, cityId),
    onSuccess: () => {
      setError(null);
      setSuccess("Город пользователя обновлен");
    },
    onError: (e) => {
      setSuccess(null);
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось обновить город пользователя");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <Card title="Смена города пользователя">
      <form onSubmit={onSubmit} className="grid">
        <Input label="User ID" value={userId} onChange={(event) => setUserId(event.target.value)} required />
        <Input label="City ID" value={cityId} onChange={(event) => setCityId(event.target.value)} required />
        {error ? <ErrorState message={error} /> : null}
        {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
        <Button type="submit" loading={mutation.isPending}>
          Обновить
        </Button>
      </form>
    </Card>
  );
}
