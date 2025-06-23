import { useAuth } from "@/hooks/use-auth";
import { startSentenceWithUppercase } from "@/lib/utils";
import { useRequestSeeds } from "@/services/mutations/seeds";
import { useCheckWhenUserCanRequestSeeds } from "@/services/queries/seeds";
import { AxiosError } from "axios";
import { formatDate } from "date-fns";
import Countdown from "react-countdown";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import LearnMoreButton from "./learn-more-button";
import { Button } from "./ui/button";

export default function NoSeeds() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const requestSeedsMtn = useRequestSeeds();
  const { data: seedAvailability } = useCheckWhenUserCanRequestSeeds(user?.id);

  const handleRequestSeeds = async () => {
    requestSeedsMtn
      .mutateAsync(user!.id)
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
    <div className="flex flex-col grow items-center justify-center gap-y-5 py-40">
      {!seedAvailability?.availableNow && seedAvailability?.timeAvailable && (
        <Countdown
          date={seedAvailability?.timeAvailable}
          className="text-4xl font-mono font-stretch-ultra-expanded font-semibold"
        />
      )}
      <div className="flex flex-col items-center gap-y-1.5 text-center">
        <h3 className="text-xl font-heading">{"No seeds?"}</h3>
        {seedAvailability?.availableNow ? (
          <p>Go ahead and request some now. I wonder what you'll get!</p>
        ) : (
          <p>That's alright! Try and request more when available.</p>
        )}
      </div>
      <div className="flex flex-col md:flex-row gap-x-5 gap-y-5">
        <Button
          disabled={!seedAvailability?.availableNow}
          onClick={handleRequestSeeds}
        >
          Request seeds
        </Button>
        <LearnMoreButton />
      </div>
    </div>
  );
}
