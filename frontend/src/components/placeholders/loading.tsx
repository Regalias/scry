import { Loader2 } from "lucide-react";

interface LoadingBlockProps {
  text?: string;
}

export const LoadingBlock = ({
  text = "One moment while we get that for you...",
}: LoadingBlockProps) => {
  return (
    <div className="flex space-x-1 justify-center place-content-center text-base text-muted-foreground size-full p-4">
      <Loader2 className="animate-spin" />
      <span>{text}</span>
    </div>
  );
};
