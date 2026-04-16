export type UserRole = "student" | "mentor" | "admin";

export type User = {
  id: string;
  full_name: string;
  birth_date?: string;
  email: string;
  role: UserRole;
  city_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};
