import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './PrivateRoute/PrivateRoute';

import Login from './components/Auth/Login';
import Logout from './components/Auth/Logout';
import Register from './components/Auth/Register';

import Header from './components/Header/Header';
import MainPage from './components/MainPage/MainPage';
import Profile from './components/Profile/Profile';
import PostPage from './components/PostPage/PostPage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          <Route
            path="/"
            element={
              <PrivateRoute>
                <Header />
                <MainPage />
              </PrivateRoute>
            }
          />

          <Route
            path="/logout"
            element={
              <PrivateRoute>
                <Logout />
              </PrivateRoute>
            }
          />

          <Route
            path="/profile"
            element={
              <PrivateRoute>
                <Header />
                <Profile />
              </PrivateRoute>
            }
          />

          <Route
            path="/profile/:username"
            element={
              <PrivateRoute>
                <Header />
                <Profile />
              </PrivateRoute>
            }
          />

          <Route
            path="/post/:postID"
            element={
              <PrivateRoute>
                <Header />
                <PostPage />
              </PrivateRoute>
            }
          />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
