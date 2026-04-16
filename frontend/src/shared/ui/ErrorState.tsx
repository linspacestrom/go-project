import { ReactNode } from "react";

export function ErrorState({
  title = "Ошибка",
  message,
  action
}: {
  title?: string;
  message?: string;
  action?: ReactNode;
}) {
  return (
    <div
      style={{
        border: "1px solid #fecaca",
        background: "#fef2f2",
        color: "#7f1d1d",
        borderRadius: 12,
        padding: 14,
        display: "grid",
        gap: 8
      }}
    >
      <strong>{title}</strong>
      {message ? <span>{message}</span> : null}
      {action}
    </div>
  );
}
