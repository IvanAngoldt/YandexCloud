import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { fetchUserProfile, fetchUserPosts, deletePost } from '../../api/api';
import PostList from '../Blog/PostList';
import '../../styles/Profile/Profile.css';

const Profile = () => {
  const { username } = useParams();
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [isLoading, setIsLoading] = useState(true); // Состояние загрузки
  const [error, setError] = useState(null);

  const isOwnProfile = username === user?.username;

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profileResponse = await fetchUserProfile(username);
        setProfile(profileResponse.data);
      } catch (err) {
        console.error('Failed to load profile:', err);
        setError('Failed to load profile.');
      } finally {
        setIsLoading(false); // Завершаем загрузку
      }

    };

    const loadPosts = async () => {
      try {
        const postsResponse = await fetchUserPosts(username);
        setPosts(postsResponse.data);
      } catch (err) {
        console.error('Failed to load posts:', err);
        setError('Failed to load posts.');
      } finally {
        setIsLoading(false); // Завершаем загрузку
      }

    };
    
    loadProfile()
    loadPosts();
  }, [username]);

  const handleDeletePost = async (postId) => {
    try {
      await deletePost(postId);
      setPosts((prev) => prev.filter((post) => post.id !== postId));
    } catch (err) {
      console.error('Failed to delete post:', err);
      alert('Failed to delete post. Please try again.');
    }
  };

  if (isLoading) {
    return <div className="spinner">Loading...</div>; // Загрузочный индикатор
  }

  if (error) {
    return <p>{error}</p>;
  }

  return (
    <div>
      <div className="profile-container">
        <div className="profile-header">
          <h1>{isOwnProfile ? 'Your Profile' : `${profile.username}'s Profile`}</h1>
        </div>
        <div className="profile-info">
          <h2>Username: {profile.username}</h2>
          <h3>Email: {profile.email}</h3>
        </div>
        <div className="profile-posts">
          <h2>{isOwnProfile ? 'My Blogs' : `${profile.username}'s Blogs`}</h2>
          <PostList
            posts={posts}
            onDeletePost={handleDeletePost}
            currentUserId={user?.id}
            canDelete={isOwnProfile}
            isOwnProfile={isOwnProfile}
          />
        </div>
      </div>
    </div>
  );
};

export default Profile;
