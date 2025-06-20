import LandingPage from "@/pages/landing/index";
import { Route, Routes } from "react-router";
import ProtectedRoute from "./components/protected";
import AuthLayout from "./layouts/auth";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import HomePage from "./pages/home";
import AllUserPlantsPage from "./pages/plants";
import IndividualPlantPage from "./pages/plants/plantId";
import ProfilePage from "./pages/profile";
import SeedsPage from "./pages/seeds";
import PlantGraveyard from "./pages/plants/graveyard";
import SettingsPage from "./pages/settings";

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
              <SettingsPage />
            </ProtectedRoute>
          }
          path="/settings"
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
              <PlantGraveyard />
            </ProtectedRoute>
          }
          path="/plants/graveyard"
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
