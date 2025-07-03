import LandingPage from "@/pages/landing/index";
import { Route, Routes } from "react-router";
import AuthLayout from "./layouts/auth";
import BaseLayout from "./layouts/base";
import ProtectedLayout from "./layouts/protected";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import GeolocationDisallowedPage from "./pages/geolocation-dissallowed";
import HomePage from "./pages/home";
import NotFoundPage from "./pages/not-found/not-found";
import UserPlantsPage from "./pages/plants";
import PlantGraveyard from "./pages/plants/graveyard";
import IndividualPlantPage from "./pages/plants/plantId";
import PublicPlantPage from "./pages/plants/plantId/public";
import ProfilePage from "./pages/profile";
import SeedsPage from "./pages/seeds";
import SettingsPage from "./pages/settings";
import GeolocationProvider from "./services/providers/geolocation-provider";

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
          element={
            <GeolocationProvider>
              <ProtectedLayout />
            </GeolocationProvider>
          }
        >
          <Route path="/home" element={<HomePage />} />
          <Route path="/seeds" element={<SeedsPage />} />
          <Route path="/plants" element={<UserPlantsPage />} />
          <Route path="/plants/:plantId" element={<IndividualPlantPage />} />
        </Route>

        <Route element={<ProtectedLayout />}>
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="/profile/:username" element={<ProfilePage />} />
          <Route path="/plants/graveyard" element={<PlantGraveyard />} />
          <Route path="/plants/:plantId/public" element={<PublicPlantPage />} />
        </Route>

        <Route
          path="/geolocation-dissallowed"
          element={<GeolocationDisallowedPage />}
        />
      </Route>

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
