import {
  createContext,
  PropsWithChildren,
  useContext,
  useEffect,
  useMemo,
  useState
} from "react";
import { User, UserRole } from "@/entities/user/model/types";
import { clearAuthSession, loadAuthSession, saveAuthSession } from "@/shared/api/auth-storage";
import { mainApi } from "@/shared/api/main-api";
import { subscribeAuthExpired } from "@/shared/lib/auth-events";

type RegisterInput = {
  email: string;
  password: string;
  fullName: string;
  role?: UserRole;
};

type AuthContextValue = {
  user: User | null;
  role: UserRole | null;
  isAuthenticated: boolean;
  isBootstrapping: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => Promise<void>;
  refreshProfile: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [user, setUser] = useState<User | null>(null);
  const [isBootstrapping, setIsBootstrapping] = useState(true);

  async function refreshProfile() {
    const profile = await mainApi.getMe();
    setUser(profile);
  }

  async function login(email: string, password: string) {
    const tokens = await mainApi.login({ email, password });
    saveAuthSession({
      accessToken: tokens.access_token,
      refreshToken: tokens.refresh_token
    });
    await refreshProfile();
  }

  async function register(input: RegisterInput) {
    const response = await mainApi.register({
      email: input.email,
      password: input.password,
      full_name: input.fullName,
      role: input.role
    });
    saveAuthSession({
      accessToken: response.tokens.access_token,
      refreshToken: response.tokens.refresh_token
    });
    setUser(response.user);
  }

  async function logout() {
    const session = loadAuthSession();
    try {
      if (session?.refreshToken) {
        await mainApi.logout(session.refreshToken);
      }
    } finally {
      clearAuthSession();
      setUser(null);
    }
  }

  useEffect(() => {
    const unsubscribe = subscribeAuthExpired(() => {
      clearAuthSession();
      setUser(null);
    });
    const session = loadAuthSession();
    if (!session?.accessToken) {
      setIsBootstrapping(false);
      return unsubscribe;
    }

    refreshProfile()
      .catch(() => {
        clearAuthSession();
        setUser(null);
      })
      .finally(() => setIsBootstrapping(false));

    return unsubscribe;
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      role: user?.role ?? null,
      isAuthenticated: Boolean(user),
      isBootstrapping,
      login,
      register,
      logout,
      refreshProfile
    }),
    [user, isBootstrapping]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}
