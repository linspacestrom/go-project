import { useMutation, useQueryClient } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { Button } from "@/shared/ui/Button";

export function MentorRequestActions({ requestId }: { requestId: string }) {
  const queryClient = useQueryClient();
  const approveMutation = useMutation({
    mutationFn: () => mainApi.approveMentorRequest(requestId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mentor-own-requests"] });
      queryClient.invalidateQueries({ queryKey: ["mentor-requests"] });
    }
  });
  const rejectMutation = useMutation({
    mutationFn: () => mainApi.rejectMentorRequest(requestId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mentor-own-requests"] });
      queryClient.invalidateQueries({ queryKey: ["mentor-requests"] });
    }
  });

  return (
    <div style={{ display: "flex", gap: 8 }}>
      <Button type="button" loading={approveMutation.isPending} onClick={() => approveMutation.mutate()}>
        Approve
      </Button>
      <Button
        type="button"
        variant="danger"
        loading={rejectMutation.isPending}
        onClick={() => rejectMutation.mutate()}
      >
        Reject
      </Button>
    </div>
  );
}
