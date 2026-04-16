import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";

export function RegisterMentorForm() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [cityId, setCityId] = useState("");
  const [description, setDescription] = useState("");
  const [title, setTitle] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () =>
      mainApi.registerMentor({
        email,
        password,
        full_name: fullName,
        city_id: cityId,
        description: description || undefined,
        title: title || undefined
      }),
    onSuccess: (data) => {
      setError(null);
      setSuccess(`Ментор ${data.full_name} создан`);
    },
    onError: (e) => {
      setSuccess(null);
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось создать ментора");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <Card title="Регистрация ментора">
      <form onSubmit={onSubmit} className="grid">
        <Input label="Email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
        <Input
          label="Пароль"
          type="password"
          value={password}
          minLength={8}
          onChange={(event) => setPassword(event.target.value)}
          required
        />
        <Input label="ФИО" value={fullName} onChange={(event) => setFullName(event.target.value)} required />
        <Input label="City ID" value={cityId} onChange={(event) => setCityId(event.target.value)} required />
        <Input
          label="Описание (optional)"
          value={description}
          onChange={(event) => setDescription(event.target.value)}
        />
        <Input label="Должность (optional)" value={title} onChange={(event) => setTitle(event.target.value)} />
        {error ? <ErrorState message={error} /> : null}
        {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
        <Button type="submit" loading={mutation.isPending}>
          Создать ментора
        </Button>
      </form>
    </Card>
  );
}
