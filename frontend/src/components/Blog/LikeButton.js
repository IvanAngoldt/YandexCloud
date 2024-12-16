import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import '../../styles/Blog/LikeButton.css';
import { toggleLike, fetchLikes } from '../../api/api';

const LikeButton = ({ postId, initialLikes = [], currentUserId, onLikeChange }) => {
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(initialLikes);
  const [showTooltip, setShowTooltip] = useState(false);
  const [tooltipData, setTooltipData] = useState([]);
  const timerRef = useRef(null);

  useEffect(() => {
    const validLikes = Array.isArray(likes) ? likes : [];
    const userHasLiked = validLikes.some((like) => like.id === currentUserId);
    setLiked(userHasLiked);
    setTooltipData(validLikes);
  }, [likes, currentUserId]);

  const handleLike = async () => {
    try {
      const updatedLikes = await toggleLike(postId, currentUserId, liked);
      setLikes(updatedLikes);

      const userHasLiked = updatedLikes.some((like) => like.id === currentUserId);
      setLiked(userHasLiked);

      onLikeChange(updatedLikes);
    } catch (error) {
      console.error('Ошибка при изменении лайка:', error);
    }
  };

  const handleMouseEnter = () => {
    if (timerRef.current) clearTimeout(timerRef.current);

    timerRef.current = setTimeout(async () => {
      try {
        const data = await fetchLikes(postId);
        setTooltipData(data || []);
        setShowTooltip(true);
      } catch (error) {
        console.error('Ошибка при загрузке лайков:', error);
      }
    }, 500);
  };

  const handleMouseLeave = () => {
    if (timerRef.current) clearTimeout(timerRef.current);
    setShowTooltip(false);
  };

  return (
    <div
      className={`like-button ${liked ? 'liked' : ''}`}
      onClick={handleLike}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={{ position: 'relative', cursor: 'pointer' }}
    >
    <span className="heart">{liked ? '❤️' : '🤍'}</span>
    <span className="like-count">{likes.length}</span>

      {showTooltip && (
        <div
          className="tooltip"
          style={{ position: 'absolute', top: '30px', left: '50%', transform: 'translateX(-50%)', zIndex: 10 }}
          onMouseEnter={() => setShowTooltip(true)} 
          onMouseLeave={() => setShowTooltip(false)}
        >
          {tooltipData.length > 0 ? (
            <>
              <h4>Лайкнули:</h4>
              <ul>
                {tooltipData.map((user) => (
                  <li key={user.id}>
                    <Link
                      to={`/profile/${user.username}`}
                      onClick={(e) => e.stopPropagation()}
                    >
                      {user.username}
                    </Link>
                  </li>
                ))}
              </ul>
            </>
          ) : (
            <p>Никто не лайкнул этот пост</p>
          )}
        </div>
      )}
    </div>
  );
};

export default LikeButton;
