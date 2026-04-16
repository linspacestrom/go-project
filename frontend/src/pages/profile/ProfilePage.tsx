import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useAuth } from "@/features/auth/model/auth-context";
import { ApiError } from "@/shared/api/types";
import { mainApi } from "@/shared/api/main-api";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { PageTitle } from "@/shared/ui/PageTitle";

export function ProfilePage() {
  const { user, refreshProfile } = useAuth();
  const [fullName, setFullName] = useState(user?.full_name ?? "");
  const [birthDate, setBirthDate] = useState(user?.birth_date?.slice(0, 10) ?? "");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () =>
      mainApi.updateMe({
        full_name: fullName || undefined,
        birth_date: birthDate ? new Date(birthDate).toISOString() : undefined
      }),
    onSuccess: async () => {
      await refreshProfile();
      setError(null);
      setSuccess("Профиль обновлен");
    },
    onError: (e) => {
      setSuccess(null);
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось обновить профиль");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <div className="grid">
      <PageTitle title="Мой профиль" subtitle="Данные текущего пользователя" />
      <Card>
        <div className="grid" style={{ gridTemplateColumns: "repeat(2, minmax(180px, 1fr))" }}>
          <div>ID: {user?.id}</div>
          <div>Email: {user?.email}</div>
          <div>Role: {user?.role}</div>
          <div>City ID: {user?.city_id ?? "—"}</div>
        </div>
      </Card>
      <Card title="Редактирование профиля">
        <form onSubmit={onSubmit} className="grid" style={{ maxWidth: 460 }}>
          <Input label="ФИО" value={fullName} onChange={(event) => setFullName(event.target.value)} />
          <Input
            label="Дата рождения"
            type="date"
            value={birthDate}
            onChange={(event) => setBirthDate(event.target.value)}
          />
          {error ? <ErrorState message={error} /> : null}
          {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
          <Button type="submit" loading={mutation.isPending}>
            Сохранить
          </Button>
        </form>
      </Card>
    </div>
  );
}
