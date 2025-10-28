import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useKillPlant } from "@/services/mutations/plants";
import { GhostIcon, SkullIcon, XIcon } from "lucide-react";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { Button } from "./ui/button";

type KillPlantButtonProps = {
  id: string;
  nickname: string;
};

export default function KillPlantButton({
  id,
  nickname,
}: KillPlantButtonProps) {
  const navigate = useNavigate();
  const killPlantMtn = useKillPlant();

  const handleKillPlant = () => {
    killPlantMtn.mutateAsync(id).then(() => navigate("/plants"));
    toast(
      <>
        RIP <i className="font-base">{nickname}</i>
      </>,
      {
        icon: <GhostIcon />,
      },
    );
  };

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button className="md:min-h-12 col-span-1 flex items-center justify-center space-x-1.5 bg-red-400 hover:bg-red-500">
          Kill <SkullIcon />
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. <strong>{nickname}</strong> would be
            gone forever.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="hover:cursor-pointer md:min-h-12 col-span-1 flex items-center justify-center space-x-1.5">
            <XIcon />
            No
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={handleKillPlant}
            className="hover:cursor-pointer md:min-h-12 col-span-1 flex items-center justify-center space-x-1.5 bg-red-400 hover:bg-red-500"
          >
            <SkullIcon />
            Yes
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
