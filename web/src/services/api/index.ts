import axios from "axios";

const ax = axios.create({
  baseURL: "/api",
  validateStatus: (status) => status < 300,
  withCredentials: true,
});

export { ax };
