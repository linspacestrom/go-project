export function Badge({
  children,
  tone = "neutral"
}: {
  children: string;
  tone?: "neutral" | "success" | "danger" | "warning";
}) {
  const palette = {
    neutral: { bg: "#e2e8f0", fg: "#0f172a" },
    success: { bg: "#dcfce7", fg: "#166534" },
    danger: { bg: "#fee2e2", fg: "#991b1b" },
    warning: { bg: "#fef3c7", fg: "#92400e" }
  }[tone];

  return (
    <span
      style={{
        display: "inline-block",
        borderRadius: 999,
        padding: "3px 8px",
        fontSize: 12,
        fontWeight: 700,
        background: palette.bg,
        color: palette.fg
      }}
    >
      {children}
    </span>
  );
}
