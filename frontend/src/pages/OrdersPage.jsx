import { useEffect, useState } from "react";
import api from "../api/api";

function OrdersPage() {
  const [orders, setOrders] = useState([]);

  const role = localStorage.getItem("role");

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
      await api.patch(`/orders/${orderId}/status`, {
        status,
      });

      alert("Статус обновлён");

      loadOrders();
    } catch (err) {
      console.error(err);
      alert("Ошибка обновления статуса");
    }
  }

  function getStatus(status) {
    switch (status) {
      case "new":
        return "Новый";
      case "processing":
        return "В обработке";
      case "completed":
        return "Завершён";
      default:
        return status;
    }
  }

  return (
    <div className="page">
      <h1>Мои заказы</h1>

      {orders.length === 0 ? (
        <p>Заказов пока нет.</p>
      ) : (
        orders.map((order) => (
          <div
            key={order.id}
            className="card"
          >
            <p>
              <strong>Заказ №</strong>
              {order.id}
            </p>

            <p>
              <strong>Статус:</strong>{" "}
              {getStatus(order.status)}
            </p>

            <p>
              <strong>Сумма:</strong>{" "}
              {order.total_price} ₽
            </p>

            {role === "admin" && (
              <>
                {order.status === "new" && (
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

                    <option value="processing">
                      В обработке
                    </option>
                  </select>
                )}

                {order.status === "processing" && (
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

                    <option value="completed">
                      Завершён
                    </option>
                  </select>
                )}
              </>
            )}
          </div>
        ))
      )}
    </div>
  );
}

export default OrdersPage;