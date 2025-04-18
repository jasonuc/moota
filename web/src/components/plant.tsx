import { Heart } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Link } from "@tanstack/react-router";

export default function Plant() {
    return (
        <Link to="/dashboard">
            <Card className="w-full min-h-fit gap-y-1.5 active:scale-[98%] transition-all duration-150 bg-secondary">
                <CardContent className="flex justify-end">
                    <div className="flex items-center gap-x-2">
                        <Heart size={15} />
                        50%
                    </div>
                </CardContent>
                <CardHeader>
                    <CardTitle className="text-xl font-bold">Sproutlet</CardTitle>
                    <CardDescription className="italic">Monstera deliciosa</CardDescription>
                </CardHeader>
            </Card >
        </Link>
    )
}