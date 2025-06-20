import Header from "@/components/header";
import { LogoWithText } from "@/components/logo";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuth } from "@/hooks/use-auth";
import { ArrowLeft, Home, Search } from "lucide-react";
import { Link, useNavigate } from "react-router";

export default function NotFoundPage() {
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();

  return (
    <div className="flex flex-col space-y-5 grow">
      {isLoggedIn ? (
        <Header />
      ) : (
        <Link to="/">
          <LogoWithText />
        </Link>
      )}

      <div className="flex flex-col items-center justify-center space-y-5 pb-10 grow">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="mx-auto mb-4 w-16 h-16 flex items-center justify-center rounded-full bg-muted">
              <Search className="w-8 h-8 text-muted-foreground" />
            </div>
            <CardTitle className="text-2xl">Page Not Found</CardTitle>
            <CardDescription>
              The page you're looking for doesn't exist or has been moved.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col md:flex-row gap-2">
              <Button
                onClick={() => navigate(-1)}
                className="flex items-center gap-2 grow"
              >
                <ArrowLeft className="w-4 h-4" />
                Go Back
              </Button>
              <Button
                onClick={() => {
                  navigate(isLoggedIn ? "/home" : "/");
                }}
                className="flex items-center gap-2 grow"
              >
                <Home className="w-4 h-4" />
                Home
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
