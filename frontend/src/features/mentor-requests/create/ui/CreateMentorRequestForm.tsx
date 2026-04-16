import { FormEvent, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { ApiError } from "@/shared/api/types";
import { Button } from "@/shared/ui/Button";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { Select } from "@/shared/ui/Select";

export function CreateMentorRequestForm() {
  const queryClient = useQueryClient();
  const [requestType, setRequestType] = useState<"category" | "skill" | "direct_mentor" | "other">("category");
  const [slotId, setSlotId] = useState("");
  const [mentorId, setMentorId] = useState("");
  const [skillId, setSkillId] = useState("");
  const [comment, setComment] = useState("");
  const [error, setError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () =>
      mainApi.createMentorRequest({
        request_type: requestType,
        slot_id: slotId || undefined,
        mentor_id: mentorId || undefined,
        skill_id: skillId || undefined,
        comment: comment || undefined
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mentor-requests"] });
      setError(null);
    },
    onError: (e) => {
      if (e instanceof ApiError) {
        setError(e.message);
      } else {
        setError("Не удалось создать mentor request");
      }
    }
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <Card title="Создать mentor request">
      <form onSubmit={onSubmit} className="grid">
        <Select
          label="Request Type"
          value={requestType}
          onChange={(event) =>
            setRequestType(event.target.value as "category" | "skill" | "direct_mentor" | "other")
          }
          options={[
            { label: "Category", value: "category" },
            { label: "Skill", value: "skill" },
            { label: "Direct mentor", value: "direct_mentor" },
            { label: "Other", value: "other" }
          ]}
        />
        <Input label="Slot ID (optional)" value={slotId} onChange={(event) => setSlotId(event.target.value)} />
        <Input
          label="Mentor ID (optional)"
          value={mentorId}
          onChange={(event) => setMentorId(event.target.value)}
        />
        <Input label="Skill ID (optional)" value={skillId} onChange={(event) => setSkillId(event.target.value)} />
        <Input label="Comment (optional)" value={comment} onChange={(event) => setComment(event.target.value)} />
        {error ? <ErrorState message={error} /> : null}
        <Button type="submit" loading={mutation.isPending}>
          Отправить
        </Button>
      </form>
    </Card>
  );
}
