import LandingPage from "@/pages/landing/index";
import { Route, Routes } from "react-router";
import ProtectedRoute from "./components/protected";
import AuthLayout from "./layouts/auth";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import HomePage from "./pages/home";
import LowGeolocationAccuracyPage from "./pages/low-geolocation-accuracy";
import NotFoundPage from "./pages/not-found/not-found";
import AllUserPlantsPage from "./pages/plants";
import PlantGraveyard from "./pages/plants/graveyard";
import IndividualPlantPage from "./pages/plants/plantId";
import PublicPlantPage from "./pages/plants/plantId/public";
import ProfilePage from "./pages/profile";
import SeedsPage from "./pages/seeds";
import SettingsPage from "./pages/settings";
import BaseLayout from "./layouts/base";

export default function App() {
  return (
    <Routes>
      <Route element={<BaseLayout />}>
        <Route index element={<LandingPage />} />

        <Route element={<AuthLayout />}>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
        </Route>

        <Route
          path="/profile/:username"
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/settings"
          element={
            <ProtectedRoute>
              <SettingsPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/low-geolocation-accuracy"
          element={<LowGeolocationAccuracyPage />}
        />

        <Route
          path="/home"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/seeds"
          element={
            <ProtectedRoute>
              <SeedsPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/plants"
          element={
            <ProtectedRoute>
              <AllUserPlantsPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/plants/graveyard"
          element={
            <ProtectedRoute>
              <PlantGraveyard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/plants/:plantId"
          element={
            <ProtectedRoute>
              <IndividualPlantPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/plants/:plantId/public"
          element={
            <ProtectedRoute>
              <PublicPlantPage />
            </ProtectedRoute>
          }
        />

        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  );
}
