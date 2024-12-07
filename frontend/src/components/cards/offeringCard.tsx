import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { models } from "@wails/go/models";
import { useMemo } from "react";
import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

import { Link } from "@tanstack/react-router";
import clsx from "clsx";
import { Price } from "../price/priceDisplay";
import "./offeringCard.css";

export interface OfferingCardProps {
  offering: models.Offering;
  selection?: models.ProductSelection;
  size?: OfferingCardSize;
  previous?: models.Offering;
  onSelect?: (offering: models.Offering) => void;
  onDeselect?: (offering: models.Offering) => void;
}

const sizeClasses = {
  small: {
    //x7
    width: "w-[196px]",
    height: "h-[273px]",
  },
  med: {
    //x8
    width: "w-[224px]",
    height: "h-[312px]",
  },
  large: {
    //x10
    width: "w-[280px]",
    height: "h-[390px]",
  },
};

const vendorColorMap = {
  gamescube: "rgb(254, 0, 0)",
  mtgmate: "rgb(249, 167, 43)",
  goodgames: "rgb(122, 137, 255)",
};

export type OfferingCardSize = keyof typeof sizeClasses;

export const OfferingCard = ({
  offering,
  selection,
  size = "small",
  onDeselect,
  onSelect,
}: OfferingCardProps) => {
  const borderType = useMemo(() => {
    if (!offering.properties.length) {
      return;
    }

    if (offering.properties.find((prop) => prop.includes("Etched"))) {
      return "Etched";
    }

    if (offering.properties.find((prop) => prop.includes("Foil"))) {
      return "Foil";
    }
    return "none";
  }, [offering.properties]);

  const selectedQty = selection?.quantity ?? 0;

  let border = "border-secondary";
  if (selection?.isFlagged) {
    border = "border-rose-500";
  } else if (selection?.isPurchased) {
    border = "border-green-500";
  } else if (selectedQty > 0) {
    border = "border-amber-500";
  }

  return (
    <div
      className={clsx(
        border,
        // sizeClasses[size].width,
        "w-[280px] m-2 rounded-xl border",
      )}
    >
      <div className="rounded-t-xl bg-stone-800 text-stone-300 mb-4 h-8 flex justify-center content-center">
        <div className="self-center flex flex-row space-x-4 items-center">
          <span className="text-lg text-stone-50">
            ${(offering.price / 100).toString()} AUD
          </span>
          <span
            className="text-sm"
            style={{
              color:
                vendorColorMap[
                  offering.vendorId as keyof typeof vendorColorMap
                ],
            }}
          >
            {offering.vendorId}
          </span>
        </div>
      </div>

      <div>
        <div className="flex justify-center">
          <CardImage
            offering={offering}
            showVendor
            size={size}
            borderType={borderType}
            onSelect={() => {
              onSelect?.(offering);
            }}
            onDeselect={() => {
              onDeselect?.(offering);
            }}
          />
        </div>
        <div>
          <div className="p-4 grid grid-cols-1 gap-2 justify-items-center">
            <span className="text-muted-foreground text-xs">
              {offering.set}
            </span>
            <Link
              className="text-base font-medium leading-none"
              onClick={(e) => {
                e.preventDefault();
                BrowserOpenURL(offering.productUri);
              }}
            >
              {offering.name}
            </Link>

            <div>
              <Badge className="mr-2" variant="destructive">
                <Price price={offering.price} />
              </Badge>
              <Badge
                variant={selectedQty > 0 ? "secondary" : "outline"}
                className={clsx(selectedQty > 0 && "border border-amber-800")}
              >
                {selectedQty} / {offering.quantity}
              </Badge>
            </div>

            <div>
              {offering.condition && (
                <Badge variant="outline">{offering.condition}</Badge>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

interface CardImageProps {
  offering: models.Offering;
  showVendor?: boolean;
  size: OfferingCardSize;
  borderType?: "Foil" | "Etched" | "none";

  onSelect?: () => void;
  onDeselect?: () => void;
}

const CardImage = ({
  offering,
  size,
  borderType = "none",
  onSelect,
  onDeselect,
}: CardImageProps) => {
  const imageUri = offering.imgUri;
  return (
    <div
      onClick={onSelect}
      onContextMenu={(e) => {
        e.preventDefault();
        onDeselect?.();
      }}
      className={clsx(
        "offering-card cursor-pointer relative",
        borderType !== "none" && "offering-card-border",
        borderType.toLowerCase(),
        sizeClasses[size].width,
        // sizeClasses[size].height,
        "aspect-[28/39]",
      )}
    >
      <div className={clsx("offering-card-contents")}>
        {imageUri && (
          <img
            src={imageUri}
            referrerPolicy="no-referrer"
            className={clsx(
              sizeClasses[size].width,
              "aspect-[28/39]",
              "transition ease-in-out hover:scale-110 duration-300",
              "offering-card-image relative rounded-xl",
            )}
          />
        )}

        {!imageUri && (
          <Skeleton
            className={clsx(sizeClasses[size].width, "aspect-[28/39]")}
          />
        )}
        {borderType !== "none" && (
          <div className="overlay absolute top-0 left-0 opacity-100 bg-transparent transition ease duration-200 w-full">
            <div className="p-2 grid grid-cols-1 gap-2">
              <Badge variant="secondary" className="justify-self-end rounded">
                {borderType}
              </Badge>
            </div>
          </div>
        )}

        {/* <div className="vendor-overlay absolute bottom-0 left-0 opacity-100 bg-transparent transition ease duration-200 w-full">
          <div className="pb-4 pr-2 grid grid-cols-1">
            {showVendor && (
              <Badge
                className="justify-self-end bg-stone-700 text-stone-300"
                variant="secondary"
              >
                {offering.vendorId}
              </Badge>
            )}
          </div>
        </div> */}
      </div>
    </div>
  );
};
