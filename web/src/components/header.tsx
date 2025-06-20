import { useCurrentUserProfile } from "@/hooks/use-current-user-profile";
import { Link } from "react-router";
import { Logo } from "./logo";
import SeedCount from "./seed-count";
import UserButton from "./user-button";

export default function Header() {
  const currentUserProfile = useCurrentUserProfile();

  return (
    <div className="flex w-full items-center justify-between">
      <Link to="/home" className="flex items-center space-x-2">
        <Logo />
      </Link>

      <div className="flex items-center space-x-5">
        <div className="flex items-center space-x-3">
          <SeedCount number={currentUserProfile?.seedCount.unused} size={40} />
        </div>

        <UserButton />
      </div>
    </div>
  );
}
