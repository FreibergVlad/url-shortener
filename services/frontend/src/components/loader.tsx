import { cn } from "@/lib/shadcn-utils";
import { LoaderCircle } from "lucide-react";

export default function Loader({ className }: React.ComponentPropsWithoutRef<"div">) {
    return (
        <div className="flex justify-center items-center h-screen w-screen">
            <LoaderCircle
                className={cn('h-16 w-16 animate-spin', className)}
            />
        </div>
      );
}