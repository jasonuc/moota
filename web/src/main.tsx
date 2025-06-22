import App from "@/app";
import { Toaster } from "@/components/ui/sonner";
import AuthProvider from "@/services/providers/auth-provider";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import "./index.css";
import GeolocationProvider from "./services/providers/geolocation-provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <GeolocationProvider>
            <App />
            <Toaster />
          </GeolocationProvider>
        </AuthProvider>
      </QueryClientProvider>
    </BrowserRouter>
  </StrictMode>
);
