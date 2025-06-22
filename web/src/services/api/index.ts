import axios from "axios";

const ax = axios.create({
  baseURL:
    process.env.NODE_ENV == "production" ? "https://api.moota.app/api" : "/api",
  validateStatus: (status) => status < 300,
  withCredentials: true,
});

export { ax };
