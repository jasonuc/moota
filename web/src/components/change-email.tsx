import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useAuth } from "@/hooks/use-auth";
import { startSentenceWithUppercase } from "@/lib/utils";
import { changeEmailFormSchema } from "@/schemas/settings";
import { useChangeEmail } from "@/services/mutations/user";
import { zodResolver } from "@hookform/resolvers/zod";
import { AxiosError } from "axios";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "./ui/button";
import { Input } from "./ui/input";

export default function ChangeEmail() {
  const { user } = useAuth();
  const changeEmailMtn = useChangeEmail();
  const form = useForm<z.infer<typeof changeEmailFormSchema>>({
    resolver: zodResolver(changeEmailFormSchema),
    defaultValues: {
      newEmail: "",
    },
  });

  function onSubmit(values: z.infer<typeof changeEmailFormSchema>) {
    if (!user) {
      toast("Error occured while trying to change email");
      return;
    }
    changeEmailMtn
      .mutateAsync({ userId: user.id, newEmail: values.newEmail })
      .then(() => {
        toast("Email changed successfully");
      })
      .catch((err: AxiosError<{ error: string }>) => {
        form.setError("newEmail", {
          message:
            typeof err.response?.data.error === "string"
              ? startSentenceWithUppercase(err.response?.data.error)
              : "Invalid Email",
        });
        toast.error("Error occurred while trying to change email");
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
          name="newEmail"
          render={({ field }) => (
            <FormItem>
              <FormLabel>New email</FormLabel>
              <FormControl>
                <Input placeholder={user?.email} {...field} />
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
