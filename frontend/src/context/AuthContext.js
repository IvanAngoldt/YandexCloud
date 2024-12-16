import React, { createContext, useContext, useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initializeAuth = () => {
      const storedUser = localStorage.getItem('user');
      const storedToken = localStorage.getItem('token');
  
      if (storedUser && storedToken) {
        try {
          const decodedToken = jwtDecode(storedToken);
  
          const isTokenValid = decodedToken.exp * 1000 > Date.now();
  
          if (isTokenValid) {
            setUser(JSON.parse(storedUser));
            setIsAuthenticated(true);
          } else {
            console.warn('Token expired, logging out...');
            logout();
          }
        } catch (error) {
          console.error('Error decoding token:', error);
          logout();
        }
      }
  
      setIsLoading(false);
    };
  
    initializeAuth();
  }, []);

  const login = ({ user: userData, token }) => {
    setIsAuthenticated(true);
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
    localStorage.setItem('token', token);
  };

  const logout = () => {
    setIsAuthenticated(false);
    setUser(null);
    localStorage.removeItem('user');
    localStorage.removeItem('token');
  };

  const setAuthToken = (token) => {
    localStorage.setItem('token', token);
    try {
      const decodedToken = jwtDecode(token);
      if (decodedToken.exp * 1000 > Date.now()) {
        setIsAuthenticated(true);
      } else {
        logout();
      }
    } catch (error) {
      console.error('Invalid token:', error);
      logout();
    }
  };

  const setAuthUser = (userData) => {
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout, setAuthToken, setAuthUser, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
