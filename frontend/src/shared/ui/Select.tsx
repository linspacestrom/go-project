import { SelectHTMLAttributes } from "react";

type Option = {
  label: string;
  value: string;
};

type Props = SelectHTMLAttributes<HTMLSelectElement> & {
  label?: string;
  options: Option[];
};

export function Select({ label, options, style, ...props }: Props) {
  return (
    <label style={{ display: "grid", gap: 6 }}>
      {label ? <span style={{ fontSize: 14, fontWeight: 600 }}>{label}</span> : null}
      <select
        {...props}
        style={{
          border: "1px solid #94a3b8",
          borderRadius: 10,
          background: "white",
          padding: "10px 12px",
          ...style
        }}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}
