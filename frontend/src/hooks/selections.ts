import { AddSelection, UpdateSelection } from "@wails/go/main/App";
import { buylist, models } from "@wails/go/models";
import { useCallback } from "react";
import { useToast } from "./use-toast";
import { useRouter } from "@tanstack/react-router";

export const useSelectionCallbacks = (card: models.Card) => {
  const { toast } = useToast();
  const router = useRouter();
  const addOffering = useCallback(
    (offering: models.Offering) => {
      const toAdd = 1;

      const existing = card.selections.find(
        (sel) => sel.offering.productUri === offering.productUri,
      );
      if (existing && existing.quantity + toAdd > offering.quantity) {
        toast({
          description: `Not enough stock for ${offering.name}`,
          variant: "destructive",
        });
        return;
      }
      const promise = existing
        ? UpdateSelection(
            new buylist.UpdateSelectionRequest({
              selectionId: existing.id,
              quantity: existing.quantity + toAdd,
              offering,
            }),
          )
        : AddSelection(card.id, offering, toAdd);

      promise
        .then(() => {
          toast({
            description: `Added ${toAdd.toString()} copy of ${offering.name}`,
          });
          void router.invalidate();
        })
        .catch((err) => {
          toast({
            title: "Error",
            description: String(err),
            variant: "destructive",
          });
        });
    },
    [card, router.invalidate, toast],
  );

  const removeOffering = useCallback(
    (offering: models.Offering) => {
      const existing = card.selections.find(
        (sel) => sel.offering.productUri === offering.productUri,
      );
      if (!existing) {
        // Do nothing
        return;
      }

      const toRemove = 1;
      const promise = UpdateSelection(
        new buylist.UpdateSelectionRequest({
          selectionId: existing.id,
          quantity: existing.quantity - toRemove,
          offering,
        }),
      );

      promise
        .then(() => {
          toast({
            description: `Removed ${toRemove.toString()} copy of ${offering.name}`,
          });
          void router.invalidate();
        })
        .catch((err) => {
          toast({
            title: "Error",
            description: String(err),
            variant: "destructive",
          });
        });
    },
    [card, router.invalidate, toast],
  );
  return { addOffering, removeOffering };
};
