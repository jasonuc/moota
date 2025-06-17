import { Link } from "react-router";
import { Logo } from "./logo";
import SeedCount from "./seed-count";
import UserButton from "./user-button";

type HeaderProps = {
  seedCount: number;
};

export default function Header({ seedCount }: HeaderProps) {
  return (
    <div className="flex w-full items-center justify-between">
      <Link to="/home" className="flex items-center space-x-2">
        <Logo />
      </Link>

      <div className="flex items-center space-x-5">
        <div className="flex items-center space-x-3">
          <SeedCount number={seedCount} size={40} />
        </div>

        <UserButton />
      </div>
    </div>
  );
}
