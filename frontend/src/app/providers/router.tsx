import { ReactNode } from "react";
import { Navigate, createBrowserRouter } from "react-router-dom";
import { useAuth } from "@/features/auth/model/auth-context";
import { AppShell } from "@/widgets/app-shell/ui/AppShell";
import { Loader } from "@/shared/ui/Loader";
import { LoginPage } from "@/pages/auth/LoginPage";
import { RegisterPage } from "@/pages/auth/RegisterPage";
import { HomePage } from "@/pages/dashboard/HomePage";
import { ProfilePage } from "@/pages/profile/ProfilePage";
import { CitiesPage } from "@/pages/cities/CitiesPage";
import { HubsPage } from "@/pages/hubs/HubsPage";
import { RoomsPage } from "@/pages/rooms/RoomsPage";
import { SlotsPage } from "@/pages/slots/SlotsPage";
import { SlotDetailsPage } from "@/pages/slots/SlotDetailsPage";
import { BookingsPage } from "@/pages/bookings/BookingsPage";
import { MentorRequestsPage } from "@/pages/mentor-requests/MentorRequestsPage";
import { NotificationsPage } from "@/pages/notifications/NotificationsPage";
import { MentorSkillsPage } from "@/pages/mentor/MentorSkillsPage";
import { MentorSlotsPage } from "@/pages/mentor/MentorSlotsPage";
import { MentorRequestsQueuePage } from "@/pages/mentor/MentorRequestsPage";
import { RegisterMentorPage } from "@/pages/admin/RegisterMentorPage";
import { AdminCitiesPage } from "@/pages/admin/AdminCitiesPage";
import { AdminHubsPage } from "@/pages/admin/AdminHubsPage";
import { AdminRoomsPage } from "@/pages/admin/AdminRoomsPage";
import { AdminUsersCityPage } from "@/pages/admin/AdminUsersCityPage";
import { AdminAnalyticsPage } from "@/pages/admin/AdminAnalyticsPage";
import { NotFoundPage } from "@/pages/not-found/NotFoundPage";

function PublicOnly({ children }: { children: ReactNode }) {
  const { isAuthenticated, isBootstrapping } = useAuth();
  if (isBootstrapping) {
    return (
      <div className="container">
        <Loader />
      </div>
    );
  }
  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

function RequireAuth({ children }: { children: ReactNode }) {
  const { isAuthenticated, isBootstrapping } = useAuth();
  if (isBootstrapping) {
    return (
      <div className="container">
        <Loader />
      </div>
    );
  }
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function RequireRole({
  children,
  roles
}: {
  children: ReactNode;
  roles: Array<"student" | "mentor" | "admin">;
}) {
  const { role } = useAuth();
  if (!role || !roles.includes(role)) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: (
      <PublicOnly>
        <LoginPage />
      </PublicOnly>
    )
  },
  {
    path: "/register",
    element: (
      <PublicOnly>
        <RegisterPage />
      </PublicOnly>
    )
  },
  {
    path: "/",
    element: (
      <RequireAuth>
        <AppShell />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <HomePage /> },
      { path: "me", element: <ProfilePage /> },
      { path: "cities", element: <CitiesPage /> },
      { path: "hubs", element: <HubsPage /> },
      { path: "rooms", element: <RoomsPage /> },
      { path: "slots", element: <SlotsPage /> },
      { path: "slots/:id", element: <SlotDetailsPage /> },
      { path: "bookings", element: <BookingsPage /> },
      { path: "mentor-requests", element: <MentorRequestsPage /> },
      { path: "notifications", element: <NotificationsPage /> },
      {
        path: "mentor/skills",
        element: (
          <RequireRole roles={["mentor"]}>
            <MentorSkillsPage />
          </RequireRole>
        )
      },
      {
        path: "mentor/slots",
        element: (
          <RequireRole roles={["mentor"]}>
            <MentorSlotsPage />
          </RequireRole>
        )
      },
      {
        path: "mentor/requests",
        element: (
          <RequireRole roles={["mentor"]}>
            <MentorRequestsQueuePage />
          </RequireRole>
        )
      },
      {
        path: "admin/register-mentor",
        element: (
          <RequireRole roles={["admin"]}>
            <RegisterMentorPage />
          </RequireRole>
        )
      },
      {
        path: "admin/cities",
        element: (
          <RequireRole roles={["admin"]}>
            <AdminCitiesPage />
          </RequireRole>
        )
      },
      {
        path: "admin/hubs",
        element: (
          <RequireRole roles={["admin"]}>
            <AdminHubsPage />
          </RequireRole>
        )
      },
      {
        path: "admin/rooms",
        element: (
          <RequireRole roles={["admin"]}>
            <AdminRoomsPage />
          </RequireRole>
        )
      },
      {
        path: "admin/users/city",
        element: (
          <RequireRole roles={["admin"]}>
            <AdminUsersCityPage />
          </RequireRole>
        )
      },
      {
        path: "admin/analytics",
        element: (
          <RequireRole roles={["admin"]}>
            <AdminAnalyticsPage />
          </RequireRole>
        )
      },
      { path: "*", element: <NotFoundPage /> }
    ]
  },
  {
    path: "*",
    element: <NotFoundPage />
  }
]);
