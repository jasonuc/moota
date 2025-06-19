import { Link } from "react-router";
import { Button } from "./ui/button";

export default function LearnMoreButton() {
  return (
    <Link to="/">
      <Button variant="neutral" className="w-full">
        <p className="text-blue-600">Learn more</p>
      </Button>
    </Link>
  );
}
