import React, { useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import { FaUser, FaSignOutAlt } from 'react-icons/fa'; // Import icons
import '../../styles/Header/Menu.css';

const Menu = ({ user }) => {
  const [menuOpen, setMenuOpen] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false); // Состояние для модального окна
  const dropdownRef = useRef(null);

  // Функция для открытия меню на наведение
  const handleMouseEnter = () => {
    setMenuOpen(true);
  };

  // Функция для закрытия меню на выходе из области
  const handleMouseLeave = () => {
    setMenuOpen(false);
  };

  // Открытие модального окна
  const openModal = () => {
    setMenuOpen(false);
    setIsModalOpen(true);
  }

  // Закрытие модального окна
  const closeModal = () => setIsModalOpen(false);

  // Обработчик подтверждения выхода
  const handleLogout = () => {
    // Здесь добавьте логику выхода, например, очистку токенов и редирект
    window.location.href = '/logout'; // Пример редиректа
  };

  const handleOverlayClick = (e) => {
    if (e.target.classList.contains('modal-overlay')) {
      closeModal();
    }
  };

  return (
    <div
      className="header-avatar"
      onMouseEnter={handleMouseEnter} // Открытие при наведении
      onMouseLeave={handleMouseLeave} // Закрытие при выходе
      ref={dropdownRef}
    >
      <div className='avatar'>
        <span>{user.username.charAt(0).toUpperCase()}</span>
      </div>
      <div className={`header-menu ${menuOpen ? 'open' : ''}`}>
      <div className="header-menu-item username">
        {user.email}
      </div>

      <Link to={`/profile/${user.username}`} className="header-menu-item">
        <FaUser /> {/* Profile icon */}
        My Profile
      </Link>

      <button className="header-menu-item" onClick={openModal}>
        <FaSignOutAlt /> {/* Logout icon */}
        Logout
      </button>
      </div>

      {/* Модальное окно с подтверждением */}
      {isModalOpen && (
        <div className="modal-overlay" onClick={handleOverlayClick}>
          <div className="modal-container">
            <div className="modal-content">
              <p>Вы уверены, что хотите выйти?</p>
              <div className="modal-actions">
                <button className="confirm-button" onClick={handleLogout}>
                  Подтвердить
                </button>
                <button className="cancel-button" onClick={closeModal}>
                  Отмена
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Menu;
