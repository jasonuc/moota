import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import Seed from "./seed";

type DashboardHeaderProps = {
    seedCount: number
};

export default function DashboardHeader({ seedCount }: DashboardHeaderProps) {
    return (
        <div className="flex w-full items-center justify-between">
            <div className="">
                <div className="flex items-center space-x-0.5">
                    <Seed number={seedCount} size={40} />
                </div>
            </div>

            <Avatar>
                <AvatarImage src={`https://api.dicebear.com/9.x/thumbs/svg?seed=jasonuc`} />
                <AvatarFallback>
                    JA
                </AvatarFallback>
            </Avatar>
        </div>
    )
}