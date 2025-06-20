import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatRelativeTime } from "@/lib/utils";
import { UserProfile } from "@/types/user";
import { GhostIcon } from "lucide-react";

type GraveyardOnProfileProps = {
  profile?: UserProfile;
};

export default function GraveyardOnProfile({
  profile,
}: GraveyardOnProfileProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <GhostIcon className="w-5 h-5 text-gray-600" />
          Plant Graveyard
        </CardTitle>
      </CardHeader>
      <CardContent>
        {profile?.deceasedPlants && profile.deceasedPlants.length > 0 ? (
          <div className="space-y-3">
            {profile.deceasedPlants.slice(0, 5).map((plant) => (
              <div
                key={plant.id}
                className="flex items-center justify-between p-3 border rounded-lg"
              >
                <div>
                  <p className="font-medium">{plant.nickname}</p>
                  <p className="text-sm text-gray-600">
                    {formatRelativeTime(plant.timeOfDeath)}
                  </p>
                </div>
                <div className="text-2xl">🪦</div>
              </div>
            ))}
            {profile.deceasedPlants.length > 5 && (
              <p className="text-sm text-gray-600 text-center italic pt-2 border-t">
                And {profile.deceasedPlants.length - 5} more...
              </p>
            )}
          </div>
        ) : (
          <p className="text-gray-600 italic text-center py-4">
            No deceased plants
          </p>
        )}
      </CardContent>
    </Card>
  );
}
