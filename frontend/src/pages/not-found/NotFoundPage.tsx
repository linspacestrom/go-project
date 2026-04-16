import { Link } from "react-router-dom";
import { Card } from "@/shared/ui/Card";
import { PageTitle } from "@/shared/ui/PageTitle";

export function NotFoundPage() {
  return (
    <div className="grid">
      <PageTitle title="404" subtitle="Страница не найдена" />
      <Card>
        <p>Похоже, такого маршрута не существует.</p>
        <Link to="/">Перейти на главную</Link>
      </Card>
    </div>
  );
}
