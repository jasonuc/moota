import Header from "@/components/header";
import { useParams } from "react-router";

export default function ProfilePage() {
  const { username } = useParams();
  const level = 4;
  const title = "The citrus mage";
  const plantsAlive = 6;
  const plantsDead = 6;
  const seedCount = 5;

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header seedCount={10} />

      <div className="flex flex-col space-y-5 w-full max-w-md mx-auto">
        <img
          src={`https://api.dicebear.com/9.x/glass/svg`}
          alt="profile"
          className="w-full h-[15rem] object-cover rounded-md border-2 border-white"
        />

        <div>
          <h1 className="text-2xl font-bold">{`@${username}`}</h1>
          {/* <p className='italic'>Untitled <span className='text-gray-600 not-italic'>{" - User has no title"}</span></p> */}
          <p className="italic">{title}</p>
          <h3 className="text-xl font-semibold mt-1.5">Level {level}</h3>
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
            <p>{plantsAlive}</p>
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
            <p>{plantsDead}</p>
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
            <p>{seedCount}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
