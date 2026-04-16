import { PropsWithChildren, ReactNode } from "react";

type Props = PropsWithChildren<{
  title?: string;
  actions?: ReactNode;
}>;

export function Card({ title, actions, children }: Props) {
  return (
    <section
      style={{
        border: "1px solid #cbd5e1",
        borderRadius: 16,
        padding: 16,
        background: "white",
        boxShadow: "0 2px 10px rgba(15, 23, 42, 0.03)"
      }}
    >
      {title || actions ? (
        <header
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: 12
          }}
        >
          <h3 style={{ margin: 0, fontSize: 18 }}>{title}</h3>
          {actions}
        </header>
      ) : null}
      {children}
    </section>
  );
}
