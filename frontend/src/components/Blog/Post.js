import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import '../../styles/Blog/Post.css';
import LikeButton from './LikeButton';

const Post = ({ post, currentUserId, canDelete, onDelete }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const isOwner = post.authorId === currentUserId;

  const openModal = () => setIsModalOpen(true);
  const closeModal = () => setIsModalOpen(false);

  const handleDelete = () => {
    onDelete(post.id);
    closeModal();
  };

  const handleOverlayClick = (e) => {
    if (e.target.classList.contains('modal-overlay')) {
      closeModal();
    }
  };

  return (
    <div className="post-card">
      <h3>
        <Link to={`/post/${post.id}`} className="post-title-link">
          {post.title}
        </Link>
      </h3>
      <p>{post.content}</p>
      <div className="post-footer">
        <div className="post-footer-left">
          <LikeButton
            postId={post.id}
            initialLikes={post.likes}
            currentUserId={currentUserId}
            onLikeChange={(updatedLikes) => {
              post.likes = updatedLikes;
            }}
          />
        </div>

        <div className="post-footer-right">
          <Link to={`/profile/${post.authorUsername}`} className="post-author">
            {isOwner ? 'By You' : `By ${post.authorUsername || 'Unknown'}`}
          </Link>

          {canDelete && isOwner && (
            <button className="delete-button" onClick={openModal}>
              ❌
            </button>
          )}
        </div>
      </div>

      {isModalOpen && (
        <div className="modal-overlay" onClick={handleOverlayClick}>
          <div className="modal-container">
            <div className="modal-content">
              <p>Вы уверены, что хотите удалить этот пост?</p>
              <div className="modal-actions">
                <button className="confirm-button" onClick={handleDelete}>
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

export default Post;
