import { loadAuthSession, saveAuthSession, clearAuthSession } from "@/shared/api/auth-storage";
import { createHttpClient } from "@/shared/api/http-client";
import { ListResponse, PaginationQuery, TokenPair } from "@/shared/api/types";
import { env } from "@/shared/config/env";
import { emitAuthExpired } from "@/shared/lib/auth-events";
import { User } from "@/entities/user/model/types";
import { City } from "@/entities/city/model/types";
import { Hub } from "@/entities/hub/model/types";
import { Room } from "@/entities/room/model/types";
import { Slot } from "@/entities/slot/model/types";
import { Booking } from "@/entities/booking/model/types";
import { MentorRequest } from "@/entities/mentor-request/model/types";

let refreshPromise: Promise<string | null> | null = null;

type RegisterPayload = {
  email: string;
  password: string;
  full_name: string;
  birth_date?: string;
  university?: string;
  course?: number;
  degree_type?: string;
  role?: string;
};

type LoginPayload = {
  email: string;
  password: string;
};

type RegisterResponse = {
  user: User;
  tokens: TokenPair;
};

export type SlotFilter = PaginationQuery & {
  city_id?: string;
  hub_id?: string;
  room_id?: string;
  mentor_id?: string;
  status?: string;
  type?: string;
  start_from?: string;
  end_to?: string;
};

export type BookingFilter = PaginationQuery & {
  status?: string;
  city_id?: string;
  room_id?: string;
  slot_id?: string;
};

export type MentorRequestFilter = PaginationQuery & {
  status?: string;
  skill_id?: string;
  city_id?: string;
  mentor_id?: string;
};

type UpdateMePayload = {
  full_name?: string;
  birth_date?: string;
};

export type CreateBookingPayload = {
  slot_id?: string;
  room_id?: string;
  booking_type: "room_only" | "room_with_mentor" | "mentor_call" | "event_seat";
  start_at: string;
  end_at: string;
  meeting_url?: string;
  seat_number?: number;
  idempotency_key?: string;
};

export type CreateMentorRequestPayload = {
  slot_id?: string;
  mentor_id?: string;
  skill_id?: string;
  request_type: "category" | "skill" | "direct_mentor" | "other";
  comment?: string;
};

export type CreateMentorSlotPayload = {
  room_id?: string;
  city_id: string;
  type: "online" | "offline";
  start_at: string;
  end_at: string;
  meeting_url?: string;
  status?: string;
  capacity?: number;
};

export type UpdateMentorSlotPayload = {
  room_id?: string;
  start_at?: string;
  end_at?: string;
  meeting_url?: string;
  status?: string;
  capacity?: number;
};

export type RegisterMentorPayload = {
  email: string;
  password: string;
  full_name: string;
  city_id: string;
  description?: string;
  title?: string;
};

export type CreateCityPayload = {
  name: string;
  is_active?: boolean;
};

export type UpdateCityPayload = {
  name?: string;
  is_active?: boolean;
};

export type CreateHubPayload = {
  city_id: string;
  name: string;
  address: string;
  is_active?: boolean;
};

export type UpdateHubPayload = {
  name?: string;
  address?: string;
  is_active?: boolean;
};

export type CreateRoomPayload = {
  hub_id: string;
  name: string;
  description?: string;
  room_type?: string;
  capacity: number;
  is_active?: boolean;
};

export type UpdateRoomPayload = {
  name?: string;
  description?: string;
  room_type?: string;
  capacity?: number;
  is_active?: boolean;
};

type AnalyticsBusinessResponse = {
  mentor_approvals: number;
  mentor_rejections: number;
  active_bookings: number;
  mentor_requests_by_skill: number;
  mentor_requests_by_category: number;
  user_cancels_after_mentor_approval: number;
};

type AnalyticsTechnicalResponse = {
  booking_conflict_count: number;
  failed_outbox_events: number;
  outbox_lag_count: number;
};

