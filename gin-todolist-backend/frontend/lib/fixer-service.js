require('dotenv').config();
const axios = require('axios');

// Axios Client declaration
const api = axios.create({
  baseURL: 'http://localhost:8080',
  params: {
    // access_key: process.env.API_KEY,
  },
  timeout: process.env.TIMEOUT || 5000,
});

// Generic GET request function
const get = async (url) => {
  const response = await api.get(url);
  const { data } = response;
  return (data);
};

module.exports = {
    getTasks: () => get(`/`),
};