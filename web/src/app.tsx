import LandingPage from "@/pages/landing/index";
import { Route, Routes } from "react-router";
import { ProtectedGeolocationAccuracyIndicatorRoute } from "./components/geolocation-accuracy-indicator";
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

export default function App() {
  return (
    <Routes>
      <Route index element={<LandingPage />} />

      <Route
        path="/low-geolocation-accuracy"
        element={<LowGeolocationAccuracyPage />}
      />

      <Route element={<AuthLayout />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Route>

      <Route
        path="/home"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <HomePage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/settings"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <SettingsPage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/profile/:username"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <ProfilePage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/seeds"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <SeedsPage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/plants"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <AllUserPlantsPage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/plants/graveyard"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <PlantGraveyard />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/plants/:plantId"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <IndividualPlantPage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route
        path="/plants/:plantId/public"
        element={
          <ProtectedGeolocationAccuracyIndicatorRoute>
            <PublicPlantPage />
          </ProtectedGeolocationAccuracyIndicatorRoute>
        }
      />

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
