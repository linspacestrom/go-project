import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/features/auth/model/auth-context";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { Select } from "@/shared/ui/Select";

export function RegisterForm() {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"student" | "mentor" | "admin">("student");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError(null);
    try {
      await register({ fullName, email, password, role });
      navigate("/", { replace: true });
    } catch (e) {
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось выполнить регистрацию");
      }
    } finally {
      setPending(false);
    }
  }

  return (
    <Card title="Регистрация">
      <form onSubmit={onSubmit} className="grid">
        <Input
          label="ФИО"
          value={fullName}
          onChange={(event) => setFullName(event.target.value)}
          required
        />
        <Input
          label="Email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />
        <Input
          label="Пароль"
          type="password"
          value={password}
          minLength={8}
          onChange={(event) => setPassword(event.target.value)}
          required
        />
        <Select
          label="Роль"
          value={role}
          onChange={(event) => setRole(event.target.value as "student" | "mentor" | "admin")}
          options={[
            { label: "Student", value: "student" },
            { label: "Mentor", value: "mentor" },
            { label: "Admin", value: "admin" }
          ]}
        />
        {error ? <ErrorState message={error} /> : null}
        <Button type="submit" loading={pending}>
          Создать аккаунт
        </Button>
      </form>
    </Card>
  );
}
