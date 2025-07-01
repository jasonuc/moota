import axios, { AxiosError, AxiosRequestConfig } from "axios";

const ax = axios.create({
  baseURL: "/api",
  validateStatus: (status) => status < 300,
  withCredentials: true,
});

ax.interceptors.response.use(
  (response) => response,
  refreshAccessTokenMiddleware
);

interface AxiosRequestConfigExtended extends AxiosRequestConfig {
  _retry?: boolean;
}

async function refreshAccessTokenMiddleware(error: AxiosError) {
  const originalRequest = error.config as AxiosRequestConfigExtended;

  const isAuthRoute = originalRequest?.url?.includes("/auth/");

  if (
    error.response?.status !== 401 ||
    originalRequest?._retry ||
    isAuthRoute ||
    !originalRequest?.url
  ) {
    return Promise.reject(error);
  }

  originalRequest._retry = true;

  try {
    await ax.post("/auth/refresh");
    return ax(originalRequest);
  } catch (refreshError) {
    window.dispatchEvent(new CustomEvent("auth:logout"));
    return Promise.reject(refreshError);
  }
}

export { ax };
