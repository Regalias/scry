import { useNavigate } from "@tanstack/react-router";
import { NewBuylist } from "@wails/go/main/App";
import { useCallback } from "react";
import { useToast } from "./use-toast";

export const useNewBuylist = () => {
  const navigate = useNavigate();
  const { toast } = useToast();

  return useCallback(
    async (args: Parameters<typeof NewBuylist>[0]) => {
      const res = await NewBuylist(args);
      toast({
        description: `Created new buylist '${res.name}'`,
      });
      // Since navigation is occuring, router cache invalidation not required
      void navigate({
        to: "/buylists/$buylistId",
        params: { buylistId: String(res.id) },
      });
    },
    [],
  );
};
