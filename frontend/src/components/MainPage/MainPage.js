import React, { useState, useEffect } from 'react';
import PostList from '../Blog/PostList';
import NewPost from '../Blog/NewPost';
import { fetchPosts, fetchPostById } from '../../api/api';
import { useAuth } from '../../context/AuthContext';
import '../../styles/MainPage/MainPage.css';

const MainPage = () => {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [isLoadingUser, setIsLoadingUser] = useState(true); // Индикатор загрузки пользователя
  const [isLoadingPosts, setIsLoadingPosts] = useState(true); // Индикатор загрузки постов
  const [error, setError] = useState(null); // Ошибки загрузки

  // Эффект для проверки готовности пользователя
  useEffect(() => {
    if (!user) {
      // Ждем, пока пользователь загрузится
      setIsLoadingUser(true);
    } else {
      setIsLoadingUser(false);
    }
  }, [user]);

  // Эффект для загрузки постов
  useEffect(() => {
    if (!user) return; // Ждем, пока пользователь будет доступен

    const loadPosts = async () => {
      try {
        setIsLoadingPosts(true); // Устанавливаем состояние загрузки постов
        const response = await fetchPosts();
        setPosts(response.data || []);
      } catch (error) {
        console.error('Failed to fetch posts:', error);
        setError('Failed to fetch posts.');
      } finally {
        setIsLoadingPosts(false); // Завершаем загрузку постов
      }
    };

    loadPosts();
  }, [user]);

  // Показываем лоадер, если данные пользователя или постов загружаются
  if (isLoadingUser || isLoadingPosts) {
    return <div className="spinner">Loading...</div>;
  }

  // Показываем ошибку, если произошла ошибка
  if (error) {
    return <p>{error}</p>;
  }

  const handlePostCreated = (newPost) => {
    setPosts((prevPosts) => [newPost, ...prevPosts]);
  };

  return (
    <div>
      <div className="main-page-container">
        <div className="main-page-new-post">
          <NewPost onPostCreated={handlePostCreated} />
        </div>
        <div className="main-page-posts">
          <PostList 
            posts={posts}
            currentUserId={user?.id}
            canDelete={false}
            isOwnProfile={false}
          />
        </div>
      </div>
    </div>
  );
};

export default MainPage;
