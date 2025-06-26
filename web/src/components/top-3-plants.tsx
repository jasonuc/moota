import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { UserProfile } from "@/types/user";
import { HeartIcon } from "lucide-react";
import { Link } from "react-router";

type Top3PlantsProps = {
  profile?: UserProfile;
};

export default function Top3Plants({ profile }: Top3PlantsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <HeartIcon className="w-5 h-5 text-green-600" />
          {Math.min(profile?.top3AlivePlants?.length || 0, 3) > 0
            ? `Top ${Math.min(profile?.top3AlivePlants?.length || 0, 3)} Plant${
                Math.min(profile?.top3AlivePlants?.length || 0, 3) > 1
                  ? "s"
                  : ""
              }`
            : "Top 3 Plants"}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {profile?.top3AlivePlants.slice(0, 3) &&
        profile?.top3AlivePlants.slice(0, 3).length > 0 ? (
          <div className="space-y-3">
            {profile.top3AlivePlants.map((plant, index) => (
              <Link to={`/plants/${plant.id}/public`} key={plant.id}>
                <div className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    <Badge>#{index + 1}</Badge>
                    <div>
                      <p className="font-medium">{plant.nickname}</p>
                      <div className="flex justify-start items-center gap-x-2">
                        <p className="text-sm text-gray-700">
                          Level: {plant.level}
                        </p>
                        <span className="text-gray-400">•</span>
                        <p className="text-sm text-gray-700">
                          HP: {plant.hp.toFixed(1)}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <p className="text-gray-600 italic text-center py-4">No plants yet</p>
        )}
      </CardContent>
    </Card>
  );
}
