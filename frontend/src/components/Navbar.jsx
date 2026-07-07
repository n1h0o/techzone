import { NavLink, useNavigate } from "react-router-dom";

import { useAuth } from "../context/AuthContext";
import { useCart } from "../context/CartContext";

function Navbar() {
  const { user, logout } = useAuth();
  const { cartCount, clearCart } = useCart();

  const navigate = useNavigate();

  function handleLogout() {
    logout();
    clearCart();
    navigate("/login");
  }

  return (
    <header className="navbar">

      <div className="navbar-left">
        <NavLink className="logo" to="/">
          🛍️ TechZone
        </NavLink>

        <NavLink to="/products">
          📦 Каталог
        </NavLink>

        {user && (
          <>
            <NavLink
              to="/cart"
              className="cart-link"
            >
              🛒 Корзина

              {cartCount > 0 && (
                <span className="cart-badge">
                  {cartCount}
                </span>
              )}
            </NavLink>

            <NavLink to="/orders">
              📋 Заказы
            </NavLink>

            <NavLink to="/notifications">
              🔔 Уведомления
            </NavLink>
          </>
        )}

        {user?.role === "admin" && (
          <NavLink to="/admin">
            ⚙️ Админ
          </NavLink>
        )}
      </div>

      <div className="navbar-right">

        {!user ? (
          <>
            <NavLink to="/login">
              Вход
            </NavLink>

            <NavLink to="/register">
              Регистрация
            </NavLink>
          </>
        ) : (
          <>
            <NavLink to="/profile">
              👤 {user.login}
            </NavLink>

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