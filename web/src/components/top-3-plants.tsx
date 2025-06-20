import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { UserProfile } from "@/types/user";
import { HeartIcon } from "lucide-react";
import { Badge } from "./ui/badge";

type Top3PlantsProps = {
  profile?: UserProfile;
};

export default function Top3Plants({ profile }: Top3PlantsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <HeartIcon className="w-5 h-5 text-green-600" />
          Top {Math.min(profile?.top3AlivePlants?.length || 0, 3)} Plants
        </CardTitle>
      </CardHeader>
      <CardContent>
        {profile?.top3AlivePlants.slice(0, 3) &&
        profile?.top3AlivePlants.slice(0, 3).length > 0 ? (
          <div className="space-y-3">
            {profile.top3AlivePlants.map((plant, index) => (
              <div
                key={plant.id}
                className="flex items-center justify-between p-3 border rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <Badge>#{index + 1}</Badge>
                  <div>
                    <p className="font-medium">{plant.nickname}</p>
                    <p className="text-sm text-gray-600">
                      HP: {plant.hp.toFixed(1)}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-600 italic text-center py-4">No plants yet</p>
        )}
      </CardContent>
    </Card>
  );
}
