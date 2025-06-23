import { Plant } from "@/types/plant";
import { EditIcon, MenuIcon, ShareIcon } from "lucide-react";
import { useState } from "react";
import ChangePlantNickname from "./change-plant-nickname";
import { Button } from "./ui/button";
import { Dialog } from "./ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { toast } from "sonner";

type MoreButtonProps = Plant;

export default function PlantPageMoreButton({ id, nickname }: MoreButtonProps) {
  const [dialogMenu, setDialogMenu] = useState("none");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleDialogMenu = () => {
    switch (dialogMenu) {
      case "changePlantNickname":
        return (
          <ChangePlantNickname
            plantId={id}
            currentNickname={nickname}
            setIsDialogOpen={setIsDialogOpen}
          />
        );
      default:
        return null;
    }
  };

  const openDialog = (menuType: string) => {
    setDialogMenu(menuType);
    setIsDialogOpen(true);
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button className="md:min-h-12 col-span-3 md:col-span-1 flex items-center justify-center space-x-1.5">
            More <MenuIcon />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onSelect={() => openDialog("changePlantNickname")}>
            <EditIcon className="size-2" />
            Change plant nickname
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              navigator.clipboard
                .writeText(window.location.href + "/public")
                .then(() => toast.success("Copied link to plant card!"));
            }}
          >
            <ShareIcon className="size-2" />
            Share plant
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        {handleDialogMenu()}
      </Dialog>
    </>
  );
}
