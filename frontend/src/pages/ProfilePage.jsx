import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function Navbar() {
  const { user, logout } = useAuth();

  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <nav
      style={{
        display: "flex",
        gap: "20px",
        padding: "20px",
        alignItems: "center",
      }}
    >
      <Link to="/products">Товары</Link>

      {user && (
        <>
          <Link to="/cart">Корзина</Link>

          <Link to="/orders">Заказы</Link>

          <Link to="/notifications">
            Уведомления
          </Link>
        </>
      )}

      {user?.role === "admin" && (
        <Link to="/admin/products">
          Добавить товар
        </Link>
      )}

      <div
        style={{
          marginLeft: "auto",
          display: "flex",
          gap: "15px",
          alignItems: "center",
        }}
      >
        {!user ? (
          <>
            <Link to="/login">Вход</Link>

            <Link to="/register">
              Регистрация
            </Link>
          </>
        ) : (
          <>
            <span>
              👤 {user.login}
            </span>

            <Link to="/profile">
              Профиль
            </Link>

            <button onClick={handleLogout}>
              Выйти
            </button>
          </>
        )}
      </div>
    </nav>
  );
}

export default Navbar;