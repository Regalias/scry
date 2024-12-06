import { Search } from "lucide-react";

import { Label } from "@/components/ui/label";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link } from "@tanstack/react-router";
import { models } from "@wails/go/models";
import clsx from "clsx";
import { useMemo, useState } from "react";

import { queryClient } from "@/query/client";
import { getAllCardOfferings } from "@/query/ipc";
import { Price } from "../price/priceDisplay";
import "./cardlist.css";

export interface CardSidebarListProps {
  buylistId: number;
  cards: models.Card[] | undefined;
}

export function CardSidebarList({ buylistId, cards }: CardSidebarListProps) {
  const [filterText, setFilterText] = useState("");
  const visibleCards = useMemo(() => {
    if (!filterText) {
      return cards;
    }
    const filtered = cards?.filter((card) => {
      return card.name.toLowerCase().includes(filterText.toLowerCase());
    });

    if (filtered?.length === 0) {
      return undefined;
    }
    return filtered;
  }, [cards, filterText]);

  return (
    <>
      <SidebarGroup>
        <form autoComplete="off">
          <SidebarGroup className="py-0">
            <SidebarGroupContent className="relative">
              <Label htmlFor="search" className="sr-only">
                Search
              </Label>
              <SidebarInput
                id="search"
                placeholder="Search..."
                className="pl-8"
                value={filterText}
                onChange={(e) => {
                  setFilterText(e.target.value);
                }}
              />
              <Search className="pointer-events-none absolute left-2 top-1/2 size-4 -translate-y-1/2 select-none opacity-50" />
            </SidebarGroupContent>
          </SidebarGroup>
        </form>
      </SidebarGroup>
      <SidebarGroup>
        <SidebarGroupLabel>Cards</SidebarGroupLabel>
        <SidebarGroupContent className="overscroll-auto">
          <SidebarMenu>
            {visibleCards?.map((card) => {
              const hasFlagged = !!card.selections.find((sel) => sel.isFlagged);

              let selectionClass = "selected-none";
              if (card.totalSelections >= card.quantity) {
                selectionClass = "selected-all";
              } else if (card.totalSelections > 0) {
                selectionClass = "selected-partial";
              }

              const prefetch = () => {
                void queryClient.prefetchQuery(getAllCardOfferings(card.name));
              };

              return (
                <SidebarMenuItem key={card.id}>
                  <SidebarMenuButton
                    asChild
                    variant="outline"
                    isActive={window.location.pathname.endsWith(
                      `/cards/${String(card.id)}`,
                    )}
                    onFocus={prefetch}
                    onMouseEnter={prefetch}
                  >
                    <Link
                      to="/buylists/$buylistId/cards/$cardId"
                      params={{
                        buylistId: String(buylistId),
                        cardId: String(card.id),
                      }}
                      className={clsx(
                        "grid grid-cols-12 cardlist",
                        selectionClass,
                        hasFlagged && "selected-flagged",
                      )}
                    >
                      <span className="col-span-8 overflow-hidden text-ellipsis whitespace-nowrap">
                        {card.name}
                      </span>
                      <span className="col-span-2">
                        <Price price={card.totalPrice} />
                      </span>
                      <span className="col-span-2">
                        {card.totalSelections} / {card.quantity}
                      </span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              );
            }) ?? (
              <SidebarMenuItem>
                <span className="pl-4">There&apos;s nothing here!</span>
              </SidebarMenuItem>
            )}
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </>
  );
}
