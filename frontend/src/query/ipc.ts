import {
  GetCardOfferingForVendor,
  ListVendors,
} from "../../wailsjs/go/main/App";
import { queryClient } from "./client";

export type CardOfferingFormat = "vendor" | "aggregate";

export const getCardOfferingsVendorQuery = (
  cardName: string,
  vendor: string,
) => ({
  queryKey: ["ipc", "cards", cardName, vendor],
  queryFn: async () => {
    const result = await GetCardOfferingForVendor(cardName, vendor);
    return result.sort((a, b) => a.price - b.price);
  },
  enabled: !!(cardName && vendor),
});

export const getVendors = () => ({
  queryKey: ["ipc", "GetVendors"],
  queryFn: () => ListVendors(),
});

export const getAllCardOfferings = (cardName: string) => ({
  queryKey: ["ipc", "cards", cardName, "all"],
  queryFn: async () => {
    const vendors = await ListVendors();
    const promises = vendors.map((vendor) => {
      return queryClient.fetchQuery(
        getCardOfferingsVendorQuery(cardName, vendor),
      );
    });

    const results = await Promise.all(promises);
    return results.flatMap((v) => v).sort((a, b) => a.price - b.price);
  },
  enabled: !!cardName,
});
