export function EmptyState({
  title = "Пока пусто",
  description
}: {
  title?: string;
  description?: string;
}) {
  return (
    <div
      style={{
        border: "1px dashed #94a3b8",
        borderRadius: 12,
        background: "#f8fafc",
        color: "#334155",
        padding: 14,
        display: "grid",
        gap: 4
      }}
    >
      <strong>{title}</strong>
      {description ? <span>{description}</span> : null}
    </div>
  );
}
