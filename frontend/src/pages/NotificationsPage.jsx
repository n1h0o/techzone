import { useEffect, useState } from "react";
import api from "../api/api";

function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    loadNotifications();
  }, []);

  async function loadNotifications() {
    try {
      const res = await api.get("/notifications");

      console.log(res.data);

      setNotifications(res.data.notifications || []);
    } catch (err) {
      console.error(err);
      alert("Не удалось загрузить уведомления");
    }
  }

  return (
    <div className="page">
      <h1>Уведомления</h1>

      {notifications.length === 0 ? (
        <p>Уведомлений пока нет</p>
      ) : (
        notifications.map((notification) => (
          <div
            key={notification.id}
            className="card"
          >
            <p>{notification.message}</p>

            <small>
              Заказ №{notification.order_id}
            </small>
          </div>
        ))
      )}
    </div>
  );
}

export default NotificationsPage;