export async function refreshAccessToken(): Promise<string | null> {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async () => {
    const session = loadAuthSession();
    if (!session?.refreshToken) {
      return null;
    }

    try {
      const response = await fetch(`${env.mainApiUrl}/api/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json"
        },
        body: JSON.stringify({ refresh_token: session.refreshToken })
      });
      if (!response.ok) {
        clearAuthSession();
        emitAuthExpired();
        return null;
      }
      const tokens = (await response.json()) as TokenPair;
      saveAuthSession({
        accessToken: tokens.access_token,
        refreshToken: tokens.refresh_token
      });

      return tokens.access_token;
    } catch {
      clearAuthSession();
      emitAuthExpired();
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

const client = createHttpClient({
  baseUrl: env.mainApiUrl,
  getAccessToken: () => loadAuthSession()?.accessToken ?? null,
  refreshAccessToken,
  onAuthFailure: () => {
    clearAuthSession();
    emitAuthExpired();
  }
});

export const mainApi = {
  register(payload: RegisterPayload) {
    return client.request<RegisterResponse, RegisterPayload>("/api/v1/auth/register", {
      method: "POST",
      body: payload,
      auth: false,
      retryOnUnauthorized: false
    });
  },
  login(payload: LoginPayload) {
    return client.request<TokenPair, LoginPayload>("/api/v1/auth/login", {
      method: "POST",
      body: payload,
      auth: false,
      retryOnUnauthorized: false
    });
  },
  logout(refreshToken: string) {
    return client.request<void, { refresh_token: string }>("/api/v1/auth/logout", {
      method: "POST",
      body: { refresh_token: refreshToken },
      auth: false,
      retryOnUnauthorized: false
    });
  },

  getMe() {
    return client.request<User>("/api/v1/me");
  },
  updateMe(payload: UpdateMePayload) {
    return client.request<User, UpdateMePayload>("/api/v1/me", {
      method: "PATCH",
      body: payload
    });
  },

  listCities(query: PaginationQuery) {
    return client.request<ListResponse<City>>("/api/v1/cities", { query });
  },
  listHubsByCity(cityId: string, query: PaginationQuery) {
    return client.request<ListResponse<Hub>>(`/api/v1/cities/${cityId}/hubs`, { query });
  },
  listRoomsByCity(cityId: string, query: PaginationQuery) {
    return client.request<ListResponse<Room>>(`/api/v1/cities/${cityId}/rooms`, { query });
  },
  listSlots(filter: SlotFilter) {
    return client.request<ListResponse<Slot>>("/api/v1/slots", { query: filter });
  },
  getSlotById(id: string) {
    return client.request<Slot>(`/api/v1/slots/${id}`);
  },

  createBooking(payload: CreateBookingPayload) {
    return client.request<Booking, CreateBookingPayload>("/api/v1/bookings", {
      method: "POST",
      body: payload
    });
  },
  cancelBooking(id: string) {
    return client.request<Booking>(`/api/v1/bookings/${id}`, {
      method: "DELETE"
    });
  },
  listBookings(filter: BookingFilter) {
    return client.request<ListResponse<Booking>>("/api/v1/bookings", { query: filter });
  },

  createMentorRequest(payload: CreateMentorRequestPayload) {
    return client.request<MentorRequest, CreateMentorRequestPayload>("/api/v1/mentor-requests", {
      method: "POST",
      body: payload
    });
  },
  listMentorRequests(filter: MentorRequestFilter) {
    return client.request<ListResponse<MentorRequest>>("/api/v1/mentor-requests", {
      query: filter
    });
  },

  subscribeMentorSkill(skillId: string) {
    return client.request<void, { skill_id: string }>("/api/v1/mentor/skills/subscribe", {
      method: "POST",
      body: { skill_id: skillId }
    });
  },
  unsubscribeMentorSkill(skillId: string) {
    return client.request<void>(`/api/v1/mentor/skills/subscribe/${skillId}`, {
      method: "DELETE"
    });
  },
  listMentorSlots() {
    return client.request<ListResponse<Slot>>("/api/v1/mentor/slots");
  },
  createMentorSlot(payload: CreateMentorSlotPayload) {
    return client.request<Slot, CreateMentorSlotPayload>("/api/v1/mentor/slots", {
      method: "POST",
      body: payload
    });
  },
  updateMentorSlot(id: string, payload: UpdateMentorSlotPayload) {
    return client.request<Slot, UpdateMentorSlotPayload>(`/api/v1/mentor/slots/${id}`, {
      method: "PATCH",
      body: payload
    });
  },
  deleteMentorSlot(id: string) {
    return client.request<void>(`/api/v1/mentor/slots/${id}`, {
      method: "DELETE"
    });
  },
  listMentorOwnRequests(filter: MentorRequestFilter) {
    return client.request<ListResponse<MentorRequest>>("/api/v1/mentor/requests", {
      query: filter
    });
  },
  approveMentorRequest(id: string) {
    return client.request<MentorRequest>(`/api/v1/mentor/requests/${id}/approve`, {
      method: "POST"
    });
  },
  rejectMentorRequest(id: string) {
    return client.request<MentorRequest>(`/api/v1/mentor/requests/${id}/reject`, {
      method: "POST"
    });
  },

  registerMentor(payload: RegisterMentorPayload) {
    return client.request<User, RegisterMentorPayload>("/api/v1/admin/mentors/register_mentor", {
      method: "POST",
      body: payload
    });
  },
  createCity(payload: CreateCityPayload) {
    return client.request<City, CreateCityPayload>("/api/v1/admin/cities", {
      method: "POST",
      body: payload
    });
  },
  updateCity(id: string, payload: UpdateCityPayload) {
    return client.request<City, UpdateCityPayload>(`/api/v1/admin/cities/${id}`, {
      method: "PATCH",
      body: payload
    });
  },
  createHub(payload: CreateHubPayload) {
    return client.request<Hub, CreateHubPayload>("/api/v1/admin/hubs", {
      method: "POST",
      body: payload
    });
  },
  updateHub(id: string, payload: UpdateHubPayload) {
    return client.request<Hub, UpdateHubPayload>(`/api/v1/admin/hubs/${id}`, {
      method: "PATCH",
      body: payload
    });
  },
  createRoom(payload: CreateRoomPayload) {
    return client.request<Room, CreateRoomPayload>("/api/v1/admin/rooms", {
      method: "POST",
      body: payload
    });
  },
  updateRoom(id: string, payload: UpdateRoomPayload) {
    return client.request<Room, UpdateRoomPayload>(`/api/v1/admin/rooms/${id}`, {
      method: "PATCH",
      body: payload
    });
  },
  deleteRoom(id: string) {
    return client.request<void>(`/api/v1/admin/rooms/${id}`, {
      method: "DELETE"
    });
  },
  updateUserCity(userId: string, cityId: string) {
    return client.request<void, { city_id: string }>(`/api/v1/admin/users/${userId}/city`, {
      method: "PATCH",
      body: { city_id: cityId }
    });
  },
  getBusinessAnalytics() {
    return client.request<AnalyticsBusinessResponse>("/api/v1/admin/analytics/business");
  },
  getTechnicalAnalytics() {
    return client.request<AnalyticsTechnicalResponse>("/api/v1/admin/analytics/technical");
  }
};
