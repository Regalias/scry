import { TriangleAlert } from "lucide-react";

interface LoadingBlockProps {
  text?: string;
}

export const EmptyBlock = ({
  text = "There's nothing here!",
}: LoadingBlockProps) => {
  return (
    <div className="flex space-x-1 justify-center place-content-center text-sm text-muted-foreground size-full p-4">
      <TriangleAlert size={18}/>
      <span>{text}</span>
    </div>
  );
};
