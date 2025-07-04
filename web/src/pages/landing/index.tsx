import Footer from "@/components/footer";
import { LogoWithText } from "@/components/logo";
import { Button } from "@/components/ui/button";
import UserButton from "@/components/user-button";
import { useAuth } from "@/hooks/use-auth";
import { useStats } from "@/hooks/use-stats";
import {
  CompassIcon,
  DumbbellIcon,
  Gamepad2Icon,
  MapPinIcon,
  TargetIcon,
  TrophyIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";

export default function LandingPage() {
  const [isVisible, setIsVisible] = useState(false);
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    setIsVisible(true);
  }, []);

  const { stats } = useStats();

  const bannerMessages = [
    "Location-Based Adventure Game",
    "Touch Grass While Tucked In",
  ];

  const features = [
    {
      icon: MapPinIcon,
      title: "Plant & Explore",
      desc: "Plant virtual seeds anywhere in the real world and watch your garden grow across your city",
    },
    {
      icon: CompassIcon,
      title: "Adventure Awaits",
      desc: "Return to your planted locations to water and evolve your virtual plants",
    },
    {
      icon: TrophyIcon,
      title: "Level Up Life",
      desc: "Turn real-world habits into epic quests - gym visits, study sessions, daily walks",
    },
  ];

  const useCases = [
    {
      icon: TargetIcon,
      title: "Daily Habits",
      desc: "Transform any routine into an adventure - work, coffee shops, study spots",
    },
    {
      icon: DumbbellIcon,
      title: "Gym Consistency",
      desc: "Plant a seed at your gym and water it during each visit to keep it alive",
    },
  ];

  return (
    <div className="min-h-screen w-full bg-background font-archivo">
      <div className="flex flex-col min-h-screen">
        <header
          className={`flex w-full justify-between items-center p-8 transition-all duration-1000 ${
            isVisible
              ? "translate-y-0 opacity-100"
              : "-translate-y-full opacity-0"
          }`}
        >
          <Link
            to="/"
            className="transform hover:scale-105 transition-transform duration-300"
          >
            <LogoWithText />
          </Link>

          <div className="flex space-x-4">
            {isLoggedIn && (
              <Link to="/home">
                <Button className="bg-secondary-background text-foreground border-2 border-border hover:bg-main hover:text-main-foreground transition-all duration-300 font-medium rounded-[10px] shadow-shadow hover:translate-x-1 hover:translate-y-1 hover:shadow-none">
                  Home
                </Button>
              </Link>
            )}

            {isLoggedIn && <UserButton />}

            {!isLoggedIn && (
              <Link to="/login">
                <Button className="bg-main text-main-foreground border-2 border-border hover:bg-secondary-background hover:text-foreground transition-all duration-300 font-medium rounded-[10px] shadow-shadow hover:translate-x-1 hover:translate-y-1 hover:shadow-none">
                  Login
                </Button>
              </Link>
            )}
          </div>
        </header>

        <main className="flex-1 flex flex-col items-center justify-center px-8 py-16">
          <div
            className={`text-center max-w-6xl transition-all duration-1000 delay-300 ${
              isVisible
                ? "translate-y-0 opacity-100"
                : "translate-y-12 opacity-0"
            }`}
          >
            <div className="inline-flex items-center space-x-2 bg-main text-main-foreground px-6 py-3 rounded-[10px] mb-8 border-2 border-border shadow-shadow font-medium">
              <Gamepad2Icon className="w-5 h-5" />
              <span>
                {
                  bannerMessages[
                    Math.floor(Math.random() * bannerMessages.length)
                  ]
                }
              </span>
            </div>

            <h1 className="text-6xl md:text-8xl font-heading text-foreground mb-6 leading-tight">
              Plant Your
              <br />
              <span className="text-main relative">
                Adventure
                <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-32 h-1 bg-main rounded-full" />
              </span>
            </h1>

            <p className="text-xl md:text-2xl text-foreground mb-12 max-w-4xl mx-auto leading-relaxed font-medium opacity-80">
              Turn the real world into your playground. Plant virtual seeds
              anywhere, return to care for them, and transform daily routines
              into epic gaming adventures.
            </p>

            <div className="text-center mb-12">
              <h2 className="text-4xl font-heading text-foreground mb-4">
                {stats.plant.alive} plants alive
              </h2>
            </div>

            <div className="flex flex-col sm:flex-row gap-6 justify-center mb-20">
              <Button
                onClick={() =>
                  isLoggedIn ? navigate("/home") : navigate("/register")
                }
                className="bg-main text-main-foreground border-2 border-border hover:bg-secondary-background hover:text-foreground transition-all duration-300 font-medium rounded-[10px] shadow-shadow hover:translate-x-1 hover:translate-y-1 hover:shadow-none px-12 py-6 text-xl"
              >
                {!isLoggedIn ? "Start Your Journey" : "Continue Your Journey"}
              </Button>
            </div>
          </div>

          <div
            className={`grid md:grid-cols-3 gap-8 max-w-6xl w-full mb-20 transition-all duration-1000 delay-700 ${
              isVisible
                ? "translate-y-0 opacity-100"
                : "translate-y-12 opacity-0"
            }`}
          >
            {features.map((feature, index) => (
              <div
                key={index}
                className="group bg-secondary-background p-8 rounded-[10px] border-2 border-border shadow-shadow hover:translate-x-1 hover:translate-y-1 hover:shadow-none transition-all duration-300 hover:bg-main hover:text-main-foreground"
              >
                <div className="w-16 h-16 bg-main text-main-foreground rounded-[10px] flex items-center justify-center mb-6 border-2 border-border shadow-shadow group-hover:bg-secondary-background group-hover:text-foreground transition-all duration-300">
                  <feature.icon className="w-8 h-8" />
                </div>
                <h3 className="text-xl font-heading mb-3">{feature.title}</h3>
                <p className="font-medium opacity-80">{feature.desc}</p>
              </div>
            ))}
          </div>

          <div
            className={`w-full max-w-6xl transition-all duration-1000 delay-900 ${
              isVisible
                ? "translate-y-0 opacity-100"
                : "translate-y-12 opacity-0"
            }`}
          >
            <div className="text-center mb-12">
              <h2 className="text-4xl font-heading text-foreground mb-4">
                Level Up Your Life
              </h2>
              <p className="text-xl text-foreground font-medium opacity-80">
                Turn everyday activities into gaming achievements
              </p>
            </div>

            <div className="grid md:grid-cols-2 gap-8">
              {useCases.map((useCase, index) => (
                <div
                  key={index}
                  className="group bg-secondary-background p-8 rounded-[10px] border-2 border-border shadow-shadow hover:translate-x-1 hover:translate-y-1 hover:shadow-none transition-all duration-300"
                >
                  <div className="w-12 h-12 bg-main text-main-foreground rounded-[10px] flex items-center justify-center mb-6 border-2 border-border shadow-shadow group-hover:scale-110 transition-transform duration-300">
                    <useCase.icon className="w-6 h-6" />
                  </div>
                  <h3 className="text-lg font-heading mb-3">{useCase.title}</h3>
                  <p className="font-medium opacity-80 text-sm">
                    {useCase.desc}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </main>
        <Footer />
      </div>
    </div>
  );
}
