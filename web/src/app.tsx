import { Route, Routes } from "react-router";
import LandingPage from "@/pages/landing/index";
import AuthLayout from "./layouts/auth";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import SeedsPage from "./pages/seeds";
import HomePage from "./pages/home";
import ProfilePage from "./pages/profile";
import AllUserPlantsPage from "./pages/plants";
import IndividualPlantPage from "./pages/plants/plantId";
import ProtectedRoute from "./components/protected";

export default function App() {
  return (
    <Routes>
      <Route index element={<LandingPage />} />

      <Route element={<AuthLayout />}>
        <Route element={<LoginPage />} path="/login" />
        <Route element={<RegisterPage />} path="/register" />
      </Route>

      <Route>
        <Route
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
          path="/home"
        />
        <Route
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
          path="/profile/:username"
        />
        <Route
          element={
            <ProtectedRoute>
              <SeedsPage />
            </ProtectedRoute>
          }
          path="/seeds"
        />
        <Route
          element={
            <ProtectedRoute>
              <AllUserPlantsPage />
            </ProtectedRoute>
          }
          path="/plants"
        />
        <Route
          element={
            <ProtectedRoute>
              <IndividualPlantPage />
            </ProtectedRoute>
          }
          path="/plants/:plantId"
        />
      </Route>
    </Routes>
  );
}
