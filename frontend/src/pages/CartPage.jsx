import { useCallback, useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";

import api from "../api/api";
import { useCart } from "../context/useCart";

function CartPage() {
  const [items, setItems] = useState([]);

  const navigate = useNavigate();

  const {
    clearCart,
    refreshCart,
  } = useCart();

  const loadCart = useCallback(async () => {
    try {
      const res = await api.get("/cart");

      setItems(res.data.items || []);

      refreshCart();
    } catch (err) {
      console.error(err);
      toast.error("Не удалось загрузить корзину");
    }
  }, [refreshCart]);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void loadCart();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [loadCart]);

  async function removeItem(item) {
    try {
      await api.delete(`/cart/items/${item.id}`);

      await refreshCart();

      toast.success("Товар удалён из корзины");

      await loadCart();
    } catch (err) {
      console.error(err);
      toast.error("Не удалось удалить товар");
    }
  }

  async function createOrder() {
    try {
      const res = await api.post("/orders");

      clearCart();

      toast.success(
        `Заказ №${res.data.order_id} успешно создан 🎉`
      );

      navigate("/orders");
    } catch (err) {
      console.error(err);

      if (err.response?.data) {
        toast.error(err.response.data);
      } else {
        toast.error("Ошибка создания заказа");
      }
    }
  }

  const totalPrice = useMemo(() => {
    return items.reduce(
      (sum, item) =>
        sum + item.price * item.quantity,
      0
    );
  }, [items]);

  return (
    <div className="container">
      <h1>🛒 Корзина</h1>

      {items.length === 0 ? (
        <div className="card">
          <h2>Корзина пуста</h2>

          <p>
            Добавьте товары из каталога.
          </p>
        </div>
      ) : (
        <>
          <div className="products-grid">
            {items.map((item) => (
              <div
                key={item.id}
                className="product-card"
              >
                <h2>{item.name}</h2>

                <p>
                  Цена:{" "}
                  <strong>
                    {Number(item.price).toLocaleString(
                      "ru-RU"
                    )}{" "}
                    ₽
                  </strong>
                </p>

                <p>
                  Количество:{" "}
                  <strong>{item.quantity}</strong>
                </p>

                <p>
                  Сумма:{" "}
                  <strong>
                    {Number(
                      item.price * item.quantity
                    ).toLocaleString("ru-RU")}{" "}
                    ₽
                  </strong>
                </p>

                <button
                  className="btn-danger"
                  onClick={() =>
                    removeItem(item)
                  }
                >
                  Удалить
                </button>
              </div>
            ))}
          </div>

          <div
            className="card"
            style={{
              marginTop: 30,
            }}
          >
            <h2>
              Итого:{" "}
              {totalPrice.toLocaleString(
                "ru-RU"
              )}{" "}
              ₽
            </h2>

            <button
              style={{
                marginTop: 20,
                width: "100%",
              }}
              onClick={createOrder}
            >
              Оформить заказ
            </button>
          </div>
        </>
      )}
    </div>
  );
}

export default CartPage;
