import { Link } from "react-router-dom";
import { RegisterForm } from "@/features/auth/ui/RegisterForm";

export function RegisterPage() {
  return (
    <div className="container" style={{ maxWidth: 520 }}>
      <RegisterForm />
      <p>
        Уже есть аккаунт? <Link to="/login">Войти</Link>
      </p>
    </div>
  );
}
