import Seed from "./seed";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Link } from "@tanstack/react-router";

type HeaderProps = {
    seedCount: number
};

export default function Header({ seedCount }: HeaderProps) {
    return (
        <div className="flex w-full items-center justify-between">
            <Link to="/dashboard" className="flex items-center space-x-2">
                <img
                    src="./moota.png"
                    alt="Moota"
                    width={45} height={45}
                />
            </Link>

            <div className="flex items-center space-x-5">
                <div className="">
                    <div className="flex items-center space-x-0.5">
                        <Seed number={seedCount} size={35} />
                    </div>
                </div>

                <Avatar>
                    <AvatarImage src={`https://api.dicebear.com/9.x/glass/svg?seed=${"jasonuc"}`} />
                    <AvatarFallback>
                        JA
                    </AvatarFallback>
                </Avatar>
            </div>
        </div>
    )
}