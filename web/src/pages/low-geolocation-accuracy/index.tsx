import { LogoWithText } from "@/components/logo";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
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
import {
  AlertTriangleIcon,
  HomeIcon,
  ListIcon,
  MapPinIcon,
  RefreshCwIcon,
  SmartphoneIcon,
  WifiIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";

export default function LowGeolocationAccuracyPage() {
  const navigate = useNavigate();
  const { isLoggedIn } = useAuth();
  const [currentAccuracy, setCurrentAccuracy] = useState<number>();
  const [isRetrying, setIsRetrying] = useState(false);

  const checkGeolocation = () => {
    setIsRetrying(true);

    if (!navigator.geolocation) {
      setCurrentAccuracy(0);
      setIsRetrying(false);
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setCurrentAccuracy(position.coords.accuracy);
        setIsRetrying(false);
      },
      () => {
        setCurrentAccuracy(0);
        setIsRetrying(false);
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0,
      }
    );
  };

  useEffect(() => {
    checkGeolocation();
  }, []);

  const getAccuracyStatus = (accuracy: number) => {
    if (accuracy === 0) return { text: "Unknown" };
    if (accuracy <= 5) return { text: "Excellent" };
    if (accuracy <= 10) return { text: "Good" };
    if (accuracy <= 50) return { text: "Poor" };
    return { text: "Very Poor" };
  };

  const accuracyStatus = getAccuracyStatus(currentAccuracy ?? 0);

  const goHome = () => navigate("/");

  return (
    <div className="flex flex-col space-y-5 grow">
      <Link to={isLoggedIn ? "/home" : "/"}>
        <LogoWithText />
      </Link>

      <Card className="w-full max-w-md mx-auto shadow-lg">
        <CardHeader className="text-center pb-4">
          <div className="mx-auto w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center mb-4">
            <MapPinIcon className="w-8 h-8 text-orange-600" />
          </div>
          <CardTitle className="text-xl font-semibold text-gray-900">
            Location Accuracy Too Low
          </CardTitle>
          <CardDescription className="text-gray-600">
            We need more precise location data to continue
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          <Alert>
            <AlertTriangleIcon className="h-4 w-4" />
            <AlertDescription>
              <strong>Required accuracy:</strong> 5 meters or better
              {currentAccuracy && (
                <div className="mt-2 flex items-center gap-2">
                  <span className="text-sm">Current accuracy:</span>
                  <Badge className="text-xs">
                    {Math.round(currentAccuracy)}m - {accuracyStatus.text}
                  </Badge>
                </div>
              )}
            </AlertDescription>
          </Alert>

          <div className="space-y-4">
            <h3 className="font-medium text-gray-900 flex items-center gap-2">
              <ListIcon className="w-4 h-4" />
              Tips to improve accuracy:
            </h3>

            <div className="space-y-3 text-sm text-gray-600">
              <div className="flex items-start gap-3">
                <SmartphoneIcon className="w-4 h-4 mt-0.5 text-blue-500 flex-shrink-0" />
                <span>
                  Try using a mobile device - phones typically have better GPS
                  accuracy than desktops
                </span>
              </div>

              <div className="flex items-start gap-3">
                <WifiIcon className="w-4 h-4 mt-0.5 text-purple-500 flex-shrink-0" />
                <span>
                  Ensure location services are enabled in your browser settings
                </span>
              </div>

              <div className="flex items-start gap-3">
                <RefreshCwIcon className="w-4 h-4 mt-0.5 text-orange-500 flex-shrink-0" />
                <span>
                  Wait a moment and try again - accuracy often improves after a
                  few seconds
                </span>
              </div>
            </div>
          </div>

          <Separator />

          <div className="space-y-3">
            <Button
              onClick={checkGeolocation}
              disabled={isRetrying}
              className="w-full"
              size="lg"
            >
              {isRetrying ? (
                <>
                  <RefreshCwIcon className="w-4 h-4 mr-2 animate-spin" />
                  Checking Location...
                </>
              ) : (
                <>
                  <RefreshCwIcon className="w-4 h-4 mr-2" />
                  Try Again
                </>
              )}
            </Button>

            <Button
              variant="default"
              onClick={goHome}
              className="w-full"
              size="lg"
            >
              <HomeIcon className="w-4 h-4 mr-2" />
              Return Home
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
