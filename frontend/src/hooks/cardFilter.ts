import { models } from "@wails/go/models";
import { useMemo } from "react";

export const useCardFilter = (
  filter: string,
  offerings: models.Offering[] | undefined,
) => {
  const filteredCards = useMemo(() => {
    if (!filter) {
      return offerings;
    }
    return offerings?.filter((offering) => {
      return Object.values(offering).find((val) =>
        String(val).toLowerCase().includes(filter.toLowerCase()),
      );
    });
  }, [offerings, filter]);

  return filteredCards;
};
