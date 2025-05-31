import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { HouseIcon, Triangle } from "lucide-react";

export default function PlantMap() {
  return (
    <Card className="h-56 relative p-0 w-full mt-10 md:mt-5">
      <CardHeader className="z-50 absolute inset-0 p-0">
        <CardTitle className="p-1.5 w-fit border-b-2 border-r-2 rounded-none text-xs flex items-center justify-center">
          <HouseIcon size={15} className="mr-2" />
          <span>home</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="size-full flex items-center justify-center">
        <Triangle />
      </CardContent>
    </Card>
  );
}
