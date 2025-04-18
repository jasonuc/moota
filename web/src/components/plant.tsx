import { Heart } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "./ui/button";

type PlantProps = {
    nickname: string
    botanicalName: string
    hp: number
}

export default function Plant({ nickname, botanicalName, hp }: PlantProps) {
    return (
        <Button asChild className="relative overflow-hidden group h-fit">
            <Card className="gap-y-1.5 bg-background flex w-full">
                <CardContent className="flex justify-end w-full p-0 pt-1.5">
                    <div className="flex items-center gap-x-2">
                        <Heart size={15} />
                        {hp}%
                    </div>
                </CardContent>
                <CardHeader className="w-full p-0 pb-3">
                    <CardTitle className="text-lg font-bold">{nickname}</CardTitle>
                    <CardDescription className="italic">{botanicalName}</CardDescription>
                </CardHeader>

                <img
                    src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${nickname}&backgroundColor=${"transparent"}`}
                    alt="avatar"
                    draggable={false}
                    className="size-20 absolute -bottom-1/4 right-0 group-hover:-bottom-5 group-active:-bottom-5 transition-all duration-100 rounded-l-md"
                />
            </Card >
        </Button>
    )
}