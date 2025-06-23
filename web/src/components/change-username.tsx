import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useAuth } from "@/hooks/use-auth";
import { useCurrentUserProfile } from "@/hooks/use-current-user-profile";
import { changeUsernameFormSchema } from "@/schemas/settings";
import { useChangeUsername } from "@/services/mutations/user";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "./ui/button";
import { Input } from "./ui/input";

export default function ChangeUsername() {
  const { user } = useAuth();
  const changeUsernameMtn = useChangeUsername();
  const currentUserProfile = useCurrentUserProfile();

  const form = useForm<z.infer<typeof changeUsernameFormSchema>>({
    resolver: zodResolver(changeUsernameFormSchema),
    defaultValues: {
      newUsername: "",
    },
  });

  function onSubmit(values: z.infer<typeof changeUsernameFormSchema>) {
    changeUsernameMtn
      .mutateAsync({
        userId: user!.id,
        currentUsername: currentUserProfile!.username,
        newUsername: values.newUsername,
      })
      .then(() => {
        toast("Username changed successfully!");
      })
      .catch(() => {
        toast("Error occured while trying to change username");
      });
    form.reset();
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col gap-y-2.5"
      >
        <FormField
          name="newUsername"
          render={({ field }) => (
            <FormItem>
              <FormLabel>New Username</FormLabel>
              <FormControl>
                <Input placeholder={currentUserProfile?.username} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="w-full flex justify-end">
          <Button disabled={!user?.id} type="submit" className="w-fit">
            Submit
          </Button>
        </div>
      </form>
    </Form>
  );
}
