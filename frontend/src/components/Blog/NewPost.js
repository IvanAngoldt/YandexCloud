import React, { useState } from 'react';
import '../../styles/Blog/NewPost.css';
import { createPost } from '../../api/api';
import MicrophoneButton from './MicrophoneButton';

const NewPost = ({ onPostCreated }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');
  const [isWaiting, setIsWaiting] = useState(false); // Флаг ожидания результата распознавания

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title || !content) {
      setErrorMessage('Both fields are required!');
      return;
    }

    try {
      const response = await createPost(title, content);
      setErrorMessage('');
      setSuccessMessage('Post created successfully!');
      onPostCreated(response.data);
      setTitle('');
      setContent('');
    } catch (error) {
      console.error('Failed to create post:', error);
      setErrorMessage('Failed to create post. Please try again later.');
    }
  };

  const handleVoiceResult = (text) => {
    setContent(text);
    setIsWaiting(false);
  };

  const handleVoiceWaiting = (waiting) => {
    setIsWaiting(waiting);
  };

  const handleContentChange = (e) => {
    if (!isWaiting) {
      setContent(e.target.value);
    }
  };

  return (
    <div className="new-post-container">
      <div className="new-post-card">
        <h2>Create New Post</h2>
        <form className="new-post-form" onSubmit={handleSubmit}>
          <input
            type="text"
            className="new-post-input"
            placeholder="Post Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <textarea
            className="new-post-textarea"
            placeholder="Post Content"
            value={isWaiting ? 'Подождите...' : content}
            onChange={handleContentChange}
            readOnly={isWaiting}
          />
          <div className="button-container">
            <div className="button-left">
              <button type="submit" className="action-button submit-button">
                Submit
              </button>
            </div>
            <div className="button-right">
              <MicrophoneButton
                onResult={handleVoiceResult}
                onWaiting={handleVoiceWaiting}
                className="action-button mic-button"
              />
            </div>
          </div>
          {errorMessage && <div className="error-message">{errorMessage}</div>}
          {successMessage && (
            <div className="success-message">{successMessage}</div>
          )}
        </form>
      </div>
    </div>
  );
};

export default NewPost;
