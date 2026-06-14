import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api/api";

function CartPage() {
  const [items, setItems] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    loadCart();
  }, []);

  async function loadCart() {
    try {
      const res = await api.get("/cart");
      setItems(res.data.items || []);
    } catch (err) {
      console.error(err);
    }
  }

  async function removeItem(id) {
    try {
      await api.delete(`/cart/items/${id}`);
      loadCart();
    } catch (err) {
      console.error(err);
    }
  }

  async function createOrder() {
    try {
      const res = await api.post("/orders");

      alert(
        `Заказ №${res.data.order_id} создан`
      );

      navigate("/orders");
    } catch (err) {
      console.error(err);
      alert("Ошибка создания заказа");
    }
  }

  return (
    <div>
      <h1>Корзина</h1>

      {items.length === 0 ? (
        <p>Корзина пуста</p>
      ) : (
        <>
          {items.map((item) => (
            <div key={item.id}>
              <h3>{item.name}</h3>

              <p>Цена: {item.price} ₽</p>

              <p>Количество: {item.quantity}</p>

              <button
                onClick={() => removeItem(item.id)}
              >
                Удалить
              </button>
            </div>
          ))}

          <button onClick={createOrder}>
            Оформить заказ
          </button>
        </>
      )}
    </div>
  );
}

export default CartPage;