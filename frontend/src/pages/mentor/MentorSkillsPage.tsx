import { MentorSkillSubscriptionForm } from "@/features/mentor/skills/ui/MentorSkillSubscriptionForm";
import { PageTitle } from "@/shared/ui/PageTitle";

export function MentorSkillsPage() {
  return (
    <div className="grid">
      <PageTitle title="Mentor Skills" subtitle="Управление подписками на skill-запросы" />
      <MentorSkillSubscriptionForm />
    </div>
  );
}
