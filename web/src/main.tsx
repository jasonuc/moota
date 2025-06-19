import App from "@/app";
import { Toaster } from "@/components/ui/sonner";
import AuthProvider from "@/services/providers/auth-provider";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <div className="flex font-archivo p-5 md:p-10 min-h-screen w-full">
          <App />
        </div>
        <Toaster />
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>
);
