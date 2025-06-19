import Header from "@/components/header";
import { getUserProfile } from "@/services/api/user";
import { UserProfile } from "@/types/user";
import { AxiosError } from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

// TODO: This page needs work maybe I add a detailed section of the user's stats
export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile>();
  const { username } = useParams();
  const navigate = useNavigate();

  const noTitleMessages = [
    "has no title",
    `is yet to learn the ways of the plant`,
    `is still finding their roots`,
    `hasn't bloomed yet`,
    `is a seedling in training`,
    `is cultivating their potential`,
    `needs more sunlight`,
    `is working on their green thumb`,
    `is branching out slowly`,
    `is photosynthesizing their skills`,
    `is still germinating`,
    `is pruning their abilities`,
    `hasn't found their soil yet`,
    `is a work in progress`,
    `is planting seeds of knowledge`,
    `is watering their dreams`,
    `is rooting for themselves`,
    `is growing at their own pace`,
    `is composting their mistakes`,
  ];

  useEffect(() => {
    getUserProfile(username!)
      .then((data) => {
        setProfile(data);
      })
      .catch((err: AxiosError<{ error: string }>) => {
        navigate("/home");
        toast.error("User does not exist", {
          description: err.response?.data.error,
          descriptionClassName: "!text-white",
        });
      });
  }, [username, navigate]);

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header />

      <div className="flex flex-col space-y-5 w-full max-w-md mx-auto">
        <img
          src={`https://api.dicebear.com/9.x/glass/svg`}
          alt="profile"
          className="w-full h-[15rem] object-cover rounded-md border-2 border-white"
        />

        <div>
          <h1 className="text-2xl font-bold">{`@${username}`}</h1>
          {!profile?.title && (
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
          <p className="italic">{profile?.title}</p>
          <h3 className="text-xl font-semibold mt-1.5">
            Level {profile?.level}
          </h3>
        </div>

        <div className="flex justify-between items-center space-y-2">
          <div className="flex flex-col items-center justify-between space-y-2">
            <div className="relative flex items-center justify-center rounded-md size-12 border-2 border-white">
              <img
                src="https://api.dicebear.com/9.x/glass/svg"
                alt="profile"
                className="absolute object-cover -z-10"
              />
              <h2>🌿</h2>
            </div>
            <p>{profile?.plantCount.alive}</p>
          </div>

          <div className="flex flex-col items-center justify-between space-y-2">
            <div className="relative flex items-center justify-center rounded-md size-12 border-2 border-white">
              <img
                src="https://api.dicebear.com/9.x/glass/svg"
                alt="profile"
                className="absolute object-cover -z-10"
              />
              <h2>🪦</h2>
            </div>
            <p>{profile?.plantCount.deceased}</p>
          </div>

          <div className="flex flex-col items-center justify-between space-y-2">
            <div className="relative flex items-center justify-center rounded-md size-12 border-2 border-white">
              <img
                src="https://api.dicebear.com/9.x/glass/svg"
                alt="profile"
                className="absolute object-cover -z-10"
              />
              <h2>🫘</h2>
            </div>
            <p>{profile?.seedCount.unused}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
