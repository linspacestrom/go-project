export type AuthSession = {
  accessToken: string;
  refreshToken: string;
};

const authStorageKey = "student_t_auth_session";

export function loadAuthSession(): AuthSession | null {
  try {
    const raw = localStorage.getItem(authStorageKey);
    if (!raw) {
      return null;
    }
    const parsed = JSON.parse(raw) as AuthSession;
    if (!parsed.accessToken || !parsed.refreshToken) {
      return null;
    }

    return parsed;
  } catch {
    return null;
  }
}

export function saveAuthSession(session: AuthSession) {
  localStorage.setItem(authStorageKey, JSON.stringify(session));
}

export function clearAuthSession() {
  localStorage.removeItem(authStorageKey);
}
