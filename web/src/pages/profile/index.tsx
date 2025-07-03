import GraveyardOnProfile from "@/components/graveyard-on-profile";
import Header from "@/components/header";
import Top3Plants from "@/components/top-3-plants";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { getDicebearGlassUrl } from "@/lib/utils";
import { useGetUserProfile } from "@/services/queries/user";
import { AxiosError } from "axios";
import { HeartIcon, SkullIcon, SproutIcon } from "lucide-react";
import { useEffect } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

const noTitleMessages = [
  "has no title",
  "is yet to learn the ways of the plant",
  "is still finding their roots",
  "hasn't bloomed yet",
  "is a seedling in training",
  "is cultivating their potential",
  "needs more sunlight",
  "is working on their green thumb",
  "is branching out slowly",
  "is photosynthesizing their skills",
  "is still germinating",
  "is pruning their abilities",
  "hasn't found their soil yet",
  "is a work in progress",
  "is planting seeds of knowledge",
  "is watering their dreams",
  "is rooting for themselves",
  "is growing at their own pace",
  "is composting their mistakes",
];

export default function ProfilePage() {
  const { username } = useParams();
  const navigate = useNavigate();
  const { data: profile, error: useGetUserProfileErr } =
    useGetUserProfile(username);

  useEffect(() => {
    if (useGetUserProfileErr) {
      navigate("/home");
      const err = useGetUserProfileErr as AxiosError<{ error: string }>;
      toast.error("User does not exist", {
        description: err.response?.data.error,
        descriptionClassName: "!text-white",
      });
    }
  }, [useGetUserProfileErr, navigate]);

  if (useGetUserProfileErr) return null;

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header />

      <div className="flex flex-col space-y-5 w-full max-w-md mx-auto px-4">
        <img
          src={getDicebearGlassUrl(username)}
          alt="profile"
          className="w-full h-[15rem] object-cover rounded-lg border"
        />

        <div className="space-y-2">
          <h1 className="text-2xl font-bold">{`@${username}`}</h1>

          {profile?.title ? (
            <p className="italic text-lg">{profile.title}</p>
          ) : (
            <p className="italic">
              Untitled{" "}
              <span className="text-gray-600 not-italic">
                {" - "}
                {`This user ${
                  noTitleMessages[
                    Math.floor(Math.random() * noTitleMessages.length)
                  ]
                }`}
              </span>
            </p>
          )}

          <div className="flex items-center gap-2">
            <Badge className="text-sm px-3 py-1">
              Level {profile?.level || 0}
            </Badge>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-3">
          <Card className="p-3.5">
            <div className="flex flex-col items-center space-y-2">
              <HeartIcon className="w-6 h-6 text-green-600" />
              <div className="text-center">
                <p className="text-xl font-bold">
                  {profile?.plantCount.alive || 0}
                </p>
                <p className="text-sm text-gray-600">Alive</p>
              </div>
            </div>
          </Card>

          <Card className="p-3.5">
            <div className="flex flex-col items-center space-y-2">
              <SkullIcon className="w-6 h-6 text-gray-600" />
              <div className="text-center">
                <p className="text-xl font-bold">
                  {profile?.plantCount.deceased || 0}
                </p>
                <p className="text-sm text-gray-600">Deceased</p>
              </div>
            </div>
          </Card>

          <Card className="p-3.5">
            <div className="flex flex-col items-center space-y-2">
              <SproutIcon className="w-6 h-6 text-amber-600" />
              <div className="text-center">
                <p className="text-xl font-bold">
                  {profile?.seedCount.unused || 0}
                </p>
                <p className="text-sm text-gray-600">Seeds</p>
              </div>
            </div>
          </Card>
        </div>

        <Top3Plants profile={profile} />

        <GraveyardOnProfile profile={profile} />
      </div>
    </div>
  );
}
