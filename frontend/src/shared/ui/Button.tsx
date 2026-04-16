import { ButtonHTMLAttributes, PropsWithChildren } from "react";

type Variant = "primary" | "secondary" | "danger" | "ghost";

type Props = PropsWithChildren<
  ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: Variant;
    loading?: boolean;
  }
>;

const colors: Record<Variant, string> = {
  primary: "#0f766e",
  secondary: "#1d4ed8",
  danger: "#b91c1c",
  ghost: "transparent"
};

export function Button({ children, variant = "primary", loading, style, disabled, ...props }: Props) {
  return (
    <button
      {...props}
      disabled={disabled || loading}
      style={{
        border: variant === "ghost" ? "1px solid #94a3b8" : "none",
        borderRadius: 10,
        padding: "10px 14px",
        cursor: "pointer",
        backgroundColor: colors[variant],
        color: variant === "ghost" ? "#0f172a" : "white",
        opacity: disabled || loading ? 0.7 : 1,
        fontWeight: 600,
        ...style
      }}
    >
      {loading ? "Загрузка..." : children}
    </button>
  );
}
