import LandingPage from "@/pages/landing/index";
import { Route, Routes } from "react-router";
import GeolocationAccuracyIndicator from "./components/geolocation-accuracy-indicator";
import ProtectedRoute from "./components/protected";
import AuthLayout from "./layouts/auth";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import HomePage from "./pages/home";
import LowGeolocationAccuracyPage from "./pages/low-geolocation-accuracy";
import AllUserPlantsPage from "./pages/plants";
import PlantGraveyard from "./pages/plants/graveyard";
import IndividualPlantPage from "./pages/plants/plantId";
import PublicPlantPage from "./pages/plants/plantId/public";
import ProfilePage from "./pages/profile";
import SeedsPage from "./pages/seeds";
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
          element={<LowGeolocationAccuracyPage />}
          path="/low-geolocation-accuracy"
        />

        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <HomePage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/home"
        />

        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <SettingsPage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/settings"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <ProfilePage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/profile/:username"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <SeedsPage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/seeds"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <AllUserPlantsPage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/plants"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <PlantGraveyard />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/plants/graveyard"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <IndividualPlantPage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/plants/:plantId"
        />
        <Route
          element={
            <ProtectedRoute>
              <GeolocationAccuracyIndicator>
                <PublicPlantPage />
              </GeolocationAccuracyIndicator>
            </ProtectedRoute>
          }
          path="/plants/:plantId/public"
        />
      </Route>
    </Routes>
  );
}
