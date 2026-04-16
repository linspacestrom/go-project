export function Loader({ label = "Загрузка..." }: { label?: string }) {
  return (
    <div
      style={{
        padding: 16,
        borderRadius: 12,
        border: "1px solid #cbd5e1",
        background: "white",
        color: "#334155"
      }}
    >
      {label}
    </div>
  );
}
