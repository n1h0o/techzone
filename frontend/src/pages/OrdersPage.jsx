import { useEffect, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { toast } from "sonner";

import api from "../api/api";
import { pay } from "../api/payment";

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
      toast.error("Не удалось загрузить заказы");
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

      toast.success("Статус заказа обновлён");

      loadOrders();
    } catch (err) {
      console.error(err);
      toast.error("Не удалось обновить статус");
    }
  }

  async function payOrder(orderId) {
    try {
      const idempotencyKey = uuidv4();

      await pay(orderId, idempotencyKey);

      toast.success("Оплата прошла успешно 💳");

      loadOrders();
    } catch (err) {
      console.error(err);

      if (err.response?.data) {
        toast.error(err.response.data);
      } else {
        toast.error("Ошибка оплаты");
      }
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

  function getPaymentStatus(status) {
    switch (status) {
      case "success":
        return "✅ Оплачено";
      case "pending":
        return "🟡 В обработке";
      case "failed":
        return "❌ Ошибка";
      default:
        return "⚪ Не оплачено";
    }
  }

  return (
    <div className="page">
      <h1>Мои заказы</h1>

      {orders.length === 0 ? (
        <p>Заказов пока нет.</p>
      ) : (
        orders.map((order) => (
          <div key={order.id} className="card">
            <p>
              <strong>Заказ №</strong> {order.id}
            </p>

            <p>
              <strong>Статус заказа:</strong>{" "}
              {getStatus(order.status)}
            </p>

            <p>
              <strong>Оплата:</strong>{" "}
              {getPaymentStatus(order.payment_status)}
            </p>

            <p>
              <strong>Сумма:</strong>{" "}
              {order.total_price} ₽
            </p>

            {role !== "admin" &&
              order.payment_status !== "success" && (
                <button
                  onClick={() => payOrder(order.id)}
                >
                  💳 Оплатить
                </button>
              )}

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