import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import api from "../api/api";

function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);

  const loadNotifications = useCallback(async () => {
    try {
      const res = await api.get("/notifications");

      setNotifications(res.data.notifications || []);
    } catch (err) {
      console.error(err);
      toast.error("Не удалось загрузить уведомления");
    }
  }, []);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void loadNotifications();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [loadNotifications]);

  return (
    <div className="container">
      <h1>Уведомления</h1>

      {notifications.length === 0 ? (
        <div className="card">
          <p>Уведомлений пока нет</p>
        </div>
      ) : (
        <div className="notifications-list">
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className="card notification-card"
            >
              <div className="notification-header">
                <span className="notification-icon">
                  🔔
                </span>

                <span className="notification-title">
                  Уведомление
                </span>
              </div>

              <p className="notification-message">
                {notification.message}
              </p>

              <small className="notification-order">
                Заказ №{notification.order_id}
              </small>

              {notification.created_at && (
                <small className="notification-date">
                  {new Date(
                    notification.created_at
                  ).toLocaleString("ru-RU")}
                </small>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default NotificationsPage;
