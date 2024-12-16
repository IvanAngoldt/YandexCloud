import React, { useEffect, useState, useRef, useCallback } from 'react';
import { fetchNotifications, markNotificationAsRead, clearNotifications } from '../../api/api';
import { ReactComponent as BellIcon } from '../../icons/bell.svg';
import '../../styles/Header/Notifications.css';

const Notifications = ({ userId }) => {
  const [notifications, setNotifications] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0); // Отдельное состояние для непрочитанных уведомлений
  const dropdownRef = useRef(null);

  // Загружаем уведомления
  useEffect(() => {
    if (userId) {
      fetchNotifications(userId)
        .then((data) => {
          setNotifications(data || []);
          const count = (data || []).filter((n) => !n.isRead).length;
          setUnreadCount(count); // Устанавливаем количество непрочитанных уведомлений
        })
        .catch((error) => console.error('Failed to fetch notifications:', error));
    }
  }, [userId]);

  const markAllAsReadOnServer = useCallback(() => {
    const unreadIds = notifications.filter((n) => !n.isRead).map((n) => n.id);
    if (unreadIds.length === 0) return;

    Promise.all(unreadIds.map((id) => markNotificationAsRead(id)))
      .then(() => {
        setNotifications((prev) =>
          prev.map((n) => ({ ...n, isRead: true }))
        );
      })
      .catch((error) => console.error('Failed to mark notifications as read on server:', error));
  }, [notifications]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        if (showDropdown) {
          markAllAsReadOnServer(); // Обновляем сервер при закрытии
          setShowDropdown(false);
        }
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showDropdown, markAllAsReadOnServer]);

  // Обработчик очистки уведомлений
  const handleClearNotifications = () => {
    clearNotifications(userId)
      .then(() => setNotifications([]))
      .catch((error) => console.error('Failed to clear notifications:', error));
  };

  // Обработчик открытия dropdown
  const toggleDropdown = () => {
    if (showDropdown) {
      markAllAsReadOnServer(); // Помечаем уведомления как прочитанные при закрытии через колокольчик
    } else {
      setUnreadCount(0); // Сбрасываем счётчик непрочитанных уведомлений при открытии
    }
    setShowDropdown((prev) => !prev);
  };

  // Рендер уведомлений
  const renderNotification = (notification) => {
    const notificationClass = notification.isRead
      ? 'notification-item read'
      : 'notification-item unread';
  
    const formattedTime = new Date(notification.createdAt).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
  
    return (
      <div className={notificationClass} key={notification.id}>
        <span className="notification-time">{formattedTime}</span>
        <span className="notification-message">
          {notification.type === 'like' ? (
            <>
              Пользователь{' '}
              <a
                href={`/profile/${notification.likerUsername}`}
                className="link"
                target="_blank"
                rel="noopener noreferrer"
                onClick={() => markNotificationAsRead(notification.id)}
              >
                {notification.likerUsername}
              </a>{' '}
              поставил лайк вашему{' '}
              <a
                href={`/post/${notification.postId}`}
                className="link"
                target="_blank"
                rel="noopener noreferrer"
                onClick={() => markNotificationAsRead(notification.id)}
              >
                посту
              </a>
              .
            </>
          ) : (
            notification.message
          )}
        </span>
      </div>
    );
  };

  return (
    <div className="notifications" ref={dropdownRef}>
      <div className="notification-icon" onClick={toggleDropdown}>
        <BellIcon />
        {unreadCount > 0 && <span className="notification-badge">{unreadCount}</span>}
      </div>
      {showDropdown && (
        <div className="notification-dropdown">
          {notifications.length > 0 && (
            <div className="dropdown-header">
              <button onClick={handleClearNotifications} className="clear-notifications-button">
                Очистить уведомления
              </button>
            </div>
          )}
          {notifications.length === 0 ? (
            <div className="no-notifications">No new notifications</div>
          ) : (
            notifications.map(renderNotification)
          )}
        </div>
      )}
    </div>
  );
};

export default Notifications;
