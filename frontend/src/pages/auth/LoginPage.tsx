import { Link } from "react-router-dom";
import { LoginForm } from "@/features/auth/ui/LoginForm";

export function LoginPage() {
  return (
    <div className="container" style={{ maxWidth: 480 }}>
      <LoginForm />
      <p>
        Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
      </p>
    </div>
  );
}
