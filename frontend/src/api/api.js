import axios from 'axios';

const AUTH_API_URL = '/api/auth';
const NOTIS_API_URL = '/api/notifications';
const POSTS_API_URL = '/api/posts';
const USERS_API_URL = '/api/users';
const RECOGNIZE_API_URL = '/api/recognize';

function getAuthHeaders() {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error('Token not found');
  }

  return { Authorization: `Bearer ${token}` };
}

export const login = async (email, password) => {
  return axios.post(`${AUTH_API_URL}/login`, { email, password });
};

export const register = async (username, email, password) => {
  return axios.post(`${AUTH_API_URL}/register`, { username, email, password });
};

export const fetchPosts = async () => {
  const headers = getAuthHeaders();

  return axios.get(`${POSTS_API_URL}/posts`, { headers });
};

export const createPost = async (title, content) => {
  const headers = getAuthHeaders();

  return axios.post(`${POSTS_API_URL}/posts`, { title, content }, { headers });
};

export const deletePost = async (postId) => {
  const headers = getAuthHeaders();

  return axios.delete(`${POSTS_API_URL}/posts/${postId}`, { headers });
};

export const fetchPostById = async (postId) => {
  const headers = getAuthHeaders();

  return axios.get(`${POSTS_API_URL}/posts/${postId}`, { headers });
};

export const fetchUserPosts = async (username) => {
  const headers = getAuthHeaders();

  return axios.get(`${POSTS_API_URL}/profile/${username}/posts`, { headers });
};

export const toggleLike = async (postId, userId, liked) => {
  const headers = getAuthHeaders();

  const response = await axios({
    method: liked ? 'delete' : 'post',
    url: `${POSTS_API_URL}/likes`,
    headers,
    data: { postId, userId },
  });

  if (response.data === null || response.data === undefined) {
    return [];
  }

  if (!Array.isArray(response.data)) {
    throw new Error('Server returned invalid response');
  }

  return response.data;
};

export const fetchLikes = async (postId) => {
  const headers = getAuthHeaders();

  const response = await axios.get(`${POSTS_API_URL}/likes?postId=${postId}`, { headers });

  if (!response.data) {
    return { users: [] };
  }

  return response.data;
};

export const fetchUserProfile = async (username) => {
  const headers = getAuthHeaders();
  
  return axios.get(`${USERS_API_URL}/users/by_username?username=${username}`, { headers });
};

export const fetchNotifications = async (userId) => {
  const headers = getAuthHeaders();

  const response = await axios.get(`${NOTIS_API_URL}/notifications?userId=${userId}`, { headers });
  return response.data;
};

export const markNotificationAsRead = async (id) => {
  const headers = getAuthHeaders();

  return axios.patch(`${NOTIS_API_URL}/notifications/read?id=${id}`, null, { headers });
};

export const clearNotifications = async (userId) => {
  const headers = getAuthHeaders();

  return axios.delete(`${NOTIS_API_URL}/notifications/${userId}/clear`, { headers });
};

export const sendAudioToServer = async (audioBlob) => {
  const buffer = await audioBlob.arrayBuffer();

  const response = await axios.post(
    `${RECOGNIZE_API_URL}/recognize`,
    new Uint8Array(buffer),
    {
      headers: {
        'Content-Type': 'audio/ogg',
        ...getAuthHeaders(),
      },
    }
  );

  return response.data;
};