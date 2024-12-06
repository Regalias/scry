import { ButtonWithConfirm } from "@/components/buttons/buttonWithConfirm";
import { Section } from "@/components/layout/section";
import { useToast } from "@/hooks/use-toast";
import { useSelectionCallbacks } from "@/hooks/selections";
import { useRouter } from "@tanstack/react-router";
import { DeleteSelectionsForCardId } from "@wails/go/main/App";
import { models } from "@wails/go/models";
import { X } from "lucide-react";
import { useCallback } from "react";
import { EmptyBlock } from "../placeholders/empty";
import { OfferingCard } from "./offeringCard";

interface SelectedOfferingsProps {
  card: models.Card;
}

export const SelectedOfferings = ({ card }: SelectedOfferingsProps) => {
  const { toast } = useToast();
  const router = useRouter();

  const clearOfferings = useCallback(() => {
    DeleteSelectionsForCardId(card.id)
      .then(() => {
        toast({
          description: `Cleared selections for '${card.name}'`,
        });
        void router.invalidate();
      })
      .catch((err) => {
        toast({ title: "Error", description: String(err) });
      });
  }, [card.id, toast, router.invalidate]);

  const { addOffering, removeOffering } = useSelectionCallbacks(card);

  return (
    <Section
      header="Selected Offerings"
      buttonGroup={
        <ButtonWithConfirm
          text="Clear Selections"
          icon={<X />}
          onClick={clearOfferings}
          disabled={!card.selections.length}
        />
      }
    >
      <div className="flex flex-wrap items-center space-x-4">
        {card.selections.length ? (
          card.selections.map((card) => {
            return (
              <OfferingCard
                offering={card.offering}
                key={card.id}
                selectedQty={card.quantity}
                onSelect={addOffering}
                onDeselect={removeOffering}
              />
            );
          })
        ) : (
          <EmptyBlock />
        )}
      </div>
    </Section>
  );
};
