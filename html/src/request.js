import axios from 'axios';

const request = axios.create({
  method: 'post',
  baseURL: window.ipConfigUrl.baseURL,
  timeout: 3000,
  responseType: 'json',
});

export default request;
