import { useAuth } from "@/features/auth/model/auth-context";
import { Card } from "@/shared/ui/Card";
import { PageTitle } from "@/shared/ui/PageTitle";

export function HomePage() {
  const { user } = useAuth();

  return (
    <div className="grid">
      <PageTitle
        title="Панель платформы Студент и Т"
        subtitle="Ролевой интерфейс для студентов, менторов и администраторов"
      />
      <Card title="Текущая сессия">
        <p style={{ margin: 0 }}>
          Вы вошли как <strong>{user?.full_name}</strong> ({user?.role})
        </p>
        <p style={{ marginBottom: 0, color: "#475569" }}>
          Используй боковое меню для навигации по сценариям платформы.
        </p>
      </Card>
    </div>
  );
}
