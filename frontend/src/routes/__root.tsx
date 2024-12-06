import { ModeToggle } from "@/components/theme/modeToggle";
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
// import { TanStackRouterDevtools } from "@tanstack/router-devtools";

import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import { Toaster } from "@/components/ui/toaster";
import { IsError, IsReady } from "@wails/go/main/App";
import { useEffect, useState } from "react";
import { ErrorBlock } from "@/components/placeholders/error";

export const Route = createRootRoute({
  component: Root,
});

function Root() {
  const [ready, setReady] = useState(false);
  const [hasError, setHasError] = useState(false);

  useEffect(() => {
    const poller = setInterval(() => {
      void IsReady().then((ready) => {
        if (ready) {
          setReady(true);
          clearInterval(poller);
        }
      });
      void IsError().then((err) => {
        if (err) {
          setHasError(true);
          clearInterval(poller);
        }
      });
    }, 100);
  }, [setReady, setHasError]);

  if (hasError) {
    return <ErrorBlock error="Something went wrong" />;
  }

  if (!ready) {
    return;
  }

  return (
    <div>
      <div className="bg-zinc-900 fg-zinc-200 border-b">
        <div className="p-2 pl-8 pr-8">
          <div className="flex flex-row gap-4">
            <span className="text-3xl select-none">&#128184;</span>
            <div>
              <NavMenu />
            </div>
            <div className="ml-auto">
              <ModeToggle />
            </div>
          </div>
        </div>
      </div>
      <div
        className="content"
        style={{
          height: "calc(100vh - 64px)",
        }}
      >
        <Outlet />
      </div>
      <Toaster />
      {/* <TanStackRouterDevtools /> */}
    </div>
  );
}

const NavMenu = () => {
  return (
    <NavigationMenu>
      <NavigationMenuList>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link to="/">Home</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link to="/buylists">Buylists</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link to="/search">Search</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  );
};
