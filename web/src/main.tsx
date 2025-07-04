import App from "@/app";
import { Toaster } from "@/components/ui/sonner";
import AuthProvider from "@/services/providers/auth-provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import InDevelopmentBanner from "./components/in-development-banner";
import "./index.css";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <InDevelopmentBanner />
          <App />
          <Toaster />
        </AuthProvider>
        {import.meta.env.DEV && <ReactQueryDevtools />}
      </QueryClientProvider>
    </BrowserRouter>
  </StrictMode>
);
