import Header from "@/components/header";
import { LogoWithText } from "@/components/logo";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import {
  AlertTriangleIcon,
  HomeIcon,
  ListIcon,
  MapPinIcon,
  SettingsIcon,
  ShieldIcon,
  SmartphoneIcon,
} from "lucide-react";
import { Link, useNavigate } from "react-router";

export default function GeolocationDisallowedPage() {
  const navigate = useNavigate();
  const { isLoggedIn } = useAuth();
  const { error } = useGeolocation();

  const goHome = () => (isLoggedIn ? navigate("/home") : navigate("/"));

  if (!error) navigate("/");

  return (
    <div className="flex flex-col space-y-5 grow">
      {isLoggedIn ? (
        <Header />
      ) : (
        <Link to="/">
          <LogoWithText />
        </Link>
      )}

      <Card className="w-full max-w-md mx-auto shadow-lg">
        <CardHeader className="text-center pb-4">
          <div className="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
            <ShieldIcon className="w-8 h-8 text-red-600" />
          </div>
          <CardTitle className="text-xl font-semibold text-gray-900">
            Location Access Denied
          </CardTitle>
          <CardDescription className="text-gray-600">
            We need location permission to run this application
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          <Alert>
            <AlertTriangleIcon className="h-4 w-4" />
            <AlertDescription>
              <strong>Location permission required:</strong> Please enable
              location access in your browser to continue using this feature.
            </AlertDescription>
          </Alert>

          <div className="space-y-4">
            <h3 className="font-medium text-gray-900 flex items-center gap-2">
              <ListIcon className="w-4 h-4" />
              How to enable location access:
            </h3>

            <div className="space-y-3 text-sm text-gray-600">
              <div className="flex items-start gap-3">
                <MapPinIcon className="w-4 h-4 mt-0.5 text-blue-500 flex-shrink-0" />
                <span>
                  Look for the location icon in your browser's address bar and
                  click "Allow"
                </span>
              </div>

              <div className="flex items-start gap-3">
                <SettingsIcon className="w-4 h-4 mt-0.5 text-purple-500 flex-shrink-0" />
                <span>
                  Check your browser settings and ensure location access is
                  enabled for this site
                </span>
              </div>

              <div className="flex items-start gap-3">
                <SmartphoneIcon className="w-4 h-4 mt-0.5 text-orange-500 flex-shrink-0" />
                <span>
                  On mobile devices, also ensure location services are enabled
                  in your device settings
                </span>
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-3">
            <Button
              onClick={() => window.location.reload()}
              className="w-full"
              size="lg"
            >
              <MapPinIcon className="w-4 h-4 mr-2" />
              Reload Page
            </Button>

            <Button onClick={goHome} className="w-full" size="lg">
              <HomeIcon className="w-4 h-4 mr-2" />
              Return Home
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
