import React from 'react';
import Post from './Post';
import '../../styles/Blog/PostList.css';

const PostList = ({ posts, onDeletePost, currentUserId, canDelete, isOwnProfile }) => {
  if (!posts || posts.length === 0) {
    return (
      <div className="post-list-empty">
        {isOwnProfile ? 'No posts yet, write the first one!' : 'No posts yet.'}
      </div>
    );
  }

  return (
    <div className="post-list-container">
      <div className="post-list">
        {posts.map((post) => {
          return (
            <Post
              key={post.id}
              post={post}
              onDelete={onDeletePost}
              currentUserId={currentUserId}
              canDelete={canDelete}
            />
          );
        })}
      </div>
    </div>
  );
  
};

export default PostList;
