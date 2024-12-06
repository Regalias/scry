import { useRouter } from "@tanstack/react-router";
import { DeleteSelection, UpdateSelection } from "@wails/go/main/App";
import { useCallback } from "react";
import { useToast } from "./use-toast";
import { models } from "@wails/go/models";

export const useSelectionStatusCallback = () => {
  const router = useRouter();
  const { toast } = useToast();
  return useCallback((req: Parameters<typeof UpdateSelection>[0]) => {
    UpdateSelection(req)
      .then(() => {
        void router.invalidate();
      })
      .catch((err) => {
        toast({ description: String(err), variant: "destructive" });
      });
  }, []);
};

export const useSelectionRemoveCallback = () => {
  const router = useRouter();
  const { toast } = useToast();
  return useCallback((selection: models.ProductSelection) => {
    DeleteSelection(selection.id)
      .then(() => {
        toast({
          description: `Removed selection ${selection.offering.name}`,
        });
        void router.invalidate();
      })
      .catch((err) => {
        toast({ description: String(err), variant: "destructive" });
      });
  }, []);
};
