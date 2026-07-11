import { Link, Navigate } from "react-router-dom";
import { useAuth } from "../context/useAuth";

function ProfilePage() {
  const { user, logout } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  function getRole(role) {
    switch (role) {
      case "admin":
        return "Администратор";

      case "client":
        return "Покупатель";

      default:
        return role;
    }
  }

  function getRoleIcon(role) {
    switch (role) {
      case "admin":
        return "🛡️";

      case "client":
        return "🛒";

      default:
        return "👤";
    }
  }

  return (
    <div className="container">
      <div className="profile-card">

        <div className="profile-avatar">
          👤
        </div>

        <h1>{user.login}</h1>

        <p className="profile-role">
          {getRoleIcon(user.role)} {getRole(user.role)}
        </p>

        <div className="profile-info">

          <div className="profile-item">
            <span>👤 Логин</span>

            <strong>{user.login}</strong>
          </div>

          <div className="profile-item">
            <span>📧 Email</span>

            <strong>{user.email}</strong>
          </div>

          <div className="profile-item">
            <span>🛡️ Роль</span>

            <strong>{getRole(user.role)}</strong>
          </div>

        </div>

        <div className="profile-actions">

          <Link
            className="profile-btn"
            to="/cart"
          >
            🛒 Корзина
          </Link>

          <Link
            className="profile-btn"
            to="/orders"
          >
            📦 Мои заказы
          </Link>

          <Link
            className="profile-btn"
            to="/notifications"
          >
            🔔 Уведомления
          </Link>

          {user.role === "admin" && (
            <Link
              className="profile-btn"
              to="/admin"
            >
              ⚙️ Админ-панель
            </Link>
          )}

        </div>

        <button
          className="logout-btn"
          onClick={logout}
        >
          Выйти из аккаунта
        </button>

      </div>
    </div>
  );
}

export default ProfilePage;
