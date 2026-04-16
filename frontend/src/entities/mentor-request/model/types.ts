export type MentorRequest = {
  id: string;
  slot_id?: string;
  mentor_id?: string;
  mentee_id: string;
  city_id: string;
  skill_id?: string;
  request_type: "category" | "skill" | "direct_mentor" | "other";
  status: string;
  comment?: string;
};
