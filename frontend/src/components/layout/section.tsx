import {
  QueryErrorResetBoundary
} from "@tanstack/react-query";
import { Suspense } from "react";
import { ErrorBoundary } from "react-error-boundary";
import { ErrorBlock } from "../placeholders/error";
import { LoadingBlock } from "../placeholders/loading";
import clsx from "clsx";

export interface SectionProps {
  header: string | React.ReactNode;
  buttonGroup?: React.ReactNode;
  margin?: boolean;
}

export const Section = ({
  header,
  buttonGroup,
  children,
  margin = false,
}: React.PropsWithChildren<SectionProps>) => {
  return (
    <div className={clsx("flex flex-col gap-4", margin && "m-4")}>
      <div className="flex flex-row items-end p-2 px-4 border-b border-secondary items-center">
        {typeof header === "string" ? (
          <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
            {header}
          </h3>
        ) : (
          header
        )}
        {buttonGroup && <div className="ml-auto">{buttonGroup}</div>}
      </div>
      <div>{children}</div>
    </div>
  );
};

export const SectionWithSuspense = ({
  header,
  buttonGroup,
  children,
}: React.PropsWithChildren<SectionProps>) => {
  return (
    <Section
      header={header}
      buttonGroup={<div className="flex flex-row space-x-1">{buttonGroup}</div>}
    >
      <QueryErrorResetBoundary>
        {({ reset }) => {
          return (
            <ErrorBoundary
              onReset={reset}
              fallbackRender={({ error, resetErrorBoundary }) => {
                return (
                  <ErrorBlock
                    error={error as unknown}
                    retry={resetErrorBoundary}
                  />
                );
              }}
            >
              <Suspense fallback={<LoadingBlock />}>{children}</Suspense>
            </ErrorBoundary>
          );
        }}
      </QueryErrorResetBoundary>
    </Section>
  );
};
