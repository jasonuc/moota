import { LogoWithText } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Link } from "react-router";

export default function LandingPage() {
  return (
    <div className="flex flex-col space-y-5 grow">
      <div className="flex w-full justify-between items-center">
        <Link to="/">
          <LogoWithText />
        </Link>

        <div className="flex space-x-2">
          <Link to="/home">
            <Button>Home</Button>
          </Link>

          <Link to="/login">
            <Button>Login</Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
