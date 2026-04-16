import { useQuery } from "@tanstack/react-query";
import { mainApi } from "@/shared/api/main-api";
import { Card } from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Loader } from "@/shared/ui/Loader";
import { PageTitle } from "@/shared/ui/PageTitle";

export function AdminAnalyticsPage() {
  const businessQuery = useQuery({
    queryKey: ["analytics-business"],
    queryFn: () => mainApi.getBusinessAnalytics()
  });
  const technicalQuery = useQuery({
    queryKey: ["analytics-technical"],
    queryFn: () => mainApi.getTechnicalAnalytics()
  });

  return (
    <div className="grid">
      <PageTitle title="Admin Analytics" subtitle="Бизнес и технические метрики платформы" />
      {(businessQuery.isLoading || technicalQuery.isLoading) ? <Loader /> : null}
      {(businessQuery.isError || technicalQuery.isError) ? (
        <ErrorState message="Не удалось загрузить аналитику" />
      ) : null}
      {businessQuery.data ? (
        <Card title="Business analytics">
          <pre style={{ margin: 0 }}>{JSON.stringify(businessQuery.data, null, 2)}</pre>
        </Card>
      ) : null}
      {technicalQuery.data ? (
        <Card title="Technical analytics">
          <pre style={{ margin: 0 }}>{JSON.stringify(technicalQuery.data, null, 2)}</pre>
        </Card>
      ) : null}
    </div>
  );
}
