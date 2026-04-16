import { RouterProvider } from "react-router-dom";
import { QueryProvider } from "@/app/providers/query-client";
import { router } from "@/app/providers/router";
import { AuthProvider } from "@/features/auth/model/auth-context";

export function App() {
  return (
    <QueryProvider>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </QueryProvider>
  );
}
