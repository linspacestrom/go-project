import { FormEvent, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";

export function MentorSkillSubscriptionForm() {
  const [skillId, setSkillId] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const subscribeMutation = useMutation({
    mutationFn: () => mainApi.subscribeMentorSkill(skillId),
    onSuccess: () => {
      setError(null);
      setSuccess("Подписка обновлена");
    },
    onError: (e) => {
      setSuccess(null);
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось оформить подписку");
      }
    }
  });

  const unsubscribeMutation = useMutation({
    mutationFn: () => mainApi.unsubscribeMentorSkill(skillId),
    onSuccess: () => {
      setError(null);
      setSuccess("Подписка отключена");
    },
    onError: (e) => {
      setSuccess(null);
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось отключить подписку");
      }
    }
  });

  function onSubscribe(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    subscribeMutation.mutate();
  }

  return (
    <Card title="Подписки на skills">
      <form onSubmit={onSubscribe} className="grid">
        <Input label="Skill ID" value={skillId} onChange={(event) => setSkillId(event.target.value)} required />
        {error ? <ErrorState message={error} /> : null}
        {success ? <div style={{ color: "#166534" }}>{success}</div> : null}
        <div style={{ display: "flex", gap: 8 }}>
          <Button type="submit" loading={subscribeMutation.isPending}>
            Подписаться
          </Button>
          <Button
            type="button"
            variant="danger"
            loading={unsubscribeMutation.isPending}
            onClick={() => unsubscribeMutation.mutate()}
          >
            Отписаться
          </Button>
        </div>
      </form>
    </Card>
  );
}
