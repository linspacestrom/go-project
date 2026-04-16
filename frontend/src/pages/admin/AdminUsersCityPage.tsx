import { UpdateUserCityForm } from "@/features/admin/manage-city/ui/UpdateUserCityForm";
import { PageTitle } from "@/shared/ui/PageTitle";

export function AdminUsersCityPage() {
  return (
    <div className="grid">
      <PageTitle
        title="Admin: User City"
        subtitle="Назначение города пользователю через админский сценарий"
      />
      <UpdateUserCityForm />
    </div>
  );
}
