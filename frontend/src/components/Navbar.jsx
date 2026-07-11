import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";

import { useAuth } from "../context/useAuth";
import { useCart } from "../context/useCart";

function Navbar() {
  const { user, logout } = useAuth();
  const { cartCount, clearCart } = useCart();

  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  function closeMenu() {
    setMenuOpen(false);
  }

  function handleLogout() {
    logout();
    clearCart();
    closeMenu();
    navigate("/login");
  }

  return (
    <>
      <header className="navbar">
        <NavLink
          to="/"
          className="logo"
          onClick={closeMenu}
        >
          🛍️ TechZone
        </NavLink>

        <nav className="navbar-desktop">
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
        </nav>

        <div className="navbar-user">
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

        <button
          className="burger"
          onClick={() =>
            setMenuOpen(!menuOpen)
          }
        >
          {menuOpen ? "✕" : "☰"}
        </button>
      </header>

      <aside
        className={`mobile-menu ${
          menuOpen ? "open" : ""
        }`}
      >
        <NavLink
          to="/products"
          onClick={closeMenu}
        >
          📦 Каталог
        </NavLink>

        {user && (
          <>
            <NavLink
              to="/cart"
              onClick={closeMenu}
            >
              🛒 Корзина
            </NavLink>

            <NavLink
              to="/orders"
              onClick={closeMenu}
            >
              📋 Заказы
            </NavLink>

            <NavLink
              to="/notifications"
              onClick={closeMenu}
            >
              🔔 Уведомления
            </NavLink>

            <NavLink
              to="/profile"
              onClick={closeMenu}
            >
              👤 Профиль
            </NavLink>

            {user.role === "admin" && (
              <NavLink
                to="/admin"
                onClick={closeMenu}
              >
                ⚙️ Админ
              </NavLink>
            )}

            <button
              className="logout-btn"
              onClick={handleLogout}
            >
              Выйти
            </button>
          </>
        )}

        {!user && (
          <>
            <NavLink
              to="/login"
              onClick={closeMenu}
            >
              Вход
            </NavLink>

            <NavLink
              to="/register"
              onClick={closeMenu}
            >
              Регистрация
            </NavLink>
          </>
        )}
      </aside>

      {menuOpen && (
        <div
          className="menu-backdrop"
          onClick={closeMenu}
        />
      )}
    </>
  );
}

export default Navbar;
