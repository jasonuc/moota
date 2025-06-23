import { Button } from "@/components/ui/button";
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useChangePlantNickname } from "@/services/mutations/plants";

type ChangePlantNicknameProps = {
  plantId: string;
  currentNickname: string;
  setIsDialogOpen: (arg0: boolean) => void;
};

export default function ChangePlantNickname({
  plantId,
  currentNickname,
  setIsDialogOpen,
}: ChangePlantNicknameProps) {
  const changePlantNicknameMtn = useChangePlantNickname();

  const handleSaveChangeClick = async (formData: FormData) => {
    const newNickname = formData.get("nickname");
    if (!newNickname) return;
    changePlantNicknameMtn.mutate({
      plantId,
      newNickname: newNickname.toString(),
    });
    setIsDialogOpen(false);
  };

  return (
    <DialogContent className="sm:max-w-[425px]">
      <form action={handleSaveChangeClick}>
        <DialogHeader>
          <DialogTitle>Edit nickname</DialogTitle>
          <DialogDescription>
            {"Make changes to your plant's nickanme here."}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-3 my-5">
          <Label htmlFor="nickname">Nickname</Label>
          <Input id="nickname" name="nickname" defaultValue={currentNickname} />
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="neutral">Cancel</Button>
          </DialogClose>
          <Button type="submit">Save</Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
