import ChangeEmail from "@/components/change-email";
import ChangePassword from "@/components/change-password";
import ChangeUsername from "@/components/change-username";
import Header from "@/components/header";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { useAuth } from "@/hooks/use-auth";

export default function SettingsPage() {
  const { logout } = useAuth();

  return (
    <div className="flex flex-col space-y-5 pb-10 w-full">
      <Header />

      <h1 className="text-3xl font-heading mb-5">Settings</h1>

      <Accordion type="single" collapsible className="w-full">
        <AccordionItem value="item-1">
          <AccordionTrigger>Change username</AccordionTrigger>
          <AccordionContent>
            <ChangeUsername />
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-2">
          <AccordionTrigger>Change email</AccordionTrigger>
          <AccordionContent>
            <ChangeEmail />
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-3">
          <AccordionTrigger>Change password</AccordionTrigger>
          <AccordionContent>
            <ChangePassword />
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-4">
          <AccordionTrigger>Request a title</AccordionTrigger>
          <AccordionContent>Coming soon</AccordionContent>
        </AccordionItem>
        <AccordionItem value="item-5">
          <AccordionTrigger
            hideChevron
            className="hover:cursor-pointer"
            onClick={() => logout()}
          >
            Logout
          </AccordionTrigger>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
