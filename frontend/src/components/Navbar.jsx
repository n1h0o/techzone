import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function Navbar() {
  const { user, logout } = useAuth();

  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/");
  }

  return (
    <header className="navbar">
      <div className="navbar-left">
        <Link className="logo" to="/">
          TechZone
        </Link>

        <Link to="/products">
          Каталог
        </Link>

        {user && (
          <>
            <Link to="/cart">
              Корзина
            </Link>

            <Link to="/orders">
              Заказы
            </Link>

            <Link to="/notifications">
              Уведомления
            </Link>
          </>
        )}

        {user?.role === "admin" && (
          <Link to="/admin">
            Админ-панель
          </Link>
        )}
      </div>

      <div className="navbar-right">
        {!user ? (
          <>
            <Link to="/login">
              Вход
            </Link>

            <Link to="/register">
              Регистрация
            </Link>
          </>
        ) : (
          <>
            <span className="username">
              👤 {user.login}
            </span>

            <Link to="/profile">
              Профиль
            </Link>

            <button
              className="logout-btn"
              onClick={handleLogout}
            >
              Выйти
            </button>
          </>
        )}
      </div>
    </header>
  );
}

export default Navbar;