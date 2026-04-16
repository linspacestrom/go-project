import { RegisterMentorForm } from "@/features/admin/register-mentor/ui/RegisterMentorForm";
import { PageTitle } from "@/shared/ui/PageTitle";

export function RegisterMentorPage() {
  return (
    <div className="grid">
      <PageTitle title="Admin: Register Mentor" subtitle="Создание mentor-профиля администратором" />
      <RegisterMentorForm />
    </div>
  );
}
