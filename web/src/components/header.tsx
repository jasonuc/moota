import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Link } from "react-router";
import { Button } from "./ui/button";
import SeedCount from "./seed-count";
import { Logo } from "./logo";
import { useAuth } from "@/hooks/use-auth";

type HeaderProps = {
  seedCount: number;
};

export default function Header({ seedCount }: HeaderProps) {
  const { logout } = useAuth();

  return (
    <div className="flex w-full items-center justify-between">
      <Link to="/home" className="flex items-center space-x-2">
        <Logo />
      </Link>

      <div className="flex items-center space-x-5">
        <div className="flex items-center space-x-3">
          <SeedCount number={seedCount} size={40} />
        </div>

        <Button
          asChild
          className="p-0 border-0"
          variant="reverse"
          onClick={logout}
        >
          <Avatar>
            <AvatarImage
              className="size-[45px]"
              src={`https://api.dicebear.com/9.x/glass/svg?seed=${"jasonuc"}`}
              draggable={false}
            />
            <AvatarFallback>JA</AvatarFallback>
          </Avatar>
        </Button>
      </div>
    </div>
  );
}
