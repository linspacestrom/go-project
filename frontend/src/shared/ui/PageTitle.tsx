export function PageTitle({
  title,
  subtitle
}: {
  title: string;
  subtitle?: string;
}) {
  return (
    <header style={{ marginBottom: 16 }}>
      <h1 style={{ margin: 0, fontSize: 28 }}>{title}</h1>
      {subtitle ? <p style={{ margin: "6px 0 0", color: "#475569" }}>{subtitle}</p> : null}
    </header>
  );
}
