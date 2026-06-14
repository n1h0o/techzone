import { useEffect, useState } from "react";
import api from "../api/api";

function OrdersPage() {
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    loadOrders();
  }, []);

  async function loadOrders() {
    try {
      const res = await api.get("/orders");
      setOrders(res.data.orders || []);
    } catch (err) {
      console.error(err);
    }
  }

  async function updateStatus(orderId, status) {
    if (!status) {
      return;
    }

    try {
      await api.patch(
        `/orders/${orderId}/status`,
        {
          status: status,
        }
      );

      alert("Статус обновлён");

      loadOrders();
    } catch (err) {
      console.error(err);
      alert("Ошибка обновления статуса");
    }
  }

  return (
    <div>
      <h1>Мои заказы</h1>

      {orders.length === 0 ? (
        <p>Заказов пока нет</p>
      ) : (
        orders.map((order) => (
          <div
            key={order.id}
            style={{
              border: "1px solid gray",
              marginBottom: "10px",
              padding: "10px",
              borderRadius: "8px",
            }}
          >
            <p>
              <strong>ID:</strong> {order.id}
            </p>

            <p>
              <strong>Статус:</strong>{" "}
              {order.status}
            </p>

            <p>
              <strong>Сумма:</strong>{" "}
              {order.total_price} ₽
            </p>

            <select
              defaultValue=""
              onChange={(e) =>
                updateStatus(
                  order.id,
                  e.target.value
                )
              }
            >
              <option value="">
                Изменить статус
              </option>

              <option value="new">
                New
              </option>

              <option value="processing">
                Processing
              </option>

              <option value="completed">
                Completed
              </option>

              <option value="cancelled">
                Cancelled
              </option>
            </select>
          </div>
        ))
      )}
    </div>
  );
}

export default OrdersPage;