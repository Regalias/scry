import { CircleX, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

export interface ErrorBlockProps {
  error?: unknown;
  isFetching?: boolean;
  retry?: () => void;
}

export const ErrorBlock = ({ isFetching, error, retry }: ErrorBlockProps) => {
  return (
    <div className="flex space-x-1 justify-center place-content-center text-base text-red-600 size-full p-4">
      <CircleX color="red" />
      <span>
        {typeof error === "string"
          ? error
          : JSON.stringify(error) || "Something went wrong"}
      </span>
      {retry && (
        <div>
          <Button onClick={retry} disabled={isFetching}>
            {isFetching && <Loader2 className="animate-spin" />}
            Retry
          </Button>
        </div>
      )}
    </div>
  );
};
