import { Check, ChevronsUpDown, GalleryVerticalEnd } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { models } from "@wails/go/models";
import { useNavigate } from "@tanstack/react-router";

export function BuylistSwitcher({
  buylists,
  selectedBuylist,
}: {
  buylists: models.Buylist[];
  selectedBuylist: models.Buylist;
}) {
  const navigate = useNavigate();
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <GalleryVerticalEnd className="size-4" />
              </div>
              <div className="flex flex-col gap-0.5 leading-none">
                <span className="font-semibold">{selectedBuylist.name}</span>
                <span>
                  {selectedBuylist.totalCards} cards (
                  {selectedBuylist.cards.length} unique)
                </span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width]"
            align="start"
          >
            {buylists.map((bl) => (
              <DropdownMenuItem
                key={bl.id}
                onSelect={() => {
                  void navigate({
                    to: "/buylists/$buylistId",
                    params: { buylistId: String(bl.id) },
                  });
                }}
              >
                {bl.name} - {bl.totalCards} cards (
                {bl.cards.length} unique)
                {bl.id === selectedBuylist.id && <Check className="ml-auto" />}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
