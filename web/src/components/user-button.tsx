import { useAuth } from "@/hooks/use-auth";
import { startSentenceWithUppercase } from "@/lib/utils";
import { requestSeeds } from "@/services/api/seeds";
import { DropdownMenuGroup } from "@radix-ui/react-dropdown-menu";
import { AxiosError } from "axios";
import { formatDate } from "date-fns";
import {
  HandCoinsIcon,
  LogOutIcon,
  SettingsIcon,
  UserIcon,
} from "lucide-react";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";

export default function UserButton() {
  const { logout, user } = useAuth();
  const navigate = useNavigate();

  const handleRequestSeeds = async () => {
    requestSeeds(user!.id)
      .then(() => {
        navigate("/seeds");
      })
      .catch(
        (
          err: AxiosError<{ error: { message: string; timeAvailable: string } }>
        ) => {
          toast.warning(
            startSentenceWithUppercase(err.response?.data.error.message ?? ""),
            {
              description: `You can request seeds again on ${formatDate(
                err.response?.data.error.timeAvailable ?? "",
                "dd/MM/yy"
              )}`,
            }
          );
        }
      );
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button asChild className="p-0 border-0" variant="reverse">
          <Avatar>
            <AvatarImage
              className="size-[45px]"
              src={`https://api.dicebear.com/9.x/glass/svg?seed=${user?.username}`}
              draggable={false}
            />
            <AvatarFallback>
              {user?.username.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem
            onClick={() => navigate(`/profile/${user?.username}`)}
          >
            <UserIcon />
            <span>Profile</span>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleRequestSeeds}>
            <HandCoinsIcon />
            <span>Request seeds</span>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <SettingsIcon />
            <span>Settings</span>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={logout}>
            <LogOutIcon />
            <span>Logout</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
