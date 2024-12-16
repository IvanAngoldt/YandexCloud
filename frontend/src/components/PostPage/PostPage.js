import React, { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useParams } from 'react-router-dom'; // Импортируем useParams для получения параметров маршрута
import { fetchPostById } from '../../api/api'; // Функция для загрузки поста

import Post from '../Blog/Post'; // Импортируем компонент Post
import '../../styles/PostPage/PostPage.css';

const PostPage = () => { // Получаем пропсы, если нужно
  const { user } = useAuth();
  const { postID } = useParams(); // Получаем postID из маршрута
  const [post, setPost] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // Загружаем пост по ID
  useEffect(() => {
    const loadPost = async () => {
      try {
        const response = await fetchPostById(postID);
        setPost(response.data); // Сохраняем ответ в состояние
      } catch (error) {
        console.error('Failed to fetch post:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadPost();
  }, [postID]);

  if (isLoading) {
    return <div className="spinner">Loading...</div>;
  }

  if (!post) {
    return <div className="post-page-container">Post not found</div>;
  }

  return (
    <div>
      <div className="post-page-container">
        {/* Передаем данные поста в компонент Post */}
        <Post 
          post={post} 
          currentUserId={user?.id} 
          canDelete={false}
        />
      </div>
    </div>
  );
};

export default PostPage;
