import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useAuth } from "@/hooks/use-auth";
import { changePasswordFormSchema } from "@/schemas/settings";
import { useChangePassword } from "@/services/mutations/user";
import { zodResolver } from "@hookform/resolvers/zod";
import { AxiosError } from "axios";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { startSentenceWithUppercase } from "@/lib/utils";

export default function ChangePassword() {
  const { user } = useAuth();
  const changePasswordMtn = useChangePassword();
  const form = useForm<z.infer<typeof changePasswordFormSchema>>({
    resolver: zodResolver(changePasswordFormSchema),
    defaultValues: {
      oldPassword: "",
      newPassword: "",
      confirmNewPassword: "",
    },
  });

  function onSubmit(values: z.infer<typeof changePasswordFormSchema>) {
    if (!user) {
      toast("Error occured while trying to change password");
      return;
    }

    if (values.newPassword !== values.confirmNewPassword) {
      form.setError(
        "confirmNewPassword",
        { message: "Passwords do not match" },
        {
          shouldFocus: true,
        }
      );
      return;
    }
    changePasswordMtn
      .mutateAsync({
        userId: user.id,
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      })
      .then(() => {
        toast("Password changed successfully!");
      })
      .catch((err: AxiosError<{ error: string }>) => {
        form.setError("oldPassword", {
          message:
            typeof err.response?.data.error === "string"
              ? startSentenceWithUppercase(err.response?.data.error)
              : "Invalid Password",
        });
        toast("Error occured while trying to change password");
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
          name="oldPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Current password</FormLabel>
              <FormControl>
                <Input placeholder={"********"} type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="newPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>New password</FormLabel>
              <FormControl>
                <Input placeholder={"********"} type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="confirmNewPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm new password</FormLabel>
              <FormControl>
                <Input placeholder={"********"} type="password" {...field} />
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
