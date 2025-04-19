import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { HomeIcon } from "lucide-react";

export default function GeneralHeader() {
    return (
        <div className="flex w-full items-center justify-between">
            <div className="">
                <div>
                    <HomeIcon />
                </div>
            </div>

            <Avatar>
                <AvatarImage src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${"jasonuc"}`} />
                <AvatarFallback>
                    JA
                </AvatarFallback>
            </Avatar>
        </div>
    )
}