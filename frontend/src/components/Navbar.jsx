import { Link } from "react-router-dom";

function Navbar() {
  return (
    <nav>
      <Link to="/products">
        Товары
      </Link>

      <Link to="/cart">
        Корзина
      </Link>

      <Link to="/orders">
        Заказы
      </Link>

      <Link to="/login">
        Вход
      </Link>

      <Link to="/register">
      Регистрация
      </Link>
    </nav>
  );
}

export default Navbar;