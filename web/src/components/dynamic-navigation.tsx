import { cn } from "@/lib/utils";
import { BeanIcon, HomeIcon, SproutIcon } from "lucide-react";
import { useLocation, useNavigate } from "react-router";
import { Button } from "./ui/button";

export default function DynamicNavigation() {
  const navigate = useNavigate();
  const { pathname } = useLocation();

  return (
    <div
      className={cn("fixed left-0 bottom-0 w-full p-5", {
        "flex gap-x-5 md:p-10": pathname === "/home",
        "flex justify-end p-5": pathname !== "/",
      })}
    >
      {pathname === "/home" && (
        <>
          <Button
            onClick={() => navigate("/plants")}
            className="hover:cursor-pointer grow md:min-h-12"
          >
            Plants
            <SproutIcon className="ml-2" />
          </Button>

          <Button
            onClick={() => navigate("/seeds")}
            className="hover:cursor-pointer grow md:min-h-12"
          >
            Seeds
            <BeanIcon className="ml-2" />
          </Button>
        </>
      )}

      {pathname !== "/home" && (
        <Button
          className="rounded-full size-12 z-50"
          onClick={() => navigate("/home")}
        >
          <HomeIcon />
        </Button>
      )}
    </div>
  );
}
