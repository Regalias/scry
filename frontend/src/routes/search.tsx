import { OfferingCard } from "@/components/cards/offeringCard";
import { EditStringForm } from "@/components/forms/editString";
import { Section } from "@/components/layout/section";
import { EmptyBlock } from "@/components/placeholders/empty";
import { ErrorBlock } from "@/components/placeholders/error";
import { LoadingBlock } from "@/components/placeholders/loading";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { BrowserOpenURL } from "@wails/runtime/runtime";
import { useState } from "react";
import { getAllCardOfferings } from "../query/ipc";

export const Route = createFileRoute("/search")({
  component: Index,
});

function Index() {
  const [submittedName, setSubmittedName] = useState("");

  const { data, isLoading, isFetching, error, refetch } = useQuery(
    getAllCardOfferings(submittedName),
  );

  return (
    <div className="p-2 container mx-auto">
      <Section header="Ad-hoc Search">
        <div className="">
          <EditStringForm
            defaultValue=""
            label="Card name"
            onSubmit={setSubmittedName}
            buttonText="Search"
          />
        </div>
      </Section>

      <div className="">
        {isLoading ? (
          <LoadingBlock />
        ) : error ? (
          <ErrorBlock
            error={error}
            isFetching={isFetching}
            retry={void refetch}
          />
        ) : data?.length ? (
          <div className="flex flex-row flex-wrap">
            {data.map((offering) => {
              return (
                <OfferingCard
                  offering={offering}
                  key={offering.productUri}
                  onSelect={() => {
                    BrowserOpenURL(offering.productUri);
                  }}
                />
              );
            })}
          </div>
        ) : submittedName ? (
          <EmptyBlock />
        ) : (
          <EmptyBlock text="Search for something" />
        )}
      </div>
    </div>
  );
}
