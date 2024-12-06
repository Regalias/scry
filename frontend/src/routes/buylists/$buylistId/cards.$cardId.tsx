import { ButtonWithConfirm } from "@/components/buttons/buttonWithConfirm";
import { CardsTabs } from "@/components/cards/cards";
import { SelectedOfferings } from "@/components/cards/selectedRow";
import { useToast } from "@/hooks/use-toast";
import {
  createFileRoute,
  useNavigate,
  useRouter,
} from "@tanstack/react-router";
import { DeleteCards, GetBuylist, ListVendors } from "@wails/go/main/App";
import { useCallback } from "react";

export const Route = createFileRoute("/buylists/$buylistId/cards/$cardId")({
  component: RouteComponent,
  loader: async (ctx) => {
    const buylist = await GetBuylist(Number(ctx.params.buylistId));
    const card = buylist.cards.find(
      (card) => card.id === Number(ctx.params.cardId),
    );
    if (!card) throw new Error("card not found in buylist!");
    const vendors = await ListVendors();
    return { buylist, card, vendors };
  },
});

function RouteComponent() {
  const { buylist, card, vendors } = Route.useLoaderData();

  const { toast } = useToast();
  const router = useRouter();
  const navigate = useNavigate();
  const removeCard = useCallback(() => {
    const promise = DeleteCards([card.id]);
    promise
      .then(() => {
        toast({
          description: `Deleted card '${card.name}' from buylist`,
        });
        void router.invalidate();
        void navigate({
          to: "/buylists/$buylistId",
          params: { buylistId: String(buylist.id) },
        });
      })
      .catch((err) => {
        toast({ title: "Error", description: String(err) });
      });
  }, [card, router.invalidate, toast]);

  return (
    <div className="p-5 flex flex-col gap-4">
      <div className="flex flex-row">
        <h2 className="mt-10 scroll-m-20 pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0">
          {card.name}
        </h2>
        <div className="ml-auto">
          <ButtonWithConfirm onClick={removeCard} />
        </div>
      </div>

      <SelectedOfferings card={card} />
      <CardsTabs card={card} vendors={vendors} />
    </div>
  );
}
