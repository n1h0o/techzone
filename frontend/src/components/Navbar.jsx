import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";

import { useAuth } from "../context/AuthContext";
import { useCart } from "../context/CartContext";

function Navbar() {
  const { user, logout } = useAuth();
  const { cartCount, clearCart } = useCart();

  const navigate = useNavigate();

  const [menuOpen, setMenuOpen] = useState(false);

  function handleLogout() {
    logout();
    clearCart();
    setMenuOpen(false);
    navigate("/login");
  }

  function closeMenu() {
    setMenuOpen(false);
  }

  return (
    <>
      <header className="navbar">
        <NavLink
          className="logo"
          to="/"
          onClick={closeMenu}
        >
          🛍️ TechZone
        </NavLink>

        <button
          className="burger"
          onClick={() => setMenuOpen(!menuOpen)}
        >
          ☰
        </button>

        <div
          className={`navbar-menu ${
            menuOpen ? "active" : ""
          }`}
        >
          <div className="navbar-left">
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
                  className="cart-link"
                  onClick={closeMenu}
                >
                  🛒 Корзина

                  {cartCount > 0 && (
                    <span className="cart-badge">
                      {cartCount}
                    </span>
                  )}
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
              </>
            )}

            {user?.role === "admin" && (
              <NavLink
                to="/admin"
                onClick={closeMenu}
              >
                ⚙️ Админ
              </NavLink>
            )}
          </div>

          <div className="navbar-right">
            {!user ? (
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
            ) : (
              <>
                <NavLink
                  to="/profile"
                  onClick={closeMenu}
                >
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
        </div>
      </header>

      {menuOpen && (
        <div
          className="navbar-overlay"
          onClick={closeMenu}
        />
      )}
    </>
  );
}

export default Navbar;