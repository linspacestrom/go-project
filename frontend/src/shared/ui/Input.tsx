import { InputHTMLAttributes } from "react";

type Props = InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export function Input({ label, error, style, ...props }: Props) {
  return (
    <label style={{ display: "grid", gap: 6 }}>
      {label ? <span style={{ fontSize: 14, fontWeight: 600 }}>{label}</span> : null}
      <input
        {...props}
        style={{
          border: `1px solid ${error ? "#dc2626" : "#94a3b8"}`,
          borderRadius: 10,
          background: "white",
          padding: "10px 12px",
          ...style
        }}
      />
      {error ? <span style={{ color: "#dc2626", fontSize: 12 }}>{error}</span> : null}
    </label>
  );
}
