import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import Notifications from './Notifications';
import Menu from './Menu';

import '../../styles/Header/Header.css';

const Header = () => {
  const { isAuthenticated, user } = useAuth();

  useEffect(() => {
    document.body.classList.add('with-header');
    return () => {
      document.body.classList.remove('with-header');
    };
  }, []);

  return (
    <header className="header">
      <nav className="header-nav">
        <div className="header-left">
          <Link to="/" className="header-link">
            Blogs
          </Link>
        </div>
        <div className="header-right">
          {isAuthenticated && user && (
            <>
              <Notifications userId={user.id} />
              <Menu user={user} />
            </>
          )}
        </div>
      </nav>
    </header>
  );
};

export default Header;
