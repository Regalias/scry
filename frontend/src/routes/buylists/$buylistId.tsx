import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { ListBuylists } from "@wails/go/main/App";

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";

import { BuylistSwitcher } from "@/components/buylistSidebar/switcher";
import { CardSidebarList } from "@/components/buylistSidebar/cardlist";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { PencilOff, ShoppingCart, TableProperties } from "lucide-react";

import "./buylist.css";

export const Route = createFileRoute("/buylists/$buylistId")({
  component: Index,
  loader: async (ctx) => {
    const id = Number(ctx.params.buylistId);
    const buylists = await ListBuylists();
    const selected = buylists.find((bl) => bl.id === id);
    if (!selected) {
      throw new Error(`no such buylist: ${String(id)}`);
    }
    return {
      buylists,
      selected,
    };
  },
  errorComponent: ({ error, reset }) => {
    return (
      <div>
        <pre>{JSON.stringify(error, null, 2)}</pre>
        <button
          onClick={() => {
            // Reset the router error boundary
            reset();
          }}
        >
          retry
        </button>
      </div>
    );
  },
});

function Index() {
  const { buylists, selected } = Route.useLoaderData();

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "22rem",
          "--sidebar-width-mobile": "22rem",
        } as React.CSSProperties
      }
    >
      <Sidebar collapsible="none" variant="sidebar" className="buylist-sidebar">
        <SidebarHeader>
          <BuylistSwitcher buylists={buylists} selectedBuylist={selected} />
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Settings</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    asChild
                    isActive={window.location.pathname.endsWith("/edit")}
                  >
                    <Link
                      to="/buylists/$buylistId/edit"
                      params={{ buylistId: String(selected.id) }}
                    >
                      <PencilOff /> Edit
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    asChild
                    isActive={window.location.pathname.endsWith("/checkout")}
                  >
                    <Link
                      to="/buylists/$buylistId/checkout"
                      params={{ buylistId: String(selected.id) }}
                    >
                      <ShoppingCart /> Summary
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
          <CardSidebarList cards={selected.cards} buylistId={selected.id} />
        </SidebarContent>
        <SidebarRail />
      </Sidebar>
      <SidebarInset className="buylist-sidebar-inset">
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <TableProperties size="18" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink asChild>
                  <Link to="/buylists">Buylists</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>{selected.name}</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <div className="buylist-sidebar-content">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
