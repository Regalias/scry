import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useCardFilter } from "@/hooks/cardFilter";
import { useSelectionCallbacks } from "@/hooks/selections";
import { queryClient } from "@/query/client";
import { getAllCardOfferings, getCardOfferingsVendorQuery } from "@/query/ipc";
import { useQueries, useQuery } from "@tanstack/react-query";
import { models } from "@wails/go/models";
import clsx from "clsx";
import { LoaderCircle, RefreshCcw } from "lucide-react";
import { useState } from "react";
import { EmptyBlock } from "../placeholders/empty";
import { ErrorBlock } from "../placeholders/error";
import { LoadingBlock } from "../placeholders/loading";
import { OfferingCard } from "./offeringCard";

export interface CardsTabsProps {
  card: models.Card;
  vendors: string[];
}

export const CardsTabs = ({ card, vendors }: CardsTabsProps) => {
  const { data, isFetching } = useQuery(getAllCardOfferings(card.name));

  const [filter, setFilter] = useState("");

  const queries = useQueries({
    queries: vendors.map((vendor) =>
      getCardOfferingsVendorQuery(card.name, vendor),
    ),
  });

  return (
    <Tabs defaultValue="all">
      <div className="flex flex-row sticky top-0 z-50 p-2 bg-background">
        <TabsList>
          <TabsTrigger value="all" className="min-w-36">
            <span>All</span>
            {isFetching ? (
              <LoaderCircle className="animate-spin ml-1" size={18} />
            ) : (
              <span className="ml-1">({data?.length.toString() ?? "0"})</span>
            )}
          </TabsTrigger>

          {vendors.map((vendor, i) => (
            <TabsTrigger key={vendor} value={vendor} className="min-w-36">
              <span>{vendor}</span>
              {queries[i].isFetching ? (
                <LoaderCircle className="animate-spin ml-1" size={18} />
              ) : (
                <span className="ml-2">
                  <Badge variant="default">
                    {queries[i].data?.length.toString() ?? "0"}
                  </Badge>
                </span>
              )}
            </TabsTrigger>
          ))}
        </TabsList>
        <Input
          placeholder="Filter cards"
          className="w-2/6 focus-visible:ring-transparent"
          value={filter}
          onChange={(e) => {
            setFilter(e.target.value);
          }}
        ></Input>
        <Button
          size="sm"
          className="ml-auto"
          disabled={isFetching}
          onClick={() => {
            void queryClient.invalidateQueries({
              queryKey: ["ipc", "cards", card.name],
            });
          }}
        >
          <RefreshCcw className={clsx(isFetching && "animate-spin")} /> Refresh
        </Button>
      </div>
      <TabsContent value="all">
        <CardList card={card} vendor="all" filter={filter} />
      </TabsContent>
      {vendors.map((vendor) => (
        <TabsContent value={vendor} key={vendor}>
          <CardList card={card} vendor={vendor} filter={filter} />
        </TabsContent>
      ))}
    </Tabs>
  );
};

interface CardListProps {
  card: models.Card;
  vendor: string;
  filter?: string;
}

const CardList = ({ card, vendor, filter = "" }: CardListProps) => {
  const { data, isLoading, isFetching, error, refetch } = useQuery(
    vendor === "all"
      ? getAllCardOfferings(card.name)
      : getCardOfferingsVendorQuery(card.name, vendor),
  );

  const { addOffering, removeOffering } = useSelectionCallbacks(card);
  const filteredCards = useCardFilter(filter, data);

  if (isLoading) return <LoadingBlock />;
  if (error)
    return (
      <ErrorBlock error={error} isFetching={isFetching} retry={void refetch} />
    );
  if (!filteredCards?.length) return <EmptyBlock />;

  return (
    <div className="flex flex-row flex-wrap">
      {filteredCards.map((offering) => {
        const selection = card.selections.find(
          (sel) => sel.offering.productUri === offering.productUri,
        );
        return (
          <OfferingCard
            offering={offering}
            key={offering.productUri}
            selection={selection}
            onSelect={addOffering}
            onDeselect={removeOffering}
          />
        );
      })}
    </div>
  );
};
